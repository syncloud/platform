<template>
  <div class="wrapper">
    <div class="content">
      <div class="headblock">
        <header class="wd12">
          <div class="logo" :class="{ onelogo: !loggedIn }">Syncloud</div>
          <div class="menulinks" v-if="loggedIn">
            <router-link to="/" id="apps" class="apps hlink" :class="{ active: activeTab === '/' }">Apps</router-link>
            <router-link to="/appcenter" id="appcenter" class="appcenter hlink"
                         :class="{ active: activeTab === '/appcenter' }">App Center
            </router-link>
            <router-link to="/settings" id="settings" class="settings hlink"
                         :class="{ active: activeTab === '/settings' }">Settings
            </router-link>
          </div>
          <div class="menuoff" v-if="loggedIn">
            <a href="#" id="logout" class="hlink" @click="logout">
              <i class="material-icons" style="vertical-align: middle">exit_to_app</i>
              <span class="button_label">Logout</span>
            </a>
            <a href="#" id="restart" class="hlink" @click="restart">
              <i class="material-icons" style="vertical-align: middle">loop</i>
              <span class="button_label">Restart</span>
            </a>
            <a href="#" id="shutdown" class="hlink" @click="shutdown">
              <i class="material-icons" style="vertical-align: middle">power_settings_new</i>
              <span class="button_label">Shutdown</span>
            </a>
          </div>
          <div id="menubutton" class="menubutton" v-if="loggedIn" @click="toggle" :class="{ menuopen: menuOpen }">
            <span></span>
            <span></span>
            <span></span>
            <span></span>
          </div>
        </header>
        <div id="menu" class="navi" v-if="loggedIn" :class="{ naviopen: menuOpen }">
          <router-link to="/" id="apps_mobile"><span style="display: block" @click="toggle">Apps</span></router-link>
          <router-link to="/appcenter" id="appcenter_mobile"><span style="display: block"
                                                                   @click="toggle">App Center</span></router-link>
          <router-link to="/settings" id="settings_mobile"><span style="display: block" @click="toggle">Settings</span>
          </router-link>
          <div class="menucolor2">
            <a href="#" id="logout_mobile" @click="logout">Log out</a>
            <a href="#" id="restart_mobile" @click="restart">Restart</a>
            <a href="#" id="shutdown_mobile" @click="shutdown">Shutdown</a>
          </div>
        </div>
      </div>
    </div>
  </div>
  <Error ref="menu_error"/>
</template>

<script>
import axios from 'axios'
import Error from '../components/Error.vue'

export default {
  props: {
    activeTab: String,
    loggedIn: Boolean,
    checkUserSession: Function
  },
  data () {
    return {
      menuOpen: false
    }
  },
  components: {
    Error
  },
  methods: {
    close: function () {
      this.menuOpen = false
    },
    toggle: function () {
      this.menuOpen = !this.menuOpen
    },
    logout: function () {
      axios.post('/rest/logout')
        .then(() => {
          this.checkUserSession()
        })
        .catch(err => {
          console.log(err)
        })
    },
    restart: function () {
      const error = this.$refs.menu_error
      axios.post('/rest/restart')
        .catch(err => error.showAxios(err))
    },
    shutdown: function () {
      const error = this.$refs.menu_error
      axios.post('/rest/shutdown')
        .catch(err => error.showAxios(err))
    }
  }
}

</script>
<style>
</style>
