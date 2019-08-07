import Vue from 'vue'
import VueRouter from 'vue-router'

import BootstrapVue from 'bootstrap-vue'

import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import 'bootswatch/dist/darkly/bootstrap.css'
import 'easymde/dist/easymde.min.css'

import App from './app.vue'
import View from './view.vue'
import Edit from './edit.vue'

Vue.use(BootstrapVue)
Vue.use(VueRouter)

const routes = [
  {
    component: View,
    name: 'view',
    path: '/:page',
  },
  {
    component: Edit,
    name: 'edit',
    path: '/:page/edit',
  },
  {
    name: 'home',
    path: '/',
    redirect: '/Home',
  },
]

const router = new VueRouter({
  mode: 'history',
  routes,
})

new Vue({
  el: '#app',
  components: { App },
  data: { },
  render: createElement => createElement('app'),
  router,
})
