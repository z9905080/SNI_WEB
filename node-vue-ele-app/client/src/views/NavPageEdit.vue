<template>
  <div class="fillcontain">
    <div class="table-container">
      <el-table v-if="tableData.length > 0" :data="tableData" style="width: 100%" max-height="1000">
        <el-table-column type="expand">
          <template>
            <div v-for="(item, index) in pageTableData" :key="index">
              <el-table height="300" style="width: 40%" align="center" v-for="(target, index) in item" :key="index">
                <el-table-column label="分頁序號" align="center" width="180">
                  {{target.page_content_id}}
                </el-table-column>
                <el-table-column label="分頁名稱" align="center" width="180">
                  {{target.page_name}}
                </el-table-column>
                <el-table-column label="操作" prop="operation" align="center" width="320">
                  <template slot-scope="scope">
                    <el-button
                      type="info"
                      icon="el-icon-edit"
                      size="small"
                      @click="handleEdit(scope.$index, scope.row)"
                    >编辑</el-button>
                    <!-- <el-button
                    type="danger"
                    icon="el-icon-delete"
                    size="small"
                    @click="handleDelete(scope.$index, scope.row)"
                    >删除</el-button>-->
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="分類序號" prop="page_group_id" align="center" width="100"></el-table-column>
        <el-table-column label="群組名稱" prop="group_name" align="center" width="300"></el-table-column>
        <el-table-column label="備註" prop="remark" align="center" width="900"></el-table-column>
        <el-table-column label="操作" prop="operation" align="center" width="350">
          <template slot-scope="scope">
            <el-button
              type="info"
              icon="el-icon-edit"
              size="small"
              @click="handleEdit(scope.$index, scope.row)"
            >编辑</el-button>
            <!-- <el-button
              type="danger"
              icon="el-icon-delete"
              size="small"
              @click="handleDelete(scope.$index, scope.row)"
            >删除</el-button>-->
          </template>
        </el-table-column>
      </el-table>
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
    </div>

    <!-- 彈出視窗 -->
    <!-- <article-dialog :dialog="dialog" :formData="formData" @update="getProfile"></article-dialog> -->
  </div>
</template>

<script>
// import ArticleDialog from "@/components/ArticleDialog.vue";

export default {
  name: "fundlist",
  components: {
    // "article-dialog": ArticleDialog
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
      pageTableData: [],
      allTableData: [],
      formData: {
        page_group_id: "",
        group_name: "",
        page_content: []
      },
      dialog: {
        show: false,
        title: "",
        option: "edit"
      }
    };
  },
  computed: {
    user() {
      return this.$store.getters.user;
    }
  },
  created() {
    this.getProfile();
  },
  methods: {
    getProfile() {
      //獲取數據
      this.$axios
        .get("https://sniweb.shouting.feedia.co/php/GetNav.php")
        .then(res => {
          this.allTableData = res.data;
          this.filterTableData = res.data;
          this.allTableData.forEach((data, index) => {
            this.pageTableData[index] = data.page_content;
          });

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
    handleEdit(index, row) {
      //編輯
      this.dialog = {
        show: true,
        title: "修改標題",
        option: "edit"
      };

      this.formData = {
        type: row.type,
        describe: row.describe,
        income: row.income,
        expend: row.expend,
        cash: row.cash,
        remark: row.remark,
        id: row._id
      };
    },
    handleDelete(index, row) {
      //刪除
      this.$axios
        .delete(`/api/acticles/delete/${row._id}`)
        .then(res => {
          this.$message("資料刪除成功");
          this.getProfile();
        })
        .catch(err => console.log(`資料刪除失敗(${err})`));
    },
    handleAdd() {
      //新增
      this.$router.push("/edit");
      // this.dialog = {
      //   show: true,
      //   title: "新增文章",
      //   option: "add"
      // };

      // this.formData = {
      //   type: "",
      //   describe: "",
      //   remark: "",
      //   id: ""
      // };
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
.pagination {
  text-align: right;
  margin-top: 10px;
}
</style>