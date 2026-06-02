<template>
  <div class="sc-page">
    <div class="sc-card" id="block1">
      <h1 class="sc-title">{{ $t('backup.title') }}</h1>
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
        <div class="backup-auto">
          <div class="backup-auto-controls">
            <el-select id="auto" v-model="auto" class="bk-sel bk-auto"
                       :placeholder="$t('backup.select')">
              <el-option id="auto-no" :label="$t('backup.autoNo')" value="no"/>
              <el-option id="auto-backup" :label="$t('backup.autoBackup')" value="backup"/>
              <el-option id="auto-restore" :label="$t('backup.autoRestore')" value="restore"/>
            </el-select>
            <el-select id="auto-day" v-model="autoDay" class="bk-sel bk-day"
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
            <el-select id="auto-hour" v-model="autoHour" class="bk-sel bk-hour"
                       :placeholder="$t('backup.select')" :disabled="auto === 'no'">
              <el-option v-for="hour in 24" :id="'auto-hour-' + (hour - 1)" :key="hour-1" :label="hour-1 + ':00'"
                         :value="(hour-1)"/>
            </el-select>
          </div>
          <el-button id="save" type="success" @click="this.saveAuto">
            {{ $t('backup.save') }}
          </el-button>
        </div>
        <div class="bk-search">
          <el-input v-model="search" size="small" :placeholder="$t('backup.typeToSearch')"/>
        </div>
        <div class="bk-list">
          <div v-for="row in filteredData" :key="row.file" class="bk-row" :data-testid="'backup-row-' + row.file">
            <span class="bk-file">{{ row.file }}</span>
            <div class="bk-actions">
              <el-button size="small" type="primary" @click="restoreConfirm(row.file)">
                {{ $t('backup.restore') }}
              </el-button>
              <el-button size="small" type="danger" @click="removeConfirm(row.file)">
                {{ $t('backup.delete') }}
              </el-button>
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
<style scoped>
.backup-auto {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 18px;
}
.backup-auto-controls {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}
.bk-auto { width: 150px; }
.bk-day { width: 112px; }
.bk-hour { width: 96px; }

.bk-search { max-width: 320px; margin-bottom: 12px; }
.bk-list { border-top: 1px solid #eef3f9; }
.bk-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 12px 6px;
  border-bottom: 1px solid #eef3f9;
}
.bk-file { word-break: break-all; text-align: left; }
.bk-actions { flex: 0 0 auto; display: flex; gap: 8px; }

@media (max-width: 600px) {
  .backup-auto {
    flex-direction: column;
    align-items: stretch;
    gap: 10px;
  }
  .backup-auto-controls {
    flex-wrap: nowrap;
    width: 100%;
    gap: 8px;
  }
  .bk-sel {
    flex: 1 1 0;
    width: auto;
    min-width: 0;
  }
  #save {
    align-self: flex-end;
  }

  .bk-search { max-width: none; }
  .bk-list { border-top: none; }
  .bk-row {
    flex-direction: column;
    align-items: stretch;
    gap: 10px;
    padding: 12px;
    border: 1px solid var(--sc-border);
    border-radius: 12px;
    background: var(--sc-field-bg);
    margin-bottom: 10px;
  }
  .bk-actions { justify-content: flex-end; }
}
</style>
