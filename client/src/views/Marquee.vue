<template>
  <div class="imageTable">
    <div class="upload">
      <el-button style="float: right;" size="small" type="success" @click="handleAdd">新增跑馬燈</el-button>
    </div>
    <el-table :data="tableData" style="width: 100%">
      <el-table-column label="編號" width="200" align="center">
        <template slot-scope="scope">
          <p>{{scope.row.index + 1}}</p>
          <span style="margin-left: 10px">{{ scope.row.date }}</span>
        </template>
      </el-table-column>
      <el-table-column label="跑馬燈內容" width="1000" align="center">
        <template slot-scope="scope">
          <p>{{ scope.row.text}}</p>
        </template>
      </el-table-column>
      <el-table-column label="操作">
        <template slot-scope="scope">
          <el-button size="mini" @click="handleEdit(scope.$index, scope.row)">编辑</el-button>
          <el-button size="mini" type="danger" @click="handleDelete(scope.$index, scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
     <!-- 彈出視窗 -->
    <marquee-dialog :dialog="dialog" :formData="formData" @update="getMarquees"></marquee-dialog>
  </div>
</template>

<script>
import MarqueeDialog from "@/components/MarqueeDialog.vue";

export default {
  components: {
    "marquee-dialog": MarqueeDialog
  },
  data() {
    return {
      tableData: [],
      dialog: {
        show: false,
        title: "",
        option: "edit"
      },
      formData: {
        marquee_id: "",
        text: "",
        color: ""
      },
    };
  },
  created() {
    this.getMarquees();
  },
  methods: {
    getMarquees() {
      this.tableData = [];
      //獲取數據
      this.$axios
        .get(
          `https://sniweb.shouting.feedia.co/php/GetMarquees.php?sid=${window.$cookies.get(
            "sid"
          )}&r=${new Date().getTime()}`
        )
        .then(res => {
          res.data.data.forEach((image, index) => {
            image["index"] = index;
            this.tableData.push(image);
          });
        })
        .catch(err => console.log(err));
    },
    handleAdd() {
       //編輯
      this.dialog = {
        show: true,
        title: "新增跑馬燈",
        option: "add"
      };

      this.formData = {
        text: "",
        color: ""
      };
    },
    handleEdit(index, row) {
      //編輯
      this.dialog = {
        show: true,
        title: "編輯跑馬燈",
        option: "edit"
      };

      this.formData = {
        marquee_id: row.id,
        text: row.text,
        color: row.color
      };
    },
    handleDelete(index, row) {
      const deleteData = {
        marquee_id: row.id
      };
      //刪除跑馬燈
      this.$axios
        .post(
          `https://sniweb.shouting.feedia.co/php/DeleteMarquee.php?sid=${window.$cookies.get(
            "sid"
          )}`,
          JSON.stringify(deleteData)
        )
        .then(res => {
          this.getMarquees();
          if (res.data.status === "Y") {
            //添加成功
            this.$message({
              message: res.data.message,
              type: "success"
            });
            this.$emit("update");
          } else {
            //添加失敗
            this.$message({
              message: res.data.message,
              type: "error"
            });
          }
        })
        .catch(err => console.log(err));
    }
  }
};
//圖片路徑
// https://sniweb.shouting.feedia.co/php/+ picture/.....
</script>

<style scoped>
.body {
  font-family: sans-serif;
  width: 100%;
}
.imageTable {
  margin-top: 3%;
  margin-left: 3%;
  width: 80%;
}
.slider {
  margin-top: 3%;
  margin-left: 3%;
  width: 80%;
}
.imageContainer {
  padding-top: 5%;
}
.image {
  width: 130px;
  height: 130px;
  background-size: cover;
  cursor: pointer;
  margin: 15%;
  border-radius: 3px;
  border: 2px solid lightgray;
  object-fit: contain;
  display: block;
}
.imageText {
  margin-left: 20%;
}
.bottom {
  margin-top: 5%;
  margin-left: 29%;
  line-height: 12px;
}
.clearfix:before,
.clearfix:after {
  display: table;
  content: "";
}
.clearfix:after {
  clear: both;
}
</style>