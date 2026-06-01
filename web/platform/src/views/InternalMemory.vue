<template>
  <div class="sc-page">
    <div class="sc-card" id="block1">
      <h1 class="sc-title">{{ $t('internalMemory.title') }}</h1>
      <div id="wrapper">
        <h3>{{ $t('internalMemory.boot') }}</h3>
        <div class="sc-row memory-line" id="block_boot_disk" v-if="boot !== undefined">
          <span class="sc-row-label">{{ $t('internalMemory.partition', { size: boot.size }) }}</span>
          <el-button
            v-if="boot.extendable"
            id="btn_boot_extend"
            type="primary"
            @click="extend"
          >{{ $t('internalMemory.extend') }}</el-button>
        </div>
      </div>
    </div>
  </div>
  <Error ref="error"/>

</template>

<script>
import Error from '../components/Error.vue'
import * as Common from '../js/common.js'
import axios from 'axios'
import { ElLoading } from 'element-plus'

export default {
  name: 'InternalMemory',
  components: {
    Error
  },
  data () {
    return {
      boot: undefined,
      loading: undefined
    }
  },
  mounted () {
    this.progressShow()
    this.reload()
  },
  methods: {
    progressShow () {
      this.loading = ElLoading.service({ lock: true, text: this.$t('common.loading'), background: 'rgba(0, 0, 0, 0.7)' })
    },
    progressHide () {
      if (this.loading) {
        this.loading.close()
      }
    },
    extend () {
      this.progressShow()
      const that = this
      const onError = err => {
        this.progressHide()
        that.$refs.error.showAxios(err)
      }
      axios.post('/rest/storage/boot_extend')
        .then(resp => {
          Common.checkForServiceError(
            resp.data,
            () => {
              Common.runAfterJobIsComplete(
                setTimeout,
                that.reload,
                onError,
                Common.JOB_STATUS_URL,
                Common.JOB_STATUS_PREDICATE)
            },
            onError)
        })
        .catch(onError)
    },
    reload () {
      axios.get('/rest/storage/boot/disk')
        .then(resp => {
          this.boot = resp.data.data
          this.progressHide()
        })
        .catch(err => {
          this.progressHide()
          this.$refs.error.showAxios(err)
        })
    }
  }
}
</script>
<style scoped>
.memory-line .el-button {
  min-width: 120px;
}
</style>
