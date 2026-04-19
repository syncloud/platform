const { config } = require('@vue/test-utils')
const { createI18n } = require('vue-i18n')
const en = require('../src/locales/en.json')

const i18n = createI18n({
  legacy: false,
  globalInjection: true,
  locale: 'en',
  fallbackLocale: 'en',
  messages: { en }
})

config.global.plugins = [i18n]
