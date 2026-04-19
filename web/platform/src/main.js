import { createApp } from 'vue'
import { createPinia } from 'pinia'
import VueApp from './VueApp.vue'
import router from './router'
import 'element-plus/dist/index.css'
import { mock } from './stub/api'
import i18n, { detectLocale, setLocale } from './i18n'

if (import.meta.env.DEV) {
  mock()
}

setLocale(detectLocale())

createApp(VueApp)
  .use(createPinia())
  .use(router)
  .use(i18n)
  .mount('#app')
