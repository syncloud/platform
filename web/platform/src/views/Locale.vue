<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12" id="block1">
        <h1>{{ $t('locale.title') }}</h1>
        <div class="locale-form">
          <div class="locale-row">
            <span class="locale-label">{{ $t('locale.language') }}</span>
            <el-select id="settings_language" v-model="language" size="default" class="locale-field">
              <el-option v-for="l in locales" :key="l.code" :label="l.name" :value="l.code" :id="'settings_lang_' + l.code"/>
            </el-select>
          </div>

          <div class="locale-row">
            <span class="locale-label">{{ $t('locale.timezone') }}</span>
            <el-select id="settings_timezone" v-model="timezone" filterable size="default" class="locale-field">
              <el-option v-for="tz in timezones" :key="tz" :label="tz" :value="tz" :id="'settings_tz_' + tz.replace(/[^a-zA-Z0-9]/g, '_')"/>
            </el-select>
          </div>

          <div class="locale-row">
            <span class="locale-label">{{ $t('locale.currentTime') }}</span>
            <span class="locale-value" id="current_time">{{ displayTime }}</span>
          </div>

          <div class="locale-actions">
            <el-button id="btn_save_timezone" type="primary" :loading="saving" @click="saveTimezone">{{ $t('locale.save') }}</el-button>
          </div>
        </div>
      </div>
    </div>
  </div>

  <Error ref="error"/>
</template>

<script>
import Error from '../components/Error.vue'
import axios from 'axios'
import { SUPPORTED_LOCALES, setLocale } from '../i18n'

const FALLBACK_TIMEZONES = [
  'UTC', 'Europe/London', 'Europe/Berlin', 'Europe/Paris', 'Europe/Moscow',
  'America/New_York', 'America/Chicago', 'America/Denver', 'America/Los_Angeles',
  'America/Sao_Paulo', 'Asia/Tokyo', 'Asia/Shanghai', 'Asia/Kolkata',
  'Asia/Dubai', 'Australia/Sydney', 'Pacific/Auckland'
]

export default {
  name: 'Locale',
  components: { Error },
  data () {
    return {
      locales: SUPPORTED_LOCALES,
      timezone: 'UTC',
      timezones: this.listTimezones(),
      serverTime: null,
      serverTimezone: 'UTC',
      tickTimer: null,
      pollTimer: null,
      saving: false,
      nowMs: Date.now()
    }
  },
  computed: {
    language: {
      get () { return this.$i18n.locale },
      set (code) { setLocale(code) }
    },
    displayTime () {
      if (!this.serverTime) return ''
      try {
        const d = new Date(this.serverTime.baseMs + (this.nowMs - this.serverTime.syncedAt))
        return new Intl.DateTimeFormat(this.$i18n.locale, {
          timeZone: this.serverTimezone,
          dateStyle: 'medium',
          timeStyle: 'medium'
        }).format(d)
      } catch {
        return ''
      }
    }
  },
  mounted () {
    this.loadTimezone()
    this.loadTime()
    this.tickTimer = setInterval(() => { this.nowMs = Date.now() }, 1000)
    this.pollTimer = setInterval(() => this.loadTime(), 30000)
  },
  unmounted () {
    if (this.tickTimer) clearInterval(this.tickTimer)
    if (this.pollTimer) clearInterval(this.pollTimer)
  },
  methods: {
    listTimezones () {
      if (typeof Intl !== 'undefined' && typeof Intl.supportedValuesOf === 'function') {
        try { return Intl.supportedValuesOf('timeZone') } catch { /* fall through */ }
      }
      return FALLBACK_TIMEZONES
    },
    loadTimezone () {
      axios.get('/rest/settings/timezone')
        .then(resp => {
          this.timezone = resp.data.data.timezone
        })
        .catch(err => {
          this.$refs.error.showAxios(err)
        })
    },
    loadTime () {
      axios.get('/rest/settings/time')
        .then(resp => {
          const d = new Date(resp.data.data.time)
          this.serverTime = { baseMs: d.getTime(), syncedAt: Date.now() }
          this.serverTimezone = resp.data.data.timezone
          this.nowMs = Date.now()
        })
        .catch(err => {
          this.$refs.error.showAxios(err)
        })
    },
    saveTimezone () {
      this.saving = true
      axios.post('/rest/settings/timezone', { timezone: this.timezone })
        .then(() => {
          this.saving = false
          this.loadTime()
        })
        .catch(err => {
          this.saving = false
          this.$refs.error.showAxios(err)
        })
    }
  }
}
</script>
<style scoped>
.locale-form {
  max-width: 520px;
  margin: 0 auto;
  padding: 24px 16px 0;
  box-sizing: border-box;
}
.locale-row {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
}
.locale-label {
  flex: 0 0 160px;
  text-align: right;
  padding-right: 14px;
  font-size: 16px;
}
.locale-field {
  width: 320px;
  max-width: 100%;
}
.locale-value {
  font-size: 16px;
}
.locale-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 24px;
}
@media (max-width: 600px) {
  .locale-label {
    flex: 0 0 110px;
    padding-right: 8px;
  }
  .locale-field {
    width: 180px;
  }
}
</style>
<style>
@import '../style/site.css';
</style>
