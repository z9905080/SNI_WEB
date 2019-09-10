<template>
<el-button type="success" icon="el-icon-check" size="small" circle @click="getCarousels()"></el-button>
</template>

<script>
import VueGallerySlideshow from "vue-gallery-slideshow";

export default {
  components: {
    vueGallerySlideshow: VueGallerySlideshow
  },
  data() {
    return {
    };
  },
  created() {
    this.getCarousels();
  },
  methods: {
    getCarousels() {
        //獲取數據
      this.$axios
        .get(
          `https://sniweb.shouting.feedia.co/php/GetCarousels.php?sid=${window.$cookies.get(
            "sid"
          )}&r=${new Date().getTime()}`
        )
        .then(res => {
            console.log(res);
            
        //   const resImages = [];
        //   this.resouseImages = [];
        //   this.checkboxGroup = [];
        //   const imageDataList = res.data.split(" <BR/>");
        //   imageDataList.forEach((image, index) => {
        //     if (!image) {
        //       resImages.splice(index, 1);
        //     } else {
        //       this.resouseImages.push(`${image}`);
        //       resImages.push(`${this.url}${image}`);
        //     }
        //   });
        //   this.images = resImages;
        })
        .catch(err => console.log(err));
    },
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
.upload {
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
.transfer-footer {
  margin-left: 20px;
  padding: 6px 5px;
}
</style>