<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <img class="login-logo" src="/logo.png" alt="Syncloud" />
        <h2>Syncloud</h2>
      </div>

      <!-- Credentials form -->
      <div v-if="step === 'credentials'">
        <el-form @submit.prevent="submitCredentials">
          <el-form-item>
            <el-input
              id="username-textfield"
              v-model="username"
              placeholder="Username"
              :disabled="loading"
              autocomplete="username"
            />
          </el-form-item>
          <el-form-item>
            <el-input
              id="password-textfield"
              v-model="password"
              type="password"
              placeholder="Password"
              :disabled="loading"
              show-password
              autocomplete="current-password"
            />
          </el-form-item>
          <el-form-item>
            <el-checkbox v-model="keepMeLoggedIn">Remember me</el-checkbox>
          </el-form-item>
          <el-form-item>
            <el-button
              id="sign-in-button"
              type="primary"
              :loading="loading"
              @click="submitCredentials"
              style="width: 100%"
            >SIGN IN</el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- TOTP setup (first time) -->
      <div v-if="step === 'totp_setup'">
        <el-alert type="info" :closable="false" show-icon style="margin-bottom: 16px">
          Two-factor authentication is required. Scan this QR code with your authenticator app.
        </el-alert>
        <div class="totp-qr">
          <img id="totp_qr" :src="totpQr" alt="TOTP QR Code" />
        </div>
        <p style="text-align: center; margin: 8px 0">Or enter this secret manually:</p>
        <div style="text-align: center; margin-bottom: 16px">
          <code id="totp_secret" style="word-break: break-all">{{ totpSecret }}</code>
        </div>
        <el-form @submit.prevent="submitTotp">
          <el-form-item>
            <el-input
              id="otp-input"
              v-model="totpCode"
              placeholder="Enter code from authenticator"
              :disabled="loading"
              autocomplete="one-time-code"
              @keyup.enter="submitTotp"
            />
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              :loading="loading"
              @click="submitTotp"
              style="width: 100%"
            >VERIFY</el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- TOTP verify (returning user) -->
      <div v-if="step === 'totp_verify'">
        <p style="text-align: center; margin-bottom: 16px">Enter the code from your authenticator app</p>
        <el-form @submit.prevent="submitTotp">
          <el-form-item>
            <el-input
              id="otp-input"
              v-model="totpCode"
              placeholder="Authentication code"
              :disabled="loading"
              autocomplete="one-time-code"
              @keyup.enter="submitTotp"
            />
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              :loading="loading"
              @click="submitTotp"
              style="width: 100%"
            >VERIFY</el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- Error display -->
      <el-alert
        v-if="errorMessage"
        class="notification"
        :title="errorMessage"
        type="error"
        show-icon
        :closable="true"
        @close="errorMessage = ''"
        style="margin-top: 8px"
      />
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
body {
  margin: 0;
  font-family: 'Roboto', sans-serif;
  background-color: #f5f5f5;
}
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
}
.login-card {
  background: white;
  border-radius: 8px;
  padding: 32px;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}
.login-header {
  text-align: center;
  margin-bottom: 24px;
}
.login-logo {
  width: 48px;
  height: 48px;
  margin-bottom: 8px;
}
.login-header h2 {
  margin: 0;
  color: #333;
}
.totp-qr {
  text-align: center;
  margin: 16px 0;
}
.totp-qr img {
  max-width: 200px;
}
</style>
