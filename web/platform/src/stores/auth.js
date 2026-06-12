import { defineStore } from 'pinia'
import axios from 'axios'

const publicRoutes = [
  '/error',
  '/login'
]

export const useAuthStore = defineStore('auth', {
  state: () => ({
    loggedIn: undefined,
    activated: true,
    admin: false
  }),
  actions: {
    checkUserSession (router, onError) {
      axios.get('/rest/user')
        .then((response) => {
          this.loggedIn = true
          this.admin = !!(response.data && response.data.data && response.data.data.admin)
          const path = router.currentRoute.value.path
          if (path === '/login' || path === '/activate') {
            router.push('/')
          }
        })
        .catch(() => {
          axios.get('/rest/activation/status')
            .then(response => {
              this.loggedIn = false
              if (!response.data.data) {
                this.activated = false
                router.push('/activate')
              } else {
                const path = router.currentRoute.value.path
                if (!publicRoutes.includes(path)) {
                  router.push('/login')
                }
              }
            })
            .catch(err => {
              if (onError) onError(err)
              router.push('/error')
            })
        })
    }
  }
})
