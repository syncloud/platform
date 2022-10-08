<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12" id="block1">
        <h1>Activation</h1>
        <div class="row-no-gutters settingsblock">

          <div class="col2">
            <div class="setline">
              <span class="span">Activated at: </span><a id="txt_device_domain" :href="url">{{ url }}</a>
            </div>

            <div class="setline">
              <span class="span">You can assign different<br>domain name to your device</span>
              <div class="spandiv">
                <button class="buttongreen bwidth smbutton" id="btn_reactivate" @click="reactivate">Reactivate</button>
              </div>
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
import Error from '../components/Error.vue'
import 'bootstrap'
import { ElLoading } from 'element-plus'

export default {
  name: 'Activation',
  props: {
    checkUserSession: Function,
    activated: Boolean
  },
  data () {
    return {
      url: 'Loading ...'
    }
  },
  components: {
    Error
  },
  mounted () {
    this.progressShow()
    axios
      .get('/rest/settings/device_url')
      .then(resp => {
        this.url = resp.data.device_url
        this.progressHide()
      })
      .catch(err => {
        this.progressHide()
        this.$refs.error.showAxios(err)
      })
  },
  methods: {
    reactivate: function () {
      axios
        .post('/rest/settings/deactivate')
        .then(_ => {
          this.checkUserSession()
        })
        .catch(err => {
          this.$refs.error.showAxios(err)
        })
    },
    progressShow () {
      this.loading = ElLoading.service({ lock: true, text: 'Loading', background: 'rgba(0, 0, 0, 0.7)' })
    },
    progressHide () {
      if (this.loading) {
        this.loading.close()
      }
    },
  }
}
</script>
<style>
@import '../style/site.css';
@import '../style/material-icons.css';
</style>
