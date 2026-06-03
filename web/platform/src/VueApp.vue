<template>
  <Menu v-if="!isStandalone" v-bind:activeTab="currentPath"/>
  <router-view/>
  <Error ref="app_error" v-if="!isStandalone"/>
</template>
<script>
import Menu from './components/Menu.vue'
import Error from './components/Error.vue'
import { useAuthStore } from './stores/auth'

export default {
  data () {
    return {
      currentPath: ''
    }
  },
  computed: {
    isStandalone () {
      return this.currentPath === '/activate'
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
  mounted () {
    this.currentPath = this.$route.path
    const auth = useAuthStore()
    auth.checkUserSession(this.$router, (err) => {
      if (this.$refs.app_error) {
        this.$refs.app_error.showAxios(err)
      }
    })
  }
}
</script>
