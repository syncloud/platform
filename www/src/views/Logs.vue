<template>

  <div class="wrapper">
    <div class="content">
      <div class="block1">
        <h1>Logs</h1>
        <div class="row-no-gutters">
          <div style="text-align: left;background-color: #3e454e; color: white; padding: 10px;max-width: 90%;margin: auto">
            <div id="logs">
              <p v-for="(log, index) in logs" :key="index" style="margin: 0; overflow-wrap: break-word">
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
import Error from '../components/Error.vue'
import { ElLoading } from 'element-plus'

export default {
  name: 'Logs',
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

    axios.get('/rest/logs')
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
@import 'material-icons/iconfont/material-icons.css';
</style>
