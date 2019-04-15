<template>
  <div class="register">
    <section class="form_container">
      <div class="manage_tip">
        <span class="title">生命之家後台管理系統</span>
        <el-form
          :model="registerUser"
          :rules="rules"
          ref="registerForm"
          label-width="100px"
          class="registerForm"
        >
          <el-form-item label="用戶名" prop="name">
            <el-input type="text" v-model="registerUser.name" placeholder="請輸入用戶名"></el-input>
          </el-form-item>
          <el-form-item label="郵件地址" prop="email">
            <el-input type="text" v-model="registerUser.email" placeholder="請輸入email"></el-input>
          </el-form-item>
          <el-form-item label="密碼" prop="password">
            <el-input type="password" v-model="registerUser.password" placeholder="請輸入密碼"></el-input>
          </el-form-item>
          <el-form-item label="確認密碼" prop="password2">
            <el-input type="password" v-model="registerUser.password2" placeholder="請確認密碼"></el-input>
          </el-form-item>
          <el-form-item label="選擇身分">
            <el-select class="identity_select" v-model="registerUser.identity" placeholder="請選擇身分">
              <el-option label="管理者" value="admin"></el-option>
              <el-option label="一般使用者" value="user"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" class="submit_btn" @click="submitForm('registerForm')">註冊</el-button>
            <!-- <el-button type="danger" class="cancel_btn" @click="resetForm('ruleForm2')">取消</el-button> -->
          </el-form-item>
        </el-form>
      </div>
    </section>
  </div>
</template>

<script>
export default {
  name: "register",
  components: {},
  data() {
    let validatePass2 = (rule, value, callback) => {
      if (value !== this.registerUser.password) {
        callback(new Error("两次输入密码不一致!"));
      } else {
        callback();
      }
    };
    return {
      registerUser: {
        name: "",
        email: "",
        password: "",
        password2: "",
        identity: ""
      },
      rules: {
        name: [
          {
            required: true,
            message: "用戶名不可為空",
            trigger: "blur"
          },
          {
            min: 2,
            max: 30,
            message: "長度在2到30個字之間",
            trigger: "blur"
          }
        ],
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
          },
          {
            min: 6,
            max: 30,
            message: "長度在6到30個字之間",
            trigger: "blur"
          }
        ],
        password2: [
          {
            required: true,
            message: "確認密碼不可為空",
            trigger: "blur"
          },
          {
            min: 6,
            max: 30,
            message: "長度在6到30個字之間",
            trigger: "blur"
          },
          {
            validator: validatePass2,
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
            this.$axios.post("/api/users/register", this.registerUser)
            .then((result) => {
                //註冊成功
                this.$message({
                    message: '帳號註冊成功!',
                    type: 'success'
                });
            })
            .catch((err) => {
                //註冊失敗
            });

            this.$router.push("/login")

        } else {
          console.log("error submit!!");
          return false;
        }
      });
    },
    resetForm(formName) {
      this.$refs[formName].resetFields();
    }
  }
};
</script>

<style scoped>
.register {
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
  top: 10%;
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
.registerForm {
  margin-top: 20px;
  background-color: #fff;
  padding: 20px 40px 20px 20px;
  border-radius: 5px;
  box-shadow: 0px 5px 10px #ccc;
}
.identity_select {
  width: 100%;
}
.submit_btn {
  width: 100%;
}
</style>
