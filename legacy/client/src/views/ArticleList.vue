<template>
  <div class="fillcontain">
    <div>
      <el-form :inline="true" ref="add_data" :model="search_data">
        <!-- 篩選 -->
        <el-form-item label="按照時間篩選">
          <el-date-picker v-model="search_data.startTime" type="datetime" placeholder="选择開始时间"></el-date-picker> --
          <el-date-picker v-model="search_data.endTime" type="datetime" placeholder="选择結束时间"></el-date-picker>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="small" icon="search" @click="handleSearch()">篩選</el-button>
        </el-form-item>
        <el-form-item class="btnRight">
          <el-button
          type="primary"
          size="small"
          icon="el-icon-plus"
          v-if="user.identity === 'admin'"
          @click="handleAdd()">新增文章</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="table-container">
      <el-table v-if="tableData.length > 0" :data="tableData" style="width: 100%" max-height="1000">
        <el-table-column type="index" label="序号" align="center" width="100"></el-table-column>
        <el-table-column prop="date" label="創建時間" width="320" align="center">
          <template slot-scope="scope">
            <i class="el-icon-time"></i>
            <span style="margin-left: 10px">{{ scope.row.date }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="文章類型" align="center" width="150"></el-table-column>
        <el-table-column prop="describe" label="標題描述" align="center" width="490"></el-table-column>
        <el-table-column prop="remark" label="備註" align="center" width="300"></el-table-column>
        <el-table-column prop="operation" align="center" label="操作" width="320" v-if="user.identity === 'admin'">
          <template slot-scope="scope">
            <el-button
              type="info"
              icon="el-icon-edit"
              size="small"
              @click="handleEdit(scope.$index, scope.row)"
            >编辑</el-button>
            <el-button
              type="danger"
              icon="el-icon-delete"
              size="small"
              @click="handleDelete(scope.$index, scope.row)"
            >删除</el-button>
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
import ArticleDialog from "@/components/ArticleDialog.vue";

export default {
  name: "fundlist",
  components: {
    "article-dialog": ArticleDialog
  },
  data() {
    return {
      search_data: {
        startTime: "",
        endTime: ""
      },
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
      formData: {
        type: "",
        describe: "",
        income: "",
        expend: "",
        cash: "",
        remark: "",
        id: ""
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
        .get("https://www.seicho-no-ie.org.tw/php/GetContent.php")
        .then(res => {
          this.allTableData = res.data;
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
    handleEdit(index, row) {
      //編輯
      this.dialog = {
        show: true,
        title: "修改資金信息",
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
      this.$router.push('/edit');
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
