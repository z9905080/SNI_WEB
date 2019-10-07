<template>
  <div class="slider-dialog">
    <el-dialog
      :title="dialog.title"
      :visible.sync="dialog.show"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      :modal-append-to-body="false"
    >
      <div class="form">
        <el-form
          ref="form"
          :model="formData"
          :rules="form_rules"
          label-width="120px"
          style="margin:10px;width:auto;"
        >
          <el-form-item prop="link_url" label="網址:" class="title">
            <el-input
              type="link_url"
              v-model="formData.link_url"
              @keydown.enter.native="onSubmit('form')"
            ></el-input>
          </el-form-item>
          <el-form-item class="text_right">
            <el-button type="primary" @click="onSubmit('form')">提交</el-button>
            <el-button @click="dialog.show = false">取消</el-button>
          </el-form-item>
        </el-form>
      </div>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: "dialog",
  data() {
    return {
      form_rules: {
        link_url: [{ required: true, message: "內容不能为空！", trigger: "blur" }]
      }
    };
  },
  props: {
    dialog: Object,
    formData: Object
  },
  methods: {
    onSubmit(form) {
      this.$refs[form].validate(valid => {
        if (valid) {
          this.formData.images.forEach((image, index) => {
            const addCarousels = {
              image_url: image.replace("\u005c", "/"),
              link_url: this.formData.link_url
            };
            // 送出選取圖片
            this.$axios
              .post(
                `https://sniweb.shouting.feedia.co/php/AddCarousel.php?sid=${window.$cookies.get(
                  "sid"
                )}`,
                JSON.stringify(addCarousels)
              )
              .then(res => {
                if (res.data.status === "Y") {
                  //添加成功
                  this.$message({
                    message: res.data.message,
                    type: "success"
                  });
                  //隱藏彈出視窗
                  this.dialog.show = false;
                  this.$emit("update");
                  this.$router.push("/carousel");
                } else {
                  //添加失敗
                  this.$message({
                    message: res.data.message,
                    type: "error"
                  });
                  //隱藏彈出視窗
                  this.dialog.show = false;
                  this.$emit("update");
                }
              })
              .catch(err => console.log(err));
          });
        }
      });
    }
  }
};
</script>

<style scoped>
.title {
  font-size: 18px;
}
</style>