<template>
  <div class="act-root">
    <div class="act-wordmark">SYNCLOUD</div>

    <div class="act-card">
      <div class="act-progress" :aria-label="'step ' + (step+1) + ' of 4'">
        <span v-for="i in 4" :key="i" class="act-dot" :class="{ active: i-1 === step, done: i-1 < step }"></span>
      </div>

      <transition name="act-step" mode="out-in">
        <section v-if="step === 0" key="welcome" class="act-step act-welcome">
          <div class="act-clock" id="activate_tz_preview">
            <div class="act-clock-time">{{ clockTime }}</div>
            <div class="act-clock-date">{{ clockDate }}</div>
            <div class="act-clock-tz">{{ timezone || browserTz }}</div>
          </div>

          <h1 class="act-hello">{{ $t('activate.welcome') }}</h1>
          <p class="act-lead">{{ $t('activate.welcomeSub') }}</p>

          <div class="act-field">
            <label for="activate_language">{{ $t('language.select') }}</label>
            <div class="act-select">
              <select id="activate_language" v-model="locale">
                <option v-for="l in locales" :key="l.code" :value="l.code" :id="'lang_' + l.code">{{ l.name }}</option>
              </select>
              <span class="act-caret">▾</span>
            </div>
          </div>

          <div class="act-field">
            <label for="activate_timezone">{{ $t('locale.timezone') }}</label>
            <div class="act-combo" :class="{ open: tzOpen }" @keydown.escape="tzOpen = false">
              <input
                id="activate_timezone"
                class="act-input act-combo-input"
                type="text"
                autocomplete="off"
                spellcheck="false"
                v-model="tzQuery"
                :placeholder="browserTz || 'UTC'"
                @focus="openTz"
                @click="openTz"
                @input="tzOpen = true"
                @blur="closeTzSoon"
                @keydown.enter.prevent="selectFirstTz"
                @keydown.down.prevent="tzOpen = true"
              />
              <span class="act-caret">▾</span>
              <div v-if="tzOpen" class="act-combo-list">
                <button
                  v-for="tz in tzFiltered"
                  :key="tz"
                  type="button"
                  class="act-combo-item"
                  :class="{ selected: tz === timezone }"
                  :id="'activate_tz_' + tz.replace(/[^a-zA-Z0-9]/g, '_')"
                  @mousedown.prevent="selectTz(tz)"
                >{{ tz }}</button>
                <div v-if="tzFiltered.length === 0" class="act-combo-empty">—</div>
              </div>
            </div>
          </div>

          <div class="act-actions act-actions-single">
            <button id="btn_welcome_next" class="act-btn act-btn-primary" @click="step = 1">{{ $t('activate.continue') }}</button>
          </div>
        </section>

        <section v-else-if="step === 1" key="type" class="act-step">
          <h1 class="act-title">{{ $t('activate.stepType') }}</h1>
          <p class="act-lead">{{ $t('activate.chooseDomainSub') }}</p>

          <div class="act-cards">
            <button id="btn_free_domain" class="act-choice" @click="selectFreeDomain">
              <div class="act-choice-icon">🏠</div>
              <div class="act-choice-title">{{ $t('activate.freeButton') }}</div>
              <div class="act-choice-body">{{ $t('activate.freeDescription', { domain: redirect_domain }) }}</div>
            </button>
            <button id="btn_premium_domain" class="act-choice" @click="selectPremiumDomain">
              <div class="act-choice-icon">🌐</div>
              <div class="act-choice-title">{{ $t('activate.premiumButton') }}</div>
              <div class="act-choice-body">{{ $t('activate.premiumDescription') }}</div>
            </button>
          </div>

          <div class="act-actions">
            <button class="act-btn act-btn-ghost" @click="step--">{{ $t('common.previous') }}</button>
          </div>
        </section>

        <section v-else-if="step === 2" key="details" class="act-step">
          <h1 class="act-title">{{ $t('activate.syncloudAccount') }}</h1>
          <p class="act-lead">{{ $t('activate.accountSub') }}</p>

          <div class="act-field">
            <label for="email">{{ $t('activate.emailPlaceholder', { domain: redirect_domain }) }}</label>
            <input id="email" class="act-input" type="text" v-model="redirectEmail" autocomplete="email">
            <div v-show="redirectEmailAlertVisible" id="email_alert" class="act-error">{{ redirectEmailAlert }}</div>
          </div>

          <div class="act-field">
            <label for="redirect_password">{{ $t('activate.passwordPlaceholder', { domain: redirect_domain }) }}</label>
            <Password id="redirect_password" v-model="redirectPassword"
                      :placeholder="$t('activate.passwordPlaceholder', { domain: redirect_domain })"
                      :show-error="redirectPasswordAlertVisible" :error="redirectPasswordAlert"/>
          </div>

          <div v-if="domainType === 'free'" class="act-hint">
            {{ $t('activate.noAccount') }}
            <a :href="'https://' + redirect_domain" target="_blank" class="act-link">{{ $t('activate.register') }}</a>
          </div>

          <div class="act-field">
            <label :for="domainType === 'free' ? 'domain_input' : 'domain_premium'">{{ $t('activate.deviceName') }}</label>
            <div v-if="domainType === 'free'" class="act-domain">
              <input id="domain_input" class="act-input act-input-flex" type="text" v-model="domain" :placeholder="$t('activate.domainNamePlaceholder')">
              <span class="act-domain-suffix">.{{ redirect_domain }}</span>
            </div>
            <input v-else id="domain_premium" class="act-input" type="text" v-model="domain" :placeholder="$t('activate.premiumDomainPlaceholder')">
            <div v-show="domainAlertVisible" id="domain_alert" class="act-error">{{ domainAlert }}</div>
          </div>

          <div class="act-actions">
            <button class="act-btn act-btn-ghost" @click="step = 1">{{ $t('common.previous') }}</button>
            <button id="btn_next" class="act-btn act-btn-primary" @click="selectDeviceName">{{ $t('common.next') }}</button>
          </div>
        </section>

        <section v-else-if="step === 3" key="creds" class="act-step">
          <h1 class="act-title">{{ $t('activate.deviceCredentials') }}</h1>
          <p class="act-lead">{{ $t('activate.credentialsSub') }}</p>

          <div class="act-field">
            <label for="device_username">{{ $t('activate.loginPlaceholder') }}</label>
            <input id="device_username" class="act-input" type="text" v-model="deviceUsername" @keyup.enter="activate">
            <div v-show="deviceUsernameAlertVisible" class="act-error">{{ deviceUsernameAlert }}</div>
          </div>

          <div class="act-field">
            <label for="device_password">{{ $t('activate.passwordInputPlaceholder') }}</label>
            <Password id="device_password" v-model="devicePassword"
                      :placeholder="$t('activate.passwordInputPlaceholder')"
                      :show-error="devicePasswordAlertVisible" :error="devicePasswordAlert"
                      @trigger="activate"/>
          </div>

          <div class="act-field">
            <label for="device_password_confirm">{{ $t('activate.confirmPasswordPlaceholder') }}</label>
            <Password id="device_password_confirm" v-model="devicePasswordConfirm"
                      :placeholder="$t('activate.confirmPasswordPlaceholder')"
                      :show-error="devicePassword !== devicePasswordConfirm"
                      :error="$t('activate.passwordsMismatch')"
                      @trigger="activate"/>
          </div>

          <div class="act-actions">
            <button class="act-btn act-btn-ghost" @click="step = 2">{{ $t('common.previous') }}</button>
            <button id="btn_activate" class="act-btn act-btn-primary" @click="activate" :disabled="!validDeviceCredentials()">{{ $t('activate.finish') }}</button>
          </div>
        </section>
      </transition>
    </div>

    <div v-if="loading" class="act-loading">
      <div class="act-spinner"></div>
      <div class="act-loading-text">{{ $t('common.loading') }}</div>
    </div>

    <transition name="act-banner">
      <div v-if="bannerVisible" class="act-banner" id="txt_error" @click="bannerVisible = false">
        <span class="act-banner-icon">⚠</span>
        <span class="act-banner-text">{{ bannerMessage }}</span>
        <span class="act-banner-close">✕</span>
      </div>
    </transition>
  </div>
</template>

<script>
import axios from 'axios'
import Password from '../components/Password.vue'
import { SUPPORTED_LOCALES, setLocale } from '../i18n'

export default {
  name: 'Activate',
  components: {
    Password
  },
  data () {
    return {
      domainType: 'free',
      loading: false,
      redirectEmail: '',
      redirectPassword: '',
      domain: '',
      redirect_domain: 'syncloud.it',
      deviceUsername: '',
      devicePassword: '',
      devicePasswordConfirm: '',
      deviceUsernameAlertVisible: false,
      deviceUsernameAlert: '',
      devicePasswordAlertVisible: false,
      devicePasswordAlert: '',
      redirectPasswordAlertVisible: false,
      redirectPasswordAlert: '',
      domainAlertVisible: false,
      domainAlert: '',
      redirectEmailAlertVisible: false,
      redirectEmailAlert: '',
      step: 0,
      locales: SUPPORTED_LOCALES,
      timezone: '',
      timezones: this.listTimezones(),
      bannerVisible: false,
      bannerMessage: '',
      nowMs: Date.now(),
      tickTimer: null,
      tzOpen: false,
      tzQuery: '',
      tzCloseTimer: null
    }
  },
  computed: {
    locale: {
      get () { return this.$i18n.locale },
      set (code) { setLocale(code) }
    },
    browserTz () {
      try {
        return Intl.DateTimeFormat().resolvedOptions().timeZone || ''
      } catch {
        return ''
      }
    },
    clockTime () {
      try {
        const opts = { hour: '2-digit', minute: '2-digit', hour12: false }
        if (this.timezone) opts.timeZone = this.timezone
        return new Intl.DateTimeFormat(this.$i18n.locale, opts).format(new Date(this.nowMs))
      } catch {
        return ''
      }
    },
    clockDate () {
      try {
        const opts = { weekday: 'long', month: 'long', day: 'numeric' }
        if (this.timezone) opts.timeZone = this.timezone
        return new Intl.DateTimeFormat(this.$i18n.locale, opts).format(new Date(this.nowMs))
      } catch {
        return ''
      }
    },
    tzFiltered () {
      const q = this.tzQuery.trim().toLowerCase()
      if (!q) return this.timezones.slice(0, 60)
      return this.timezones.filter(t => t.toLowerCase().includes(q)).slice(0, 60)
    }
  },
  mounted () {
    this.tickTimer = setInterval(() => { this.nowMs = Date.now() }, 1000)
    this.loading = true
    axios
      .get('/rest/redirect_info')
      .then(response => {
        this.redirect_domain = response.data.data.domain
        this.loading = false
      })
      .catch(err => {
        this.loading = false
        if (err.response && err.response.status !== 502) {
          this.showBanner(err)
        }
      })
  },
  unmounted () {
    if (this.tickTimer) clearInterval(this.tickTimer)
    if (this.tzCloseTimer) clearTimeout(this.tzCloseTimer)
  },
  methods: {
    openTz () {
      if (this.tzCloseTimer) {
        clearTimeout(this.tzCloseTimer)
        this.tzCloseTimer = null
      }
      this.tzOpen = true
    },
    closeTzSoon () {
      this.tzCloseTimer = setTimeout(() => {
        this.tzOpen = false
        if (this.timezones.indexOf(this.tzQuery) === -1) {
          this.tzQuery = this.timezone
        }
      }, 120)
    },
    selectTz (tz) {
      this.timezone = tz
      this.tzQuery = tz
      this.tzOpen = false
    },
    selectFirstTz () {
      if (this.tzFiltered.length > 0) {
        this.selectTz(this.tzFiltered[0])
      }
    },
    listTimezones () {
      if (typeof Intl !== 'undefined' && typeof Intl.supportedValuesOf === 'function') {
        try { return Intl.supportedValuesOf('timeZone') } catch { /* fall through */ }
      }
      return ['UTC']
    },
    activateRequestBody (domain) {
      const body = {
        redirect_email: this.redirectEmail,
        redirect_password: this.redirectPassword,
        domain,
        device_username: this.deviceUsername,
        device_password: this.devicePassword
      }
      if (this.timezone) {
        body.timezone = this.timezone
      }
      return body
    },
    validDeviceCredentials () {
      if (this.deviceUsername === '') return false
      if (this.devicePassword === '') return false
      if (this.devicePassword !== this.devicePasswordConfirm) return false
      return true
    },
    activate (event) {
      if (!this.validDeviceCredentials()) return
      if (event && event.preventDefault) event.preventDefault()
      this.loading = true
      this.hideAlerts()
      if (this.domainType === 'premium') {
        this.activatePremiumDomain()
      } else {
        this.activateFreeDomain()
      }
    },
    forceCertificateRecheck () {
      window.location = '/?t=' + (new Date()).getTime()
    },
    hideAlerts () {
      this.deviceUsernameAlertVisible = false
      this.devicePasswordAlertVisible = false
      this.redirectEmailAlertVisible = false
      this.redirectPasswordAlertVisible = false
      this.domainAlertVisible = false
      this.bannerVisible = false
    },
    showBanner (err) {
      if (err && err.response && err.response.status === 401) {
        this.$router.push('/login')
        return
      }
      let message = this.$t('common.serverError')
      if (err && err.response && err.response.data && err.response.data.message) {
        message = err.response.data.message
      }
      this.bannerMessage = message
      this.bannerVisible = true
    },
    showRedirectAlert (err) {
      if (err.response && err.response.data) {
        const data = err.response.data
        if (data.parameters_messages) {
          for (const pm of data.parameters_messages) {
            const message = pm.messages.join('\n')
            if (pm.parameter === 'redirect_password') {
              this.redirectPasswordAlertVisible = true
              this.redirectPasswordAlert = message
            }
            if (pm.parameter === 'email') {
              this.redirectEmailAlertVisible = true
              this.redirectEmailAlert = message
            }
            if (pm.parameter === 'domain') {
              this.domainAlertVisible = true
              this.domainAlert = message
            }
          }
        } else {
          this.showBanner(err)
        }
      }
    },
    showActivateAlert (err) {
      if (err.response && err.response.data) {
        const data = err.response.data
        if (data.parameters_messages) {
          for (const pm of data.parameters_messages) {
            const message = pm.messages.join(', ')
            if (pm.parameter === 'device_username') {
              this.deviceUsernameAlertVisible = true
              this.deviceUsernameAlert = message
            }
            if (pm.parameter === 'device_password') {
              this.devicePasswordAlertVisible = true
              this.devicePasswordAlert = message
            }
          }
        } else {
          this.showBanner(err)
        }
      }
    },
    activateFreeDomain () {
      axios
        .post('/rest/activate/managed', this.activateRequestBody(this.fullDomain()))
        .then(this.forceCertificateRecheck)
        .catch(err => {
          this.loading = false
          this.showActivateAlert(err)
        })
    },
    activatePremiumDomain () {
      axios
        .post('/rest/activate/managed', this.activateRequestBody(this.domain))
        .then(this.forceCertificateRecheck)
        .catch(err => {
          this.loading = false
          this.showActivateAlert(err)
        })
    },
    selectPremiumDomain () {
      this.hideAlerts()
      this.domainType = 'premium'
      this.step = 2
    },
    selectFreeDomain () {
      this.hideAlerts()
      this.domainType = 'free'
      this.step = 2
    },
    selectDeviceName () {
      this.hideAlerts()
      this.domainAvailability()
    },
    fullDomain () {
      if (this.domainType === 'free') {
        return this.domain + '.' + this.redirect_domain
      }
      return this.domain
    },
    domainAvailability () {
      this.loading = true
      axios
        .post('/rest/redirect/domain/availability',
          {
            email: this.redirectEmail,
            password: this.redirectPassword,
            domain: this.fullDomain()
          })
        .then(() => {
          this.step = 3
          this.loading = false
        })
        .catch(err => {
          this.loading = false
          this.showRedirectAlert(err)
        })
    }
  }
}
</script>

<style scoped>
.act-root {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-start;
  padding: 56px 16px 40px;
  min-height: 100vh;
  background: linear-gradient(140deg, #f7fafe 0%, #eaf2fb 55%, #dbe7f4 100%);
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
  color: #1a2a3a;
  box-sizing: border-box;
}
.act-wordmark {
  font-size: 13px;
  font-weight: 700;
  letter-spacing: 6px;
  color: #2b7bd6;
  margin: 8px 0 24px;
}
.act-clock {
  text-align: center;
  margin: 0 0 20px;
}
.act-clock-time {
  font-size: 64px;
  font-weight: 200;
  letter-spacing: -2px;
  color: #1a2a3a;
  line-height: 1;
  font-variant-numeric: tabular-nums;
}
.act-clock-date {
  font-size: 14px;
  color: #5a6b80;
  margin-top: 6px;
  text-transform: capitalize;
}
.act-clock-tz {
  font-size: 11px;
  color: #8796a8;
  letter-spacing: 1.5px;
  text-transform: uppercase;
  margin-top: 4px;
}
.act-card {
  width: 100%;
  max-width: 460px;
  background: #fff;
  border-radius: 22px;
  box-shadow: 0 24px 60px -20px rgba(22, 50, 92, 0.18), 0 2px 6px rgba(22, 50, 92, 0.04);
  padding: 24px 28px 26px;
  box-sizing: border-box;
}
.act-progress {
  display: flex;
  justify-content: center;
  gap: 8px;
  margin-bottom: 20px;
}
.act-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #d5dde8;
  transition: background 0.25s ease, width 0.25s ease;
}
.act-dot.active {
  background: #2b7bd6;
  width: 24px;
  border-radius: 6px;
}
.act-dot.done {
  background: #8fb7e6;
}
.act-step {
  display: flex;
  flex-direction: column;
}
.act-hello {
  font-size: 30px;
  font-weight: 700;
  margin: 0 0 4px;
  text-align: center;
  letter-spacing: -0.5px;
}
.act-title {
  font-size: 24px;
  font-weight: 700;
  margin: 0 0 4px;
  text-align: center;
  letter-spacing: -0.3px;
}
.act-lead {
  font-size: 14px;
  color: #5a6b80;
  text-align: center;
  margin: 0 0 18px;
  line-height: 1.45;
}
.act-field {
  margin-bottom: 16px;
}
.act-field label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: #3c5373;
  margin-bottom: 6px;
  padding-left: 2px;
  letter-spacing: 0.2px;
}
.act-input {
  width: 100%;
  height: 46px;
  padding: 0 16px;
  border-radius: 12px;
  border: 1px solid #d5dde8;
  background: #f6f9fd;
  font-size: 16px;
  color: #1a2a3a;
  transition: border-color 0.2s ease, background 0.2s ease, box-shadow 0.2s ease;
  box-sizing: border-box;
}
.act-input:focus {
  outline: none;
  background: #fff;
  border-color: #2b7bd6;
  box-shadow: 0 0 0 4px rgba(43, 123, 214, 0.12);
}
.act-input-flex {
  flex: 1 1 auto;
  min-width: 0;
}
.act-domain {
  display: flex;
  align-items: center;
  gap: 8px;
}
.act-domain-suffix {
  color: #5a6b80;
  font-size: 15px;
  white-space: nowrap;
}
.act-select {
  position: relative;
}
.act-select select {
  width: 100%;
  height: 46px;
  padding: 0 40px 0 16px;
  border-radius: 12px;
  border: 1px solid #d5dde8;
  background: #f6f9fd;
  font-size: 16px;
  color: #1a2a3a;
  appearance: none;
  -webkit-appearance: none;
  cursor: pointer;
  transition: border-color 0.2s ease, background 0.2s ease, box-shadow 0.2s ease;
}
.act-select select:focus {
  outline: none;
  background: #fff;
  border-color: #2b7bd6;
  box-shadow: 0 0 0 4px rgba(43, 123, 214, 0.12);
}
.act-caret {
  position: absolute;
  right: 16px;
  top: 23px;
  transform: translateY(-50%);
  pointer-events: none;
  color: #5a6b80;
  font-size: 12px;
}
.act-combo {
  position: relative;
}
.act-combo-input {
  padding-right: 40px;
}
.act-combo-list {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  right: 0;
  max-height: 220px;
  overflow-y: auto;
  background: #fff;
  border: 1px solid #d5dde8;
  border-radius: 12px;
  box-shadow: 0 12px 24px -8px rgba(22, 50, 92, 0.18);
  z-index: 20;
  padding: 4px;
}
.act-combo-item {
  display: block;
  width: 100%;
  text-align: left;
  padding: 9px 12px;
  border: none;
  background: transparent;
  border-radius: 8px;
  font: inherit;
  font-size: 14px;
  color: #1a2a3a;
  cursor: pointer;
}
.act-combo-item:hover, .act-combo-item:focus {
  background: #f0f6fe;
  outline: none;
}
.act-combo-item.selected {
  background: #e7f1fc;
  color: #1d6ec7;
  font-weight: 600;
}
.act-combo-empty {
  padding: 10px 12px;
  color: #8796a8;
  font-size: 13px;
  text-align: center;
}
.act-hint {
  font-size: 13px;
  color: #5a6b80;
  text-align: right;
  margin: -10px 2px 16px 0;
}
.act-link {
  color: #2b7bd6;
  text-decoration: none;
  font-weight: 600;
}
.act-link:hover {
  text-decoration: underline;
}
.act-error {
  font-size: 13px;
  color: #d9363e;
  margin-top: 6px;
  padding-left: 2px;
  white-space: pre-line;
}
.act-cards {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin-bottom: 20px;
}
.act-choice {
  background: #f6f9fd;
  border: 1px solid #d5dde8;
  border-radius: 16px;
  padding: 20px 14px;
  text-align: center;
  cursor: pointer;
  transition: transform 0.15s ease, border-color 0.2s ease, box-shadow 0.2s ease, background 0.2s ease;
  font-family: inherit;
  color: inherit;
}
.act-choice:hover {
  border-color: #2b7bd6;
  background: #fff;
  box-shadow: 0 10px 24px -12px rgba(43, 123, 214, 0.35);
  transform: translateY(-2px);
}
.act-choice:focus {
  outline: none;
  border-color: #2b7bd6;
  box-shadow: 0 0 0 4px rgba(43, 123, 214, 0.12);
}
.act-choice-icon {
  font-size: 32px;
  margin-bottom: 8px;
  line-height: 1;
}
.act-choice-title {
  font-size: 16px;
  font-weight: 700;
  margin-bottom: 6px;
  color: #1a2a3a;
}
.act-choice-body {
  font-size: 13px;
  color: #5a6b80;
  line-height: 1.4;
}
.act-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-top: auto;
  padding-top: 16px;
}
.act-actions-single {
  justify-content: center;
}
.act-btn {
  font: inherit;
  font-size: 15px;
  font-weight: 600;
  padding: 12px 22px;
  border-radius: 12px;
  border: none;
  cursor: pointer;
  transition: transform 0.1s ease, box-shadow 0.2s ease, background 0.2s ease;
  min-width: 110px;
}
.act-btn:active {
  transform: translateY(1px);
}
.act-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.act-btn-primary {
  background: linear-gradient(135deg, #2b7bd6 0%, #1d6ec7 100%);
  color: #fff;
  box-shadow: 0 6px 16px -6px rgba(43, 123, 214, 0.5);
}
.act-btn-primary:hover:not(:disabled) {
  box-shadow: 0 10px 24px -8px rgba(43, 123, 214, 0.55);
}
.act-btn-ghost {
  background: transparent;
  color: #5a6b80;
  padding: 12px 18px;
}
.act-btn-ghost:hover:not(:disabled) {
  color: #2b7bd6;
}
.act-loading {
  position: fixed;
  inset: 0;
  background: rgba(10, 20, 40, 0.45);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  gap: 14px;
  backdrop-filter: blur(4px);
}
.act-spinner {
  width: 42px;
  height: 42px;
  border: 3px solid rgba(255, 255, 255, 0.25);
  border-top-color: #fff;
  border-radius: 50%;
  animation: act-spin 0.9s linear infinite;
}
.act-loading-text {
  color: #fff;
  font-size: 14px;
  letter-spacing: 0.3px;
}
@keyframes act-spin {
  to { transform: rotate(360deg); }
}
.act-banner {
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  background: #fff;
  border-radius: 12px;
  border-left: 4px solid #d9363e;
  padding: 12px 40px 12px 16px;
  box-shadow: 0 12px 32px -10px rgba(10, 20, 40, 0.25);
  display: flex;
  align-items: center;
  gap: 10px;
  max-width: 90vw;
  cursor: pointer;
  z-index: 2100;
}
.act-banner-icon { color: #d9363e; font-size: 16px; }
.act-banner-text { font-size: 14px; color: #1a2a3a; }
.act-banner-close {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: #8796a8;
  font-size: 13px;
}
.act-banner-enter-active, .act-banner-leave-active {
  transition: opacity 0.25s ease, transform 0.25s ease;
}
.act-banner-enter-from { opacity: 0; transform: translate(-50%, -20px); }
.act-banner-leave-to { opacity: 0; transform: translate(-50%, -20px); }

.act-step-enter-active, .act-step-leave-active {
  transition: opacity 0.25s ease, transform 0.25s ease;
}
.act-step-enter-from {
  opacity: 0;
  transform: translateY(8px);
}
.act-step-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
@media (max-width: 520px) {
  .act-root { padding: 36px 10px 24px; }
  .act-wordmark { margin: 4px 0 16px; letter-spacing: 4px; }
  .act-clock-time { font-size: 52px; }
  .act-clock { margin-bottom: 16px; }
  .act-card { padding: 18px 18px 20px; border-radius: 16px; }
  .act-progress { margin-bottom: 14px; }
  .act-hello { font-size: 24px; }
  .act-title { font-size: 20px; }
  .act-lead { margin-bottom: 12px; font-size: 14px; }
  .act-field { margin-bottom: 12px; }
  .act-input, .act-select select { height: 44px; font-size: 15px; }
  .act-cards { grid-template-columns: 1fr; gap: 10px; }
  .act-choice { padding: 16px 14px; }
  .act-choice-icon { font-size: 28px; }
}
</style>
