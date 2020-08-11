import Vue from "vue";
import Vuex from 'vuex'
import App from "./App.vue";
import router from "./router";
import "./styles/global.scss"
import { BootstrapVue, BootstrapVueIcons } from "bootstrap-vue";
import VueProgressBar from 'vue-progressbar'

import "bootstrap-vue/dist/bootstrap-vue.css";
import "bootstrap/dist/css/bootstrap.css";

Vue.use(Vuex)
Vue.use(BootstrapVue);
Vue.use(BootstrapVueIcons);
Vue.use(VueProgressBar, {
  color: '#4590d3',
  failedColor: 'red',
  height: '2px'
});

Vue.config.productionTip = false;

new Vue({
  router,
  render: h => h(App)
}).$mount("#app");
