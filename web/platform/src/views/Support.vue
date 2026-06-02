<template>

  <div class="sc-page">
    <div class="sc-card" id="block1">
      <h1 class="sc-title">{{ $t('support.title') }}</h1>
      <div class="sc-row">
        <span class="sc-row-label">{{ $t('support.copyToSupport') }}</span>
        <s-switch id="switch" size="large" v-model="includeSupport"
                   :active-text="$t('support.yes')" :inactive-text="$t('support.no')" inline-prompt/>
      </div>
      <div class="sc-row">
        <span class="sc-row-label">{{ $t('support.reportIssue') }}</span>
        <s-button id="send"
                   :loading="loading"
                   type="primary"
                   @click="sendLogs">
          {{ $t('support.sendLogs') }}
        </s-button>
      </div>
    </div>
  </div>

  <Error ref="error" :enable-logs="false"/>

</template>

<script>
import axios from 'axios'
import Error from '../components/Error.vue'

export default {
  name: 'Support',
  components: {
    Error
  },
  data () {
    return {
      includeSupport: false,
      loading: false
    }
  },
  mounted () {
  },
  methods: {
    sendLogs () {
      this.loading = true
      axios
        .post('/rest/logs/send', null, { params: { include_support: this.includeSupport } })
        .then(() => {
          this.loading = false
        })
        .catch(err => {
          this.loading = false
          this.$refs.error.showAxios(err)
        })
    }
  }
}
</script>
