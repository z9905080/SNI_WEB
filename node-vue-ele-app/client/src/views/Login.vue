<template>
  <div class="login">
    <section class="form_container">
      <div class="manage_tip">
        <span class="title">生命之家後台管理系統</span>
        <el-form
          :model="loginUser"
          :rules="rules"
          ref="loginForm"
          label-width="80px"
          class="loginForm"
        >
          <el-form-item label="信箱" prop="email">
            <el-input type="text" v-model="loginUser.email" placeholder="請輸入email"></el-input>
          </el-form-item>
          <el-form-item label="密碼" prop="password">
            <el-input type="password" v-model="loginUser.password" placeholder="請輸入密碼"></el-input>
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
import jwt_decode from 'jwt-decode'

export default {
  name: "login",
  components: {},
  data() {
    return {
      loginUser: {
        email: "",
        password: ""
      },
      rules: {
        email: [
          {
            type: "email",
            required: true,
            message: "信箱格式不正確",
            trigger: "blur"
          }
        ],
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
          this.$axios.post("/api/users/login", this.loginUser).then(res => {
            // token
            const { token } = res.data;

            // 儲存至localStorage
            localStorage.setItem('eleToken', token);

            // 解析token
            const decoded = jwt_decode(token);

            // token儲存到vuex中
            this.$store.dispatch('setAuthenticated', !this.isEmpty(decoded));
            this.$store.dispatch('setUser', decoded);

          });

          this.$router.push("/index");
        }
      });
    },
    isEmpty(value) {
        return (
            value === undefined ||
            value === null ||
            (typeof value === 'object' && Object.keys(value).length === 0) ||
            (typeof value === 'string' && value.trim().length === 0)
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
