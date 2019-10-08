<template>
  <div class="imageTable">
    <el-table :data="tableData" style="width: 100%">
      <el-table-column label="編號" width="300" align="center">
        <template slot-scope="scope">
          <p>{{scope.row.index + 1}}</p>
          <span style="margin-left: 10px">{{ scope.row.date }}</span>
        </template>
      </el-table-column>
      <el-table-column label="縮圖" width="500" align="center">
        <template slot-scope="scope">
          <el-popover trigger="hover" placement="top">
            <!-- <p>姓名: {{ scope.row.name }}</p>
            <p>住址: {{ scope.row.address }}</p>-->
            <!-- <p>住址: {{ scope.row.imageUrl }}</p> -->
            <!-- <p>住址: {{ scope.row }}</p> -->
            <!-- <img class="image" :src="scope.row.imageUrl" @click="index = i" /> -->
            <div class="demo-image__preview">
              <el-image
                style="width: 200px; height: 200px"
                :src="scope.row.imageUrl"
                :preview-src-list="images"
              ></el-image>
            </div>
            <div slot="reference" class="name-wrapper">
              <el-tag size="medium" class="el-icon-picture">{{ scope.row.name }}</el-tag>
            </div>
          </el-popover>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center">
        <template slot-scope="scope">
          <el-button size="medium" @click="handleEdit(scope.$index, scope.row)">编辑</el-button>
          <el-button size="medium" type="danger" @click="handleDelete(scope.$index, scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <!-- <vueGallerySlideshow :images="images" :index="index" @close="index = null"></vueGallerySlideshow> -->
  </div>
</template>

<script>
import VueGallerySlideshow from "vue-gallery-slideshow";

export default {
  components: {
    vueGallerySlideshow: VueGallerySlideshow
  },
  data() {
    return {
      tableData: [],
      images: [],
      url: "https://sniweb.shouting.feedia.co/php"
    };
  },
  created() {
    this.getCarousels();
  },
  methods: {
    getCarousels() {
      this.tableData = [];
      this.images = [];
      //獲取數據
      this.$axios
        .get(
          `https://sniweb.shouting.feedia.co/php/GetCarousels.php?sid=${window.$cookies.get(
            "sid"
          )}&r=${new Date().getTime()}`
        )
        .then(res => {
          res.data.data.forEach((image, index) => {
            image["imageUrl"] = `${this.url}${image.image}`;
            image["index"] = index;
            this.images.push(image["imageUrl"]);
            this.tableData.push(image);
          });
        })
        .catch(err => console.log(err));
    },
    handleEdit(index, row) {
      localStorage.link_url = `${row.url}`;
      this.$router.push(`/carousel/selectcarousel/${row.id}`);
    },
    handleDelete(index, row) {
      const deleteData = {
        carousel_id: row.id
      };
      //送出刪除圖片
      this.$axios
        .post(
          `https://sniweb.shouting.feedia.co/php/DeleteCarousel.php?sid=${window.$cookies.get(
            "sid"
          )}`,
          JSON.stringify(deleteData)
        )
        .then(res => {
          this.getCarousels();
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