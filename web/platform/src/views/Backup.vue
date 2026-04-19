<template>
  <div>
    <div>
      <div>
        <div class="block1 wd12" style="max-width: 500px">
          <h1>{{ $t('backup.title') }}</h1>
          <div v-if="progress" id="progress">
            <el-row>
              <el-col :span="4"></el-col>
              <el-col :span="16" style="min-height: 30px " id="progress_summary">
                {{ progressSummary }}
              </el-col>
              <el-col :span="4"></el-col>
            </el-row>
            <el-row>
              <el-col :span="4"></el-col>
              <el-col :span="16">
                <el-progress :show-text="false" :percentage="progressPercentage"
                             :indeterminate="progressIndeterminate"/>
              </el-col>
              <el-col :span="4"></el-col>
            </el-row>
          </div>
          <div v-if="!progress">
            <div>
              <div style="padding-left: 10px; padding-top:10px; padding-bottom: 10px; display: inline-block">
                <el-select id="auto" v-model="auto" class="m-2" style="width: 140px; padding-right: 10px"
                           :placeholder="$t('backup.select')">
                  <el-option id="auto-no" :label="$t('backup.autoNo')" value="no"/>
                  <el-option id="auto-backup" :label="$t('backup.autoBackup')" value="backup"/>
                  <el-option id="auto-restore" :label="$t('backup.autoRestore')" value="restore"/>
                </el-select>
                <el-select id="auto-day" v-model="autoDay" class="m-2" style="width: 100px; padding-right: 10px"
                           :placeholder="$t('backup.select')" :disabled="auto === 'no'">
                  <el-option id="auto-day-every" :label="$t('backup.daily')" :value="0"/>
                  <el-option id="auto-day-monday" :label="$t('backup.mon')" :value="1"/>
                  <el-option :label="$t('backup.tue')" :value="2"/>
                  <el-option :label="$t('backup.wed')" :value="3"/>
                  <el-option :label="$t('backup.thu')" :value="4"/>
                  <el-option :label="$t('backup.fri')" :value="5"/>
                  <el-option :label="$t('backup.sat')" :value="6"/>
                  <el-option :label="$t('backup.sun')" :value="7"/>
                </el-select>
                <el-select id="auto-hour" v-model="autoHour" class="m-2" style="width: 90px; padding-right: 10px"
                           :placeholder="$t('backup.select')" :disabled="auto === 'no'">
                  <el-option v-for="hour in 24" :id="'auto-hour-' + (hour - 1)" :key="hour-1" :label="hour-1 + ':00'"
                             :value="(hour-1)"/>
                </el-select>
              </div>
              <div
                style="padding-top:10px; padding-bottom: 10px; padding-right: 10px; display: inline-block; float: right">
                <el-button id="save" type="success" @click="this.saveAuto">
                  {{ $t('backup.save') }}
                </el-button>
              </div>
            </div>
            <div class="row-no-gutters settingsblock">
              <el-table :data="filteredData"
                        style="width: 100%" table-layout="fixed">
                <el-table-column :label="$t('backup.name')" prop="file" :sortable="true"/>
                <el-table-column align="right" width="200px">
                  <template #header>
                    <el-input v-model="search" size="small" :placeholder="$t('backup.typeToSearch')"/>
                  </template>
                  <template #default="scope">
                    <el-button size="small" type="primary" @click="this.restoreConfirm(scope.row.file)">
                      {{ $t('backup.restore') }}
                    </el-button>
                    <el-button size="small" type="danger" @click="this.removeConfirm(scope.row.file)">
                      {{ $t('backup.delete') }}
                    </el-button>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <Dialog :visible="confirmationVisible" id="confirmation" @confirm="submit"
          @cancel="confirmationVisible = false">
    <template v-slot:title>
      <span v-if="action === 'restore'">{{ $t('backup.restoreTitle') }}</span>
      <span v-if="this.action === 'remove'">{{ $t('backup.removeTitle') }}</span>
    </template>
    <template v-slot:text>
      <div class="bodymod">
        <div class="btext">
          <span v-if="action === 'restore'">{{ $t('backup.restorePrompt') }}<br>{{ file }}?</span>
          <span v-if="action === 'remove'">{{ $t('backup.removePrompt') }}<br>{{ file }}?</span>
        </div>
      </div>
    </template>
  </Dialog>

</template>

<script>
import axios from 'axios'
import * as Common from '../js/common.js'
import Dialog from '../components/Dialog.vue'
import Notification from '../components/Notification.vue'

export default {
  name: 'Backup',
  data() {
    return {
      file: '',
      action: '',
      confirmationVisible: false,
      data: [],
      search: '',
      auto: 'no',
      autoDay: 0,
      autoHour: 0,
      progressSummary: '',
      progress: true,
      progressPercentage: 20,
      progressIndeterminate: true,
    }
  },
  computed: {
    filteredData() {
      return this.data.filter((v) => !this.search || v.file.toLowerCase().includes(this.search.toLowerCase()))
    }
  },
  components: {
    Dialog
  },
  mounted() {
    this.progressShow(this.$t('backup.loading'))
    this.status()
  },
  methods: {
    progressShow(summary) {
      this.progressSummary = summary
      this.progress = true
    },
    progressHide() {
      this.progressSummary = ''
      this.progress = false
    },
    removeConfirm(file) {
      this.file = file
      this.action = 'remove'
      this.confirmationVisible = true
    },
    restoreConfirm(file) {
      this.file = file
      this.action = 'restore'
      this.confirmationVisible = true
    },
    submit() {
      this.confirmationVisible = false
      switch (this.action) {
        case 'restore':
          this.restore()
          break
        case 'remove':
          this.remove()
          break
      }
    },
    remove() {
      axios.post('/rest/backup/remove', {file: this.file})
        .then(() => {
          this.reload()
        })
        .catch(this.showError)
    },
    showError(error) {
      this.progressHide()
      Notification.error(error)
    },
    restore() {
      this.progressShow(this.$t('backup.restoring', { file: this.file }))
      axios
        .post('/rest/backup/restore', {file: this.file})
        .then(() => {
          Common.runAfterJobIsComplete(
            setTimeout,
            () => {
              this.progressHide()
              this.reload()
            },
            Notification.error,
            Common.JOB_STATUS_URL,
            Common.JOB_STATUS_PREDICATE)
        })
        .catch(this.showError)
    },
    reload() {
      axios.get('/rest/backup/list')
        .then((response) => {
          if (response.data.data) {
            this.data = response.data.data
          } else {
            this.data = []
          }
          this.progressHide()
        })
        .catch(this.showError)
      axios.get('/rest/backup/auto')
        .then((response) => {
          this.auto = response.data.data.auto
          this.autoDay = response.data.data.day
          this.autoHour = response.data.data.hour
          this.progressHide()
        })
        .catch(this.showError)
    },
    status() {
      const error = this.$refs.error
      const onError = (err) => {
        this.progressHide()
        error.showAxios(err)
      }
      Common.runAfterJobIsComplete(
        setTimeout,
        () => {
          this.progressHide()
          this.reload()
        },
        onError,
        Common.JOB_STATUS_URL,
        (response) => {
          if (!response.data.data.status) {
            return false
          }
          this.progressSummary = response.data.data.status + ' (' + response.data.data.name + ')'
          return response.data.data.status !== 'Idle'
        },
      )
    },
    saveAuto() {
      this.progressShow(this.$t('backup.savingAuto'))
      axios.post('/rest/backup/auto',
        {auto: this.auto, day: this.autoDay, hour: this.autoHour})
        .then(() => {
          this.progressHide()
        })
        .catch(this.showError)
    }
  }
}
</script>
<style>
@import '../style/site.css';
@import 'material-icons/iconfont/material-icons.css';
</style>
