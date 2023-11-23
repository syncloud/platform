import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import inject from '@rollup/plugin-inject'

export default defineConfig({
  plugins: [
    inject({
      $: 'jquery',
      jQuery: 'jquery'
    }),
    vue()
  ]
})
