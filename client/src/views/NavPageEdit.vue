<template>
  <div class="fillcontain">
    <el-button
      type="warning"
      icon="el-icon-edit"
      size="medium"
      class="btnSortSend"
      v-if="tableDataChangeSort"
      @click="handleSortSend()"
    >確定排序</el-button>
    <el-button
      type="danger"
      icon="el-icon-delete"
      size="medium"
      class="btnSortCanel"
      v-if="tableDataChangeSort"
      @click="handleSortCanel()"
    >取消排序</el-button>
    <el-button
      class="btnRight"
      type="primary"
      size="medium"
      icon="el-icon-plus"
      @click="handleAdd()"
    >添加頁籤</el-button>
    <el-table-draggable>
      {{checkNavSort}}
      <el-table :data="tableData" style="width: 100%" class="nav">
        <el-table-column type="expand">
          <template slot-scope="props">
            <div class="contentTable">
              <el-button
                class="btnRight"
                type="success"
                icon="el-icon-plus"
                size="medium"
                @click="handleArticleAdd(props.row.page_group_id)"
              >新增內文</el-button>
              <el-table-draggable>
                {{dragContent(props.row.page_content, props.row.page_group_id)}}
                <el-table
                  :data="props.row.page_content"
                  max-height="1000"
                  style="width: 100%"
                  align="center"
                  class="content"
                >
                  <el-table-column prop="page_name" label="分頁名稱" align="center" style="width: 30%"></el-table-column>
                  <el-table-column label="備註 " prop="remark" align="center" style="width: 50%"></el-table-column>
                  <el-table-column label="操作" prop="operation" align="center" style="width: 20%">
                    <template slot-scope="scope">
                      <el-button
                        type="primary"
                        icon="el-icon-edit"
                        size="medium"
                        @click="handleArticleEdit(scope.$index, scope.row.page_content_id, scope.row.page_name)"
                      >內文编輯</el-button>
                      <el-button
                        type="danger"
                        icon="el-icon-delete"
                        size="medium"
                        @click="handleDelete(scope.$index, scope.row)"
                      >删除</el-button>
                    </template>
                  </el-table-column>
                </el-table>
              </el-table-draggable>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="群組名稱" prop="group_name"></el-table-column>
        <el-table-column label="備註" prop="remark"></el-table-column>
        <el-table-column label="操作" prop="operation">
          <template slot-scope="scope">
            <el-button
              type="info"
              icon="el-icon-edit"
              size="medium"
              @click="handleEdit(scope.$index, scope.row.page_group_id, scope.row.group_name)"
            >编輯</el-button>
            <el-button
              type="danger"
              icon="el-icon-delete"
              size="medium"
              @click="handleDelete(scope.$index, scope.row)"
            >删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-table-draggable>
    <!-- 分頁 -->
    <el-row>
      <el-col :span="24">
        <div class="pagination">
          <el-pagination
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
            :current-page.sync="paginations.page_index"
            :page-sizes="paginations.page_sizes"
            :page-size="paginations.page_size"
            :layout="paginations.layout"
            :total="paginations.total"
          ></el-pagination>
        </div>
      </el-col>
    </el-row>

    <!-- 彈出視窗 -->
    <navpage-dialog :dialog="dialog" :formData="formData" @update="getProfile"></navpage-dialog>
  </div>
</template>

<script>
import NavPageDialog from "@/components/NavPageDialog.vue";
import ElTableDraggable from "element-ui-el-table-draggable";
window['Lodash'] = require('lodash');

export default {
  name: "navpagelist",
  components: {
    "navpage-dialog": NavPageDialog,
    ElTableDraggable
  },
  data() {
    return {
      filterTableData: [],
      paginations: {
        page_index: 1, //當前位於哪頁
        total: 0, //總數
        page_size: 15, //一頁顯示多少條
        page_sizes: [15, 20, 25, 30], //每頁顯示多少條
        layout: "total, sizes, prev, pager, next, jumper" //翻頁屬性
      },
      tableData: [],
      allTableData: [],
      backupTableData: [],
      allContentData: {},
      navSort: [],
      contentSort: {},
      tableDataChangeSort: false,
      formData: {
        page_group_id: "",
        group_name: "",
        page_content: []
      },
      dialog: {
        show: false,
        title: "",
        option: "edit"
      },
      editorData: {
        page_id: "",
        page_name: ""
      },
      editor: {
        show: false,
        option: "edit"
      }
    };
  },
  computed: {
    checkNavSort() {
      if (this.navSort.length) {
        this.navSort = [];
      }
      let isChange = false;
      for (let i = 0; i < this.tableData.length; i++) {
        let newId = this.tableData[i].page_group_id;
        let oldId = this.backupTableData[i].page_group_id;
        if (newId !== oldId) {
          isChange = true;
        }
        this.navSort.push(newId);
      }
      if (isChange) {
        this.tableDataChangeSort = true;
      } else {
        this.tableDataChangeSort = false;
        this.navSort = [];
      }
    },
    user() {
      return this.$store.getters.user;
    }
  },
  created() {
    this.getProfile();
  },
  methods: {
    dragContent(tableData, navId) {
      if (this.contentSort[navId].length) {
        this.contentSort[navId] = [];
      }
      let isChange = false;
      let oldTableData = this.allContentData[navId];
      for (let i = 0; i < tableData.length; i++) {
        let newId = tableData[i].page_content_id;
        let oldId = oldTableData[i].page_content_id;
        if (newId !== oldId) {
          isChange = true;
        }
        this.contentSort[navId].push(newId);
      }

      if (this.contentSort[navId].length) {
        this.tableDataChangeSort = true;
      } else {
        this.tableDataChangeSort = this.navSort.length ? true : false;
        this.contentSort[navId] = [];
      }
    },
    getProfile() {
      //獲取數據
      this.$axios
        .get(
          `https://www.seicho-no-ie.org.tw/php/GetNav.php?&r=${new Date().getTime()}`
        )
        .then(res => {
          this.allTableData = res.data;
          this.backupTableData = Lodash.cloneDeep(this.allTableData);
          for (let i = 0; i < this.backupTableData.length; i++) {
            this.allContentData[this.backupTableData[i].page_group_id] = this.backupTableData[i].page_content;
            this.contentSort[this.backupTableData[i].page_group_id] = [];
          }

          this.filterTableData = res.data;
          // 設置分頁數據
          this.setPaginations();
        })
        .catch(err => console.log(err));
    },
    setPaginations() {
      // 分頁屬性設置
      this.paginations.total = this.allTableData.length;
      this.paginations.page_index = 1;
      this.paginations.page_size = 15;
      // 設置默認的分頁數據
      this.tableData = this.allTableData.filter((item, index) => {
        return index < this.paginations.page_size;
      });
    },
    handleEdit(index, id, name) {
      //編輯
      this.dialog = {
        show: true,
        title: "修改頁籤",
        option: "edit"
      };

      this.formData = {
        page_group_id: id,
        page_group_name: name
      };
    },
    handleDelete(index, row) {
      const deleteType = row.page_group_id ? "DeleteNav" : "DeleteContent";
      let deleteData =
        deleteType === "DeleteNav"
          ? { page_group_id: row.page_group_id }
          : { page_id: row.page_content_id };
      //刪除
      this.$axios
        .post(
          `https://www.seicho-no-ie.org.tw/php/${deleteType}.php?sid=${window.$cookies.get(
            "sid"
          )}`,
          JSON.stringify(deleteData)
        )
        .then(res => {
          if (res.data.status === "Y") {
            this.$message({
              message: res.data.message,
              type: "warning"
            });

            this.getProfile();
          } else {
            this.$message({
              message: res.data.message,
              type: "error"
            });
          }
        });
    },
    handleAdd() {
      //新增
      this.dialog = {
        show: true,
        title: "新增頁籤",
        option: "add"
      };

      this.formData = {
        page_group_id: "",
        page_group_name: ""
      };
    },
    handleSizeChange(page_size) {
      //切換size
      this.paginations.page_index = 1;
      this.paginations.page_size = page_size;
      this.tableData = this.allTableData.filter((item, index) => {
        return index < page_size;
      });
    },
    handleCurrentChange(page) {
      //獲取當前頁
      let index = this.paginations.page_size * (page - 1);
      //數據的總數
      let nums = this.paginations.page_size * page;
      //容器
      let tables = [];

      for (let i = index; i < nums; i++) {
        if (this.allTableData[i]) {
          tables.push(this.allTableData[i]);
        }
        this.tableData = tables;
      }
    },
    handleSearch() {
      //篩選
      if (!this.search_data.startTime || !this.search_data.endTime) {
        this.$message({
          type: "warning",
          message: "請選擇時間區間"
        });
        this.getProfile();
        return;
      }

      const sTime = this.search_data.startTime.getTime();
      const eTime = this.search_data.endTime.getTime();

      this.allTableData = this.filterTableData.filter(item => {
        let date = new Date(item.date);
        let time = date.getTime();
        return time >= sTime && time <= eTime;
      });

      // 分頁數據的調用
      this.setPaginations();
    },
    handleArticleEdit(index, id, name) {
      //內文編輯
      this.$router.push(`/navpageedit/edit/${id}`);
    },
    handleArticleAdd(page_group_id) {
      //新增內文
      this.$router.push({
        name: "add",
        params: {
          page_group_id: page_group_id
        }
      });
    },
    handleSortSend() {
      if (this.navSort.length) {
        const navSortData = {
          page_group_sort: `[${this.navSort}]`
        };
        //送出數據
        this.$axios
          .post(
            `https://www.seicho-no-ie.org.tw/php/EditNavSort.php?sid=${window.$cookies.get(
              "sid"
            )}`,
            JSON.stringify(navSortData)
          )
          .then(res => {
            if (res.data.status === "Y") {
              //添加成功
              this.$message({
                message: res.data.message,
                type: "success"
              });
              this.getProfile();
              this.$emit("update");
            } else {
              //添加失敗
              this.$message({
                message: res.data.message,
                type: "error"
              });
            }
          });
      }

      if (Object.keys(this.contentSort).length) {
        for (let navId in this.contentSort) {
          if (this.contentSort[navId].length) {
            const contentSortData = {
              page_group_id: Number(navId),
              page_sort: `[${this.contentSort[navId]}]`
            };
            //送出數據
            this.$axios
              .post(
                `https://www.seicho-no-ie.org.tw/php/EditContentSort.php?sid=${window.$cookies.get(
                  "sid"
                )}`,
                JSON.stringify(contentSortData)
              )
              .then(res => {
                if (res.data.status === "Y") {
                  //添加成功
                  this.$message({
                    message: res.data.message,
                    type: "success"
                  });
                  this.getProfile();
                  this.$emit("update");
                } else {
                  //添加失敗
                  this.$message({
                    message: res.data.message,
                    type: "error"
                  });
                }
              });
          }
        }
      }
    },
    handleSortCanel() {
      this.getProfile();
      this.$emit("update");
    }
  }
};
</script>

<style scoped>
.fillcontain {
  width: 100%;
  height: 100%;
  padding: 16px;
  box-sizing: border-box;
}
.btnRight {
  float: right;
}
.btnSortSend {
  float: left;
  margin-left: 65%;
}
.btnSortCanel {
  float: left;
  /* margin-left: 5% */
}
.pagination {
  text-align: right;
  margin-top: 10px;
}
</style>
