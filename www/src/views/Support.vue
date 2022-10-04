<template>

  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12" id="block1">
        <h1>Support</h1>
        <div class="row-no-gutters settingsblock">
          <div class="col2">
            <div class="setline">
              <span class="span">Send a copy to support</span>
              <Switch
                :checked="includeSupport"
                @toggle="includeSupport = !includeSupport"
                on-label="Yes"
                off-label="No"
              />
            </div>
            <div class="setline">
              <span class="span">Report Issue:</span>
              <button
                id="send"
                class="buttonblue bwidth smbutton"
                @click="sendLogs"
                data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Sending...">Send logs
              </button>
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
import Error from '@/components/Error.vue'
import Switch from '@/components/Switch.vue'

export default {
  name: 'Support',
  components: {
    Error,
    Switch
  },
  props: {
    checkUserSession: Function,
    activated: Boolean
  },
  data () {
    return {
      includeSupport: false
    }
  },
  mounted () {
  },
  methods: {
    sendLogs () {
      const btn = $('#send')
      btn.button('loading')
      axios
        .post('/rest/send_log', null, { params: { include_support: this.includeSupport } })
        .then(_ => {
          btn.button('reset')
        })
        .catch(this.$refs.error.showAxios)
    }
  }
}
</script>
<style>
@import '../style/site.css';
@import '../style/material-icons.css';
</style>
