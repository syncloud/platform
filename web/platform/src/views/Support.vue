<template>

  <div class="sc-page">
    <div class="sc-card" id="block1">
      <h1 class="sc-title" data-testid="support-title">{{ $t('support.title') }}</h1>
      <p class="sc-lead support-hint" data-testid="support-hint">{{ includeSupport ? $t('support.hintWithSupport') : $t('support.hintOwnerOnly') }}</p>
      <div class="sc-row">
        <span class="sc-row-label">{{ $t('support.copyToSupport') }}</span>
        <s-switch id="switch" data-testid="support-include" size="large" v-model="includeSupport"
                   :active-text="$t('support.yes')" :inactive-text="$t('support.no')" inline-prompt/>
      </div>
      <div class="sc-row">
        <span class="sc-row-label">{{ $t('support.reportIssue') }}</span>
        <s-button id="send"
                   data-testid="support-send"
                   :loading="loading"
                   type="primary"
                   @click="sendLogs">
          {{ $t('support.sendLogs') }}
        </s-button>
      </div>
      <s-alert v-if="sent" class="support-sent" data-testid="support-sent" type="success" :closable="false" show-icon>
        {{ sentToSupport ? $t('support.sentWithSupport') : $t('support.sentOwnerOnly') }}
      </s-alert>
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
      loading: false,
      sent: false,
      sentToSupport: false
    }
  },
  methods: {
    sendLogs () {
      this.loading = true
      this.sent = false
      const includeSupport = this.includeSupport
      axios
        .post('/rest/logs/send', null, { params: { include_support: includeSupport } })
        .then(() => {
          this.loading = false
          this.sentToSupport = includeSupport
          this.sent = true
        })
        .catch(err => {
          this.loading = false
          this.$refs.error.showAxios(err)
        })
    }
  }
}
</script>

<style scoped>
.support-hint { margin-bottom: 16px; }
.support-sent { margin-top: 16px; }
</style>
