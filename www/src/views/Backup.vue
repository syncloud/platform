<template>
  <div>
    <div>
      <div>
        <div class="block1 wd12">
          <h1>Backup</h1>
          <div style="padding-top:10px; padding-bottom: 10px">
            Auto:
            <el-select v-model="auto" class="m-2" style="width: 100px; padding-right: 10px" placeholder="Select">
              <el-option label="No" value="no"/>
              <el-option label="Backup" value="backup"/>
              <el-option label="Restore" value="restore"/>
            </el-select>
            <el-select v-model="autoDay" class="m-2" style="width: 120px; padding-right: 10px" placeholder="Select">
              <el-option label="Every day" value="0"/>
              <el-option label="Monday" value="1"/>
              <el-option label="Tuesday" value="2"/>
              <el-option label="Wednesday" value="3"/>
              <el-option label="Thursday" value="4"/>
              <el-option label="Friday" value="5"/>
              <el-option label="Saturday" value="6"/>
              <el-option label="Sunday" value="7"/>
            </el-select>
            <el-select v-model="autoTime" class="m-2" style="width: 90px; padding-right: 10px" placeholder="Select">
              <el-option v-for="hour in 24" :label="hour-1 + ':00'" :value="hour-1"/>
            </el-select>
            <el-button type="success" @click="this.save">
              Save
            </el-button>
          </div>
          <div class="row-no-gutters settingsblock">
            <el-table :data="filteredData" style="width: 100%" table-layout="fixed">
              <el-table-column label="File" prop="file"/>
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

  <Error ref="error"/>

</template>

<script>
import Error from '../components/Error.vue'
import toastr from 'toastr'
import axios from 'axios'
import * as Common from '../js/common.js'
import Confirmation from '../components/Confirmation.vue'

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
      autoDay: '0',
      autoTime: 1
    }
  },
  computed: {
    filteredData () {
      return this.data.filter((v) => !this.search || v.file.toLowerCase().includes(this.search.toLowerCase()))
    }
  },
  components: {
    Error,
    Confirmation
  },
  mounted () {
    this.reload()
  },
  methods: {
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
        .catch(err => this.$refs.error.showToast(err))
    },
    restore () {
      const that = this
      axios
        .post('/rest/backup/restore', { file: this.file })
        .then(_ => {
          toastr.info('Restoring an app from a backup')

          Common.runAfterJobIsComplete(
            setTimeout,
            () => {
              toastr.info('Backup restore has finished')
              this.reload()
            },
            err => that.$refs.error.showToast(err),
            Common.JOB_STATUS_URL,
            Common.JOB_STATUS_PREDICATE)
        })
        .catch(err => {
          this.$refs.error.showToast(err)
        })
    },
    reload () {
      axios.get('/rest/backup/list')
        .then((response) => {
          this.data = response.data.data
        })
    },
    save () {
      axios.post('/rest/backup/save',
        { auto: this.auto, day: this.autoDay, time: this.autoTime })
        .catch(err => {
          this.$refs.error.showToast(err)
        })
    }
  }
}
</script>
<style>
@import '../style/site.css';
@import 'material-icons/iconfont/material-icons.css';
@import 'toastr/build/toastr.css';
</style>
