<template>
  <el-dialog v-model="visible" style="min-width: 300px; max-width: 500px">
    <template #header>
      <h4 class="modal-title">
        <slot name="title">Error</slot>
      </h4>
    </template>
    <slot name="text">
      <div class="bodymod">
        <div class="btext" id="txt_error">{{ message }}</div>
      </div>
    </slot>
    <template #footer>
      <el-button @click="close">Close</el-button>
      <el-button v-if="enableLogs" id="btn_error_send_logs" @click="sendLogs">Send logs</el-button>
    </template>
  </el-dialog>
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
      message: 'Some error happened!',
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
      let message = 'Server Error'
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
