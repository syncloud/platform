<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12">
        <h1>Internal Memory</h1>
        <div class="row-no-gutters settingsblock" id="wrapper">
          <div class="col2">
            <div class="setline" style="margin-top: 20px;">
              <span class="span" style="font-weight: bold;">Boot</span>
            </div>
            <div class="setline memory-line" id="block_boot_disk" v-if="boot !== undefined">
              <span class="span">Partition - {{ boot.size }}</span>
              <el-button
                v-if="boot.extendable"
                id="btn_boot_extend"
                type="primary"
                @click="extend"
              >Extend</el-button>
            </div>
          </div>
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
      this.loading = ElLoading.service({ lock: true, text: 'Loading', background: 'rgba(0, 0, 0, 0.7)' })
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
<style>
@import '../style/site.css';
@import 'material-icons/iconfont/material-icons.css';
.memory-line {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.memory-line .el-button {
  min-width: 120px;
}
</style>
