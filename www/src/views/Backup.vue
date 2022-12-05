<template>
  <div>
    <div>
      <div>
        <div class="block1 wd12" style="max-width: 500px">
          <h1>Backup</h1>
          <div :style="{ visibility: visibility }">
            <div>
              <div style="padding-left: 10px; padding-top:10px; padding-bottom: 10px; display: inline-block">
                <span style="padding-right: 10px">Auto:</span>
                <el-select id="auto" v-model="auto" class="m-2" style="width: 100px; padding-right: 10px"
                           placeholder="Select">
                  <el-option id="auto-no" label="No" value="no"/>
                  <el-option id="auto-backup" label="Backup" value="backup"/>
                  <el-option id="auto-restore" label="Restore" value="restore"/>
                </el-select>
                <el-select id="auto-day" v-model="autoDay" class="m-2" style="width: 130px; padding-right: 10px"
                           placeholder="Select" :disabled="auto === 'no'">
                  <el-option id="auto-day-every" label="Every day" :value="0"/>
                  <el-option id="auto-day-monday" label="Monday" :value="1"/>
                  <el-option label="Tuesday" :value="2"/>
                  <el-option label="Wednesday" :value="3"/>
                  <el-option label="Thursday" :value="4"/>
                  <el-option label="Friday" :value="5"/>
                  <el-option label="Saturday" :value="6"/>
                  <el-option label="Sunday" :value="7"/>
                </el-select>
                <el-select id="auto-hour" v-model="autoHour" class="m-2" style="width: 90px; padding-right: 10px"
                           placeholder="Select" :disabled="auto === 'no'">
                  <el-option v-for="hour in 24" :id="'auto-hour-' + (hour - 1)" :key="hour-1" :label="hour-1 + ':00'"
                             :value="(hour-1)"/>
                </el-select>
              </div>
              <div
                style="padding-top:10px; padding-bottom: 10px; padding-right: 10px; display: inline-block; float: right">
                <el-button id="save" type="success" @click="this.saveAuto">
                  Save
                </el-button>
              </div>
            </div>
            <div class="row-no-gutters settingsblock">
              <el-table :data="filteredData"
                        style="width: 100%" table-layout="fixed">
                <el-table-column label="Name" prop="file" :sortable="true"/>
                <el-table-column align="right" width="200px">
                  <template #header>
                    <el-input v-model="search" size="small" placeholder="Type to search"/>
                  </template>
                  <template #default="scope">
                    <el-button size="small" type="primary" @click="this.restoreConfirm(scope.row.file)">
                      Restore
                    </el-button>
                    <el-button size="small" type="danger" @click="this.removeConfirm(scope.row.file)">
                      Delete
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

  <Confirmation :visible="confirmationVisible" id="confirmation" @confirm="submit"
                @cancel="confirmationVisible = false">
    <template v-slot:title>
      <span v-if="action === 'restore'">Restore</span>
      <span v-if="this.action === 'remove'">Remove</span>
    </template>
    <template v-slot:text>
      <div class="bodymod">
        <div class="btext">
          <span v-if="action === 'restore'">Do you want to restore:<br>{{ file }}?</span>
          <span v-if="action === 'remove'">Do you want to remove:<br>{{ file }}?</span>
        </div>
      </div>
    </template>
  </Confirmation>

</template>

<script>
import axios from 'axios'
import * as Common from '../js/common.js'
import Confirmation from '../components/Confirmation.vue'
import Notification from '../components/Notification.vue'
import { ElLoading } from 'element-plus'

export default {
  name: 'Backup',
  props: {
    checkUserSession: Function,
    activated: Boolean
  },
  data () {
    return {
      file: '',
      action: '',
      confirmationVisible: false,
      data: [],
      search: '',
      auto: 'no',
      autoDay: 0,
      autoHour: 0,
      visibility: 'hidden'
    }
  },
  computed: {
    filteredData () {
      return this.data.filter((v) => !this.search || v.file.toLowerCase().includes(this.search.toLowerCase()))
    }
  },
  components: {
    Confirmation
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
      this.visibility = 'visible'
      this.loading.close()
    },
    removeConfirm (file) {
      this.file = file
      this.action = 'remove'
      this.confirmationVisible = true
    },
    restoreConfirm (file) {
      this.file = file
      this.action = 'restore'
      this.confirmationVisible = true
    },
    submit () {
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
    remove () {
      axios.post('/rest/backup/remove', { file: this.file })
        .then(_ => {
          this.reload()
        })
        .catch(this.showError)
    },
    showError (error) {
      this.progressHide()
      Notification.error(error)
    },
    restore () {
      axios
        .post('/rest/backup/restore', { file: this.file })
        .then(_ => {
          Notification.info('Restoring an app from a backup')
          Common.runAfterJobIsComplete(
            setTimeout,
            () => {
              Notification.success('Backup restore has finished')
              this.reload()
            },
            Notification.error,
            Common.JOB_STATUS_URL,
            Common.JOB_STATUS_PREDICATE)
        })
        .catch(this.showError)
    },
    reload () {
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
    saveAuto () {
      this.progressShow()
      axios.post('/rest/backup/auto',
        { auto: this.auto, day: this.autoDay, hour: this.autoHour })
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
