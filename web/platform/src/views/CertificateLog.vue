<template>
  <div class="sc-page">
    <div class="sc-card sc-card-wide" id="block1">
      <h1 class="sc-title">{{ $t('certificateLog.title') }}</h1>
      <div class="sc-console" id="logs">
        <p v-for="(log, index) in logs" :key="index" style="margin: 0">
          {{ log }}
        </p>
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
  name: 'CertificateLog',
  components: {
    Error
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
      this.loading = ElLoading.service({ lock: true, text: this.$t('common.loading'), background: 'rgba(0, 0, 0, 0.7)' })
    },
    progressHide () {
      if (this.loading) {
        this.loading.close()
      }
    }
  }
}
</script>
