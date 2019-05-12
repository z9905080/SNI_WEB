import Vue from 'vue'
import App from './App.vue';
// import VueCookies from 'vue-cookies';
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';

import axios from './http';
import router from './router'
import store from './store'

Vue.config.productionTip = false;
Vue.use(ElementUI);
// Vue.use(VueCookies);
Vue.prototype.$axios = axios;
// Vue.prototype.$cookies = VueCookies;

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
