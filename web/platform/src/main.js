import { createApp } from 'vue'
import { createPinia } from 'pinia'
import VueApp from './VueApp.vue'
import router from './router'
import 'font-awesome/css/font-awesome.css'
import 'roboto-fontface/css/roboto/roboto-fontface.css'
import 'material-icons/iconfont/material-icons.css'
import './style/design.css'
import ui from './ui'
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
  .use(ui)
  .mount('#app')
