import { createApp } from 'vue'
import VueApp from './VueApp.vue'
import router from './router'
import './jQuery'
import 'bootstrap'
import 'element-plus/dist/index.css'
import { mock } from './stub/api'

if (import.meta.env.DEV) {
  mock()
}

createApp(VueApp)
  .use(router)
  .mount('#app')
