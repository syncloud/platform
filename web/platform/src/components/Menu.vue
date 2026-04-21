<template>
  <div class="wrapper">
    <div class="content">
      <div class="headblock">
        <header class="wd12">
          <div class="logo" :class="{ onelogo: !auth.loggedIn }">Syncloud</div>
          <div class="menulinks" v-if="auth.loggedIn">
            <router-link to="/" id="apps" class="apps hlink" :class="{ active: activeTab === '/' }">{{ $t('menu.apps') }}</router-link>
            <router-link to="/appcenter" id="appcenter" class="appcenter hlink"
                         :class="{ active: activeTab === '/appcenter' }">{{ $t('menu.appCenter') }}
            </router-link>
            <router-link to="/settings" id="settings" class="settings hlink"
                         :class="{ active: activeTab === '/settings' }">{{ $t('menu.settings') }}
            </router-link>
          </div>
          <div class="menuoff" v-if="auth.loggedIn">
            <a href="#" id="logout" class="hlink" @click="logout">
              <i class="material-icons" style="vertical-align: middle">exit_to_app</i>
              <span class="button_label">{{ $t('menu.logout') }}</span>
            </a>
          </div>
          <div id="menubutton" class="menubutton" v-if="auth.loggedIn" @click="toggle" :class="{ menuopen: menuOpen }">
            <span></span>
            <span></span>
            <span></span>
            <span></span>
          </div>
        </header>
        <div id="menu" class="navi" v-if="auth.loggedIn" :class="{ naviopen: menuOpen }">
          <router-link to="/" id="apps_mobile" @click="close">{{ $t('menu.apps') }}</router-link>
          <router-link to="/appcenter" id="appcenter_mobile" @click="close">{{ $t('menu.appCenter') }}</router-link>
          <router-link to="/settings" id="settings_mobile" @click="close">{{ $t('menu.settings') }}</router-link>
          <div class="menucolor2">
            <a href="#" id="logout_mobile" @click="logout(); close()">{{ $t('menu.logout') }}</a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { useAuthStore } from '../stores/auth'

export default {
  props: {
    activeTab: String
  },
  data () {
    return {
      menuOpen: false,
      auth: useAuthStore()
    }
  },
  methods: {
    close: function () {
      this.menuOpen = false
    },
    toggle: function () {
      this.menuOpen = !this.menuOpen
    },
    logout: function () {
      window.location.href = '/rest/logout'
    }
  }
}

</script>
<style>
</style>
