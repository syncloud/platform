import { createApp } from 'vue'
import LoginApp from './LoginApp.vue'

async function start () {
  if (import.meta.env.VITE_STUB) {
    const { mock } = await import('./stub/api')
    mock()
  }
  createApp(LoginApp).mount('#app')
}

start()
