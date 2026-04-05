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

          <div v-if="enabled" class="setline" style="max-width: 500px; margin: 10px auto">
            <el-alert type="info" :closable="false" show-icon>
              Users will be prompted to set up their authenticator app during login.
            </el-alert>
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
  data () {
    return {
      enabled: false,
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
          this.loading = false
          this.enabled = true
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
