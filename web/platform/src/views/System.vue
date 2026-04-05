<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12" id="block1">
        <h1>System</h1>
        <div class="row-no-gutters settingsblock">
          <div class="col2">
            <div class="setline system-action">
              <span class="span">Restart the device</span>
              <el-button id="restart" type="primary" @click="restartConfirmVisible = true">Restart</el-button>
            </div>
            <div class="setline system-action">
              <span class="span">Shut down the device</span>
              <el-button id="shutdown" type="danger" @click="shutdownConfirmVisible = true">Shutdown</el-button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <Dialog :visible="restartConfirmVisible" @confirm="restart" @cancel="restartConfirmVisible = false">
    <template #title>Restart</template>
    <template #text>Are you sure you want to restart the device?</template>
  </Dialog>

  <Dialog :visible="shutdownConfirmVisible" @confirm="shutdown" @cancel="shutdownConfirmVisible = false">
    <template #title>Shutdown</template>
    <template #text>Are you sure you want to shut down the device? You will need physical access to turn it back on.</template>
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
