<template>
  <div v-if="visible" class="s-modal-overlay" @click.self="close">
    <div class="s-modal syncloud-dialog" role="dialog">
      <h4 class="modal-title"><slot name="title">{{ $t('common.error') }}</slot></h4>
      <slot name="text">
        <div class="bodymod">
          <div class="btext" id="txt_error">{{ message }}</div>
        </div>
      </slot>
      <div class="s-modal-footer">
        <button class="sc-btn sc-btn-ghost" type="button" @click="close">{{ $t('common.close') }}</button>
        <button v-if="enableLogs" id="btn_error_send_logs" class="sc-btn sc-btn-ghost" type="button" @click="sendLogs">{{ $t('error.sendLogs') }}</button>
      </div>
    </div>
  </div>
</template>
<script>
import axios from 'axios'

export default {
  name: 'Error',
  props: {
    enableLogs: { type: Boolean, default: true }
  },
  data () {
    return {
      message: '',
      visible: false
    }
  },
  methods: {
    sendLogs () {
      axios
        .post('/rest/logs/send', null, { params: { include_support: true } })
        .catch(err => {
          console.log(err)
        })
      this.visible = false
    },
    showAxios (err) {
      let status = 500
      if (err.response !== undefined) {
        status = err.response.status
      }
      if (status === 401) {
        this.$router.push('/login')
        return
      }
      if (status === 501) {
        this.$router.push('/activate')
        return
      }
      if (status === 0) {
        console.log('user navigated away from the page')
        return
      }
      let message = this.$t('common.serverError')
      if (err.response !== undefined && err.response.data !== undefined) {
        const error = err.response.data
        if (error.message !== undefined) {
          message = error.message
        }
      }
      this.message = message
      this.visible = true
    },
    close () {
      this.visible = false
      this.$emit('cancel')
    }
  }
}
</script>
