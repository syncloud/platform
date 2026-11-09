<template>
  <div class="sc-page">
    <div class="sc-card" id="block1">
      <h1 class="sc-title">{{ $t('storage.title') }}</h1>
      <div class="sc-row" style="border: none; justify-content: flex-start; gap: 14px;">
        <s-switch size="large" id="multi" :active-text="$t('storage.disks')" :inactive-text="$t('storage.partitions')" v-model="multiMode"/>
        <button @click="helpVisible = true" type=button class="control" style="background:transparent; border: none; cursor: pointer;">
          <i class='fa fa-question-circle fa-lg' style="color: var(--sc-faint)"></i>
        </button>
      </div>

      <span id="no_disks" class="sc-lead" v-if="disks.length === 0">{{ $t('storage.noExternalDisks') }}</span>

      <!--Single disk-->
      <div v-if="!multiMode">
        <s-radio-group v-model="activeSinglePartition">
          <template v-for="(disk, index) in disks" :key="index">
            <s-radio v-for="(partition, pindex) in disk.partitions" :key="pindex"
                      :id="'partition_' + index + '_' + pindex" :label="partition.device" size="large"
                      border class="disk">
              <span style="white-space: normal;">
                {{ disk.name }}  - {{ partition.size }}
              </span>
            </s-radio>
          </template>
          <s-radio id="none" label="none" size="large" border class="disk" v-if="disks.length !== 0">
            <span>{{ $t('storage.none') }}</span>
          </s-radio>
        </s-radio-group>
      </div>

      <!--Multi disk-->
      <div v-if="multiMode">
        <s-checkbox-group v-model="activeMultiDisks" size="large" class="disk-list">
          <div v-for="(disk, index) in disks" :key="index" class="disk-row">
            <s-checkbox :id="'disk_' + index" class="disk-choice" size="large" border :label="disk.device">
              <span style="white-space: normal;">
                {{ disk.name }} - {{ disk.size }}
                <span v-if="disk.raid">({{ disk.raid }})</span>
              </span>
            </s-checkbox>
            <i v-if="status(disk)" class="material-icons-outlined disk-status" :class="'disk-status-' + status(disk)"
               :title="$t('storage.status_' + status(disk))" :data-testid="'disk-status-' + index"
               :data-status="status(disk)">{{ statusIcon[status(disk)] }}</i>
            <router-link :to="{ path: '/storage/disk', query: { device: disk.device } }" class="disk-details-link"
                         :title="$t('storage.details')" :data-testid="'disk-details-' + index">
              <i class="material-icons-outlined">chevron_right</i>
            </router-link>
          </div>
        </s-checkbox-group>
      </div>

      <router-link v-if="disks.length !== 0" to="/storage/checks" class="checks-link" data-testid="storage-checks-link"
                   :class="{ 'checks-link-bad': scrubBad }">
        <span class="checks-link-text">
          <span class="checks-link-title">{{ $t('storage.checks') }}</span>
          <span class="sc-muted" data-testid="storage-checks-summary">{{ checksSummary }}</span>
        </span>
        <i class="material-icons-outlined">chevron_right</i>
      </router-link>

      <!--Save-->
      <div class="sc-actions">
        <button class="sc-btn sc-btn-success" id="btn_save" type="submit"
                data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Working..."
                @click="confirmationVisible = true">{{ $t('storage.save') }}
        </button>
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
        <div v-show="activeMultiDisks.length !== 0" class="format-row">
          <span>{{ $t('storage.initializeQuestion') }}</span>
          <s-switch size="large" id="format" v-model="format"/>
        </div>
      </div>
      <div style="display: grid" v-if="!multiMode">
        <span v-if="activeSinglePartition !== 'none'" style="font-weight: bold;">
          {{ descriptionByPartition(activeSinglePartition) }}
        </span>
        <div v-show="activeSinglePartition !== 'none'" class="format-row">
          <span>{{ $t('storage.initializeQuestion') }}</span>
          <s-switch size="large" id="format" v-model="format"/>
        </div>
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
import Loading from '../util/loading'
import notify from '../util/notify'
import { diskStatus, healthOf, scrubText } from '../js/storage'
import { useStorageAlertsStore } from '../stores/storageAlerts'

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
      helpVisible: false,
      health: { disks: [], scrub: null },
      statusIcon: { ok: 'check_circle', warning: 'warning', error: 'error' }
    }
  },
  computed: {
    scrubBad () {
      return !!(this.health.scrub && (this.health.scrub.errors || this.health.scrub.uncorrectable > 0))
    },
    checksSummary () {
      if (this.health.scrub) {
        return this.$t('storage.scrub') + ': ' + scrubText((...args) => this.$t(...args), this.health.scrub)
      }
      return this.$t('storage.checksSummary')
    }
  },
  mounted () {
    this.progressShow()
    this.reload()
    useStorageAlertsStore().seen().catch(err => this.$refs.error.showAxios(err))
  },
  methods: {
    progressShow () {
      this.loading = Loading.service({ lock: true, text: this.$t('common.loading'), background: 'rgba(0, 0, 0, 0.7)' })
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
      this.reloadHealth()
    },
    reloadHealth () {
      axios.get('/rest/storage/health')
        .then(resp => { this.health = resp.data.data })
        .catch(err => this.$refs.error.showAxios(err))
    },
    status (disk) {
      return diskStatus(disk, healthOf(this.health, disk))
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
            notify({
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
          notify({
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
<style scoped>
.disk {
  min-width: 300px;
  max-width: 300px
}
.disk-row {
  display: flex;
  align-items: center;
  gap: 6px;
  width: 100%;
}
.disk-list {
  display: flex;
  width: 100%;
}
.disk-choice {
  flex: 1 1 auto;
  min-width: 0;
}
.disk-choice :deep(.s-choice-body) {
  min-width: 0;
  overflow-wrap: anywhere;
}
.disk-status,
.disk-details-link {
  flex-shrink: 0;
}
.disk-status {
  font-size: 22px !important;
}
.disk-status-ok {
  color: var(--sc-success);
}
.disk-status-warning {
  color: var(--sc-warn);
}
.disk-status-error {
  color: var(--sc-danger);
}
.disk-details-link {
  display: flex;
  color: var(--sc-muted);
  text-decoration: none;
}
.checks-link {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-top: 16px;
  padding: 12px 14px;
  border: 1px solid var(--sc-border);
  border-radius: var(--sc-control-radius);
  color: var(--sc-ink);
  text-decoration: none;
}
.checks-link-bad {
  border-color: var(--sc-danger);
}
.checks-link-bad .sc-muted {
  color: var(--sc-danger);
}
.checks-link-text {
  display: grid;
  gap: 2px;
}
.checks-link-title {
  font-weight: 600;
}
.format-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-top: 14px;
}
</style>
