import Vue from 'vue'
import Router from 'vue-router'
import Index from '@/views/index.vue'
import Register from '@/views/Register.vue'
import Login from '@/views/Login.vue'
import NotFound from '@/views/404.vue'
import Home from '@/views/Home.vue'
import InfoShow from '@/views/InfoShow.vue'
import ArticleList from '@/views/ArticleList.vue'
import FundList from '@/views/FundList.vue'
import Edit from '@/views/Edit.vue'
import NavPageEdit from '@/views/NavPageEdit.vue'

Vue.use(Router)

const router = new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      redirect: '/index'
    },
    {
      path: '/index',
      name: 'index',
      component: Index,
      children: [
        {path: '', component: Home},
        {path: '/home', name: 'home', component: Home},
        {path: '/infoshow', name: 'infoshow', component: InfoShow},
        {path: '/articlelist', name: 'articlelist', component: ArticleList},
        {path: '/navpageedit', name: 'navpageedit', component: NavPageEdit},
        {path: '/fundlist', name: 'fundlist', component: FundList},
        {path: '/edit', name: 'edit', component: Edit},
      ]
    },
    {
      path: '/register',
      name: 'register',
      component: Register
    },
    {
      path: '/login',
      name: 'login',
      component: Login
    },
    {
      path: '*',
      name: '/404',
      component: NotFound
    }
  ]
});

// 路由守衛(router guard)
router.beforeEach((to, from, next) => {
  const isLogin = localStorage.eleToken ? true : false;
  if (to.path === '/login' || to.path === '/register') {
    next();
  } else {
    isLogin ? next() : next('/login');
  }
})

export default router;
