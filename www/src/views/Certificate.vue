<template>

  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12" id="block1">
        <h1>Certificate</h1>
        <div class="row-no-gutters settingsblock">
          <div class="col2">
            <div class="setline">
              <p>{{log}}</p>
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
import Error from '@/components/Error'
import Switch from '@/components/Switch'
import * as Common from '@/js/common'

export default {
  name: 'Certificate',
  components: {
    Error
  },
  props: {
    checkUserSession: Function,
  },
  data () {
    return {
      log: ''
    }
  },
  mounted () {
    this.progressShow()

    axios.get('/rest/certificate/log')
        .then((resp) => {
          this.log = resp.data
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
    },
  }
}
</script>
<style>
@import '../style/site.css';
@import '../style/material-icons.css';
</style>
