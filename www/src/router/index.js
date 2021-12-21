import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', name: 'Apps', component: () => import('../views/Apps.vue') },
  { path: '/login', name: 'Login', component: () => import('../views/Login.vue') },
  { path: '/app', name: 'App', component: () => import('../views/App.vue') },
  { path: '/appcenter', name: 'AppCenter', component: () => import('../views/AppCenter.vue') },
  { path: '/settings', name: 'Settings', component: () => import('../views/Settings.vue') },
  { path: '/activation', name: 'Activation', component: () => import('../views/Activation.vue') },
  { path: '/activate', name: 'Activate', component: () => import('../views/Activate.vue') },
  { path: '/backup', name: 'Backup', component: () => import('../views/Backup.vue') },
  { path: '/network', name: 'Network', component: () => import('../views/Network.vue') },
  { path: '/access', name: 'Access', component: () => import('../views/Access.vue') },
  { path: '/storage', name: 'Storage', component: () => import('../views/Storage.vue') },
  { path: '/internalmemory', name: 'InternalMemory', component: () => import('../views/InternalMemory.vue') },
  { path: '/updates', name: 'Updates', component: () => import('../views/Updates.vue') },
  { path: '/support', name: 'Support', component: () => import('../views/Support.vue') },
  { path: '/certificate', name: 'Certificate', component: () => import('../views/Certificate.vue') },
  { path: '/certificate/log', name: 'Certificate Log', component: () => import('../views/CertificateLog.vue') },
  { path: '/:catchAll(.*)', redirect: '/' }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
