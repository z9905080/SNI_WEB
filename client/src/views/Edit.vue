<template>
  <div>
    <div class="form">
      <el-form ref="form" label-width="120px" style="margin:30px;width:auto;">
        <h2 style="font-weight:bold; margin-bottom: 10px; font-size:20px">內文標題:</h2>
        <el-input
          class="title"
          v-model="page_content.page_name"
          placeholder="請輸入標題"
          required="true"
        ></el-input>
        <!-- 给editor加key是因为给tinymce keep-alive以后组件切换时tinymce编辑器会显示异常，
        在activated钩子里改变key的值可以让编辑器重新创建-->
        <editor
          id="tinymceEditor"
          :init="tinymceInit"
          v-model="page_content.html_context"
          :key="tinymceFlag"
        ></editor>
        <el-upload
          class="upload-demo"
          ref="imageUpload"
          :action="imageUrl"
          :on-preview="handlePreview"
          :on-remove="handleRemove"
          :file-list="fileList"
          :auto-upload="true"
          multiple
          :on-success="insertImage"
          style="display:none"
        ></el-upload>
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
import "tinymce/plugins/image";
import "tinymce/plugins/media";
import "tinymce/plugins/link";
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
      url: "",
      imageUrl: `http://www.seicho-no-ie.org.tw/php/PictureUpload.php?sid=${window.$cookies.get(
        "sid"
      )}`,
      tinymceFlag: 1,
      id: "",
      tinymceInit: {},
      page_content: {
        html_context: "",
        id: "",
        page_group_id: "",
        page_name: ""
      },
      fileList: []
    };
  },
  methods: {
    // 插入图片至编辑器
    insertImage(res, file) {
      let src = `http://www.seicho-no-ie.org.tw/php/${file.response.url}`; // 图片存储地址
      tinymce.execCommand("mceInsertContent", false, `<img src=${src}>`);
    },
    getPageContent() {
      // 拿到網址的id
      this.url = location.href;
      //判斷是否為edit
      if (this.url.indexOf("/navpageedit/edit/") !== -1) {
        //之後去分割字串把分割後的字串放進陣列中
        const ary1 = this.url.split("/navpageedit/edit/");
        this.id = ary1[ary1.length - 1];
        //獲取數據
        this.$axios
          .post(
            `https://www.seicho-no-ie.org.tw/php/GetContent.php?page_id=${this.id}`
          )
          .then(res => {
            this.page_content = res.data;
            this.page_content.page_id = res.data.id;
          })
          .catch(err => console.log(err));
      } else {
        this.page_content.page_group_id = this.$route.params.page_group_id;
      }
    },
    onSubmit() {
      const isEdit = this.url.indexOf("/navpageedit/edit/") !== -1;
      const apiType = isEdit ? "EditContent" : "AddContent";

      fetch(`https://www.seicho-no-ie.org.tw/php/${apiType}.php?sid=${window.$cookies.get("sid")}`, {
        body: JSON.stringify(this.page_content), // must match 'Content-Type' header
        cache: "no-cache", // *default, no-cache, reload, force-cache, only-if-cached
        headers: {
          "content-type": "application/json"
        },
        method: "POST", // *GET, POST, PUT, DELETE, etc.
        redirect: "follow", // manual, *follow, error
        referrer: "no-referrer" // *client, no-referrer
      }).then(response => {
          return response.json();
        })
        .then(data=>{
            if (data.status === "Y") {
            //添加成功
            this.$message({
              message: data.message,
              type: "success"
            });
            this.$emit("update");
            this.$router.push("/navpageedit");
          } else {
            //添加失敗
            this.$message({
              message: data.message,
              type: "error"
            });
          }
        }); // 輸出成 json

      // //送出數據
      // this.$axios
      //   .post(
      //     `http://www.seicho-no-ie.org.tw/php/${apiType}.php?sid=${window.$cookies.get(
      //       "sid"
      //     )}`,
      //     JSON.stringify(this.page_content)
      //   )
      //   .then(res => {
      //     if (res.data.status === "Y") {
      //       //添加成功
      //       this.$message({
      //         message: res.data.message,
      //         type: "success"
      //       });
      //       this.$emit("update");
      //       this.$router.push("/navpageedit");
      //     } else {
      //       //添加失敗
      //       this.$message({
      //         message: res.data.message,
      //         type: "error"
      //       });
      //     }
      //   });
    },
    onCancel() {
      this.$router.push("/navpageedit");
    },
    handleRemove(file, fileList) {},
    handlePreview(file) {},
    handleExceed(files, fileList) {
      this.$message.warning(
        `当前限制选择 3 个文件，本次选择了 ${
          files.length
        } 个文件，共选择了 ${files.length + fileList.length} 个文件`
      );
    },
    beforeRemove(file, fileList) {
      return this.$confirm(`确定移除 ${file.name}？`);
    }
  },
  created() {
    this.tinymceInit = {
      skin_url: "/client/dist/tinymce/skins/ui/oxide",
      language_url: `/client/dist/tinymce/langs/zh_CN.js`,
      language: "zh_CN",
      height: document.body.offsetHeight - 300,
      browser_spellcheck: true, // 拼写检查
      branding: false, // 去水印
      // elementpath: false,  //禁用编辑器底部的状态栏
      statusbar: false, // 隐藏编辑器底部的状态栏
      paste_data_images: true, // 允许粘贴图像
      menubar: false, // 隐藏最上方menu
      // menubar: "insert", // 隐藏最上方menu
      plugins: "advlist table lists paste preview fullscreen image media link",
      toolbar:
        "fontselect fontsizeselect forecolor backcolor bold italic underline strikethrough link | alignleft aligncenter alignright alignjustify | imageUpload media quicklink h2 h3 blockquote table numlist bullist preview fullscreen",
      images_upload_url: "postAcceptor.php",
      automatic_uploads: false,
      file_browser_callback: function(field_name, url, type, win) {
        win.document.getElementById(field_name).value = "my browser value";
      },
      file_browser_callback_types: "file image media",
      /**
       * 下面方法是为tinymce添加自定义插入图片按钮
       * 借助ElementUI的Upload组件,将图片先上传到存储云上，再将图片的存储地址放入编辑器内容
       */
      setup: editor => {
        editor.ui.registry.addButton("imageUpload", {
          tooltip: "插入图片",
          icon: "image",
          onAction: () => {
            let upload = this.$refs.imageUpload.$refs["upload-inner"];
            upload.handleClick();
          }
        });
      }
    };
    this.getPageContent();
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
.title {
  margin-bottom: 2%;
}
</style>
