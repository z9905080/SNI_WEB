/*!

 =========================================================
 * Vue Now UI Kit - v1.1.0
 =========================================================

 * Product Page: https://www.creative-tim.com/product/now-ui-kit
 * Copyright 2019 Creative Tim (http://www.creative-tim.com)

 * Designed by www.invisionapp.com Coded by www.creative-tim.com

 =========================================================

 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

 */
import Vue from 'vue';
import Vuex from 'vuex'
import App from './App.vue';
import router from './router';
import NowUiKit from './plugins/now-ui-kit';
import VueAnalytics from 'vue-analytics'
import storeParam from './store'

Vue.config.productionTip = false;

Vue.use(Vuex);
Vue.use(NowUiKit);
Vue.use(VueAnalytics, {
  id: 'UA-154926984-1',
  router,
  autoTracking: {
    pageviewOnLoad: false
  }
})

const store = new Vuex.Store(storeParam);

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app');
