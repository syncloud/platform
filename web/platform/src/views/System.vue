<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12" id="block1">
        <h1>{{ $t('system.title') }}</h1>
        <div class="row-no-gutters settingsblock">
          <div class="col2">
            <div class="setline system-action">
              <span class="span">{{ $t('system.restartLabel') }}</span>
              <el-button id="restart" type="primary" @click="restartConfirmVisible = true">{{ $t('system.restart') }}</el-button>
            </div>
            <div class="setline system-action">
              <span class="span">{{ $t('system.shutdownLabel') }}</span>
              <el-button id="shutdown" type="danger" @click="shutdownConfirmVisible = true">{{ $t('system.shutdown') }}</el-button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <Dialog :visible="restartConfirmVisible" @confirm="restart" @cancel="restartConfirmVisible = false">
    <template #title>{{ $t('system.restart') }}</template>
    <template #text>{{ $t('system.restartConfirm') }}</template>
  </Dialog>

  <Dialog :visible="shutdownConfirmVisible" @confirm="shutdown" @cancel="shutdownConfirmVisible = false">
    <template #title>{{ $t('system.shutdown') }}</template>
    <template #text>{{ $t('system.shutdownConfirm') }}</template>
  </Dialog>

  <Error ref="error"/>
</template>

<script>
import axios from 'axios'
import Error from '../components/Error.vue'
import Dialog from '../components/Dialog.vue'

export default {
  name: 'System',
  components: {
    Error,
    Dialog
  },
  data () {
    return {
      restartConfirmVisible: false,
      shutdownConfirmVisible: false
    }
  },
  methods: {
    restart () {
      this.restartConfirmVisible = false
      axios.post('/rest/restart')
        .catch(err => this.$refs.error.showAxios(err))
    },
    shutdown () {
      this.shutdownConfirmVisible = false
      axios.post('/rest/shutdown')
        .catch(err => this.$refs.error.showAxios(err))
    }
  }
}
</script>
<style>
@import '../style/site.css';
@import 'material-icons/iconfont/material-icons.css';
.system-action {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.system-action .el-button {
  min-width: 120px;
}
</style>
