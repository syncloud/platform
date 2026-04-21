<template>
  <el-config-provider :locale="elementLocale">
    <Menu v-if="!isStandalone" v-bind:activeTab="currentPath"/>
    <router-view/>
    <Error ref="app_error" v-if="!isStandalone"/>
  </el-config-provider>
</template>
<script>
import Menu from './components/Menu.vue'
import Error from './components/Error.vue'
import { useAuthStore } from './stores/auth'
import { elementLocale } from './i18n'

export default {
  data () {
    return {
      currentPath: '',
      elementLocale
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
<style>
@import 'bootstrap/dist/css/bootstrap.css';
@import 'font-awesome/css/font-awesome.css';
@import 'roboto-fontface/css/roboto/roboto-fontface.css';
</style>
