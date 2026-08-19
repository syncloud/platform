<template>
  <div class="sc-page">
    <div class="sc-card" id="block1">
      <h1 class="sc-title">{{ $t('system.title') }}</h1>
      <div class="sc-row">
        <span class="sc-row-label">{{ $t('system.restartLabel') }}</span>
        <s-button id="restart" type="primary" @click="restartConfirmVisible = true">{{ $t('system.restart') }}</s-button>
      </div>
      <div class="sc-row">
        <span class="sc-row-label">{{ $t('system.shutdownLabel') }}</span>
        <s-button id="shutdown" type="danger" @click="shutdownConfirmVisible = true">{{ $t('system.shutdown') }}</s-button>
      </div>
    </div>
  </div>

  <Dialog :visible="restartConfirmVisible" @confirm="restart" @cancel="restartConfirmVisible = false">
    <template #title>{{ $t('system.restart') }}</template>
    <template #text>{{ $t('system.restartConfirm') }}</template>
  </Dialog>

  <Dialog :visible="shutdownConfirmVisible" @confirm="shutdown" @cancel="shutdownConfirmVisible = false">
    <template #title>{{ $t('system.shutdown') }}</template>
    <template #text>{{ $t('system.shutdownConfirm') }}</template>
  </Dialog>

  <div v-if="progressVisible" class="s-modal-overlay" data-testid="progress">
    <div class="s-modal" role="dialog">
      <h4 class="modal-title" data-testid="progress-title">{{ progressTitle }}</h4>
      <div class="s-modal-body progress-body">
        <div v-if="waiting" class="sc-progress-spinner" data-testid="progress-spinner"></div>
        <span data-testid="progress-text">{{ progressText }}</span>
      </div>
      <div class="s-modal-footer">
        <button v-if="stage === 'timeout'" id="btn_reload" data-testid="btn_reload" class="sc-btn sc-btn-primary"
                type="button" @click="reloadPage">{{ $t('system.reload') }}</button>
        <button v-if="!waiting" id="btn_close" data-testid="btn_close" class="sc-btn sc-btn-ghost"
                type="button" @click="progressVisible = false">{{ $t('common.close') }}</button>
      </div>
    </div>
  </div>

  <Error ref="error"/>
</template>

<script>
import axios from 'axios'
import Error from '../components/Error.vue'
import Dialog from '../components/Dialog.vue'

const ONLINE_CONFIRMATIONS = 2

export default {
  name: 'System',
  components: {
    Error,
    Dialog
  },
  props: {
    pollIntervalMs: {
      type: Number,
      default: 2000
    },
    pollAttempts: {
      type: Number,
      default: 300
    }
  },
  data () {
    return {
      restartConfirmVisible: false,
      shutdownConfirmVisible: false,
      progressVisible: false,
      action: 'restart',
      stage: 'stopping',
      uptimeBefore: undefined,
      stopped: false
    }
  },
  computed: {
    waiting () {
      return this.stage === 'stopping' || this.stage === 'starting' || this.stage === 'done'
    },
    progressTitle () {
      return this.action === 'restart' ? this.$t('system.restarting') : this.$t('system.shuttingDown')
    },
    progressText () {
      return this.$t(`system.stage.${this.action}.${this.stage}`)
    }
  },
  beforeUnmount () {
    this.stopped = true
  },
  methods: {
    restart () {
      this.restartConfirmVisible = false
      return this.run('restart', '/rest/restart', () => this.waitForRestart())
    },
    shutdown () {
      this.shutdownConfirmVisible = false
      return this.run('shutdown', '/rest/shutdown', () => this.waitForShutdown())
    },
    async run (action, url, wait) {
      this.action = action
      this.stage = 'stopping'
      this.progressVisible = true
      this.uptimeBefore = await this.uptime()
      try {
        await axios.post(url)
      } catch (err) {
        if (err.response) {
          this.progressVisible = false
          this.$refs.error.showAxios(err)
          return
        }
      }
      await wait()
    },
    async waitForRestart () {
      const offline = await this.waitOffline()
      if (this.stopped) {
        return
      }
      if (offline === 'timeout') {
        this.stage = 'timeout'
        return
      }
      if (offline === 'offline') {
        this.stage = 'starting'
        const online = await this.waitOnline()
        if (this.stopped) {
          return
        }
        if (!online) {
          this.stage = 'timeout'
          return
        }
      }
      this.stage = 'done'
      this.reloadPage()
    },
    async waitForShutdown () {
      const offline = await this.waitOffline()
      if (this.stopped) {
        return
      }
      this.stage = offline === 'timeout' ? 'timeout' : 'off'
    },
    async waitOffline () {
      for (let attempt = 0; attempt < this.pollAttempts; attempt++) {
        if (this.stopped) {
          return 'stopped'
        }
        const probe = await this.probe()
        if (!probe.reachable) {
          return 'offline'
        }
        if (this.restarted(probe.uptime)) {
          return 'restarted'
        }
        await this.sleep(this.pollIntervalMs)
      }
      return 'timeout'
    },
    async waitOnline () {
      let reachable = 0
      for (let attempt = 0; attempt < this.pollAttempts; attempt++) {
        if (this.stopped) {
          return false
        }
        const probe = await this.probe()
        reachable = probe.reachable ? reachable + 1 : 0
        if (reachable >= ONLINE_CONFIRMATIONS) {
          return true
        }
        await this.sleep(this.pollIntervalMs)
      }
      return false
    },
    restarted (uptime) {
      return uptime !== undefined && this.uptimeBefore !== undefined && uptime < this.uptimeBefore
    },
    probe () {
      return axios.get('/rest/uptime', { timeout: 5000 })
        .then(response => ({ reachable: true, uptime: response.data.data }))
        .catch(err => {
          const status = err.response ? err.response.status : 0
          return { reachable: status > 0 && status < 500, uptime: undefined }
        })
    },
    uptime () {
      return this.probe().then(probe => probe.uptime)
    },
    sleep (ms) {
      return new Promise(resolve => setTimeout(resolve, ms))
    },
    reloadPage () {
      window.location.reload()
    }
  }
}
</script>

<style scoped>
.progress-body {
  display: flex;
  align-items: center;
  gap: 12px;
}
.sc-progress-spinner {
  flex: none;
  width: 22px;
  height: 22px;
  border: 3px solid var(--sc-faint);
  border-top-color: var(--sc-primary);
  border-radius: 50%;
  animation: sc-spin 0.9s linear infinite;
}
</style>
