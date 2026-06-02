<template>
  <div class="sc-page">
    <div class="sc-card sc-card-wide" id="block1">
      <h1 class="sc-title">{{ $t('certificateLog.title') }}</h1>
      <label class="logs-wrap">
        <input type="checkbox" id="logs_wrap" data-testid="logs-wrap" v-model="wrap">
        {{ $t('logs.wordWrap') }}
      </label>
      <div class="sc-console" id="logs" :class="{ nowrap: !wrap }">
        <p v-for="(log, index) in logs" :key="index" class="logs-line">
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
      loading: undefined,
      wrap: true
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
<style scoped>
.logs-wrap {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
  font-size: 14px;
  color: var(--sc-muted);
  margin-bottom: 10px;
  cursor: pointer;
}
.logs-line {
  margin: 0;
  white-space: pre-wrap;
  overflow-wrap: break-word;
  word-break: break-word;
}
.sc-console.nowrap .logs-line {
  white-space: pre;
  overflow-wrap: normal;
  word-break: normal;
}
</style>
