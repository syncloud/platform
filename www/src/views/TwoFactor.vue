<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12" id="block1">
        <h1>Two-Factor Authentication</h1>
        <div class="setblock">
          <div class="setline">
            <span class="setname">Status: </span>
            <span class="setvalue" id="twofa_status">{{ enabled ? 'Enabled' : 'Disabled' }}</span>
          </div>

          <div class="setline">
            <el-button
              v-if="!enabled"
              id="btn_enable_2fa"
              type="success"
              :loading="loading"
              @click="enableTwoFactor"
            >Enable 2FA</el-button>
            <el-button
              v-else
              id="btn_disable_2fa"
              type="danger"
              :loading="loading"
              @click="disableTwoFactor"
            >Disable 2FA</el-button>
          </div>

          <div v-if="totpQr" class="setline">
            <el-alert type="warning" :closable="false" show-icon>
              Scan this QR code now. It will not be shown again.
              If you lose access to your authenticator app, see
              <a href="https://github.com/syncloud/platform/wiki/Two-Factor-Authentication#recovery" target="_blank">recovery instructions</a>.
            </el-alert>
            <p>Scan this QR code with your authenticator app:</p>
            <img id="totp_qr" :src="totpQr" alt="TOTP QR Code" />
            <p>Or enter this secret manually:</p>
            <code id="totp_secret">{{ totpSecret }}</code>
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
import QRCode from 'qrcode'

export default {
  name: 'TwoFactor',
  data () {
    return {
      enabled: false,
      loading: false,
      totpQr: null,
      totpSecret: null
    }
  },
  components: {
    Error
  },
  mounted () {
    this.load()
  },
  methods: {
    load () {
      const error = this.$refs.error
      axios.get('/rest/settings/2fa')
        .then(resp => {
          this.enabled = resp.data.data.enabled
        })
        .catch(err => {
          error.showAxios(err)
        })
    },
    enableTwoFactor () {
      this.loading = true
      const error = this.$refs.error
      axios.post('/rest/settings/2fa', { enabled: true })
        .then(() => {
          this.enabled = true
          return axios.post('/rest/settings/2fa/totp')
        })
        .then(resp => {
          this.loading = false
          const uri = resp.data.data.uri
          const params = new URL(uri).searchParams
          this.totpSecret = params.get('secret')
          return QRCode.toDataURL(uri)
        })
        .then(qr => {
          this.totpQr = qr
        })
        .catch(err => {
          this.loading = false
          error.showAxios(err)
        })
    },
    disableTwoFactor () {
      this.loading = true
      const error = this.$refs.error
      axios.post('/rest/settings/2fa', { enabled: false })
        .then(() => {
          this.loading = false
          this.enabled = false
          this.totpQr = null
          this.totpSecret = null
        })
        .catch(err => {
          this.loading = false
          error.showAxios(err)
        })
    }
  }
}
</script>
<style>
@import '../style/site.css';
@import 'material-icons/iconfont/material-icons.css';
</style>
