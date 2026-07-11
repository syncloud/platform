<template>
  <div class="login-root">
    <div class="login-wordmark">SYNCLOUD</div>

    <div class="login-card">
      <img class="login-logo" src="/logo.png" alt="Syncloud" />

      <!-- Credentials form -->
      <form v-if="step === 'credentials'" class="login-form" @submit.prevent="submitCredentials">
        <div class="field">
          <label for="username-textfield">Username</label>
          <input
            id="username-textfield"
            class="input"
            v-model="username"
            placeholder="Username"
            :disabled="loading"
            autocomplete="username">
        </div>
        <div class="field">
          <label for="password-textfield">Password</label>
          <div class="password-wrap">
            <input
              id="password-textfield"
              class="input"
              :type="showPassword ? 'text' : 'password'"
              v-model="password"
              placeholder="Password"
              :disabled="loading"
              autocomplete="current-password"
              @keyup.enter="submitCredentials">
            <button type="button" class="pw-toggle" @click="showPassword = !showPassword">{{ showPassword ? 'Hide' : 'Show' }}</button>
          </div>
        </div>
        <label class="remember">
          <input type="checkbox" v-model="keepMeLoggedIn"> Remember me
        </label>
        <button id="sign-in-button" type="submit" class="btn btn-primary" :disabled="loading">
          <span v-if="loading" class="spinner"></span><span v-else>Sign in</span>
        </button>
      </form>

      <!-- TOTP setup (first time) -->
      <div v-if="step === 'totp_setup'">
        <div class="alert alert-info">
          Two-factor authentication is required. Scan this QR code with your authenticator app.
        </div>
        <div class="totp-qr">
          <img id="totp_qr" :src="totpQr" alt="TOTP QR Code" />
        </div>
        <p class="muted center">Or enter this secret manually:</p>
        <div class="center secret"><code id="totp_secret">{{ totpSecret }}</code></div>
        <form class="login-form" @submit.prevent="submitTotp">
          <div class="field">
            <input
              id="otp-input"
              class="input"
              v-model="totpCode"
              placeholder="Enter code from authenticator"
              :disabled="loading"
              autocomplete="one-time-code"
              @keyup.enter="submitTotp">
          </div>
          <button id="verify-button" type="submit" class="btn btn-primary" :disabled="loading">
            <span v-if="loading" class="spinner"></span><span v-else>Verify</span>
          </button>
        </form>
      </div>

      <!-- TOTP verify (returning user) -->
      <div v-if="step === 'totp_verify'">
        <p class="muted center">Enter the code from your authenticator app</p>
        <form class="login-form" @submit.prevent="submitTotp">
          <div class="field">
            <input
              id="otp-input"
              class="input"
              v-model="totpCode"
              placeholder="Authentication code"
              :disabled="loading"
              autocomplete="one-time-code"
              @keyup.enter="submitTotp">
          </div>
          <button id="verify-button" type="submit" class="btn btn-primary" :disabled="loading">
            <span v-if="loading" class="spinner"></span><span v-else>Verify</span>
          </button>
        </form>
      </div>

      <!-- Error display -->
      <div v-if="errorMessage" class="alert alert-error notification" id="login_error">
        {{ errorMessage }}
        <button type="button" class="alert-close" @click="errorMessage = ''">✕</button>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
import QRCode from 'qrcode'

export default {
  name: 'LoginApp',
  data() {
    return {
      step: 'credentials',
      username: '',
      password: '',
      showPassword: false,
      keepMeLoggedIn: false,
      totpCode: '',
      totpQr: null,
      totpSecret: null,
      loading: false,
      errorMessage: '',
      targetURL: '',
      flowID: ''
    }
  },
  async mounted() {
    // Handle logout path - clear Authelia session first
    if (window.location.pathname === '/logout') {
      try {
        await axios.post('/api/logout')
      } catch {
        // Session may already be cleared
      }
      window.history.replaceState(null, '', '/')
      return
    }
    const params = new URLSearchParams(window.location.search)
    this.targetURL = params.get('rd') || ''
    this.flowID = params.get('flow_id') || ''
    this.checkState()
  },
  methods: {
    async checkState() {
      console.log('login: checkState, flowID=' + this.flowID + ', targetURL=' + this.targetURL)
      // For OIDC flows, always show login form (don't auto-complete)
      if (this.flowID) {
        console.log('login: OIDC flow detected, showing login form')
        return
      }
      try {
        const resp = await axios.get('/api/state')
        if (resp.data && resp.data.data) {
          const level = resp.data.data.authentication_level
          if (level >= 2) {
            await this.completeFlow()
            return
          }
          if (level === 1) {
            await this.handleSecondFactor()
            return
          }
        }
      } catch {
        // Not authenticated, show login form
      }
    },
    async submitCredentials() {
      if (!this.username || !this.password) {
        this.showError('Please enter username and password')
        return
      }
      this.loading = true
      this.errorMessage = ''
      try {
        await axios.post('/api/firstfactor', {
          username: this.username,
          password: this.password,
          keepMeLoggedIn: this.keepMeLoggedIn,
          targetURL: this.targetURL
        })
        await this.handleSecondFactor()
      } catch (err) {
        this.loading = false
        this.showError(this.extractError(err, 'Incorrect username or password'))
      }
    },
    async handleSecondFactor() {
      const stateResp = await axios.get('/api/state')
      const level = stateResp.data.data.authentication_level
      console.log('login: authentication_level=' + level)
      if (level >= 2) {
        // Fully authenticated (2FA complete or not required)
        await this.completeFlow()
      } else if (level === 1) {
        // First factor done - try OIDC consent first (works for one_factor policy)
        if (this.flowID) {
          try {
            const consentInfo = await axios.get('/api/oidc/consent?flow_id=' + this.flowID)
            const clientID = consentInfo.data.data.client_id
            const resp = await axios.post('/api/oidc/consent', {
              flow_id: this.flowID,
              client_id: clientID,
              consent: true
            })
            if (resp.data && resp.data.data && resp.data.data.redirect_uri) {
              window.location.href = resp.data.data.redirect_uri
              return
            }
          } catch {
            // Consent requires higher auth level, proceed to TOTP setup
            console.log('login: OIDC consent needs 2FA, proceeding to TOTP setup')
          }
        }
        // No OIDC flow or consent failed - need second factor
        if (!this.flowID) {
          // Direct login to auth domain without OIDC flow, redirect to target or home
          await this.completeFlow()
          return
        }
        // Check if user already has TOTP configured (returning user)
        try {
          const statusResp = await axios.post('/login/totp/status', { username: this.username, password: this.password })
          if (statusResp.data.data.configured) {
            this.step = 'totp_verify'
            this.loading = false
            return
          }
        } catch {
          console.log('login: TOTP status check failed, proceeding to setup')
        }
        await this.setupTotp()
      } else {
        this.loading = false
        this.showError('Authentication failed')
      }
    },
    async setupTotp() {
      this.loading = true
      try {
        const resp = await axios.post('/login/totp/setup', { username: this.username, password: this.password })
        const uri = resp.data.data.uri
        const params = new URL(uri).searchParams
        this.totpSecret = params.get('secret')
        this.totpQr = await QRCode.toDataURL(uri)
        this.step = 'totp_setup'
        this.loading = false
      } catch (err) {
        this.loading = false
        this.showError(this.extractError(err, 'Failed to set up two-factor authentication'))
      }
    },
    async submitTotp() {
      if (!this.totpCode) {
        this.showError('Please enter the authentication code')
        return
      }
      this.loading = true
      this.errorMessage = ''
      try {
        await axios.post('/api/secondfactor/totp', {
          token: this.totpCode,
          targetURL: this.targetURL
        })
        await this.completeFlow()
      } catch (err) {
        this.loading = false
        this.showError(this.extractError(err, 'Invalid authentication code'))
      }
    },
    async completeFlow() {
      if (this.flowID) {
        try {
          const consentInfo = await axios.get('/api/oidc/consent?flow_id=' + this.flowID)
          const clientID = consentInfo.data.data.client_id
          const resp = await axios.post('/api/oidc/consent', {
            flow_id: this.flowID,
            client_id: clientID,
            consent: true
          })
          if (resp.data && resp.data.data && resp.data.data.redirect_uri) {
            window.location.href = resp.data.data.redirect_uri
            return
          }
        } catch {
          // Consent failed, fall through to default redirect
        }
      }
      if (this.targetURL) {
        window.location.href = this.targetURL
      } else {
        // Redirect to main domain (strip 'auth.' prefix from auth.domain)
        const mainDomain = window.location.hostname.replace(/^auth\./, '')
        window.location.href = 'https://' + mainDomain
      }
    },
    showError(message) {
      this.errorMessage = message
      setTimeout(() => { this.errorMessage = '' }, 5000)
    },
    extractError(err, fallback) {
      if (err.response && err.response.data && err.response.data.data && err.response.data.data.message) {
        return err.response.data.data.message
      }
      if (err.response && err.response.data && err.response.data.message) {
        return err.response.data.message
      }
      return fallback
    }
  }
}
</script>

<style>
:root {
  --lg-primary: #2b7bd6;
  --lg-primary-dark: #1d6ec7;
  --lg-ink: #1a2a3a;
  --lg-ink-2: #3c5373;
  --lg-muted: #5a6b80;
  --lg-faint: #8796a8;
  --lg-field-bg: #f6f9fd;
  --lg-border: #d5dde8;
  --lg-danger: #d9363e;
}
* { box-sizing: border-box; }
body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
  color: var(--lg-ink);
  background: linear-gradient(140deg, #f7fafe 0%, #eaf2fb 55%, #dbe7f4 100%);
  min-height: 100vh;
  min-height: 100dvh;
}
.login-root {
  min-height: 100vh;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 32px 16px 160px;
}
.login-wordmark {
  font-size: 14px;
  font-weight: 700;
  letter-spacing: 6px;
  color: var(--lg-primary);
  margin-top: auto;
  margin-bottom: 24px;
}
.login-card {
  width: 100%;
  max-width: 420px;
  margin-bottom: auto;
  background: #fff;
  border-radius: 22px;
  box-shadow: 0 24px 60px -20px rgba(22, 50, 92, 0.18), 0 2px 6px rgba(22, 50, 92, 0.04);
  padding: 34px 32px 36px;
}
.login-logo {
  display: block;
  width: 56px;
  height: 56px;
  margin: 0 auto 22px;
}
.login-form { display: block; }
.field { margin-bottom: 16px; }
.field label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: var(--lg-ink-2);
  margin-bottom: 6px;
}
.input {
  width: 100%;
  height: 46px;
  padding: 0 16px;
  border-radius: 12px;
  border: 1px solid var(--lg-border);
  background: var(--lg-field-bg);
  font-size: 16px;
  color: var(--lg-ink);
  transition: border-color 0.2s ease, background 0.2s ease, box-shadow 0.2s ease;
}
.input:focus {
  outline: none;
  background: #fff;
  border-color: var(--lg-primary);
  box-shadow: 0 0 0 4px rgba(43, 123, 214, 0.12);
}
.password-wrap { position: relative; }
.pw-toggle {
  position: absolute;
  right: 8px;
  top: 50%;
  transform: translateY(-50%);
  background: transparent;
  border: none;
  color: var(--lg-faint);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  padding: 6px 8px;
}
.pw-toggle:hover { color: var(--lg-primary); }
.remember {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: var(--lg-muted);
  margin-bottom: 18px;
  cursor: pointer;
}
.remember input { width: 16px; height: 16px; }
.btn {
  width: 100%;
  height: 48px;
  border: none;
  border-radius: 12px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.1s ease, box-shadow 0.2s ease, filter 0.15s ease;
}
.btn:active { transform: translateY(1px); }
.btn:disabled { opacity: 0.6; cursor: not-allowed; }
.btn-primary {
  background: linear-gradient(135deg, var(--lg-primary) 0%, var(--lg-primary-dark) 100%);
  color: #fff;
  box-shadow: 0 6px 16px -6px rgba(43, 123, 214, 0.5);
}
.btn-primary:hover:not(:disabled) { filter: brightness(1.05); }
.spinner {
  display: inline-block;
  width: 18px;
  height: 18px;
  border: 2px solid rgba(255, 255, 255, 0.4);
  border-top-color: #fff;
  border-radius: 50%;
  animation: lg-spin 0.8s linear infinite;
}
@keyframes lg-spin { to { transform: rotate(360deg); } }
.muted { color: var(--lg-muted); font-size: 14px; }
.center { text-align: center; }
.secret { margin-bottom: 16px; }
.secret code { word-break: break-all; font-size: 14px; color: var(--lg-ink-2); }
.totp-qr { text-align: center; margin: 16px 0; }
.totp-qr img { max-width: 200px; border-radius: 12px; }
.alert {
  border-radius: 12px;
  padding: 12px 14px;
  font-size: 14px;
  line-height: 1.45;
}
.alert-info {
  background: #e7f1fc;
  color: var(--lg-primary-dark);
  margin-bottom: 16px;
}
.alert-error {
  position: relative;
  background: #fdeced;
  color: var(--lg-danger);
  margin-top: 14px;
  padding-right: 34px;
}
.alert-close {
  position: absolute;
  right: 10px;
  top: 50%;
  transform: translateY(-50%);
  background: transparent;
  border: none;
  color: var(--lg-danger);
  cursor: pointer;
  font-size: 13px;
}
</style>
