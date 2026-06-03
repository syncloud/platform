<template>

  <div class="sc-page">
    <div class="sc-card" id="block1">
      <h1 class="sc-title">{{ $t('updates.title') }}</h1>
      <div id="block_updates">
        <div class="setline update-line">
          <span class="sc-row-label">{{ $t('updates.check') }}</span>
          <s-button id="btn_check_updates" type="primary" :disabled="busy" @click="check">{{ $t('updates.checkButton') }}</s-button>
        </div>
        <div class="setline">
          <span class="sc-row-label">{{ $t('updates.note') }}</span>
        </div>

        <div class="setline update-line">
          <span class="sc-row-label">{{ $t('updates.system') }} <span id="txt_platform_version">{{ platformVersion }}</span></span>
          <s-button
            v-if="!platform.progress && platformVersion !== platformVersionAvailable"
            id="btn_platform_upgrade"
            type="success"
            :disabled="busy"
            @click="upgradePlatform"
          >
            {{ $t('updates.upgradeTo', { version: platformVersionAvailable }) }}
          </s-button>
          <div class="update-progress" id="platform_progress" v-if="platform.progress">
            <span class="update-progress-summary" id="platform_progress_summary">{{ platform.summary }}</span>
            <s-progress
              :percentage="platform.percentage"
              :indeterminate="platform.indeterminate"
              :show-text="false"
              :stroke-width="6"
            />
          </div>
        </div>

        <div class="setline update-line">
          <span class="sc-row-label">{{ $t('updates.installer') }} <span id="txt_installer_version">{{ installerVersion }}</span></span>
          <s-button
            v-if="!installer.progress && installerVersion !== installerVersionAvailable"
            id="btn_installer_upgrade"
            type="success"
            :disabled="busy"
            @click="upgradeInstaller"
          >
            {{ $t('updates.upgradeTo', { version: installerVersionAvailable }) }}
          </s-button>
          <div class="update-progress" id="installer_progress" v-if="installer.progress">
            <span class="update-progress-summary" id="installer_progress_summary">{{ installer.summary }}</span>
            <s-progress
              :percentage="installer.percentage"
              :indeterminate="installer.indeterminate"
              :show-text="false"
              :stroke-width="6"
            />
          </div>
        </div>
      </div>
    </div>
  </div>

  <Error ref="error"/>

</template>

<script>
import axios from 'axios'
import * as Common from '../js/common.js'
import Error from '../components/Error.vue'
import Loading from '../util/loading'

export default {
  name: 'Updates',
  components: {
    Error
  },
  data () {
    return {
      platformVersion: undefined,
      platformVersionAvailable: undefined,
      installerVersion: undefined,
      installerVersionAvailable: undefined,
      platform: { progress: false, summary: '', percentage: 0, indeterminate: true },
      installer: { progress: false, summary: '', percentage: 0, indeterminate: true },
      loading: undefined
    }
  },
  computed: {
    busy () {
      return this.platform.progress || this.installer.progress
    }
  },
  mounted () {
    this.progressShow()
    this.versions(() => this.resume())
  },
  methods: {
    progressShow () {
      this.loading = Loading.service({ lock: true, text: this.$t('common.loading'), background: 'rgba(0, 0, 0, 0.7)' })
    },
    progressHide () {
      if (this.loading) {
        this.loading.close()
      }
    },
    onError (err) {
      this.platform.progress = false
      this.installer.progress = false
      this.progressHide()
      this.$refs.error.showAxios(err)
    },
    resume () {
      axios.get(Common.INSTALLER_STATUS_URL)
        .then(resp => {
          const progress = resp.data.data.progress
          if (progress && progress.platform) {
            this.trackPlatform()
          }
        })
        .catch(() => {})
      axios.get(Common.JOB_STATUS_URL)
        .then(resp => {
          const job = resp.data.data
          if (job.status !== 'Idle' && job.name === 'installer.upgrade') {
            this.trackInstaller()
          }
        })
        .catch(() => {})
    },
    trackPlatform () {
      this.platform.progress = true
      Common.runAfterJobIsComplete(
        setTimeout,
        () => { this.platform.progress = false; this.versions() },
        this.onError,
        Common.INSTALLER_STATUS_URL,
        (response) => {
          const all = response.data.data.progress
          const progress = all ? all.platform : undefined
          if (!progress) {
            return false
          }
          this.platform.progress = true
          this.platform.summary = progress.summary
          if (progress.indeterminate) {
            this.platform.indeterminate = true
            this.platform.percentage = 20
          } else {
            this.platform.indeterminate = false
            this.platform.percentage = progress.percentage
          }
          return true
        }
      )
    },
    trackInstaller () {
      this.installer.progress = true
      this.installer.indeterminate = true
      this.installer.summary = this.$t('updates.upgrading')
      Common.runAfterJobIsComplete(
        setTimeout,
        () => { this.installer.progress = false; this.versions() },
        this.onError,
        Common.JOB_STATUS_URL,
        (response) => {
          const job = response.data.data
          if (job.status !== 'Idle' && job.name === 'installer.upgrade') {
            this.installer.progress = true
            this.installer.indeterminate = true
            this.installer.summary = this.$t('updates.upgrading')
            return true
          }
          return false
        }
      )
    },
    upgradePlatform () {
      this.platform.progress = true
      this.platform.indeterminate = true
      this.platform.percentage = 0
      this.platform.summary = this.$t('updates.upgrading')
      axios.post('/rest/app/upgrade', { app_id: 'platform' })
        .then((resp) => {
          Common.checkForServiceError(resp.data, this.trackPlatform, this.onError)
        })
        .catch(this.onError)
    },
    upgradeInstaller () {
      this.installer.progress = true
      this.installer.indeterminate = true
      this.installer.percentage = 0
      this.installer.summary = this.$t('updates.upgrading')
      axios.post('/rest/installer/upgrade')
        .then((resp) => {
          Common.checkForServiceError(resp.data, this.trackInstaller, this.onError)
        })
        .catch(this.onError)
    },
    check () {
      this.progressShow()
      this.versions()
    },
    versions (onComplete) {
      Promise.all([
        axios.get('/rest/app', { params: { app_id: 'platform' } }),
        axios.get('/rest/installer/version')]
      )
        .then((results) => {
          this.platformVersion = results[0].data.data.installed_version
          this.platformVersionAvailable = results[0].data.data.current_version
          this.installerVersion = results[1].data.data.installed_version
          this.installerVersionAvailable = results[1].data.data.store_version
          this.progressHide()
          if (onComplete) {
            onComplete()
          }
        })
        .catch(err => {
          this.progressHide()
          this.$refs.error.showAxios(err)
        })
    }
  }
}
</script>
<style scoped>
.update-line {
  display: flex;
  justify-content: space-between;
  align-items: center;
  min-height: 40px;
}
.update-line .el-button {
  min-width: 120px;
}
.update-progress {
  flex: 0 0 auto;
  width: 120px;
}
.update-progress-summary {
  display: block;
  margin-bottom: 4px;
  color: var(--sc-muted);
  font-size: 12px;
  text-align: center;
}
</style>
