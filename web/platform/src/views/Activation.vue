<template>
  <div class="sc-page">
    <div class="sc-card" id="block1">
      <h1 class="sc-title">{{ $t('activation.title') }}</h1>

      <div class="sc-field">
        <span class="sc-label">{{ $t('activation.activatedAt') }}</span>
        <a id="txt_device_domain" :href="url">{{ url }}</a>
      </div>

      <div class="sc-field">
        <span class="sc-label" style="white-space: pre-line; font-weight: 400; color: var(--sc-muted)">{{ $t('activation.reassign') }}</span>
      </div>

      <div class="sc-actions">
        <button class="sc-btn sc-btn-success" id="btn_reactivate" @click="reactivate">{{ $t('activation.reactivate') }}</button>
      </div>
    </div>
  </div>

  <Error ref="error"/>

</template>

<script>
import axios from 'axios'
import Error from '../components/Error.vue'
import { ElLoading } from 'element-plus'
import { useAuthStore } from '../stores/auth'

export default {
  name: 'Activation',
  data () {
    return {
      url: ''
    }
  },
  components: {
    Error
  },
  mounted () {
    this.progressShow()
    axios
      .get('/rest/device/url')
      .then(resp => {
        this.url = resp.data.data
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
        .post('/rest/deactivate')
        .then(() => {
          const auth = useAuthStore()
          auth.checkUserSession(this.$router)
        })
        .catch(err => {
          this.$refs.error.showAxios(err)
        })
    },
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
