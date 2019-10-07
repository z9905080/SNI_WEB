<template>
  <div class="body">
    <div class="upload">
      <el-upload
        class="upload-demo"
        ref="upload"
        :action="uploadUrl"
        :on-preview="handlePreview"
        :on-remove="handleRemove"
        :file-list="fileList"
        :auto-upload="false"
        :on-success="uploadSuccess"
      >
        <el-button slot="trigger" size="medium" type="primary">從電腦尋找圖片</el-button>
        <el-button
          style="margin-left: 10px;"
          size="medium"
          type="success"
          @click="submitUpload"
        >上傳到伺服器</el-button>
        <el-button style="float: right;" size="medium" type="success" @click="sendSelectImage">送出選取圖片</el-button>
        <div slot="tip" class="el-upload__tip">只能上傳jpg/png/gif檔案，且不超過2MB，已上傳圖片會在下方顯示</div>
      </el-upload>
    </div>
    <div class="transfer"></div>
    <div class="slider">
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
            <!-- <div class="demo-image__preview">
              <el-image style="width: 100px; height: 100px" :src="image" :preview-src-list="images"></el-image>
            </div> -->
            <div style="padding: 10px;">
              <!-- <span class="imageText">{{ `${checkboxText}${i + 1}` }}</span> -->
              <div class="bottom clearfix">
                <el-checkbox-group v-model="checkboxGroup" size="medium" class="imageText">
                  <el-checkbox :label="`${checkboxText}${i + 1}`" circle></el-checkbox>
                </el-checkbox-group>
                <!-- <el-button type="success" icon="el-icon-check" size="small" circle></el-button> -->
                <el-button
                  type="danger"
                  icon="el-icon-delete"
                  plain
                  size="medium"
                  class="deleteBtn"
                  @click="deleteImage(i)"
                >刪除</el-button>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
      <vueGallerySlideshow :images="images" :index="index" @close="index = null"></vueGallerySlideshow>
    </div>
    <!-- 彈出視窗 -->
    <slider-dialog :dialog="dialog" :formData="formData" @update="getImages"></slider-dialog>
  </div>
</template>

<script>
import VueGallerySlideshow from "vue-gallery-slideshow";
import SliderDialog from "@/components/SliderDialog.vue";

export default {
  components: {
    "vueGallerySlideshow": VueGallerySlideshow,
    "slider-dialog": SliderDialog
  },
  data() {
    return {
      uploadUrl: `https://sniweb.shouting.feedia.co/php/PictureUpload.php?sid=${window.$cookies.get(
        "sid"
      )}&r=${new Date().getTime()}`,
      url: `https://sniweb.shouting.feedia.co/php`,
      checkboxText: "圖片",
      index: null,
      fileList: [],
      images: [],
      resouseImages: [],
      checkboxGroup: [],
      dialog: {
        show: false,
        title: ""
      },
      formData: {
        link_url: "",
        images: ""
      },
    };
  },
  created() {
    this.getImages();
  },
  methods: {
    getImages() {
      //獲取全部圖片
      this.$axios
        .get(
          `https://sniweb.shouting.feedia.co/php/GetGallery.php?sid=${window.$cookies.get(
            "sid"
          )}&r=${new Date().getTime()}`
        )
        .then(res => {
          if (res.data.status && res.data.status === 'N') {
            this.$message({
              message: `${res.data.message}`,
              type: "error",
            });
            return;
          }

          const resImages = [];
          this.resouseImages = [];
          this.checkboxGroup = [];
          const imageDataList = res.data.split(" <BR/>");
          imageDataList.forEach((image, index) => {
            if (!image) {
              resImages.splice(index, 1);
            } else {
              this.resouseImages.push(`${image}`);
              resImages.push(`${this.url}${image}`);
            }
          });
          this.images = resImages;
        })
        .catch(err => console.log(err));
    },
    uploadSuccess(file, fileList) {
      this.getImages();
    },
    sendSelectImage(index) {
      const imageSelectData = this.getSelectImage();
      if (!imageSelectData.length) {
        this.$message({
          message: '請選取圖片',
          type: "warning",
        });
        return;
      }

      //編輯
      this.dialog = {
        show: true,
        title: "新增跑馬燈網址"
      };

      this.formData = {
        link_url: "",
        images: imageSelectData
      };

      // imageSelectData.forEach((image, index) => {
      //   const addCarousels = {
      //     image_url: image.replace("\u005c", "/"),
      //     link_url: this.url
      //   };
      //   // 送出選取圖片
      //   this.$axios
      //     .post(
      //       `https://sniweb.shouting.feedia.co/php/AddCarousel.php?sid=${window.$cookies.get(
      //         "sid"
      //       )}`,
      //       JSON.stringify(addCarousels)
      //     )
      //     .then(res => {
      //       if (res.data.status === "Y") {
      //         //添加成功
      //         this.$message({
      //           message: res.data.message,
      //           type: "success"
      //         });
      //         this.$emit("update");
      //         this.$router.push("/carousel");
      //       } else {
      //         //添加失敗
      //         this.$message({
      //           message: res.data.message,
      //           type: "error"
      //         });
      //       }
      //     })
      //     .catch(err => console.log(err));
      // });
    },
    getSelectImage() {
      const imageSelectList = [];
      this.checkboxGroup.forEach(obj => {
        const select = Number(obj.replace(`${this.checkboxText}`, "")) - 1;
        const imagePath = this.resouseImages[select].replace(`${this.url}`, "");
        imageSelectList.push(imagePath);
      });
      return imageSelectList;
    },
    deleteImage(index) {
      if (index === undefined || null) {
        return;
      }
      const image = this.resouseImages[index].replace("\u005c", "/");
      const deleteFile = {
        file_path: `${image}`
      };

      //送出刪除圖片
      this.$axios
        .post(
          `https://sniweb.shouting.feedia.co/php/PictureDelete.php?sid=${window.$cookies.get(
            "sid"
          )}`,
          JSON.stringify(deleteFile)
        )
        .then(res => {
          this.getImages();
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
    },
    submitUpload() {
      this.$refs.upload.submit();
    },
    handleRemove(file, fileList) {
    },
    handlePreview(file) {
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
.bottom {
  margin-top: 3%;
  margin-left: 25%;
  line-height: 12px;
}
.deleteBtn {
  margin-top: 8%;
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