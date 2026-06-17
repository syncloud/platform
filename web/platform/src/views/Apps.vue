<template>
  <div class="sc-page">
    <div class="sc-card sc-card-wide" id="block1">
      <h1 class="sc-title">{{ $t('apps.title') }}</h1>
      <div id="block_apps">
        <div v-if="apps.length === 0" class="sc-empty">
          <p class="sc-lead">{{ $t('apps.emptyHeading') }}</p>
          <router-link v-if="auth.admin" to="/appcenter" class="sc-link">{{ $t('apps.appCenterLink') }}</router-link>
        </div>
        <div v-else class="sc-grid">
          <router-link v-for="(app, index) in apps" :key="index" :to="'/app?id=' + app.id" class="sc-tile"
                       :data-testid="'app-tile-' + app.id">
            <img :src="app.icon" class="appimg" :alt="app.name" @error="(e) => e.target.src = defaultIcon">
            <div class="sc-tile-name">{{ app.name }}</div>
          </router-link>
        </div>
      </div>
    </div>
  </div>
  <Error ref="error"/>

</template>

<script>
import axios from 'axios'
import Error from '../components/Error.vue'
import Loading from '../util/loading'
import { useAuthStore } from '../stores/auth'

export default {
  name: 'Apps',
  components: {
    Error
  },
  data () {
    return {
      apps: Array,
      loading: undefined,
      defaultIcon: '/images/default-app.svg',
      auth: useAuthStore()
    }
  },
  mounted () {
    this.progressShow()
    axios.get('/rest/apps/installed')
      .then(resp => {
        if (resp.data.data == null) {
          this.apps = []
        } else {
          this.apps = resp.data.data
        }
        this.progressHide()
      })
      .catch(err => {
        this.progressHide()
        this.$refs.error.showAxios(err)
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
