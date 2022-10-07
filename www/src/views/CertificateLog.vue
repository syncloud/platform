<template>

  <div class="wrapper">
    <div class="content">
      <div class="block1">
        <h1>Certificate Log</h1>
        <div class="row-no-gutters">
          <div style="text-align: left;background-color: #3e454e; color: white; padding: 10px;max-width: 90%;margin: auto">
            <div class="setline" id="logs">
              <p v-for="(log, index) in logs" :key="index" style="margin: 0px">
                {{ log }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <Error ref="error"/>

</template>

<script>
import axios from 'axios'
import $ from 'jquery'
import 'bootstrap'
import Error from '../components/Error.vue'
import 'gasparesganga-jquery-loading-overlay'

export default {
  name: 'CertificateLog',
  components: {
    Error
  },
  props: {
    checkUserSession: Function,
    activated: Boolean
  },
  data () {
    return {
      logs: Array
    }
  },
  mounted () {
    this.progressShow()

    axios.get('/rest/certificate/log')
      .then((resp) => {
        this.logs = resp.data.data
      })
      .catch(err => {
        this.$refs.error.showAxios(err)
        this.progressHide()
      })
  },
  methods: {
    progressShow () {
      $('#block_updates').LoadingOverlay('show', { background: 'rgb(0,0,0,0)' })
    },
    progressHide () {
      $('#block_updates').LoadingOverlay('hide')
    }
  }
}
</script>
<style>
@import '../style/site.css';
@import '../style/material-icons.css';
</style>
