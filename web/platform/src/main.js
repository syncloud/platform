import { createApp } from 'vue'
import { createPinia } from 'pinia'
import VueApp from './VueApp.vue'
import router from './router'
import 'element-plus/dist/index.css'
import { mock } from './stub/api'

if (import.meta.env.DEV) {
  mock()
}

createApp(VueApp)
  .use(createPinia())
  .use(router)
  .mount('#app')
