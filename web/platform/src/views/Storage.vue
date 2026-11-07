<template>
  <div class="sc-page">
    <div class="sc-card sc-card-wide" id="block1">
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
        <s-checkbox-group v-model="activeMultiDisks" size="large">
          <div v-for="(disk, index) in disks" :key="index" class="disk-card" :data-testid="'disk-card-' + index">
            <div class="disk-head">
              <s-checkbox :id="'disk_' + index" class="disk" size="large" border :label="disk.device">
                <span style="white-space: normal;">
                  {{ disk.name }} - {{ disk.size }}
                  <span v-if="disk.raid">({{ disk.raid }})</span>
                </span>
              </s-checkbox>
              <i v-if="disk.has_errors" class="material-icons-outlined disk-error-icon" :data-testid="'disk-has-errors-' + index">error</i>
            </div>
            <dl class="disk-details">
              <template v-if="disk.serial">
                <dt>{{ $t('storage.serial') }}</dt>
                <dd :data-testid="'disk-serial-' + index">{{ disk.serial }}</dd>
              </template>
              <template v-if="portOf(disk)">
                <dt>{{ $t('storage.port') }}</dt>
                <dd :data-testid="'disk-port-' + index">{{ portOf(disk) }}</dd>
              </template>
              <template v-if="disk.devid">
                <dt>{{ $t('storage.devid') }}</dt>
                <dd :data-testid="'disk-devid-' + index">{{ disk.devid }}</dd>
              </template>
              <template v-if="disk.raid">
                <dt>{{ $t('storage.profile') }}</dt>
                <dd :data-testid="'disk-profile-' + index">{{ $t('storage.profileValue', { data: disk.raid, metadata: disk.raid_metadata || '-' }) }}</dd>
              </template>
              <template v-if="disk.errors">
                <dt>{{ $t('storage.errorCounters') }}</dt>
                <dd :data-testid="'disk-errors-' + index" :class="{ 'disk-bad': disk.has_errors }">{{ $t('storage.counters', disk.errors) }}</dd>
              </template>
              <template v-if="healthOf(disk)">
                <dt>{{ $t('storage.smart') }}</dt>
                <dd :data-testid="'disk-smart-' + index" :class="{ 'disk-bad': smartBad(healthOf(disk).smart) }">{{ smartText(healthOf(disk).smart) }}</dd>
                <template v-if="healthOf(disk).smart.available">
                  <dt>{{ $t('storage.selfTest') }}</dt>
                  <dd :data-testid="'disk-selftest-' + index" class="disk-selftest">
                    <span>{{ selfTestText(healthOf(disk).smart.self_test) }}</span>
                    <button class="sc-btn sc-btn-ghost sc-btn-small" type="button" :data-testid="'disk-smart-test-' + index"
                            :disabled="health.busy || healthOf(disk).smart.self_test.running"
                            @click.prevent="startSmartTest(disk.device)">{{ $t('storage.smartTestStart') }}</button>
                  </dd>
                </template>
                <dt>{{ $t('storage.linkEvents') }}</dt>
                <dd :data-testid="'disk-link-' + index" :class="{ 'disk-bad': healthOf(disk).link.count > 0 }">{{ eventsText(healthOf(disk).link) }}</dd>
                <dt>{{ $t('storage.ioErrors') }}</dt>
                <dd :data-testid="'disk-io-' + index" :class="{ 'disk-bad': healthOf(disk).io.count > 0 }">{{ eventsText(healthOf(disk).io) }}</dd>
              </template>
            </dl>
          </div>
        </s-checkbox-group>

        <div v-if="health.scrub" class="disk-card" data-testid="array-section">
          <div class="disk-head">
            <span class="array-title">{{ $t('storage.array') }}</span>
          </div>
          <dl class="disk-details">
            <dt>{{ $t('storage.scrub') }}</dt>
            <dd data-testid="scrub-status" :class="{ 'disk-bad': !!health.scrub.errors }">{{ scrubText(health.scrub) }}</dd>
          </dl>
          <div class="check-actions">
            <button class="sc-btn sc-btn-ghost sc-btn-small" type="button" data-testid="scrub-start"
                    :disabled="health.busy || health.scrub.state === 'running'"
                    @click="startScrub">{{ $t('storage.scrubStart') }}</button>
            <button v-if="anyErrors" class="sc-btn sc-btn-ghost sc-btn-small" type="button" data-testid="errors-reset"
                    @click="resetVisible = true">{{ $t('storage.resetCounters') }}</button>
          </div>
          <div class="schedule-row">
            <span>{{ $t('storage.scheduleScrub') }}</span>
            <s-switch data-testid="schedule-scrub" :model-value="schedule.scrub" @update:model-value="value => saveSchedule('scrub', value)"/>
          </div>
        </div>

        <div v-if="anySmart" class="schedule-row" data-testid="schedule-smart-row">
          <span>{{ $t('storage.scheduleSmart') }}</span>
          <s-switch data-testid="schedule-smart" :model-value="schedule.smart" @update:model-value="value => saveSchedule('smart', value)"/>
        </div>
        <p v-if="health.busy" class="sc-muted" data-testid="checks-busy">{{ $t('storage.checksBusy') }}</p>
      </div>

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

  <Dialog :visible="resetVisible" id="reset-confirmation" data-testid="reset-confirmation" @confirm="resetErrors"
          @cancel="resetVisible = false">
    <template v-slot:title>{{ $t('storage.resetTitle') }}</template>
    <template v-slot:text>
      <div class="btext">{{ $t('storage.resetText') }}</div>
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
      resetVisible: false,
      health: { disks: [], scrub: null, schedule: { scrub: false, smart: false }, busy: false },
      schedule: { scrub: false, smart: false },
      healthTimer: undefined
    }
  },
  computed: {
    anyErrors () {
      return this.disks.some(d => d.has_errors)
    },
    anySmart () {
      return this.health.disks.some(h => h.smart.available)
    },
    checksRunning () {
      return (this.health.scrub && this.health.scrub.state === 'running') ||
        this.health.disks.some(h => h.smart.self_test && h.smart.self_test.running)
    }
  },
  mounted () {
    this.progressShow()
    this.reload()
  },
  unmounted () {
    clearTimeout(this.healthTimer)
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
      clearTimeout(this.healthTimer)
      return axios.get('/rest/storage/health')
        .then(resp => {
          this.health = resp.data.data
          this.schedule = { ...this.health.schedule }
          if (this.checksRunning) {
            this.healthTimer = setTimeout(this.reloadHealth, 10000)
          }
        })
        .catch(err => this.$refs.error.showAxios(err))
    },
    healthOf (disk) {
      return this.health.disks.find(h => h.device === disk.device)
    },
    portOf (disk) {
      const health = this.healthOf(disk)
      const port = health && health.port
      if (port && disk.hctl) {
        return port + ' (' + disk.hctl + ')'
      }
      return port || disk.hctl
    },
    value (number) {
      return number === undefined || number === null ? '-' : number
    },
    smartBad (smart) {
      return smart.available && (smart.passed === false || [smart.reallocated, smart.pending, smart.offline_uncorrectable, smart.crc_errors].some(v => v > 0))
    },
    smartText (smart) {
      if (!smart.available) {
        return this.$t('storage.smartNotAvailable')
      }
      return this.$t(smart.passed ? 'storage.smartPassed' : 'storage.smartFailed', {
        hours: this.value(smart.power_on_hours),
        temperature: this.value(smart.temperature),
        reallocated: this.value(smart.reallocated),
        pending: this.value(smart.pending),
        uncorrectable: this.value(smart.offline_uncorrectable),
        crc: this.value(smart.crc_errors)
      })
    },
    selfTestText (test) {
      if (test.running) {
        return this.$t('storage.selfTestRunning', { percent: test.remaining_percent })
      }
      if (!test.last_status) {
        return this.$t('storage.selfTestNone')
      }
      return [test.last_type, test.last_status].filter(v => v).join(': ')
    },
    eventsText (events) {
      if (events.count === 0) {
        return this.$t('storage.eventsNone')
      }
      return this.$t('storage.eventsCount', { count: events.count, last: events.last })
    },
    scrubText (scrub) {
      let text
      if (scrub.state === 'never') {
        text = this.$t('storage.scrubNever')
      } else if (scrub.state === 'running') {
        text = this.$t('storage.scrubRunning', { progress: scrub.progress || '0%' })
      } else if (scrub.state === 'finished') {
        text = this.$t('storage.scrubFinished', { started: scrub.started, duration: scrub.duration })
      } else {
        text = this.$t('storage.scrubStopped', { started: scrub.started })
      }
      if (scrub.errors) {
        return text + '. ' + this.$t('storage.scrubErrors', { errors: scrub.errors, corrected: scrub.corrected, uncorrectable: scrub.uncorrectable })
      }
      return text
    },
    checkAction (url, request) {
      return axios.post(url, request)
        .then(() => this.reloadHealth())
        .catch(err => this.$refs.error.showAxios(err))
    },
    startScrub () {
      return this.checkAction('/rest/storage/check/scrub')
    },
    startSmartTest (device) {
      return this.checkAction('/rest/storage/check/smart', { device: device })
    },
    resetErrors () {
      this.resetVisible = false
      return axios.post('/rest/storage/errors/reset')
        .then(() => this.reload())
        .catch(err => this.$refs.error.showAxios(err))
    },
    saveSchedule (check, value) {
      this.schedule = { ...this.schedule, [check]: value }
      return axios.post('/rest/storage/check/schedule', this.schedule)
        .then(resp => { this.schedule = { ...resp.data.data } })
        .catch(err => this.$refs.error.showAxios(err))
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
.disk-card {
  border: 1px solid var(--sc-border);
  border-radius: var(--sc-control-radius);
  padding: 12px 14px;
  margin-bottom: 12px;
}
.disk-head {
  display: flex;
  align-items: center;
  gap: 8px;
}
.disk-error-icon {
  color: var(--sc-danger);
  font-size: 20px !important;
}
.disk-details {
  display: grid;
  grid-template-columns: max-content 1fr;
  gap: 4px 14px;
  margin: 10px 0 0;
  font-size: 14px;
}
.disk-details dt {
  color: var(--sc-muted);
}
.disk-details dd {
  margin: 0;
  overflow-wrap: anywhere;
}
.disk-bad {
  color: var(--sc-danger);
}
.disk-selftest {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}
.sc-btn-small {
  padding: 4px 10px;
  font-size: 13px;
}
.array-title {
  font-weight: 600;
}
.check-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}
.schedule-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-top: 12px;
}
.format-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-top: 14px;
}
</style>
