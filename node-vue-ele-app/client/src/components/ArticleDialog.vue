<template>
  <div class="article-dialog">
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
          <el-form-item label="文章類型:">
            <el-select v-model="formData.type" placeholder="文章類型">
              <el-option
                v-for="(formtype, index) in format_type_list"
                :key="index"
                :label="formtype"
                :value="formtype"
              ></el-option>
            </el-select>
          </el-form-item>
          <el-form-item prop="describe" label="標題描述:">
            <el-input type="describe" v-model="formData.describe"></el-input>
          </el-form-item>
          <el-form-item prop="remark" label="備註:">
            <el-input type="remark" v-model="formData.remark"></el-input>
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
      format_type_list: [
        "運動",
        "財金",
        "娛樂",
        "政治",
        "社會",
        "國際"
      ],
      form_rules: {
        describe: [
          { required: true, message: "標題描述不能为空！", trigger: "blur" }
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
          const url = this.dialog.option === 'add' ? 'add' : `edit/${this.formData.id}`;
          this.$axios.post(`/api/acticles/${url}`, this.formData).then(res => {
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