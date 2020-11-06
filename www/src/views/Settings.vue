<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12" id="block1">
        <h1>Settings</h1>
        <div class="row-no-gutters settingsblock">
          <div class="colapp app">
            <router-link to="/activation" id="activation" class="app-link">
              <i class="material-icons settings-icon">perm_identity</i>
              <div class="setting-name"><span class="withline">Activation</span></div>
            </router-link>
          </div>

          <div class="colapp app">
            <router-link to="/network" id="network" class="app-link">
              <i class="material-icons settings-icon">router</i>
              <div class="setting-name"><span class="withline">Network</span></div>
            </router-link>
          </div>

          <div class="colapp app">
            <router-link to="/access" id="access" class="app-link">
              <i class="material-icons settings-icon">cloud_queue</i>
              <div class="setting-name"><span class="withline">Access</span></div>
            </router-link>
          </div>

          <div class="colapp app">
            <router-link to="/internalmemory" id="internalmemory" class="app-link">
              <i class="material-icons settings-icon">memory</i>
              <div class="setting-name"><span class="withline">Internal Memory</span></div>
            </router-link>
          </div>

          <div class="colapp app">
            <router-link to="/storage" id="storage" class="app-link">
              <i class="material-icons settings-icon">storage</i>
              <div class="setting-name"><span class="withline">Storage</span></div>
            </router-link>
          </div>

          <div class="colapp app">
            <router-link to="/updates" id="updates" class="app-link">
              <i class="material-icons settings-icon">update</i>
              <div class="setting-name"><span class="withline">Updates</span></div>
            </router-link>
          </div>

          <div class="colapp app">
            <router-link to="/support" id="support" class="app-link">
              <i class="material-icons settings-icon">help_outline</i>

              <div class="setting-name"><span class="withline">Support</span></div>
            </router-link>
          </div>

          <div class="colapp app">
            <router-link to="/backup" id="backup" class="app-link">
              <i class="material-icons settings-icon">archive</i>
              <div class="setting-name"><span class="withline">Backup</span></div>
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
@import '../style/material-icons.css';
</style>
