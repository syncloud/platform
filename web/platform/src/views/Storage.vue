<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1" id="block1">
        <h1>{{ $t('storage.title') }}</h1>
        <div>
          <div class="col2">
            <div class="setline">

              <div class="setline" style="margin-top: 20px;">
                <div class="spandiv" style="font-weight: bold; margin-right: 10px;">
                  <el-switch size="large" id="multi" :active-text="$t('storage.disks')" :inactive-text="$t('storage.partitions')" v-model="multiMode" style="--el-switch-on-color: #36ad40;"/>
                </div>
                <button @click="helpVisible = true" type=button
                        class="control" style="background:transparent;">
                  <i class='fa fa-question-circle fa-lg'></i>
                </button>
              </div>

              <div>
                <span id="no_disks" class="span" v-if="disks.length === 0">{{ $t('storage.noExternalDisks') }}</span>

                <!--Single disk-->
                <div v-if="!multiMode">
                  <el-radio-group v-model="activeSinglePartition" style="display: table;">
                    <div v-for="(disk, index) in disks" :key="index">
                      <div v-for="(partition, pindex) in disk.partitions" :key="pindex">
                        <el-radio :id="'partition_' + index + '_' + pindex" :label="partition.device" size="large"
                                  border class="disk">
                        <span style="white-space: normal;">
                          {{ disk.name }}  - {{ partition.size }}
                        </span>
                        </el-radio>
                      </div>
                    </div>
                    <el-radio id="none" label="none" size="large" border class="disk" v-if="disks.length !== 0">
                      <span>{{ $t('storage.none') }}</span>
                    </el-radio>
                  </el-radio-group>
                </div>

                <!--Multi disk-->
                <div v-if="multiMode">
                  <el-checkbox-group v-model="activeMultiDisks" size="large">
                    <div v-for="(disk, index) in disks" :key="index" style="display: flex">
                      <el-checkbox :id="'disk_' + index" class="disk" size="large" border :label="disk.device">
                        <span style="white-space: normal;">
                          {{ disk.name }} - {{ disk.size }}
                          <span v-if="disk.raid">({{ disk.raid }})</span>
                        </span>
                      </el-checkbox>
                      <i v-if="disk.has_errors" class="material-icons-outlined" style="color: red; padding-top: 8px; font-size: 20px !important;">error</i>
                    </div>
                  </el-checkbox-group>
                </div>

                <!--Save-->
                <div class="setline">
                  <div class="spandiv">
                    <button class="submit buttongreen control" id="btn_save" type="submit"
                            data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Working..."
                            style="width: 150px" @click="confirmationVisible = true">{{ $t('storage.save') }}
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

  <Dialog :visible="confirmationVisible" id="confirmation" @confirm="diskAction"
                @cancel="diskActionCancel">
    <template v-slot:title>
      <div v-if="multiMode">
        <span v-if="activeMultiDisks.length !== 0">{{ $t('storage.activateDisks') }}</span>
        <span v-if="activeMultiDisks.length === 0">{{ $t('storage.deactivateDisk') }}</span>
      </div>
      <div v-if="!multiMode">
        <span v-if="activeSinglePartition !== 'none'">{{ $t('storage.activatePartition') }}</span>
        <span v-if="activeSinglePartition === 'none'">{{ $t('storage.deactivateDisk') }}</span>
      </div>
    </template>
    <template v-slot:text>
      <div style="display: grid" v-if="multiMode">
        <span style="padding-bottom: 10px">{{ $t('storage.permissionsWarning') }}</span>
        <span style="font-weight: bold;" v-for="(device, index) in activeMultiDisks" :key="index">
          {{ descriptionByDisk(device) }}
        </span>
        <el-row v-show="activeMultiDisks.length !== 0" style="align-items: center;">
          <el-col :span="24" style="min-height: 20px"></el-col>
          <el-col :span="20" style="padding-right: 10px">
            {{ $t('storage.initializeQuestion') }}
          </el-col>
          <el-col :span="4" style="text-align: right;">
            <el-switch size="large" id="format" v-model="format" style="--el-switch-on-color: Tomato;"/>
          </el-col>
        </el-row>
      </div>
      <div style="display: grid" v-if="!multiMode">
        <span v-if="activeSinglePartition !== 'none'" style="font-weight: bold;">
          {{ descriptionByPartition(activeSinglePartition) }}
        </span>
        <el-row v-show="activeSinglePartition !== 'none'" style="align-items: center;">
          <el-col :span="24" style="min-height: 20px"></el-col>
          <el-col :span="20" style="padding-right: 10px">
            {{ $t('storage.initializeQuestion') }}
          </el-col>
          <el-col :span="4" style="text-align: right;">
            <el-switch size="large" id="format" v-model="format" style="--el-switch-on-color: Tomato;"/>
          </el-col>
        </el-row>
      </div>
    </template>
  </Dialog>

  <Dialog :visible="helpVisible" @cancel="helpVisible=false" :confirm-enabled="false" :cancel-text="$t('common.close')" >
    <template v-slot:title>{{ $t('storage.externalDiskTitle') }}</template>
    <template v-slot:text>
      <div class="btext">
        {{ $t('storage.externalDiskText') }}
      </div>
    </template>
  </Dialog>
</template>

<script>
import axios from 'axios'
import Error from '../components/Error.vue'
import * as Common from '../js/common.js'
import Dialog from '../components/Dialog.vue'
import { ElLoading, ElNotification } from 'element-plus'

export default {
  name: 'Storage',
  components: {
    Dialog,
    Error
  },
  data () {
    return {
      disks: [],
      multiMode: true,
      confirmationVisible: false,
      loading: undefined,
      format: false,
      activeSinglePartition: 'none',
      activeMultiDisks: [],
      helpVisible: false
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
    reload () {
      this.uiCheckDisks()
      this.checkProgress()
      this.checkError()
    },
    uiCheckDisks () {
      axios.get('/rest/storage/disks')
        .then(resp => {
          this.activeMultiDisks = []
          this.format = false
          this.disks = resp.data.data || []
          const activeDisks = this.disks.filter(d => d.active).map(d => d.device)
          if (activeDisks && activeDisks.length > 0) {
            this.activeMultiDisks = activeDisks
          } else {
            const activePartition = this.disks.flatMap(d => d.partitions).find(p => p.active)
            if (activePartition) {
              this.activeSinglePartition = activePartition.device
              this.multiMode = false
            }
          }
          this.progressHide()
        })
        .catch(err => {
          this.progressHide()
          this.$refs.error.showAxios(err)
        })
    },
    checkProgress () {
      axios.get('/rest/job/status')
        .then(resp => {
          if (resp.data.data.name.startsWith('storage.')) {
            ElNotification({
              title: this.$t('storage.currentChange'),
              message: this.$t('storage.inProgress'),
              type: 'warning'
            })
          }
        })
        .catch(err => {
          this.progressHide()
          this.$refs.error.showAxios(err)
        })
    },
    checkError () {
      axios.get('/rest/storage/error/last')
        .catch(err => {
          ElNotification({
            title: this.$t('storage.previousChange'),
            message: err.response.data.message,
            type: 'error',
            duration: 0,
            onClose: this.clearLastError
          })
        })
    },
    clearLastError () {
      axios.post('/rest/storage/error/clear')
    },
    descriptionByDisk (device) {
      const disk = this.disks.find(d => d.device === device)
      return disk.name + ' - ' + disk.size
    },
    descriptionByPartition (device) {
      const disk = this.disks.find(d => d.partitions.some(p => p.device === device))
      return disk.name + ' - ' + disk.partitions.find(p => p.device === device).size
    },
    diskActionCancel () {
      this.confirmationVisible = false
      this.format = false
    },
    diskAction () {
      this.confirmationVisible = false
      this.progressShow()
      const error = this.$refs.error
      const that = this
      const onError = (err) => {
        this.progressHide()
        error.showAxios(err)
        this.reload()
      }
      let request = {}
      let mode = 'deactivate'
      if (this.multiMode) {
        if (this.activeMultiDisks.length !== 0) {
          mode = 'activate/disk'
          request = {
            devices: this.activeMultiDisks,
            format: this.format
          }
        }
      } else {
        if (this.activeSinglePartition !== 'none') {
          mode = 'activate/partition'
          request = {
            device: this.activeSinglePartition,
            format: this.format
          }
        }
      }
      axios
        .post('/rest/storage/' + mode, request)
        .then(resp => {
          Common.checkForServiceError(resp.data, function () {
            Common.runAfterJobIsComplete(
              setTimeout,
              that.reload,
              onError,
              Common.JOB_STATUS_URL,
              Common.JOB_STATUS_PREDICATE)
          }, onError)
        })
        .catch(onError)
    }
  }
}
</script>
<style>
@import '../style/site.css';
@import 'material-icons/iconfont/material-icons.css';

.disk {
  min-width: 300px;
  max-width: 300px
}
</style>
