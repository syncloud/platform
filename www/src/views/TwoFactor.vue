<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12" id="block1">
        <h1>Two-Factor Authentication</h1>
        <div class="setblock">
          <div class="setline">
            <span class="setname">Status</span>
            <span class="setvalue" id="twofa_status">{{ enabled ? 'Enabled' : 'Disabled' }}</span>
          </div>

          <div v-if="!enabled" class="setline">
            <p>Before enabling 2FA, register your authenticator device in Authelia settings.</p>
            <el-button id="btn_authelia_settings" type="primary" @click="openAutheliaSettings">
              Open Authelia Settings
            </el-button>
          </div>

          <div class="setline" style="margin-top: 20px;">
            <el-button
              v-if="!enabled"
              id="btn_enable_2fa"
              type="success"
              :loading="loading"
              @click="setTwoFactor(true)"
            >Enable 2FA</el-button>
            <el-button
              v-else
              id="btn_disable_2fa"
              type="danger"
              :loading="loading"
              @click="setTwoFactor(false)"
            >Disable 2FA</el-button>
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

export default {
  name: 'TwoFactor',
  props: {
    checkUserSession: Function,
    activated: Boolean
  },
  data () {
    return {
      enabled: false,
      autheliaUrl: '',
      loading: false
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
          this.autheliaUrl = resp.data.data.authelia_url
        })
        .catch(err => {
          error.showAxios(err)
        })
    },
    openAutheliaSettings () {
      window.open(this.autheliaUrl + '/settings', '_blank')
    },
    setTwoFactor (enabled) {
      this.loading = true
      const error = this.$refs.error
      axios.post('/rest/settings/2fa', { enabled: enabled })
        .then(() => {
          this.loading = false
          this.enabled = enabled
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
