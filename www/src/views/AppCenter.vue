<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12" id="block1">
        <h1 class="bigh1">App Center</h1>
        <div class="row-no-gutters appcenterlist" id="block_apps" style="min-height: 200px">
          <router-link v-for="(app, index) in apps" :key="index" :to="'/app?id=' + app.id" class="colapp app">
            <img :src="app.icon" class="appimg" alt="">
            <div class="appname"><span class="withline">{{ app.name }}</span></div>
            <div class="appdesc"></div>
          </router-link>
        </div>
      </div>
    </div>
  </div>
  <Error ref="error"/>
</template>

<script>
import axios from 'axios'
import 'bootstrap'
import 'bootstrap-switch'
import * as Common from '../js/common.js'
import Error from '@/components/Error'
import 'gasparesganga-jquery-loading-overlay'
import $ from 'jquery'

export default {
  name: 'AppCenter',
  props: {
    checkUserSession: Function,
    activated: Boolean
  },
  data () {
    return {
      apps: undefined
    }
  },
  components: {
    Error
  },
  mounted () {
    this.progressShow()
    const error = this.$refs.error
    const that = this
    const onError = (err) => {
      error.showAxios(err)
      that.progressHide()
    }
    axios.get('/rest/apps/available')
      .then(
        (resp) => {
          Common.checkForServiceError(
            resp.data,
            () => {
              that.apps = resp.data.data
              that.progressHide()
            },
            onError)
        })
      .catch(onError)
  },
  methods: {
    progressShow () {
      $('#block_apps').LoadingOverlay('show', { background: 'rgb(0,0,0,0)' })
    },
    progressHide () {
      $('#block_apps').LoadingOverlay('hide')
    }
  }
}
</script>
<style>
@import '../style/site.css';
@import '../style/material-icons.css';
</style>
