<template>
  <div class="slider">
    <el-upload
      class="upload-demo"
      ref="upload"
      :action="uploadUrl"
      :on-preview="handlePreview"
      :on-remove="handleRemove"
      :file-list="fileList"
      :auto-upload="false"
    >
      <el-button slot="trigger" size="small" type="primary">從電腦尋找圖片</el-button>
      <el-button style="margin-left: 10px;" size="small" type="success" @click="submitUpload">上傳到伺服器</el-button>
      <el-button style="float: right;" size="small" type="success" @click="sendSelectImage">送出選取圖片</el-button>
      <div slot="tip" class="el-upload__tip">只能上傳jpg/png/gif檔案，且不超過2MB，已上傳圖片會在下方顯示</div>
    </el-upload>
    <el-row>
      <el-col
        :span="4"
        :offset="2"
        class="imageContainer"
        v-for="(image, i) in images"
        :key="i"
        @click="index = i"
      >
        <el-card :body-style="{ padding: '15px' }">
          <img class="image" :src="image" @click="index = i" />
          <div style="padding: 10px;">
            <div class="bottom clearfix">
              <el-checkbox-group v-model="checkboxGroup" size="small">
                <el-checkbox :label="`${checkboxText}${i + 1}`" circle></el-checkbox>
                <el-button type="danger" icon="el-icon-delete" size="small" circle @click="deleteImage(i)"></el-button>
              </el-checkbox-group>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    <vueGallerySlideshow :images="images" :index="index" @close="index = null"></vueGallerySlideshow>
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
      uploadUrl: `https://sniweb.shouting.feedia.co/php/PictureUpload.php?sid=${`${window.$cookies.get(
        "sid"
      )}`}`,
      url: `https://sniweb.shouting.feedia.co/php`,
      checkboxText: "圖片",
      index: null,
      fileList: [],
      images: [],
      checkboxGroup: []
    };
  },
  created() {
    this.getImages();
  },
  methods: {
    getImages() {
      //獲取數據
      this.$axios
        .get(
          `https://sniweb.shouting.feedia.co/php/GetGallery.php?sid=${`${window.$cookies.get(
            "sid"
          )}`}&r=${new Date().getTime()}`
        )
        .then(res => {
          const resImages = [];
          res.data.split("<BR/>").forEach((image, index) => {
            if (!image.length) {
              resImages.splice(index, 1);
            } else {
              resImages.push(`${this.url}${image}`);
            }
          });
          this.images = resImages;
        })
        .catch(err => console.log(err));
    },
    sendSelectImage() {
      const imageSelectData = this.getSelectImage();
      if (!imageSelectData.length) {
        return;
      }
      console.log(imageSelectData);

      //送出選取圖片
      // this.$axios
      //   .post(
      //     `https://sniweb.shouting.feedia.co/php/AddCarousels.php?sid=${`${window.$cookies.get(
      //       "sid"
      //     )}`}`,
      //     JSON.stringify(imageSelectData)
      //   )
      //   .then(res => {
      //     console.log(res);
      //   })
      //   .catch(err => console.log(err));
    },
    getSelectImage() {
      const imageSelectList = [];
      this.checkboxGroup.forEach(obj => {
        const select = Number(obj.replace(`${this.checkboxText}`, "")) - 1;
        const imagePath = this.images[select].replace(
          `${this.url}`,
          ""
        );
        imageSelectList.push(imagePath);
      });
      // let imageData;
      // imageSelectList.forEach(image => {
      //   imageData = imageData ? `${imageData}<BR/>${image}` : `${image}`;
      // });
      // return imageData;
      return imageSelectList;
    },
    deleteImage(index) {
      if (!index) {
        return;
      }

      // const image = this.images[index].replace(`${this.url}`, '');
      const image = this.images[index];
      console.log(index);
      console.log(image);

      //送出刪除圖片
      this.$axios
        .post(
          `https://sniweb.shouting.feedia.co/php/PictureDelete.php?sid=${`${window.$cookies.get(
            "sid"
          )}`}`,
          JSON.stringify(image)
        )
        .then(res => {
          console.log(res);
          this.getImages();
        })
        .catch(err => console.log(err));
    },
    submitUpload() {
      this.$refs.upload.submit();
    },
    handleRemove(file, fileList) {
      console.log(file, fileList);
    },
    handlePreview(file) {
      console.log(file);
    }
  }
};
//圖片路徑
// https://sniweb.shouting.feedia.co/php/+ picture/.....
</script>

<style scoped>
body {
  font-family: sans-serif;
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
.bottom {
  margin-top: 5%;
  margin-left: 15%;
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