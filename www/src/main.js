import { createApp } from 'vue'
import VueApp from './VueApp.vue'
import router from './router'
import 'jquery'
import 'bootstrap'
import 'bootstrap-switch'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import { mock } from "./stub/api"

if (import.meta.env.DEV) {
  mock()
}

createApp(VueApp)
  .use(router)
  .use(ElementPlus)
  .mount('#app')
