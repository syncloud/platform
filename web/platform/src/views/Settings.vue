<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12" id="block1">
        <h1>{{ $t('settings.title') }}</h1>
        <div class="row-no-gutters settingsblock">
          <div class="colapp app">
            <router-link to="/activation" id="activation" class="app-link">
              <i class="material-icons settings-icon">perm_identity</i>
              <div class="setting-name"><span class="withline">{{ $t('settings.activation') }}</span></div>
            </router-link>
          </div>

          <div class="colapp app">
            <router-link to="/network" id="network" class="app-link">
              <i class="material-icons settings-icon">router</i>
              <div class="setting-name"><span class="withline">{{ $t('settings.network') }}</span></div>
            </router-link>
          </div>

          <div class="colapp app">
            <router-link to="/access" id="access" class="app-link">
              <i class="material-icons settings-icon">cloud_queue</i>
              <div class="setting-name"><span class="withline">{{ $t('settings.access') }}</span></div>
            </router-link>
          </div>

          <div class="colapp app">
            <router-link to="/internalmemory" id="internalmemory" class="app-link">
              <i class="material-icons settings-icon">memory</i>
              <div class="setting-name"><span class="withline">{{ $t('settings.internalMemory') }}</span></div>
            </router-link>
          </div>

          <div class="colapp app">
            <router-link to="/storage" id="storage" class="app-link">
              <i class="material-icons settings-icon">storage</i>
              <div class="setting-name"><span class="withline">{{ $t('settings.storage') }}</span></div>
            </router-link>
          </div>

          <div class="colapp app">
            <router-link to="/updates" id="updates" class="app-link">
              <i class="material-icons settings-icon">update</i>
              <div class="setting-name"><span class="withline">{{ $t('settings.updates') }}</span></div>
            </router-link>
          </div>

          <div class="colapp app">
            <router-link to="/support" id="support" class="app-link">
              <i class="material-icons settings-icon">help_outline</i>

              <div class="setting-name"><span class="withline">{{ $t('settings.support') }}</span></div>
            </router-link>
          </div>

          <div class="colapp app">
            <router-link to="/backup" id="backup" class="app-link">
              <i class="material-icons settings-icon">archive</i>
              <div class="setting-name"><span class="withline">{{ $t('settings.backup') }}</span></div>
            </router-link>
          </div>

          <div class="colapp app">
            <router-link to="/certificate" id="certificate" class="app-link">
              <i class="material-icons settings-icon">security</i>
              <div class="setting-name"><span class="withline">{{ $t('settings.certificate') }}</span></div>
            </router-link>
          </div>

          <div class="colapp app">
            <router-link to="/twofactor" id="twofactor" class="app-link">
              <i class="material-icons settings-icon">verified_user</i>
              <div class="setting-name"><span class="withline">{{ $t('settings.twoFactor') }}</span></div>
            </router-link>
          </div>

          <div class="colapp app">
            <router-link to="/logs" id="logs" class="app-link">
              <i class="material-icons-outlined settings-icon">feed</i>
              <div class="setting-name"><span class="withline">{{ $t('settings.logs') }}</span></div>
            </router-link>
          </div>

          <div class="colapp app">
            <router-link to="/customproxy" id="customproxy" class="app-link">
              <i class="material-icons settings-icon">swap_horiz</i>
              <div class="setting-name"><span class="withline">{{ $t('settings.customProxy') }}</span></div>
            </router-link>
          </div>

          <div class="colapp app">
            <router-link to="/system" id="system" class="app-link">
              <i class="material-icons settings-icon">settings_power</i>
              <div class="setting-name"><span class="withline">{{ $t('settings.system') }}</span></div>
            </router-link>
          </div>

          <div class="colapp app">
            <router-link to="/locale" id="locale" class="app-link">
              <i class="material-icons settings-icon">translate</i>
              <div class="setting-name"><span class="withline">{{ $t('settings.locale') }}</span></div>
            </router-link>
          </div>

        </div>
      </div>
    </div>
  </div>

</template>

<script>
export default {
  name: 'Settings'
}
</script>
<style>
@import '../style/site.css';
@import 'material-icons/iconfont/material-icons.css';
.settings-icon {
  font-size: 72px !important;
  padding-top: 20px !important;
  padding-bottom: 10px !important;
}

</style>
