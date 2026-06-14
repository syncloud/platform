import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  { path: '/', name: 'Apps', component: () => import('../views/Apps.vue') },
  { path: '/login', name: 'Login', component: () => import('../views/Login.vue') },
  { path: '/app', name: 'App', component: () => import('../views/App.vue') },
  { path: '/appcenter', name: 'AppCenter', component: () => import('../views/AppCenter.vue') },
  { path: '/settings', name: 'Settings', component: () => import('../views/Settings.vue') },
  { path: '/activation', name: 'Activation', component: () => import('../views/Activation.vue') },
  { path: '/activate', name: 'Activate', component: () => import('../views/Activate.vue') },
  { path: '/backup', name: 'Backup', component: () => import('../views/Backup.vue'), meta: { admin: true } },
  { path: '/network', name: 'Network', component: () => import('../views/Network.vue') },
  { path: '/access', name: 'Access', component: () => import('../views/Access.vue'), meta: { admin: true } },
  { path: '/storage', name: 'Storage', component: () => import('../views/Storage.vue'), meta: { admin: true } },
  { path: '/internalmemory', name: 'InternalMemory', component: () => import('../views/InternalMemory.vue'), meta: { admin: true } },
  { path: '/updates', name: 'Updates', component: () => import('../views/Updates.vue'), meta: { admin: true } },
  { path: '/support', name: 'Support', component: () => import('../views/Support.vue') },
  { path: '/certificate', name: 'Certificate', component: () => import('../views/Certificate.vue'), meta: { admin: true } },
  { path: '/certificate/log', name: 'Certificate Log', component: () => import('../views/CertificateLog.vue'), meta: { admin: true } },
  { path: '/twofactor', name: 'TwoFactor', component: () => import('../views/TwoFactor.vue') },
  { path: '/logs', name: 'Logs', component: () => import('../views/Logs.vue') },
  { path: '/customproxy', name: 'CustomProxy', component: () => import('../views/CustomProxy.vue') },
  { path: '/users', name: 'Users', component: () => import('../views/UsersList.vue'), meta: { admin: true } },
  { path: '/useredit', name: 'UserEdit', component: () => import('../views/UserEdit.vue'), meta: { admin: true } },
  { path: '/system', name: 'System', component: () => import('../views/System.vue') },
  { path: '/locale', name: 'Locale', component: () => import('../views/Locale.vue') },
  { path: '/health', name: 'Health', component: () => import('../views/Health.vue'), meta: { admin: true } },
  { path: '/:catchAll(.*)', redirect: '/' }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

router.beforeEach((to) => {
  if (to.meta.admin) {
    const auth = useAuthStore()
    if (auth.loggedIn && !auth.admin) {
      return '/'
    }
  }
})

export default router
