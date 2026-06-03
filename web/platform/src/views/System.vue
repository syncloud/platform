<template>
  <div class="sc-page">
    <div class="sc-card" id="block1">
      <h1 class="sc-title">{{ $t('system.title') }}</h1>
      <div class="sc-row">
        <span class="sc-row-label">{{ $t('system.restartLabel') }}</span>
        <s-button id="restart" type="primary" @click="restartConfirmVisible = true">{{ $t('system.restart') }}</s-button>
      </div>
      <div class="sc-row">
        <span class="sc-row-label">{{ $t('system.shutdownLabel') }}</span>
        <s-button id="shutdown" type="danger" @click="shutdownConfirmVisible = true">{{ $t('system.shutdown') }}</s-button>
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
