<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12" id="block1">
        <h1>Storage</h1>

        <div class="row-no-gutters settingsblock" id="block_storage">

          <div class="col2">
            <div class="setline">

              <div class="setline" style="margin-top: 20px;">
                <span class="span" style="font-weight: bold; margin-right: 10px">External disks</span>
                <div class="spandiv" style="margin-right: 10px;">
                  Multi disk <el-switch size="large" v-model="multiMode" style="--el-switch-on-color: #36ad40;" />
                </div>
                <button data-toggle="modal" data-target="#help_external_disk" type=button
                        class="control" style="background:transparent;">
                  <i class='fa fa-question-circle fa-lg'></i>
                </button>
              </div>

              <div>
                <span class="span" v-if="disks.length === 0">No external disks found</span>

                <div class="setline" v-if="disks.length !== 0">
                  <input type="radio" id="disk_none" v-model="activeSingleDisk"
                         v-bind:value="{ name: 'None', partition: {device: 'none', active: false} }"
                         :checked="activeSingleDisk.name === 'None'" style="margin-right: 10px;">
                  <span class="span">None</span>
                </div>

                <div v-for="(disk, index) in disks" :key="index">
                  <div class="setline" style="margin-top: 20px;">
                    <span class="span" style="font-weight: bold;" :id="'disk_name_' + index">
                      {{ disk.name }} - {{ disk.size }}
                    </span>
                    <div class="spandiv" v-if="!disk.active">
                      <button class="buttonred bwidth smbutton btn-lg"
                              :id="'format_' + index"
                              data-type="format"
                              data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> "
                              @click.stop="diskFormatConfirm(index, disk.device, disk.name)">Format
                      </button>
                    </div>
                  </div>
                  <div v-for="(partition, pindex) in disk.partitions" :key="pindex">
                    <div class="setline" v-if="partition.mountable || partition.active">

                      <input v-if="multiMode" type="checkbox" v-model="partition.active" style="margin-right: 10px;">
                      <input v-if="!multiMode" type="radio" :id="'disk_' + index + '_' + pindex" v-model="activeSingleDisk"
                             v-bind:value="{ name: disk.name, partition: partition }"
                             :checked="partition.active" style="margin-right: 10px;">

                      <span class="span" :id="'partition_name_' + index + '_' + pindex">
                        Partition - {{ partition.size }}
                      </span>
                    </div>
                  </div>
                </div>

                <div class="setline">
                  <br>
                  <div class="spandiv">
                    <button class="submit buttongreen control" id="btn_save" type="submit"
                            data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Working..."
                            style="width: 150px" @click="diskActionConfirm">Save
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

  <Confirmation :visible="partitionConfirmationVisible" id="partition_confirmation" @confirm="diskAction" @cancel="partitionConfirmationVisible=false">
    <template v-slot:title>
      <span v-if="activeSingleDisk.partition.active">Activate disk</span>
      <span v-if="!activeSingleDisk.partition.active">Deactivate disk</span>
    </template>
    <template v-slot:text>
      Your existing data will not be touched or moved<br>
      <span v-if="activeSingleDisk.partition.active" style="font-weight: bold;">{{ activeSingleDisk.name }}</span>
      <br>
      Are you sure?
    </template>
  </Confirmation>

  <Confirmation :visible="diskFormatConfirmationVisible" id="disk_format_confirmation" @confirm="diskFormat"
                @cancel="diskFormatConfirmationVisible = false">
    <template v-slot:title>Disk format</template>
    <template v-slot:text>
      This will destroy all the data on this disk!<br>
      <span id="disk_name" style="font-weight: bold;">{{ deviceToFormatName }}</span><br>
      Are you sure?
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
              Currently you can activate only one storage at a time.
              When activating a disk partition existing data will not be copied to the selected disk.<br><br>
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
      deviceToFormat: undefined,
      deviceToFormatIndex: undefined,
      deviceToFormatName: undefined,
      partitionActionDiskName: undefined,
      partitionActionDevice: undefined,
      partitionAction: undefined,
      multiMode: false,
      activeSingleDisk: { name: 'None', partition: { device: 'none', active: false } },
      activeMultiDisks: [],
      partitionConfirmationVisible: false,
      diskFormatConfirmationVisible: false,
      loading: undefined
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
    diskFormatConfirm (index, device, name) {
      this.deviceToFormatIndex = index
      this.deviceToFormat = device
      this.deviceToFormatName = name
      this.diskFormatConfirmationVisible = true
    },
    diskFormat () {
      this.diskFormatConfirmationVisible = false
      this.progressShow()
      const error = this.$refs.error
      const that = this
      const onError = (err) => {
        this.progressHide()
        error.showAxios(err)
        this.uiCheckDisks()
      }

      axios.post('/rest/storage/disk_format', { device: this.deviceToFormat })
        .then(function (resp) {
          Common.checkForServiceError(resp.data, function () {
            Common.runAfterJobIsComplete(
              setTimeout,
              that.uiCheckDisks,
              onError,
              Common.JOB_STATUS_URL,
              Common.JOB_STATUS_PREDICATE)
          }, onError)
        })
        .catch(onError)
    },
    uiCheckDisks () {
      axios.get('/rest/storage/disks')
        .then(resp => {
          this.disks = resp.data.data
          const activeDisk = this.disks.find(d => d.active)
          if (activeDisk) {
            this.activeSingleDisk = {
              name: activeDisk.name,
              partition: activeDisk.partitions.find(p => p.active)
            }
          }
          this.progressHide()
        })
        .catch(err => {
          console.debug(err)

          this.progressHide()
          this.$refs.error.showAxios(err)
        })
    },
    diskActionConfirm () {
      this.partitionActionDiskName = this.activeSingleDisk.name
      this.partitionActionDevice = this.activeSingleDisk.partition.device
      this.partitionAction = !this.activeSingleDisk.partition.active
      this.partitionConfirmationVisible = true
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
      const mode = this.partitionAction ? 'activate' : 'deactivate'
      axios.post('/rest/storage/disk/' + mode, { device: this.activeSingleDisk.partition.device })
        .then(resp => {
          Common.checkForServiceError(
            resp.data,
            that.uiCheckDisks,
            onError)
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
