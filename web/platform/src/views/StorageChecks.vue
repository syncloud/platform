<template>
  <div class="sc-page">
    <div class="sc-card sc-card-wide" id="block1">
      <router-link to="/storage" class="sc-link back-link" data-testid="storage-back">{{ $t('storage.back') }}</router-link>
      <h1 class="sc-title">{{ $t('storage.checks') }}</h1>

      <div v-if="health.scrub" class="check-block" data-testid="array-section">
        <div class="check-head">
          <span class="check-name">{{ $t('storage.scrub') }}</span>
          <button class="sc-btn sc-btn-ghost sc-btn-small" type="button" data-testid="scrub-start"
                  :disabled="health.busy || health.scrub.state === 'running'"
                  @click="startScrub">{{ $t('storage.scrubStart') }}</button>
        </div>
        <p class="check-text" data-testid="scrub-status" :class="{ 'disk-bad': !!health.scrub.errors }">{{ scrubText(t, health.scrub) }}</p>
        <div class="schedule-row">
          <span>{{ $t('storage.scheduleScrub') }}</span>
          <s-switch data-testid="schedule-scrub" :model-value="schedule.scrub" @update:model-value="value => saveSchedule('scrub', value)"/>
        </div>
      </div>

      <div v-if="health.scrub && anyErrors" class="check-block" data-testid="errors-section">
        <div class="check-head">
          <span class="check-name">{{ $t('storage.errorCounters') }}</span>
          <button class="sc-btn sc-btn-ghost sc-btn-small" type="button" data-testid="errors-reset"
                  @click="resetVisible = true">{{ $t('storage.resetCounters') }}</button>
        </div>
        <p class="check-text">{{ $t('storage.resetText') }}</p>
      </div>

      <div class="check-block" data-testid="smart-section">
        <div class="check-head">
          <span class="check-name">{{ $t('storage.smart') }}</span>
        </div>
        <div class="schedule-row">
          <span>{{ $t('storage.scheduleSmart') }}</span>
          <s-switch data-testid="schedule-smart" :model-value="schedule.smart" @update:model-value="value => saveSchedule('smart', value)"/>
        </div>
      </div>

      <p v-if="health.busy" class="sc-muted" data-testid="checks-busy">{{ $t('storage.checksBusy') }}</p>
    </div>
  </div>
  <Error ref="error"/>

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
import { scrubText } from '../js/storage'

export default {
  name: 'StorageChecks',
  components: {
    Dialog,
    Error
  },
  data () {
    return {
      health: { disks: [], scrub: null, schedule: { scrub: false, smart: false }, busy: false },
      schedule: { scrub: false, smart: false },
      anyErrors: false,
      resetVisible: false,
      timer: undefined
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
    reload () {
      clearTimeout(this.timer)
      return Promise.all([axios.get('/rest/storage/health'), axios.get('/rest/storage/disks')])
        .then(([health, disks]) => {
          this.health = health.data.data
          this.schedule = { ...this.health.schedule }
          this.anyErrors = (disks.data.data || []).some(d => d.has_errors)
          if (this.health.scrub && this.health.scrub.state === 'running') {
            this.timer = setTimeout(this.reload, 10000)
          }
        })
        .catch(err => this.$refs.error.showAxios(err))
    },
    startScrub () {
      return axios.post('/rest/storage/check/scrub')
        .then(() => this.reload())
        .catch(err => this.$refs.error.showAxios(err))
    },
    resetErrors () {
      this.resetVisible = false
      return axios.post('/rest/storage/errors/reset')
        .then(() => this.reload())
        .catch(err => this.$refs.error.showAxios(err))
    },
    saveSchedule (check, value) {
      this.schedule = { ...this.schedule, [check]: value }
      return axios.post('/rest/storage/check/schedule', this.schedule)
        .then(resp => { this.schedule = { ...resp.data.data } })
        .catch(err => this.$refs.error.showAxios(err))
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
  font-weight: 600;
}
.check-text {
  margin: 8px 0 0;
  color: var(--sc-muted);
}
.disk-bad {
  color: var(--sc-danger);
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
