<template>
  <header class="sc-header">
    <div class="sc-header-inner">
      <router-link to="/" class="sc-brand" :class="{ centered: !auth.loggedIn }">SYNCLOUD</router-link>

      <nav class="sc-nav" v-if="auth.loggedIn">
        <router-link to="/" id="apps" class="sc-nav-link" :class="{ active: activeTab === '/' }">{{ $t('menu.apps') }}</router-link>
        <router-link v-if="auth.admin" to="/appcenter" id="appcenter" class="sc-nav-link" :class="{ active: activeTab === '/appcenter' }">{{ $t('menu.appCenter') }}</router-link>
        <router-link v-if="auth.admin" to="/settings" id="settings" class="sc-nav-link" :class="{ active: activeTab === '/settings' }">{{ $t('menu.settings') }}</router-link>
      </nav>

      <a href="#" id="logout" class="sc-logout" v-if="auth.loggedIn" @click="logout">
        <i class="material-icons">exit_to_app</i>
        <span class="button_label">{{ $t('menu.logout') }}</span>
      </a>

      <button
        type="button"
        id="theme_toggle"
        data-testid="theme-toggle"
        class="sc-theme-toggle"
        :title="theme.isDark ? $t('menu.lightMode') : $t('menu.darkMode')"
        :aria-label="theme.isDark ? $t('menu.lightMode') : $t('menu.darkMode')"
        @click="theme.toggle()">
        <i class="material-icons">{{ theme.isDark ? 'light_mode' : 'dark_mode' }}</i>
      </button>

      <div id="menubutton" class="sc-burger" v-if="auth.loggedIn" @click="toggle" :class="{ menuopen: menuOpen }">
        <span></span><span></span><span></span><span></span>
      </div>
    </div>

    <div id="menu" class="sc-mobile-nav" v-if="auth.loggedIn" :class="{ naviopen: menuOpen }">
      <router-link to="/" id="apps_mobile" @click="close">{{ $t('menu.apps') }}</router-link>
      <router-link v-if="auth.admin" to="/appcenter" id="appcenter_mobile" @click="close">{{ $t('menu.appCenter') }}</router-link>
      <router-link v-if="auth.admin" to="/settings" id="settings_mobile" @click="close">{{ $t('menu.settings') }}</router-link>
      <a href="#" id="logout_mobile" class="sc-mobile-logout" @click="logout(); close()">{{ $t('menu.logout') }}</a>
    </div>

    <router-link
      v-if="spaceWarningVisible"
      to="/internalmemory"
      id="disk_space_warning"
      data-testid="disk-space-warning"
      class="sc-header-alert">
      <i class="material-icons">warning</i>
      <span>{{ $t(spaceMessage, { free: freeText }) }}</span>
    </router-link>
  </header>
</template>

<script>
import { useAuthStore } from '../stores/auth'
import { useThemeStore } from '../stores/theme'
import { useSpaceStore } from '../stores/space'

const SPACE_INTERVAL_MS = 5 * 60 * 1000

export default {
  props: {
    activeTab: String
  },
  data () {
    return {
      menuOpen: false,
      auth: useAuthStore(),
      theme: useThemeStore(),
      space: useSpaceStore(),
      spaceTimer: undefined
    }
  },
  computed: {
    spaceWarningVisible () {
      return this.auth.admin && this.space.low && this.space.lowest !== undefined
    },
    spaceMessage () {
      return this.space.lowest.kind === 'data' ? 'space.lowData' : 'space.lowSystem'
    },
    freeText () {
      const kb = this.space.lowest.free_kb
      if (kb >= 1024 * 1024) {
        return (kb / (1024 * 1024)).toFixed(1) + ' GB'
      }
      return Math.round(kb / 1024) + ' MB'
    }
  },
  watch: {
    'auth.admin' (admin) {
      if (admin) {
        this.space.load()
      }
    }
  },
  mounted () {
    if (this.auth.admin) {
      this.space.load()
    }
    this.spaceTimer = setInterval(() => {
      if (this.auth.admin) {
        this.space.load()
      }
    }, SPACE_INTERVAL_MS)
  },
  beforeUnmount () {
    clearInterval(this.spaceTimer)
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

<style scoped>
.sc-header {
  position: relative;
  background: var(--sc-header-bg);
  backdrop-filter: saturate(180%) blur(12px);
  border-bottom: 1px solid var(--sc-border-soft);
  z-index: 30;
}
.sc-header-inner {
  max-width: 1100px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  padding: 0 24px;
  height: 64px;
  gap: 28px;
}
.sc-brand {
  font-size: 15px;
  font-weight: 700;
  letter-spacing: 5px;
  color: var(--sc-primary);
  text-decoration: none;
}
.sc-brand.centered { margin: 0 auto; }
.sc-nav {
  display: flex;
  gap: 8px;
  flex: 1;
}
.sc-nav-link {
  font-size: 15px;
  font-weight: 600;
  color: var(--sc-muted);
  text-decoration: none;
  padding: 8px 14px;
  border-radius: 10px;
  transition: color 0.2s ease, background 0.2s ease;
}
.sc-nav-link:hover { color: var(--sc-primary); background: var(--sc-primary-soft); }
.sc-nav-link.active { color: var(--sc-primary); background: var(--sc-primary-soft); }
.sc-logout {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--sc-muted);
  font-size: 14px;
  font-weight: 600;
  text-decoration: none;
}
.sc-logout:hover { color: var(--sc-primary); }
.sc-logout .material-icons { font-size: 20px; vertical-align: middle; }

.sc-theme-toggle {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  margin-left: auto;
  width: 38px;
  height: 38px;
  padding: 0;
  border: 1px solid var(--sc-border);
  border-radius: 50%;
  background: var(--sc-surface);
  color: var(--sc-muted);
  cursor: pointer;
  transition: color 0.2s ease, border-color 0.2s ease, background 0.2s ease;
}
.sc-theme-toggle:hover { color: var(--sc-primary); border-color: var(--sc-primary); }
.sc-theme-toggle .material-icons { font-size: 20px; }

.sc-burger {
  display: none;
  position: relative;
  width: 28px;
  height: 22px;
  margin-left: auto;
  cursor: pointer;
}
.sc-burger span {
  display: block;
  position: absolute;
  height: 3px;
  width: 100%;
  background: var(--sc-ink-2);
  border-radius: 3px;
  left: 0;
  transition: 0.25s ease;
}
.sc-burger span:nth-child(1) { top: 0; }
.sc-burger span:nth-child(2), .sc-burger span:nth-child(3) { top: 9px; }
.sc-burger span:nth-child(4) { top: 18px; }
.sc-burger.menuopen span:nth-child(1) { top: 9px; width: 0; left: 50%; }
.sc-burger.menuopen span:nth-child(2) { transform: rotate(45deg); }
.sc-burger.menuopen span:nth-child(3) { transform: rotate(-45deg); }
.sc-burger.menuopen span:nth-child(4) { top: 9px; width: 0; left: 50%; }

.sc-mobile-nav {
  position: absolute;
  top: 64px;
  left: 0;
  right: 0;
  background: var(--sc-surface);
  box-shadow: 0 12px 24px -8px rgba(22, 50, 92, 0.18);
  opacity: 0;
  visibility: hidden;
  transition: opacity 0.25s ease;
  z-index: 29;
  overflow: hidden;
}
.sc-mobile-nav.naviopen { opacity: 1; visibility: visible; }
.sc-mobile-nav a {
  display: block;
  padding: 14px 24px;
  text-align: center;
  color: var(--sc-ink-2);
  text-decoration: none;
  font-weight: 600;
  border-top: 1px solid var(--sc-border-soft);
}
.sc-mobile-nav a:hover { background: var(--sc-primary-soft); color: var(--sc-primary); }
.sc-mobile-logout { color: var(--sc-muted) !important; }

.sc-header-alert {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 8px 16px;
  background: var(--sc-danger-soft);
  color: var(--sc-danger);
  font-size: 13px;
  font-weight: 600;
  text-align: center;
  text-decoration: none;
  border-top: 1px solid var(--sc-border-soft);
}
.sc-header-alert i { font-size: 18px; }

@media (max-width: 850px) {
  .sc-nav, .sc-logout { display: none; }
  .sc-burger { display: block; }
  .sc-brand { margin: 0; }
}
</style>
