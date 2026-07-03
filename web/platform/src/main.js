import { createApp } from 'vue'
import { createPinia } from 'pinia'
import VueApp from './VueApp.vue'
import router from './router'
import 'font-awesome/css/font-awesome.css'
import 'roboto-fontface/css/roboto/roboto-fontface.css'
import 'material-icons/iconfont/material-icons.css'
import './style/design.css'
import ui from './ui'
import i18n, { detectLocale, setLocale } from './i18n'
import { useThemeStore } from './stores/theme'

async function start () {
  if (import.meta.env.VITE_STUB) {
    const { mock } = await import('./stub/api')
    mock()
  }

  setLocale(detectLocale())

  const pinia = createPinia()

  createApp(VueApp)
    .use(pinia)
    .use(router)
    .use(i18n)
    .use(ui)
    .mount('#app')

  useThemeStore(pinia).init()
}

start()
