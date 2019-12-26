<template>
  <div class="login">
    <section class="form_container">
      <div class="manage_tip">
        <span class="title">生長之家後台管理系統</span>
        <el-form
          :model="loginUser"
          :rules="rules"
          ref="loginForm"
          label-width="80px"
          class="loginForm"
        >
          <el-form-item label="帳號" prop="account">
            <el-input type="text" v-model="loginUser.account" placeholder="請輸入帳號" @keydown.enter.native="submitForm('loginForm')"></el-input>
          </el-form-item>
          <el-form-item label="密碼" prop="password">
            <el-input type="password" v-model="loginUser.password" placeholder="請輸入密碼" @keydown.enter.native="submitForm('loginForm')"></el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" class="submit_btn" @click="submitForm('loginForm')">登入</el-button>
          </el-form-item>
          <div class="tiparea">
            <p>
              還沒有帳號? 現在
              <router-link to="/register">註冊</router-link>
            </p>
          </div>
        </el-form>
      </div>
    </section>
  </div>
</template>

<script>
export default {
  name: "login",
  components: {},
  data() {
    return {
      loginUser: {
        account: "", //test123
        password: "" //test123
      },
      rules: {
        password: [
          {
            required: true,
            message: "密碼不可為空",
            trigger: "blur"
          }
        ]
      }
    };
  },
  methods: {
    submitForm(formName) {
      this.$refs[formName].validate(valid => {
        if (valid) {
          this.$axios
            .post(
              "http://www.seicho-no-ie.org.tw/php/Login.php",
              JSON.stringify(this.loginUser)
            )
            .then(res => {
              if (res.data) {
                // get token
                const { token, expire_time } = res.data;
                
                //利用regular expression  將'-' 代換成 '/'
                const newDate = new Date(Date.parse(expire_time .replace(/-/g, "/")));
                // set cookie
                this.$cookies.set('sid', token, newDate);
                this.$router.push("/index");
                this.$message({
                message: '歡迎進入系統',
                type: "success"
                });
              } else {
                this.$message({
                message: '帳號或密碼輸入錯誤，請重新輸入',
                type: "error"
              });
              }
            });
        }
      });
    },
    isEmpty(value) {
      return (
        value === undefined ||
        value === null ||
        (typeof value === "object" && Object.keys(value).length === 0) ||
        (typeof value === "string" && value.trim().length === 0)
      );
    }
  }
};
</script>

<style scoped>
.login {
  position: relative;
  width: 100%;
  height: 100%;
  background: url(../assets/bg.jpg) no-repeat center center;
  background-size: 100% 100%;
}
.form_container {
  width: 400px;
  height: 250px;
  position: absolute;
  top: 20%;
  left: 34%;
  padding: 25px;
  border-radius: 5px;
  text-align: center;
}
.form_container .manage_tip .title {
  font-family: "Microsoft YaHei";
  font-weight: bold;
  font-size: 30px;
  color: #fff;
}
.loginForm {
  margin-top: 20px;
  background-color: #fff;
  padding: 20px 40px 20px 20px;
  border-radius: 5px;
  box-shadow: 0px 5px 10px #ccc;
}
.submit_btn {
  width: 100%;
}
.tiparea {
  text-align: right;
  font-size: 12px;
  color: #333;
}
.tiparea p a {
  color: #409eff;
}
</style>
