<template>
  <div class="sc-page" id="block_app">
    <div class="sc-card" v-if="info !== undefined">
      <div class="app-head">
        <img :src="info.app.icon" class="appimg" alt="" @error="(e) => e.target.src = defaultIcon">

        <div class="app-actions" v-if="!progress">
          <button id="btn_open" :data-url="info.app.url" class="app-btn app-btn-primary"
                  @click="open"
                  v-if="info.installed_version !== null">
            {{ $t('app.open') }}
          </button>
          <button id="btn_install" class="app-btn app-btn-primary"
                  data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Installing..."
                  @click="install"
                  v-if="info.installed_version === null && !info.local_install">
            {{ $t('app.install', { version: info.current_version }) }}
          </button>
          <button id="btn_upgrade" class="app-btn app-btn-upgrade"
                  data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Upgrading..."
                  @click="upgrade"
                  v-if="info.installed_version !== null && !info.local_install && info.installed_version !== info.current_version">
            {{ $t('app.upgrade', { version: info.current_version }) }}
          </button>
          <button id="btn_remove" data-testid="btn_remove" class="app-btn app-btn-danger"
                  data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Removing..."
                  @click="remove"
                  v-if="info.installed_version !== null">
            {{ $t('app.remove') }}
          </button>
          <button id="btn_backup" class="app-btn app-btn-tonal"
                  data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Creating backup..."
                  @click="backupConfirm"
                  v-if="info.installed_version !== null">
            {{ $t('app.backup') }}
          </button>
        </div>
      </div>

      <h1 id="app_name" data-testid="app_name" class="app-name">{{ info.app.name }}</h1>
      <div class="app-meta">
        <span v-if="info.installed_version !== null && !progress" class="app-version">{{ $t('app.version') }} {{ info.installed_version }}</span>
        <span v-if="info.local_install" id="local_install_badge" data-testid="local_install_badge" class="app-badge">{{ $t('app.localInstall') }}</span>
      </div>

      <div v-if="progress" id="progress" class="app-progress">
        <div id="progress_summary" class="app-progress-summary">{{ progressSummary }}</div>
        <el-progress :show-text="false" :percentage="progressPercentage" :indeterminate="progressIndeterminate" :stroke-width="8"/>
      </div>

      <p class="app-description" v-if="showDescription">{{ info.app.description }}</p>
    </div>
  </div>

  <Dialog :visible="appActionConfirmationVisible" id="app_confirmation" @confirm="confirm"
                @cancel="appActionConfirmationVisible = false">
    <template v-slot:title>
      <span id="confirm_caption">{{ action }}</span>
    </template>
    <template v-slot:text>
      <div class="bodymod">
        <div class="btext">
          {{ $t('app.areYouSure') }}
        </div>
      </div>
    </template>
  </Dialog>

  <Dialog :visible="backupConfirmationVisible" id="backup_confirmation" @confirm="backup"
                @cancel="backupConfirmationVisible = false">
    <template v-slot:title>
      <span id="confirm_caption">{{ $t('app.backupTitle') }}</span>
    </template>
    <template v-slot:text>
      <div class="bodymod">
        <div class="btext">
          {{ $t('app.backupConfirmText') }}
        </div>
      </div>
    </template>
  </Dialog>

  <Error ref="error"/>

</template>

<script>
import axios from 'axios'
import Error from '../components/Error.vue'
import * as Common from '../js/common.js'
import Dialog from '../components/Dialog.vue'

export default {
  name: 'App',
  data () {
    return {
      defaultIcon: '/images/default-app.svg',
      info: undefined,
      appId: undefined,
      action: '',
      progress: true,
      progressPercentage: 0,
      progressSummary: '',
      progressIndeterminate: true,
      appActionConfirmationVisible: false,
      backupConfirmationVisible: false
    }
  },
  components: {
    Dialog,
    Error
  },
  computed: {
    showDescription () {
      const d = (this.info && this.info.app && this.info.app.description ? this.info.app.description : '').trim()
      const name = (this.info && this.info.app && this.info.app.name ? this.info.app.name : '').trim()
      return d !== '' && d.toLowerCase() !== name.toLowerCase()
    }
  },
  mounted () {
    this.progressShow()
    this.appId = this.$route.query.id
    this.loadApp()
      .then(() => this.status())
  },
  methods: {
    progressShow () {
      this.progress = true
    },
    progressHide () {
      this.progressSummary = ''
      this.progress = false
    },
    async loadApp() {
      try {
        let resp = await axios
          .get('/rest/app', {params: {app_id: this.appId}})
        this.info = resp.data.data
      } catch (err) {
        this.$refs.error.showAxios(err)
      }
    },
    open () {
      window.location.href = this.info.app.url
    },
    install () {
      this.action = this.$t('app.installAction')
      this.actionUrl = '/rest/app/install'
      this.appActionConfirmationVisible = true
    },
    upgrade () {
      this.action = this.$t('app.upgradeAction')
      this.actionUrl = '/rest/app/upgrade'
      this.appActionConfirmationVisible = true
    },
    remove () {
      this.action = this.$t('app.removeAction')
      this.actionUrl = '/rest/app/remove'
      this.appActionConfirmationVisible = true
    },
    backupConfirm () {
      this.backupConfirmationVisible = true
    },
    backup () {
      this.backupConfirmationVisible = false
      this.progressShow()
      this.progressSummary = this.$t('app.creatingBackup')

      const error = this.$refs.error
      const onError = (err) => {
        this.progressHide()
        error.showAxios(err)
      }
      axios.post('/rest/backup/create', { app: this.appId })
        .then((resp) => {
          Common.checkForServiceError(resp.data, () => {
            Common.runAfterJobIsComplete(
              setTimeout,
              () => { this.loadApp().then(() => this.progressHide()) },
              onError,
              Common.JOB_STATUS_URL,
              Common.JOB_STATUS_PREDICATE)
          }, onError)
        })
        .catch(onError)
    },
    status () {
      const error = this.$refs.error
      const onError = (err) => {
        this.progressHide()
        error.showAxios(err)
      }
      Common.runAfterJobIsComplete(
        setTimeout,
        () => { this.loadApp().then(() => this.progressHide()) },
        onError,
        Common.INSTALLER_STATUS_URL,
        (response) => {
          if (!response.data.data.progress) {
            return false
          }
          const progress = response.data.data.progress[this.appId]
          if (progress) {
            this.progressShow()
            this.progressSummary = progress.summary
            if (progress.indeterminate) {
              this.progressIndeterminate = true
              this.progressPercentage = 20
            } else {
              this.progressIndeterminate = false
              this.progressPercentage = progress.percentage
            }
            return true
          } else {
            return false
          }
        }
      )
    },
    confirm () {
      this.appActionConfirmationVisible = false
      this.progressShow()
      const error = this.$refs.error
      const onError = (err) => {
        this.progressHide()
        error.showAxios(err)
      }
      axios.post(this.actionUrl, { app_id: this.appId })
        .then(resp => {
          Common.checkForServiceError(resp.data, this.status, onError)
        })
        .catch(onError)
    }
  }
}
</script>
<style scoped>
.app-head {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 18px;
}
.appimg {
  width: 88px;
  height: 88px;
  border-radius: 22px;
  flex-shrink: 0;
}
.app-actions {
  flex: 1;
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  justify-content: flex-end;
  align-items: center;
  min-width: 0;
}

.app-name {
  font-size: 24px;
  font-weight: 700;
  letter-spacing: -0.3px;
  text-align: left;
  margin: 0 0 4px;
  color: var(--sc-ink);
}
.app-meta {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
  margin-bottom: 14px;
}
.app-version {
  font-size: 14px;
  color: var(--sc-muted);
  font-variant-numeric: tabular-nums;
}
.app-badge {
  font-size: 12px;
  font-weight: 600;
  color: var(--sc-primary);
  background: var(--sc-primary-soft);
  padding: 3px 10px;
  border-radius: 999px;
}
.app-description {
  text-align: left;
  font-size: 15px;
  line-height: 1.6;
  color: var(--sc-ink-2);
  margin: 4px 0 0;
  white-space: pre-line;
}
.app-progress {
  margin: 6px 0 18px;
}
.app-progress-summary {
  text-align: center;
  font-size: 13px;
  color: var(--sc-muted);
  margin-bottom: 8px;
}

/* modern app-store style pill buttons */
.app-btn {
  font-family: var(--sc-font);
  font-weight: 600;
  font-size: 14px;
  padding: 9px 22px;
  border-radius: 999px;
  border: none;
  cursor: pointer;
  white-space: nowrap;
  transition: transform 0.1s ease, box-shadow 0.2s ease, filter 0.15s ease, background 0.2s ease;
}
.app-btn:active { transform: translateY(1px); }
.app-btn:disabled { opacity: 0.5; cursor: not-allowed; }
.app-btn-primary {
  background: linear-gradient(135deg, var(--sc-primary) 0%, var(--sc-primary-dark) 100%);
  color: #fff;
  box-shadow: 0 6px 16px -8px rgba(43, 123, 214, 0.6);
}
.app-btn-primary:hover { filter: brightness(1.06); }
.app-btn-upgrade {
  background: linear-gradient(135deg, #34b566 0%, #2faa5d 100%);
  color: #fff;
  box-shadow: 0 6px 16px -8px rgba(47, 170, 93, 0.6);
}
.app-btn-upgrade:hover { filter: brightness(1.04); }
.app-btn-tonal {
  background: var(--sc-primary-soft);
  color: var(--sc-primary);
}
.app-btn-tonal:hover { background: #d9e9fb; }
.app-btn-danger {
  background: #fdeced;
  color: var(--sc-danger);
}
.app-btn-danger:hover { background: #fbdcde; }

@media (max-width: 600px) {
  .appimg { width: 72px; height: 72px; border-radius: 18px; }
  .app-btn { padding: 8px 16px; font-size: 13px; }
}
</style>
