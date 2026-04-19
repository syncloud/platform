<template>

  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12" id="block1">
        <h1>{{ $t('support.title') }}</h1>
        <div class="row-no-gutters settingsblock">
          <div class="col2">
            <div class="setline">
              <span class="span">{{ $t('support.copyToSupport') }}</span>
              <el-switch id="switch" size="large" v-model="includeSupport" style="--el-switch-on-color: #36ad40;"
                         :active-text="$t('support.yes')" :inactive-text="$t('support.no')" inline-prompt/>
            </div>
            <div class="setline">
              <span class="span">{{ $t('support.reportIssue') }}</span>
              <el-button id="send"
                         :loading="loading"
                         type="primary"
                         @click="sendLogs">
                {{ $t('support.sendLogs') }}
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <Error ref="error" :enable-logs="false"/>

</template>

<script>
import axios from 'axios'
import Error from '../components/Error.vue'

export default {
  name: 'Support',
  components: {
    Error
  },
  data () {
    return {
      includeSupport: false,
      loading: false
    }
  },
  mounted () {
  },
  methods: {
    sendLogs () {
      this.loading = true
      axios
        .post('/rest/logs/send', null, { params: { include_support: this.includeSupport } })
        .then(() => {
          this.loading = false
        })
        .catch(err => {
          this.loading = false
          this.$refs.error.showAxios(err)
        })
    }
  }
}
</script>
<style>
@import '../style/site.css';
@import 'material-icons/iconfont/material-icons.css';
</style>
