import { defineStore } from 'pinia'

const STORAGE_KEY = 'syncloud-theme'

function preferredTheme () {
  const saved = localStorage.getItem(STORAGE_KEY)
  if (saved === 'dark' || saved === 'light') {
    return saved
  }
  if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
    return 'dark'
  }
  return 'light'
}

function apply (theme) {
  document.documentElement.setAttribute('data-theme', theme)
}

export const useThemeStore = defineStore('theme', {
  state: () => ({
    theme: preferredTheme()
  }),
  getters: {
    isDark: (state) => state.theme === 'dark'
  },
  actions: {
    init () {
      apply(this.theme)
    },
    set (theme) {
      this.theme = theme === 'dark' ? 'dark' : 'light'
      localStorage.setItem(STORAGE_KEY, this.theme)
      apply(this.theme)
    },
    toggle () {
      this.set(this.isDark ? 'light' : 'dark')
    }
  }
})
