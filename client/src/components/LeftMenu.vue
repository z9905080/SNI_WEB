<template>
    <el-row class="menu_page">
         <el-col>
             <el-menu
                mode="vertical"
                background-color="#324057"
                text-color="#fff"
                active-text-color="#409eff"
                class="el-menu-vertical-demo">
                <router-link to="/home">
                    <el-menu-item index="0">
                        <i class="fa fa-margin fa-server"></i>
                        <span slot="title" class="title">首頁</span>
                    </el-menu-item>
                </router-link>
                <template  v-for="item in items" >
                    <el-submenu v-if="item.children" :index="item.path" :key="item.path">
                        <template slot="title">
                            <i :class="'fa fa-margin '+item.icon"></i>
                            <span slot="title" class="title">{{item.name}}</span>
                        </template>
                        <router-link v-for="(citem,cindex) in item.children"
                            :to="citem.path" :key="cindex">
                            <el-menu-item
                                :index='citem.path'>
                                <span slot="title">{{citem.name}}</span>
                            </el-menu-item>
                        </router-link>
                    </el-submenu>

                </template>
             </el-menu>
         </el-col>
    </el-row>
</template>
<script>
export default {
  name: "leftmenu",
  data() {
    return {
      items: [
        // {
        //   icon: "fa-money",
        //   name: "資金管理",
        //   path: "fund",
        //   children: [{ path: "fundlist", name: "資金流水" }]
        // },
        // {
        //   icon: "fa-asterisk",
        //   name: "信息管理",
        //   path: "info",
        //   children: [{ path: "infoshow", name: "个人信息" }]
        // }
        {
          icon: "far fa-bookmark",
          name: "頁籤&內文管理",
          path: "navpage",
          children: [{ path: "/navpageedit", name: "頁籤列表&編輯" }]
        },
        // {
        //   icon: "el-icon-document",
        //   name: "文章管理",
        //   path: "article",
        //   children: [{ path: "/articlelist", name: "文章列表" }]
        // },
        {
          icon: "fas fa-image",
          name: "輪播圖設定",
          path: "slider",
          children: [{ path: "/slider", name: "圖片管理" }, { path: "/carousel", name: "輪播圖列表" }]
        },
        {
          icon: "fas fa-audio-description",
          name: "跑馬燈設定",
          path: "marquee",
          children: [{ path: "/marquee", name: "跑馬燈設定" }]
        },
      ]
    };
  }
};
</script>
<style scoped>
.menu_page {
  position: fixed;
  top: 71px;
  left: 0;
  min-height: 100%;
  background-color: #324057;
  z-index: 99;
}
.el-menu {
  border: none;
  font-size: 18px;
}
.fa-margin {
  margin-right: 5px;
}
.el-menu-vertical-demo:not(.el-menu--collapse) {
  width: 200px;
  min-height: 400px;
}
.el-menu-vertical-demo {
  width: 35px;
}
.el-submenu .el-menu-item {
  min-width: 200px;
  font-size: 18px;
}
.title {
  font-size: 18px;
}

.hiddenDropdown,
.hiddenDropname {
  display: none;
}
a {
  text-decoration: none;
}
</style>
