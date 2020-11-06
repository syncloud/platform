<template>
  <div :id="'block_' + name" class="modal fade bs-are-use-sure" tabindex="-1" role="dialog"
       aria-labelledby="mySmallModalLabel">
    <div class="modal-dialog" role="document">
      <div class="modal-content">
        <div class="modal-header">
          <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span
            aria-hidden="true">&times;</span></button>
          <h4 class="modal-title">Error</h4>
        </div>
        <div class="modal-body">
          <div class="bodymod">
            <div :id="'txt_' + name" class="btext">Some error happened!</div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn buttonlight bwidth smbutton" data-dismiss="modal">Close</button>
            <button
              v-if="enableLogs"
              id="btn_error_send_logs"
              type="button"
              @click="sendLogs"
              data-dismiss="modal"
              class="btn buttonblue bwidth smbutton">Send logs
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div v-if="testing">
    <label for="test_parameter1"></label>
    <input id="test_parameter1"/>
  </div>
</template>
<script>
import $ from 'jquery'
import toastr from 'toastr'
import axios from 'axios'

function showFieldError (field, error) {
  const txtFieldSelector = '#' + field
  const errorBlockId = getErrorBlockId(field)
  const errorBlockSelector = '#' + errorBlockId
  const errorHtml = '<div class=\'alert alert-danger alert90\' id=\'' + errorBlockId + '\'><b>' + error + '</b></div>'
  $(errorHtml).insertAfter(txtFieldSelector)
  $(txtFieldSelector).bind('keyup change', function () {
    $(errorBlockSelector).remove()
  })
}

function getErrorBlockId (field) {
  return field + '_alert'
}

export default {
  name: 'Error',
  props: {
    name: { type: String, default: 'error' },
    enableLogs: { type: Boolean, default: true },
    testing: { type: Boolean, default: false }
  },
  methods: {
    sendLogs () {
      axios
        .post('/rest/send_log', null, { params: { include_support: true } })
        .catch(err => {
          console.log(err)
        })
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
        if (error.parameters_messages !== undefined) {
          for (let i = 0; i < error.parameters_messages.length; i++) {
            const pm = error.parameters_messages[i]
            const message = pm.messages.join('\n')
            showFieldError(pm.parameter, message)
          }
          return
        }
        if (error.message !== undefined) {
          message = error.message
        }
      }
      $('#txt_' + this.name).text(message)
      $('#block_' + this.name).modal()
    },
    showToast (error) {
      const status = error.response.status
      if (status === 401) {
        this.$router.push('/login')
      } else if (status === 0) {
        console.log('user navigated away from the page')
      } else {
        let message = 'Server Error'
        if ('data' in error.response && 'message' in error.response.data) {
          message = error.response.data.message
        }
        toastr.error(message)
      }
    }
  }
}
</script>
