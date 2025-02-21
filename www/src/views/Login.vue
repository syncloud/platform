<template>

  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12" id="block1">
        <h1>Log in</h1>
        <div class="formblock">
          <form id="form-login">
            <input placeholder="Login" class="nameinput" id="username" type="text" required="" v-model="username"
                   v-on:keyup.enter="login">
            <input placeholder="Password" class="passinput" id="password" type="password" required=""
                   v-model="password" v-on:keyup.enter="login">
            <el-button class="submit control" id="btn_login"
                       style="width: 100%; height: 40px;"
                       :loading="loading"
                       type="success"
                       @click="login"
            >Log in
            </el-button>
          </form>
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
      this.loading = true
      const error = this.$refs.error
      axios.post('/rest/login', { username: this.username, password: this.password })
        .then(_ => {
          this.loading = false
          this.checkUserSession()
          this.$router.push('/')
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

input[type="text"], input[type="password"] {
  width: 100%;
  padding: 0 50px 0 20px;
  border-radius: 3px;
  border: 1px solid #dcdee0;
  margin-bottom: 10px;
  background-color: #fff !important;
  background-size: 14px 14px;
  -webkit-transition: all .3s ease-out;
  -moz-transition: all .3s ease-out;
  -o-transition: all .3s ease-out;
  -ms-transition: all .3s ease-out;
}

input:required {
  box-shadow: none;
}
</style>
