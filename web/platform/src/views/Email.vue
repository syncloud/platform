<template>
  <div class="sc-page">
    <div class="sc-card" id="block1">
      <h1 class="sc-title">{{ $t('email.title') }}</h1>
      <div :style="{ visibility: visibility }">

        <div class="access-block" data-testid="mail-relay-section">
          <div class="access-head">
            <div class="access-heading">
              <span class="access-name">{{ $t('email.relay') }}</span>
            </div>
            <s-switch id="tgl_mail_relay" data-testid="mail-relay-toggle" size="large" v-model="relayEnabled"
                      style="--el-switch-on-color: #2faa5d"/>
          </div>
          <p class="detail-desc">{{ $t('email.relayDescription') }}</p>
        </div>

        <div class="sc-actions">
          <button class="sc-btn sc-btn-success" id="btn_save" data-testid="mail-relay-save" type="submit"
                  @click="save">{{ $t('email.save') }}
          </button>
        </div>
      </div>
    </div>
  </div>
  <Error ref="error"/>
</template>

<script>
import Error from '../components/Error.vue'
import axios from 'axios'
import Loading from '../util/loading'

export default {
  name: 'Email',
  data () {
    return {
      relayEnabled: undefined,
      visibility: 'hidden',
      loading: undefined
    }
  },
  components: {
    Error
  },
  mounted () {
    this.progressShow()
    this.reload()
  },
  methods: {
    progressShow (text) {
      this.loading = Loading.service({ lock: true, text: text || this.$t('common.loading'), background: 'rgba(0, 0, 0, 0.7)' })
    },
    progressHide () {
      this.visibility = 'visible'
      this.loading.close()
    },
    reload () {
      axios.get('/rest/mail_relay')
        .then(resp => {
          this.relayEnabled = resp.data.data.enabled
          this.progressHide()
        })
        .catch(e => {
          this.$refs.error.showAxios(e)
          this.progressHide()
        })
    },
    save () {
      this.progressShow()
      axios.post('/rest/mail_relay', { enabled: this.relayEnabled })
        .then(() => {
          this.progressHide()
          this.reload()
        })
        .catch(e => {
          this.$refs.error.showAxios(e)
          this.progressHide()
        })
    }
  }
}
</script>
<style scoped>
.access-block {
  border: 1px solid var(--sc-border);
  border-radius: 16px;
  padding: 20px 22px;
  margin-bottom: 16px;
  background: var(--sc-surface);
}
.access-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}
.access-heading {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}
.access-name {
  font-size: 17px;
  font-weight: 700;
  color: var(--sc-ink);
}
.detail-desc {
  color: var(--sc-muted);
  font-size: 14px;
  margin: 8px 0 0;
  max-width: 460px;
  line-height: 1.5;
}
</style>
