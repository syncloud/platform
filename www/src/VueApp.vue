<template>
  <Menu v-bind:activeTab="currentPath" v-bind:checkUserSession="checkUserSession" v-bind:loggedIn="loggedIn"/>
  <router-view v-bind:checkUserSession="checkUserSession" :activated="activated"/>
  <Error ref="app_error"/>
</template>
<script>
import axios from 'axios'
import Menu from './components/Menu.vue'
import Error from './components/Error.vue'

const publicRoutes = [
  '/error',
  '/login'
]

export default {
  data () {
    return {
      currentPath: '',
      loggedIn: undefined,
      email: '',
      activated: true
    }
  },
  name: 'VueApp',
  components: {
    Menu,
    Error
  },
  watch: {
    $route (to, from) {
      console.log('route change from ' + from.path + ' to ' + to.path)
      this.currentPath = to.path
    }
  },
  methods: {
    checkUserSession: function () {
      axios.get('/rest/user')
        .then(() => {
          this.loggedIn = true
          if (this.currentPath === '/login' || this.currentPath === '/activate') {
            this.$router.push('/')
          }
        })
        .catch(() => {
          axios.get('/rest/activation/status')
            .then(response => {
              this.loggedIn = false
              if (!response.data.data) {
                this.activated = false
                this.$router.push('/activate')
              } else {
                if (!publicRoutes.includes(this.currentPath)) {
                  this.$router.push('/login')
                }
              }
            })
            .catch(err => {
              this.$refs.app_error.showAxios(err)
              this.$router.push('/error')
            })
        })
    }
  },
  mounted () {
    this.currentPath = this.$route.path
    this.checkUserSession()
  }
}
</script>
<style>
@import 'bootstrap/dist/css/bootstrap.css';
@import 'font-awesome/css/font-awesome.css';
@import 'roboto-fontface/css/roboto/roboto-fontface.css';
</style>
