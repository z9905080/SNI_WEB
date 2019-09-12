<template>
  <div class="marquee-dialog">
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
          <!-- <el-form-item prop="page_group_id" label="序號:">
            <label type="page_group_id" v-text="formData.page_group_id"></label>
          </el-form-item>-->
          <el-form-item prop="text" label="文字內容:">
            <el-input type="text" v-model="formData.text" @keydown.enter.native="onSubmit('form')" :style="style"></el-input>
            <el-color-picker v-model="formData.color"></el-color-picker>
          </el-form-item>
          <!-- <el-form-item prop="remark" label="備註:">
            <el-input type="remark" v-model="formData.remark"></el-input>
          </el-form-item>-->
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
        text: [{ required: true, message: "內容不能为空！", trigger: "blur" }]
      }
    };
  },
  props: {
    dialog: Object,
    formData: Object
  },
  computed: {
      style () {
        return 'color: ' + this.formData.color;
      }
    },
  methods: {
    onSubmit(form) {
      const apiType =
        this.dialog.option === "edit" ? "EditMarquee" : "AddMarquee";

      this.$refs[form].validate(valid => {
        if (valid) {
          console.log(this.formData);

          this.$axios
            .post(
              `https://sniweb.shouting.feedia.co/php/${apiType}.php?sid=${window.$cookies.get(
                "sid"
              )}`,
              JSON.stringify(this.formData)
            )
            .then(res => {
              console.log(res);

              if (res.data.status === "Y") {
                //添加成功
                this.$message({
                  message: res.data.message,
                  type: "success"
                });
                //隱藏彈出視窗
                this.dialog.show = false;
                this.$emit("update");
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
            });
        }
      });
    }
  }
};
</script>

<style scoped>
</style>