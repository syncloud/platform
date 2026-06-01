<template>
  <div class="sc-page" id="block_app">
    <div class="sc-card" v-if="info !== undefined">
      <div class="app-head">
        <img :src="info.app.icon" class="appimg" alt="" @error="(e) => e.target.src = defaultIcon">

        <div class="app-actions" v-if="!progress">
          <button id="btn_open" :data-url="info.app.url" class="app-btn app-btn-primary"
                  @click="open" :title="$t('app.open')" :aria-label="$t('app.open')"
                  v-if="info.installed_version !== null">
            <i class="material-icons app-btn-glyph">launch</i><span class="app-btn-label">{{ $t('app.open') }}</span>
          </button>
          <button id="btn_install" class="app-btn app-btn-primary"
                  data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Installing..."
                  @click="install" :title="$t('app.install', { version: info.current_version })" :aria-label="$t('app.install', { version: info.current_version })"
                  v-if="info.installed_version === null && !info.local_install">
            <i class="material-icons app-btn-glyph">download</i><span class="app-btn-label">{{ $t('app.install', { version: info.current_version }) }}</span>
          </button>
          <button id="btn_upgrade" class="app-btn app-btn-upgrade"
                  data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Upgrading..."
                  @click="upgrade" :title="$t('app.upgrade', { version: info.current_version })" :aria-label="$t('app.upgrade', { version: info.current_version })"
                  v-if="info.installed_version !== null && !info.local_install && info.installed_version !== info.current_version">
            <i class="material-icons app-btn-glyph">upgrade</i><span class="app-btn-label">{{ $t('app.upgrade', { version: info.current_version }) }}</span>
          </button>
          <button id="btn_remove" data-testid="btn_remove" class="app-btn app-btn-danger"
                  data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Removing..."
                  @click="remove" :title="$t('app.remove')" :aria-label="$t('app.remove')"
                  v-if="info.installed_version !== null">
            <i class="material-icons app-btn-glyph">delete</i><span class="app-btn-label">{{ $t('app.remove') }}</span>
          </button>
          <button id="btn_backup" class="app-btn app-btn-tonal"
                  data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Creating backup..."
                  @click="backupConfirm" :title="$t('app.backup')" :aria-label="$t('app.backup')"
                  v-if="info.installed_version !== null">
            <i class="material-icons app-btn-glyph">backup</i><span class="app-btn-label">{{ $t('app.backup') }}</span>
          </button>
        </div>

        <div v-if="progress" id="progress" class="app-progress" :title="progressSummary">
          <svg class="app-ring" :class="{ spin: progressIndeterminate }" viewBox="0 0 44 44">
            <circle class="app-ring-track" cx="22" cy="22" r="19"/>
            <circle class="app-ring-fill" cx="22" cy="22" r="19"
                    :stroke-dasharray="ringCircumference"
                    :stroke-dashoffset="progressIndeterminate ? ringCircumference * 0.72 : ringCircumference * (1 - progressPercentage / 100)"/>
          </svg>
          <span id="progress_summary" class="sr-only">{{ progressSummary }}</span>
        </div>
      </div>

      <h1 id="app_name" data-testid="app_name" class="app-name">{{ info.app.name }}</h1>
      <div class="app-meta">
        <span v-if="info.installed_version !== null && !progress" class="app-version">{{ $t('app.version') }} {{ info.installed_version }}</span>
        <span v-if="info.local_install" id="local_install_badge" data-testid="local_install_badge" class="app-badge">{{ $t('app.localInstall') }}</span>
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
    ringCircumference () {
      return 2 * Math.PI * 19
    },
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
  flex: 1;
  min-width: 0;
  align-self: center;
  display: flex;
  justify-content: center;
  align-items: center;
}

/* round (circular) progress, no words */
.app-ring {
  width: 56px;
  height: 56px;
}
.app-ring.spin {
  animation: app-ring-rotate 0.9s linear infinite;
}
.app-ring-track {
  fill: none;
  stroke: var(--sc-primary-soft);
  stroke-width: 4;
}
.app-ring-fill {
  fill: none;
  stroke: var(--sc-primary);
  stroke-width: 4;
  stroke-linecap: round;
  transform: rotate(-90deg);
  transform-origin: 50% 50%;
  transition: stroke-dashoffset 0.3s ease;
}
@keyframes app-ring-rotate {
  to { transform: rotate(360deg); }
}

.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border: 0;
}

/* modern app-store style buttons: glyph + label pills, icon-only on mobile */
.app-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-family: var(--sc-font);
  font-weight: 600;
  font-size: 14px;
  padding: 9px 20px;
  border-radius: 999px;
  border: none;
  cursor: pointer;
  white-space: nowrap;
  transition: transform 0.1s ease, box-shadow 0.2s ease, filter 0.15s ease, background 0.2s ease;
}
.app-btn:active { transform: translateY(1px); }
.app-btn:disabled { opacity: 0.5; cursor: not-allowed; }
.app-btn-glyph { font-size: 18px; line-height: 1; }
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
  .app-actions { gap: 12px; }
  .app-btn {
    width: 48px;
    height: 48px;
    padding: 0;
    justify-content: center;
    border-radius: 50%;
  }
  .app-btn-label { display: none; }
  .app-btn-glyph { font-size: 22px; }
  .app-ring { width: 48px; height: 48px; }
}
</style>
