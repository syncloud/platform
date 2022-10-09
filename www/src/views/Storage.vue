<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1" id="block1">
        <h1>Storage</h1>
        <div>
          <div class="col2">
            <div class="setline">

              <div class="setline" style="margin-top: 20px;">
                <div class="spandiv" style="font-weight: bold; margin-right: 10px;">
                  Multi disk
                  <el-switch size="large" v-model="multiMode" style="--el-switch-on-color: #36ad40;"/>
                </div>
                <button data-toggle="modal" data-target="#help_external_disk" type=button
                        class="control" style="background:transparent;">
                  <i class='fa fa-question-circle fa-lg'></i>
                </button>
              </div>

              <div>
                <span class="span" v-if="disks.length === 0">No external disks found</span>

                <!--Single disk-->
                <div v-if="!multiMode">
                  <el-radio-group v-model="activeSinglePartition" style="display: table;">
                    <div v-for="(disk, index) in disks" :key="index">
                      <el-radio v-for="(partition, pindex) in disk.partitions" :key="pindex" :label="partition.device" size="large" border style="min-width: 300px">
                        <span class="span">
                          {{ disk.name }}  - {{ partition.size }}
                        </span>
                      </el-radio>
                    </div>
                    <el-radio label="none" size="large" border style="min-width: 300px" v-if="disks.length !== 0">
                      <span class="span">None</span>
                    </el-radio>
                  </el-radio-group>
                </div>

                <!--Multi disk-->
                <div v-if="multiMode">
                  <el-checkbox-group v-model="activeMultiDisks" size="large">
                    <div v-for="(disk, index) in disks" :key="index">
                      <el-checkbox style="min-width: 300px" size="large" border :label="disk.device">
                        <span class="span">
                          {{ disk.name }} - {{ disk.size }}
                        </span>
                      </el-checkbox>
                    </div>
                  </el-checkbox-group>
                </div>

                <!--Save-->
                <div class="setline">
                  <div class="spandiv">
                    <button class="submit buttongreen control" id="btn_save" type="submit"
                            data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Working..."
                            style="width: 150px" @click="partitionConfirmationVisible = true">Save
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

  <Confirmation :visible="partitionConfirmationVisible" id="partition_confirmation" @confirm="diskAction"
                @cancel="diskActionCancel">
    <template v-slot:title>
      <div v-if="multiMode">
        <span v-if="activeMultiDisks.length !== 0">Activate multiple disks</span>
        <span v-if="activeMultiDisks.length === 0">Deactivate disk</span>
      </div>
      <div v-if="!multiMode">
        <span v-if="activeSinglePartition !== 'none'">Activate disk</span>
        <span v-if="activeSinglePartition === 'none'">Deactivate disk</span>
      </div>
    </template>
    <template v-slot:text>
      <div style="display: grid" v-if="multiMode">
        <span style="font-weight: bold;" v-for="(device, index) in activeMultiDisks" :key="index">
          {{ descriptionByDisk(device) }}
        </span>
        <span v-if="activeMultiDisks.length !== 0" style="color: Tomato;">
          It will remove all data on them!
        </span>
        <span>Are you sure?</span>
      </div>
      <div style="display: grid" v-if="!multiMode">
        <span v-if="activeSinglePartition !== 'none'" style="font-weight: bold;">
          {{ descriptionByPartition(activeSinglePartition) }}
        </span>
        <span v-if="activeSinglePartition !== 'none'">
          Initialize disk by removing all data on it?
          <el-switch size="large" v-model="format" style="--el-switch-on-color: Tomato;"/>
        </span>
        <span>Are you sure?</span>
      </div>
    </template>
  </Confirmation>

  <div id="help_external_disk" class="modal fade bs-are-use-sure" tabindex="-1" role="dialog"
       aria-labelledby="mySmallModalLabel">
    <div class="modal-dialog" role="document">
      <div class="modal-content">
        <div class="modal-header">
          <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span
            aria-hidden="true">&times;</span></button>
          <h4 class="modal-title">External disk</h4>
        </div>
        <div class="modal-body">
          <div class="bodymod">
            <div class="btext">
              Every app is configured to use storage provided by the system (which is available at /data).
              This setting screen allows you to choose which attached disk to use for that storage.<br>
              When activating a disk existing data is not copied to the selected disk.<br><br>
              You can initialize a disk by formatting it to clear all the data or to make it compatible with the system.
            </div>

          </div>
          <div class="modal-footer">
            <button type="button" class="btn buttonlight bwidth smbutton" data-dismiss="modal">Close</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
import 'bootstrap'
import Error from '../components/Error.vue'
import * as Common from '../js/common.js'
import Confirmation from '../components/Confirmation.vue'
import { ElLoading } from 'element-plus'

export default {
  name: 'Storage',
  components: {
    Confirmation,
    Error
  },
  props: {
    checkUserSession: Function,
    activated: Boolean
  },
  data () {
    return {
      disks: [],
      multiMode: false,
      partitionConfirmationVisible: false,
      loading: undefined,
      format: false,
      activeSinglePartition: undefined,
      activeMultiDisks: []
    }
  },
  mounted () {
    this.progressShow()
    this.uiCheckDisks()
  },
  methods: {
    progressShow () {
      this.loading = ElLoading.service({ lock: true, text: 'Loading', background: 'rgba(0, 0, 0, 1)' })
    },
    progressHide () {
      if (this.loading) {
        this.loading.close()
      }
    },
    uiCheckDisks () {
      axios.get('/rest/storage/disks')
        .then(resp => {
          this.format = false;
          this.disks = resp.data.data
          const activeDisks = this.disks.filter(d => d.active).map(d => d.device)
          if (activeDisks.length > 0) {
            this.activeMultiDisks = activeDisks
            this.multiMode = true
          } {
            this.activeSinglePartition = this.disks.flatMap(d => d.partitions).find(p => p.active).device
            this.multiMode = false
          }
          this.progressHide()
        })
        .catch(err => {
          this.progressHide()
          this.$refs.error.showAxios(err)
        })
    },
    descriptionByDisk(device) {
      let disk = this.disks.find(d => d.device === device)
      return disk.name + " - " + disk.size
    },
    descriptionByPartition(device) {
      let disk = this.disks.find(d => d.partitions.some(p => p.device === device))
      return disk.name + " - " + disk.partitions.find(p => p.device === device).size
    },
    diskActionCancel () {
      this.partitionConfirmationVisible=false
      this.format = false;
    },
    diskAction () {
      this.partitionConfirmationVisible = false
      this.progressShow()
      const error = this.$refs.error
      const that = this
      const onError = (err) => {
        this.progressHide()
        error.showAxios(err)
        this.uiCheckDisks()
      }
      let request = {}
      let mode = 'deactivate'
      if (this.multiMode) {
        if (this.activeMultiDisks.length !== 0) {
          mode = 'activate_multi'
          request = {
            devices: this.activeMultiDisks
          }
        }
      } else {
        if (this.activeSinglePartition !== 'none') {
          mode = 'activate'
          request = {
            device: this.activeSinglePartition,
            format: this.format
          }
        }
      }
      axios
        .post('/rest/storage/disk/' + mode, request)
        .then(resp => {
          Common.checkForServiceError(resp.data, that.uiCheckDisks, onError)
        })
        .catch(onError)
    }
  }
}
</script>
<style>
@import '../style/site.css';
@import '../style/material-icons.css';
</style>
