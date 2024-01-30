import 'bootstrap/dist/css/bootstrap.css'
import 'easymde/dist/easymde.min.css'
import './easymde.css'

import { createApp, h } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'

import App from './app.vue'
import Edit from './edit.vue'
import View from './view.vue'

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

const router = createRouter({
  history: createWebHistory(),
  mode: 'history',
  routes,
})

const app = createApp({
  name: 'WikiMain',
  render() {
    return h(App)
  },

  router,
})

app.use(router)
app.mount('#app')
