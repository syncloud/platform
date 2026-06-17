<template>
  <div class="sc-page">
    <div class="sc-card sc-card-wide" id="block1">
      <h1 class="sc-title">{{ $t('settings.title') }}</h1>
      <div class="sc-grid">
        <router-link v-if="isAdmin" to="/activation" id="activation" class="sc-tile">
          <i class="material-icons sc-tile-icon">perm_identity</i>
          <div class="sc-tile-name">{{ $t('settings.activation') }}</div>
        </router-link>
        <router-link v-if="isAdmin" to="/users" id="users" class="sc-tile">
          <i class="material-icons sc-tile-icon">people</i>
          <div class="sc-tile-name">{{ $t('settings.users') }}</div>
        </router-link>
        <router-link v-if="isAdmin" to="/network" id="network" class="sc-tile">
          <i class="material-icons sc-tile-icon">router</i>
          <div class="sc-tile-name">{{ $t('settings.network') }}</div>
        </router-link>
        <router-link v-if="isAdmin" to="/access" id="access" class="sc-tile">
          <i class="material-icons sc-tile-icon">cloud_queue</i>
          <div class="sc-tile-name">{{ $t('settings.access') }}</div>
        </router-link>
        <router-link v-if="isAdmin" to="/internalmemory" id="internalmemory" class="sc-tile">
          <i class="material-icons sc-tile-icon">memory</i>
          <div class="sc-tile-name">{{ $t('settings.internalMemory') }}</div>
        </router-link>
        <router-link v-if="isAdmin" to="/storage" id="storage" class="sc-tile">
          <i class="material-icons sc-tile-icon">storage</i>
          <div class="sc-tile-name">{{ $t('settings.storage') }}</div>
        </router-link>
        <router-link v-if="isAdmin" to="/updates" id="updates" class="sc-tile">
          <i class="material-icons sc-tile-icon">update</i>
          <div class="sc-tile-name">{{ $t('settings.updates') }}</div>
        </router-link>
        <router-link v-if="isAdmin" to="/support" id="support" class="sc-tile">
          <i class="material-icons sc-tile-icon">help_outline</i>
          <div class="sc-tile-name">{{ $t('settings.support') }}</div>
        </router-link>
        <router-link v-if="isAdmin" to="/backup" id="backup" class="sc-tile">
          <i class="material-icons sc-tile-icon">archive</i>
          <div class="sc-tile-name">{{ $t('settings.backup') }}</div>
        </router-link>
        <router-link v-if="isAdmin" to="/certificate" id="certificate" class="sc-tile">
          <i class="material-icons sc-tile-icon">security</i>
          <div class="sc-tile-name">{{ $t('settings.certificate') }}</div>
        </router-link>
        <router-link v-if="isAdmin" to="/twofactor" id="twofactor" class="sc-tile">
          <i class="material-icons sc-tile-icon">verified_user</i>
          <div class="sc-tile-name">{{ $t('settings.twoFactor') }}</div>
        </router-link>
        <router-link v-if="isAdmin" to="/logs" id="logs" class="sc-tile">
          <i class="material-icons-outlined sc-tile-icon">feed</i>
          <div class="sc-tile-name">{{ $t('settings.logs') }}</div>
        </router-link>
        <router-link v-if="isAdmin" to="/customproxy" id="customproxy" class="sc-tile">
          <i class="material-icons sc-tile-icon">swap_horiz</i>
          <div class="sc-tile-name">{{ $t('settings.customProxy') }}</div>
        </router-link>
        <router-link v-if="isAdmin" to="/system" id="system" class="sc-tile">
          <i class="material-icons sc-tile-icon">settings_power</i>
          <div class="sc-tile-name">{{ $t('settings.system') }}</div>
        </router-link>
        <router-link v-if="isAdmin" to="/locale" id="locale" class="sc-tile">
          <i class="material-icons sc-tile-icon">translate</i>
          <div class="sc-tile-name">{{ $t('settings.locale') }}</div>
        </router-link>
        <router-link v-if="isAdmin" to="/health" id="health" data-testid="nav-health" class="sc-tile">
          <i class="material-icons sc-tile-icon">favorite</i>
          <div class="sc-tile-name">{{ $t('settings.health') }}</div>
        </router-link>
      </div>
    </div>
  </div>

</template>

<script>
import { useAuthStore } from '../stores/auth'

export default {
  name: 'Settings',
  computed: {
    isAdmin () {
      return useAuthStore().admin
    }
  }
}
</script>
