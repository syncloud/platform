<template>
  <div class="sc-page">
    <div class="sc-card sc-card-wide" id="block1">
      <router-link to="/storage" class="sc-link back-link" data-testid="storage-back">{{ $t('storage.back') }}</router-link>
      <h1 class="sc-title">{{ $t('storage.checks') }}</h1>

      <div class="check-block" data-testid="checks-section">
        <div class="check-head">
          <span class="check-name">{{ $t('storage.diskCheck') }}
            <button @click="helpVisible = true" type="button" class="help-button" data-testid="checks-help">
              <i class='fa fa-question-circle fa-lg'></i>
            </button>
          </span>
          <button class="sc-btn sc-btn-ghost sc-btn-small" type="button" data-testid="checks-start"
                  :disabled="health.busy || running"
                  @click="startChecks">{{ $t('storage.checkNow') }}</button>
        </div>
        <div class="check-results">
          <span class="check-label">{{ $t('storage.selfTest') }}</span>
          <span data-testid="selftest-status" :class="{ 'disk-bad': selfTestsBad(health) }">{{ selfTestsText(t, health) }}</span>
          <template v-if="health.scrub">
            <span class="check-label">{{ $t('storage.scrub') }}</span>
            <span data-testid="scrub-status" :class="{ 'disk-bad': !!health.scrub.errors }">{{ scrubText(t, health.scrub) }}</span>
          </template>
        </div>
        <div class="schedule-row">
          <span>{{ $t('storage.scheduleChecks') }}</span>
          <s-switch data-testid="schedule-checks" v-model="schedule.enabled"/>
        </div>
        <p class="sc-muted schedule-note" data-testid="schedule-time">{{ $t('storage.scheduleTime') }}</p>
      </div>

      <div v-if="health.scrub && anyErrors" class="check-block" data-testid="errors-section">
        <div class="check-head">
          <span class="check-name">{{ $t('storage.errorCounters') }}</span>
          <button class="sc-btn sc-btn-ghost sc-btn-small" type="button" data-testid="errors-reset"
                  @click="resetVisible = true">{{ $t('storage.resetCounters') }}</button>
        </div>
        <p class="check-text">{{ $t('storage.resetText') }}</p>
      </div>

      <p v-if="health.busy" class="sc-muted" data-testid="checks-busy">{{ $t('storage.checksBusy') }}</p>

      <div class="sc-actions">
        <button class="sc-btn sc-btn-success" id="btn_save" data-testid="checks-save" type="button"
                @click="saveSchedule">{{ $t('storage.save') }}</button>
      </div>
    </div>
  </div>
  <Error ref="error"/>

  <Dialog :visible="helpVisible" id="checks-help-dialog" data-testid="checks-help-dialog" @cancel="helpVisible = false"
          :confirm-enabled="false" :cancel-text="$t('common.close')">
    <template v-slot:title>{{ $t('storage.checksHelpTitle') }}</template>
    <template v-slot:text>
      <div class="btext">{{ $t('storage.checksHelpText') }}</div>
    </template>
  </Dialog>

  <Dialog :visible="resetVisible" id="reset-confirmation" data-testid="reset-confirmation" @confirm="resetErrors"
          @cancel="resetVisible = false">
    <template v-slot:title>{{ $t('storage.resetTitle') }}</template>
    <template v-slot:text>
      <div class="btext">{{ $t('storage.resetText') }}</div>
    </template>
  </Dialog>
</template>

<script>
import axios from 'axios'
import Error from '../components/Error.vue'
import Dialog from '../components/Dialog.vue'
import Loading from '../util/loading'
import { checksRunning, scrubText, selfTestsBad, selfTestsText } from '../js/storage'

export default {
  name: 'StorageChecks',
  components: {
    Dialog,
    Error
  },
  data () {
    return {
      health: { disks: [], scrub: null, schedule: { enabled: false }, busy: false },
      schedule: { enabled: false },
      anyErrors: false,
      resetVisible: false,
      helpVisible: false,
      loading: undefined,
      timer: undefined
    }
  },
  computed: {
    running () {
      return checksRunning(this.health)
    }
  },
  mounted () {
    this.reload()
  },
  unmounted () {
    clearTimeout(this.timer)
  },
  methods: {
    t (...args) {
      return this.$t(...args)
    },
    scrubText,
    selfTestsText,
    selfTestsBad,
    reload () {
      clearTimeout(this.timer)
      return Promise.all([axios.get('/rest/storage/health'), axios.get('/rest/storage/disks')])
        .then(([health, disks]) => {
          this.health = health.data.data
          this.schedule = { ...this.health.schedule }
          this.anyErrors = (disks.data.data || []).some(d => d.has_errors)
          if (this.running) {
            this.timer = setTimeout(this.reload, 10000)
          }
        })
        .catch(err => this.$refs.error.showAxios(err))
    },
    startChecks () {
      return axios.post('/rest/storage/check')
        .then(() => this.reload())
        .catch(err => this.$refs.error.showAxios(err))
    },
    resetErrors () {
      this.resetVisible = false
      return axios.post('/rest/storage/errors/reset')
        .then(() => this.reload())
        .catch(err => this.$refs.error.showAxios(err))
    },
    saveSchedule () {
      this.loading = Loading.service({ lock: true, text: this.$t('common.loading'), background: 'rgba(0, 0, 0, 0.7)' })
      return axios.post('/rest/storage/check/schedule', this.schedule)
        .then(() => this.reload())
        .catch(err => this.$refs.error.showAxios(err))
        .finally(() => this.loading.close())
    }
  }
}
</script>
<style scoped>
.back-link {
  display: inline-block;
  margin-bottom: 8px;
}
.check-block {
  border: 1px solid var(--sc-border);
  border-radius: var(--sc-control-radius);
  padding: 12px 14px;
  margin-bottom: 12px;
}
.check-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}
.check-name {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
}
.help-button {
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 0;
  color: var(--sc-faint);
}
.check-text {
  margin: 8px 0 0;
  color: var(--sc-muted);
}
.disk-bad {
  color: var(--sc-danger);
}
.check-results {
  display: grid;
  grid-template-columns: max-content minmax(0, 1fr);
  gap: 6px 12px;
  margin-top: 10px;
  color: var(--sc-muted);
  overflow-wrap: anywhere;
}
.check-label {
  color: var(--sc-faint);
}
.schedule-note {
  margin: 8px 0 0;
}
.schedule-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-top: 12px;
}
.sc-btn-small {
  padding: 4px 10px;
  font-size: 13px;
}
</style>
