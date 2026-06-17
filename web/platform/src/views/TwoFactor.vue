<template>
  <div class="sc-page">
    <div class="sc-card" id="block1">
      <h1 class="sc-title">{{ $t('twoFactor.title') }}</h1>
      <div class="sc-row">
        <span class="sc-row-label">{{ $t('twoFactor.status') }}</span>
        <span id="twofa_status" style="font-weight: 600">{{ enabled ? $t('twoFactor.enabled') : $t('twoFactor.disabled') }}</span>
      </div>

      <div class="sc-actions">
        <s-button
          v-if="!enabled"
          id="btn_enable_2fa"
          type="success"
          :loading="loading"
          @click="enableTwoFactor"
        >{{ $t('twoFactor.enable') }}</s-button>
        <s-button
          v-else
          id="btn_disable_2fa"
          type="danger"
          :loading="loading"
          @click="disableTwoFactor"
        >{{ $t('twoFactor.disable') }}</s-button>
      </div>

      <div v-if="enabled" style="margin-top: 16px">
        <s-alert type="info" :closable="false" show-icon>
          {{ $t('twoFactor.enabledNote') }}
        </s-alert>
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
      axios.get('/rest/2fa')
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
      axios.post('/rest/2fa', { enabled: true })
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
      axios.post('/rest/2fa', { enabled: false })
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
