<template>

  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12" id="block1">
        <h1>Log in</h1>
        <div class="formblock">
          <form id="form-login">
            <input placeholder="Login" class="nameinput" id="username" type="text" required="" v-model="username">
            <input placeholder="Password" class="passinput" id="password" type="password" required=""
                   v-model="password">
            <button class="submit buttongreen control" id="btn_login" type="submit"
                    data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Logging in..."
                    @click="login"
            >Log in
            </button>
          </form>
        </div>
      </div>
    </div>
  </div>

  <div id="block_error" class="modal fade bs-are-use-sure" tabindex="-1" role="dialog"
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
            <div id="txt_error" class="btext">Some error happened!</div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn buttonlight bwidth smbutton" data-dismiss="modal">Close</button>
            <button id="btn_error_send_logs" type="button" class="btn buttonblue bwidth smbutton">Send logs</button>
          </div>
        </div>
      </div>
    </div>
  </div>

  <Error ref="error"/>

</template>

<script>
import $ from 'jquery'
import Error from '@/components/Error'
import axios from 'axios'
import 'bootstrap'

export default {
  name: 'Login',
  props: {
    checkUserSession: Function,
    activated: Boolean
  },
  data () {
    return {
      username: '',
      password: '',
      loading: false
    }
  },
  components: {
    Error
  },
  methods: {
    login: function (event) {
      const error = this.$refs.error
      event.preventDefault()
      const btn = $('#btn_login')
      btn.button('loading')
      $('#form-login input').prop('disabled', true)
      $('#form-login .alert').remove()
      axios.post('/rest/login', { username: this.username, password: this.password })
        .then(_ => {
          btn.button('reset')
          $('#form-login input').prop('disabled', false)
          this.checkUserSession()
          this.$router.push('/')
        })
        .catch(err => {
          btn.button('reset')
          $('#form-login input').prop('disabled', false)
          error.showAxios(err)
        })
    }
  }
}
</script>
<style>
@import '../style/site.css';
@import '../style/material-icons.css';

input:required {
  box-shadow: none;
}
</style>
