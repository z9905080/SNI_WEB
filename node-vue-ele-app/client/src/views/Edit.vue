<template>
  <div>
    <div class="form">
      <el-form
        ref="form"
        label-width="120px"
        style="margin:10px;width:auto;"
      >
        <!--
    给editor加key是因为给tinymce keep-alive以后组件切换时tinymce编辑器会显示异常，
    在activated钩子里改变key的值可以让编辑器重新创建
        -->
        <editor id="tinymceEditor" :init="tinymceInit" v-model="content" :key="tinymceFlag"></editor>
        <el-form-item class="text_right">
          <el-button type="primary" @click="onSubmit()">提交</el-button>
          <el-button @click="onCancel()">取消</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>
<script>
import tinymce from "tinymce/tinymce";
import "tinymce/themes/silver/theme";
import Editor from "@tinymce/tinymce-vue";

import "tinymce/plugins/textcolor";
import "tinymce/plugins/advlist";
import "tinymce/plugins/table";
import "tinymce/plugins/lists";
import "tinymce/plugins/paste";
import "tinymce/plugins/preview";
import "tinymce/plugins/fullscreen";

export default {
  name: "TinymceEditor",
  components: {
    editor: Editor
  },
  data() {
    return {
      tinymceFlag: 1,
      tinymceInit: {},
      content: ""
    };
  },
  methods: {
    // 插入图片至编辑器
    insertImage(res, file) {
      let src = ""; // 图片存储地址
      tinymce.execCommand("mceInsertContent", false, `<img src=${src}>`);
    },
    onSubmit() {
      this.$refs[form].validate(valid => {
        if (valid) {
          // const url = this.dialog.option === 'add' ? 'add' : `edit/${this.formData.id}`;
          this.$axios.post(`/api/acticles/${url}`, this.formData).then(res => {
            //添加成功
            this.$message({
              // message: this.dialog.option === 'add' ? "資料添加成功" : "資料編輯成功",
              type: "success"
            });
            this.$emit("update");
            this.$router.push('/articlelist');
          });
        }
      });
    },
    onCancel() {
      this.$router.push('/articlelist');
    }
  },
  created() {
    const that = this;
    this.tinymceInit = {
      skin_url: "/tinymce/skins/ui/oxide",
      language_url: `/tinymce/langs/zh_CN.js`,
      language: "zh_CN",
      height: document.body.offsetHeight - 300,
      browser_spellcheck: true, // 拼写检查
      branding: false, // 去水印
      // elementpath: false,  //禁用编辑器底部的状态栏
      statusbar: false, // 隐藏编辑器底部的状态栏
      paste_data_images: true, // 允许粘贴图像
      menubar: false, // 隐藏最上方menu
      plugins: "advlist table lists paste preview fullscreen",
      toolbar:
        "fontselect fontsizeselect forecolor backcolor bold italic underline strikethrough | alignleft aligncenter alignright alignjustify | imageUpload quicklink h2 h3 blockquote table numlist bullist preview fullscreen",
      /**
       * 下面方法是为tinymce添加自定义插入图片按钮
       * 借助iview的Upload组件,将图片先上传到存储云上，再将图片的存储地址放入编辑器内容
       */
      setup: editor => {
        editor.ui.registry.addButton("imageUpload", {
          // text: '插入图片',
          tooltip: "插入图片",
          icon: "image",
          onAction: () => {}
        });
      }
    };
  },
  activated() {
    this.tinymceFlag++;
  },
  mounted() {}
};
</script>

<style scoped>
.text_right {
  float: right;
  margin-top: 20px;
}
</style>
