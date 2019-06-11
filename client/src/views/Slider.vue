<template>
  <div class="slider">
    <el-upload
      class="upload-demo"
      ref="upload"
      :action="url"
      :on-preview="handlePreview"
      :on-remove="handleRemove"
      :file-list="fileList"
      :auto-upload="false"
      list-type="picture"
    >
      <el-button slot="trigger" size="small" type="primary">選取圖片</el-button>
      <el-button style="margin-left: 10px;" size="small" type="success" @click="submitUpload">上傳到伺服器</el-button>
      <div slot="tip" class="el-upload__tip">只能上傳jpg/png檔案，且不超過2MB</div>
    </el-upload>
  </div>
</template>

<script>
export default {
  data() {
    return {
      url: `https://sniweb.shouting.feedia.co/php/Upload.php?sid=${`${window.$cookies.get(
        "sid"
      )}`}`,
      fileList: []
    };
  },
  methods: {
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
// https://sniweb.shouting.feedia.co/php/+ picture/.....
</script>

<style scoped>
.slider {
  margin-top: 3%;
  margin-left: 3%;
  width: 80%;
}
</style>