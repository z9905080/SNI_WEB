<template>
  <div class="navpage-dialog">
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
          </el-form-item> -->
          <el-form-item prop="page_group_name" label="名稱:">
            <el-input type="page_group_name" v-model="formData.page_group_name" @keydown.enter.native="onSubmit('form')"></el-input>
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
        page_group_name: [
          { required: true, message: "名稱不能为空！", trigger: "blur" }
        ]
      }
    };
  },
  props: {
    dialog: Object,
    formData: Object
  },
  methods: {
    onSubmit(form) {
      const apiType = this.dialog.option === "edit" ? "EditNav" : "AddNav";

      this.$refs[form].validate(valid => {
        if (valid) {
          this.$axios
            .post(
              `https://sniweb.shouting.feedia.co/php/${apiType}.php?sid=${window.$cookies.get(
                "sid"
              )}`,
              JSON.stringify(this.formData)
            )
            .then(res => {
              if (res.data.status === "Y" && res.data) {
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