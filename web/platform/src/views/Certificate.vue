<template>
  <div class="sc-page">
    <div class="sc-card" id="block1">
      <h1 class="sc-title">{{ $t('certificate.title') }}</h1>

      <div class="sc-row">
        <span class="sc-row-label">{{ $t('certificate.valid') }}</span>
        <i v-if="valid" class="material-icons icon-good" id="valid_good">check_circle</i>
        <i v-if="!valid" class="material-icons icon-bad" id="valid_bad">error</i>
      </div>
      <div class="sc-row">
        <span class="sc-row-label">{{ $t('certificate.validDays') }}</span>
        <span id="valid_days">{{ validDays }}</span>
      </div>
      <div class="sc-row">
        <span class="sc-row-label">{{ $t('certificate.real') }}</span>
        <i v-if="real" class="material-icons icon-good" id="real_good">check_circle</i>
        <i v-if="!real" class="material-icons icon-bad" id="real_bad">error</i>
      </div>
      <div class="sc-row">
        <span class="sc-row-label">{{ $t('certificate.moreDetails') }}</span>
        <router-link to="/certificate/log">
          <button class="sc-btn sc-btn-primary">{{ $t('certificate.log') }}</button>
        </router-link>
      </div>
    </div>
  </div>

  <Error ref="error"/>

</template>

<script>
import axios from 'axios'
import Error from '../components/Error.vue'
import Loading from '../util/loading'

export default {
  name: 'Certificate',
  components: {
    Error
  },
  data () {
    return {
      valid: false,
      real: false,
      validDays: 0,
      loading: undefined
    }
  },
  mounted () {
    this.progressShow()
    axios.get('/rest/certificate')
      .then((resp) => {
        const data = resp.data.data
        this.valid = data.is_valid
        this.real = data.is_real
        this.validDays = data.valid_for_days
        this.progressHide()
      })
      .catch(err => {
        this.$refs.error.showAxios(err)
        this.progressHide()
      })
  },
  methods: {
    progressShow () {
      this.loading = Loading.service({ lock: true, text: this.$t('common.loading'), background: 'rgba(0, 0, 0, 0.7)' })
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
.icon-good { font-size: 22px; color: var(--sc-success); }
.icon-bad { font-size: 22px; color: var(--sc-danger); }
</style>
