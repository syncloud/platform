<template>
  <div class="sc-page">
    <div class="sc-card sc-card-wide" id="block1">
      <h1 class="sc-title">{{ $t('appCenter.title') }}</h1>
      <div class="appcenterfilter">
        <el-input
          v-model="filter"
          id="appcenter_filter"
          data-testid="appcenter_filter"
          size="large"
          :placeholder="$t('appCenter.filter')"/>
      </div>
      <div class="sc-grid" id="block_apps" style="min-height: 200px">
        <router-link v-for="(app, index) in filteredApps" :key="index" :to="'/app?id=' + app.id" class="sc-tile">
          <img :src="app.icon" class="appimg" alt="" @error="(e) => e.target.src = defaultIcon">
          <div class="sc-tile-name">{{ app.name }}</div>
        </router-link>
      </div>
    </div>
  </div>
  <Error ref="error"/>
</template>

<script>
import axios from 'axios'
import * as Common from '../js/common.js'
import Error from '../components/Error.vue'
import { ElLoading } from 'element-plus'

export default {
  name: 'AppCenter',
  data () {
    return {
      apps: undefined,
      filter: '',
      loading: undefined,
      defaultIcon: '/images/default-app.svg'
    }
  },
  components: {
    Error
  },
  computed: {
    filteredApps () {
      if (!this.apps) return []
      const q = this.filter.trim().toLowerCase()
      if (!q) return this.apps
      return this.apps.filter(a =>
        (a.name || '').toLowerCase().includes(q) ||
        (a.id || '').toLowerCase().includes(q) ||
        (a.description || '').toLowerCase().includes(q)
      )
    }
  },
  mounted () {
    this.progressShow()
    const error = this.$refs.error
    const that = this
    const onError = (err) => {
      error.showAxios(err)
      that.progressHide()
    }
    axios.get('/rest/apps/available')
      .then(
        (resp) => {
          Common.checkForServiceError(
            resp.data,
            () => {
              that.apps = resp.data.data
              that.progressHide()
            },
            onError)
        })
      .catch(onError)
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
