<template>

  <div class="wrapper">
    <div class="content">
      <div class="block1">
        <h1>Certificate Log</h1>
        <div class="row-no-gutters">
          <div style="text-align: left;background-color: #3e454e; color: white; padding: 10px;max-width: 90%;margin: auto">
            <div class="setline" id="logs">
              <p v-for="(log, index) in logs" :key="index" style="margin: 0">
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
import 'bootstrap'
import Error from '../components/Error.vue'
import { ElLoading } from 'element-plus'

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
      logs: Array,
      loading: undefined
    }
  },
  mounted () {
    this.progressShow()

    axios.get('/rest/certificate/log')
      .then((resp) => {
        this.logs = resp.data.data
        this.progressHide()
      })
      .catch(err => {
        this.$refs.error.showAxios(err)
        this.progressHide()
      })
  },
  methods: {
    progressShow () {
      this.loading = ElLoading.service({ lock: true, text: 'Loading', background: 'rgba(0, 0, 0, 0.7)' })
    },
    progressHide () {
      if (this.loading) {
        this.loading.close()
      }
    }
  }
}
</script>
<style>
@import '../style/site.css';
@import '../style/material-icons.css';
</style>
