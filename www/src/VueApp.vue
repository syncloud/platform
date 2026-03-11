<template>
  <Menu v-bind:activeTab="currentPath"/>
  <router-view/>
  <Error ref="app_error"/>
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
      this.$refs.app_error.showAxios(err)
    })
  }
}
</script>
<style>
@import 'bootstrap/dist/css/bootstrap.css';
@import 'font-awesome/css/font-awesome.css';
@import 'roboto-fontface/css/roboto/roboto-fontface.css';
</style>
