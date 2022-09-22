<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12" id="block1">
        <h1 class="bigh1">Applications</h1>
          <div class="row-no-gutters appcenterlist" id="block_apps">
            <div v-if="apps.length === 0">
              <h2 class="bh2">You don't have any apps installed yet. You can install one from App Center</h2>
              <router-link to="/appcenter" class="appcenterh">App Center</router-link>
            </div>
            <router-link v-for="(app, index) in apps" :key="index" :to="'/app?id=' + app.id" class="colapp app">
              <img :src="app.icon" class="appimg" :alt="app.name">
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
import Error from '@/components/Error'
import 'gasparesganga-jquery-loading-overlay'
import $ from 'jquery'

export default {
  name: 'Apps',
  props: {
    checkUserSession: Function,
    activated: Boolean
  },
  components: {
    Error
  },
  data () {
    return {
      apps: Array
    }
  },
  mounted () {
    this.progressShow()
    axios.get('/rest/apps/installed')
      .then(resp => {
        if (resp.data.data == null) {
            this.apps = []
        } else {
            this.apps = resp.data.data
        }
        this.progressHide()
      })
      .catch(err => {
        this.progressHide()
        this.$refs.error.showAxios(err)
      })
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
