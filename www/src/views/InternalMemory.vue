<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12">
        <h1>Internal Memory</h1>
        <div class="row-no-gutters settingsblock" id="wrapper">
          <div class="col2">
            <div class="setline">
              <div class="setline" style="margin-top: 20px;">
                <span class="span" style="font-weight: bold;">Boot</span>
              </div>

              <div class="setline">
                <div id="block_boot_disk" v-if="boot !== undefined">
                  <span class="span">Partition - {{ boot.size }}</span>
                  <div class="spandiv" v-if="boot.extendable">
                    <button class="buttongreen bwidth smbutton btn-lg"
                            @click="extend"
                            id="btn_boot_extend"
                            data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Extending...">Extend
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <Error ref="error"/>

</template>

<script>
import $ from 'jquery'
import 'bootstrap'
import 'bootstrap-switch'
import Error from '@/components/Error.vue'
import * as Common from '../js/common.js'
import axios from 'axios'
import 'gasparesganga-jquery-loading-overlay'

export default {
  name: 'InternalMemory',
  components: {
    Error
  },
  props: {
    checkUserSession: Function,
    activated: Boolean
  },
  data () {
    return {
      boot: undefined
    }
  },
  mounted () {
    this.progressShow()
    this.reload()
  },
  methods: {
    progressShow () {
      $('#wrapper').LoadingOverlay('show', { background: 'rgb(0,0,0,0)' })
    },
    progressHide () {
      $('#wrapper').LoadingOverlay('hide')
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
@import '../style/material-icons.css';
</style>
