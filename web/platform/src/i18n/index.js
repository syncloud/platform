import { createI18n } from 'vue-i18n'
import { ref } from 'vue'
import en from '../locales/en.json'
import enLocale from 'element-plus/dist/locale/en.mjs'

export const SUPPORTED_LOCALES = [
  { code: 'en', name: 'English' },
  { code: 'zh-CN', name: '中文' },
  { code: 'es', name: 'Español' },
  { code: 'hi', name: 'हिन्दी' },
  { code: 'ar', name: 'العربية' },
  { code: 'pt', name: 'Português' },
  { code: 'ru', name: 'Русский' },
  { code: 'ja', name: '日本語' },
  { code: 'de', name: 'Deutsch' },
  { code: 'fr', name: 'Français' }
]

const RTL_LOCALES = ['ar']

const ELEMENT_LOCALE_FILES = {
  en: () => import('element-plus/dist/locale/en.mjs'),
  'zh-CN': () => import('element-plus/dist/locale/zh-cn.mjs'),
  es: () => import('element-plus/dist/locale/es.mjs'),
  hi: () => import('element-plus/dist/locale/hi.mjs'),
  ar: () => import('element-plus/dist/locale/ar.mjs'),
  pt: () => import('element-plus/dist/locale/pt.mjs'),
  ru: () => import('element-plus/dist/locale/ru.mjs'),
  ja: () => import('element-plus/dist/locale/ja.mjs'),
  de: () => import('element-plus/dist/locale/de.mjs'),
  fr: () => import('element-plus/dist/locale/fr.mjs')
}

const APP_LOCALE_FILES = {
  en: () => Promise.resolve({ default: en }),
  'zh-CN': () => import('../locales/zh-CN.json'),
  es: () => import('../locales/es.json'),
  hi: () => import('../locales/hi.json'),
  ar: () => import('../locales/ar.json'),
  pt: () => import('../locales/pt.json'),
  ru: () => import('../locales/ru.json'),
  ja: () => import('../locales/ja.json'),
  de: () => import('../locales/de.json'),
  fr: () => import('../locales/fr.json')
}

export const STORAGE_KEY = 'syncloud.locale'

function codes () {
  return SUPPORTED_LOCALES.map(l => l.code)
}

export function detectLocale () {
  try {
    const stored = typeof localStorage !== 'undefined' && localStorage.getItem(STORAGE_KEY)
    if (stored && codes().includes(stored)) return stored
  } catch { /* localStorage may be disabled */ }

  const nav = typeof navigator !== 'undefined' ? (navigator.language || 'en') : 'en'
  if (codes().includes(nav)) return nav
  const base = nav.split('-')[0]
  const match = codes().find(c => c === base || c.startsWith(base + '-'))
  return match || 'en'
}

export const elementLocale = ref(enLocale)

const i18n = createI18n({
  legacy: false,
  globalInjection: true,
  locale: 'en',
  fallbackLocale: 'en',
  messages: { en }
})

export async function setLocale (code) {
  if (!codes().includes(code)) code = 'en'

  if (!i18n.global.availableLocales.includes(code)) {
    const mod = await APP_LOCALE_FILES[code]()
    i18n.global.setLocaleMessage(code, mod.default || mod)
  }
  i18n.global.locale.value = code

  const elMod = await ELEMENT_LOCALE_FILES[code]()
  elementLocale.value = elMod.default || elMod

  try { localStorage.setItem(STORAGE_KEY, code) } catch { /* ignore */ }

  if (typeof document !== 'undefined') {
    document.documentElement.lang = code
    document.documentElement.dir = RTL_LOCALES.includes(code) ? 'rtl' : 'ltr'
  }
}

export default i18n
