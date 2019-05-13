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
          <el-form-item prop="page_group_id" label="分類序號:">
            <label type="page_group_id" v-text="formData.page_group_id"></label>
          </el-form-item>
          <el-form-item prop="group_name" label="群組名稱:">
            <el-input type="group_name" v-model="formData.group_name"></el-input>
          </el-form-item>
          <!-- <el-form-item prop="remark" label="備註:">
            <el-input type="remark" v-model="formData.remark"></el-input>
          </el-form-item> -->
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
        // page_group_id: [
        //   { required: true, message: "分類序號不能为空！", trigger: "blur" }
        // ],
        group_name: [
          { required: true, message: "群組名稱不能为空！", trigger: "blur" }
        ],
        page_name:[
          { required: true, message: "分頁名稱不能为空！", trigger: "blur" }
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
      this.$refs[form].validate(valid => {
        if (valid) {
        //   const url = this.dialog.option === 'add' ? 'add' : `edit/${this.formData.page_group_id}`;
          const url = `edit/${this.formData.page_group_id}`;
          this.$axios.post(`https://sniweb.shouting.feedia.co/php/${url}`, this.formData).then(res => {
            //添加成功
            this.$message({
              message: this.dialog.option === 'add' ? "資料添加成功" : "資料編輯成功",
              type: "success"
            });
            //隱藏彈出視窗
            this.dialog.show = false;
            this.$emit("update");
          });
        }
      });
    }
  }
};
</script>

<style scoped>
</style>