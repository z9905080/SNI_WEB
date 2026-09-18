<template>
  <header class="head-nav">
    <el-row>
      <el-col :span="6" class="logo-container">
        <!-- <img src="../assets/logo.png" class="logo"> -->
        <span class="title">生命之家後台管理系統</span>
      </el-col>
      <el-col :span="6" class="user">
        <div class="userinfo">
          <!-- <img src="https://www.gravatar.com/avatar/anything?s=200&d=mm" class="avatar"> -->
          <div class="welcome">
              <p class="name comename">歡迎</p>
              <p class="name avatarname">Hulk</p>
          </div>
          <span class="username">

          </span>
        </div>
      </el-col>
    </el-row>
  </header>
</template>

<script>
export default {
  name: "haed-nav"
};
</script>

<style scoped>
.head-nav {
  width: 100%;
  height: 60px;
  min-width: 600px;
  padding: 5px;
  background: #324057;
  color: #fff;
  border-bottom: 1px solid #1f2d3d;
}
.logo-container {
    line-height: 60px;
    min-width: 400px;
}
.welcome {
    display: inline-block;
    width: auto;
    vertical-align: middle;
    padding: 0 5px;
}
.name {
    line-height: 20px;
    text-align: center;
    font-size: 14px;
}
.comename {
    font-size: 12px;
}
.avatarname {
    color: #409eff;
    font-weight: bolder;
}
.username {
    cursor: pointer;
    margin-right: 5px;
}
</style>