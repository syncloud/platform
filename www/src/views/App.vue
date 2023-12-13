<template>
  <div class="wrapper" id="block_app">
    <div class="content">
      <div class="block1 wd12">
        <div class="appblock" v-if="info !== undefined">
          <div>
            <div>
              <img :src="info.app.icon" class="appimg" alt="">
            </div>
            <div class="appinfo">
              <h1 id="app_name">{{ info.app.name }}</h1>
              <div v-if="info.installed_version !== null && !progress">
                <b>Version:</b> {{ info.installed_version }}<br>
              </div>
            </div>
          </div>
          <div>
            <div v-if="progress" id="progress">
            <el-row >
              <el-col :span="8"></el-col>
              <el-col :span="8" style="min-height: 30px " id="progress_summary" >
                {{ progressSummary }}
              </el-col>
              <el-col :span="8"></el-col>
            </el-row>
            <el-row >
              <el-col :span="8"></el-col>
              <el-col :span="8">
                <el-progress :show-text="false" :percentage="progressPercentage" :indeterminate="progressIndeterminate"/>
              </el-col>
              <el-col :span="8"></el-col>
            </el-row>
            </div>
            <div class="buttonblock" v-if="!progress">
              <button id="btn_open" :data-url="info.app.url" class="buttonblue bwidth smbutton"
                      @click="open"
                      v-if="info.installed_version !== null">
                Open
              </button>
              <button id="btn_install" class="buttonblue bwidth smbutton"
                      data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Installing..."
                      @click="install"
                      v-if="info.installed_version === null">
                Install v{{ info.current_version }}
              </button>
              <button id="btn_upgrade" class="buttongreen bwidth smbutton"
                      data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Upgrading..."
                      @click="upgrade"
                      v-if="info.installed_version !== null && info.installed_version !== info.current_version">
                Upgrade v{{ info.current_version }}
              </button>
              <button id="btn_remove" class="buttongrey bwidth smbutton"
                      data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Removing..."
                      @click="remove"
                      v-if="info.installed_version !== null">
                Remove
              </button>
              <button id="btn_backup" class="buttonblue bwidth smbutton"
                      data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Creating backup..."
                      @click="backupConfirm"
                      v-if="info.installed_version !== null">
                Backup
              </button>
            </div>
            <div class="btext">{{ info.app.description }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <Dialog :visible="appActionConfirmationVisible" id="app_confirmation" @confirm="confirm"
                @cancel="appActionConfirmationVisible = false">
    <template v-slot:title>
      <span id="confirm_caption">{{ action }}</span>
    </template>
    <template v-slot:text>
      <div class="bodymod">
        <div class="btext">
          Are you sure?
        </div>
      </div>
    </template>
  </Dialog>

  <Dialog :visible="backupConfirmationVisible" id="backup_confirmation" @confirm="backup"
                @cancel="backupConfirmationVisible = false">
    <template v-slot:title>
      <span id="confirm_caption">Backup</span>
    </template>
    <template v-slot:text>
      <div class="bodymod">
        <div class="btext">
          This will backup app settings excluding files uploaded to the disk storage.<br>
          Later you can restore it from Settings - Backup<br>
          Are you sure?
        </div>
      </div>
    </template>
  </Dialog>

  <Error ref="error"/>

</template>

<script>
import axios from 'axios'
import Error from '../components/Error.vue'
import * as Common from '../js/common.js'
import Dialog from '../components/Dialog.vue'

export default {
  name: 'App',
  props: {
    checkUserSession: Function,
    activated: Boolean
  },
  data () {
    return {
      info: undefined,
      appId: undefined,
      action: '',
      loading: undefined,
      progress: true,
      progressPercentage: 20,
      progressSummary: '',
      progressIndeterminate: true,
      appActionConfirmationVisible: false,
      backupConfirmationVisible: false
    }
  },
  components: {
    Dialog,
    Error
  },
  mounted () {
    this.progressShow()
    this.appId = this.$route.query.id
    this.loadApp()
      .then(() => this.status())
  },
  methods: {
    progressShow () {
      this.progress = true
    },
    progressHide () {
      this.progressSummary = ''
      this.progress = false
    },
    loadApp () {
      return axios
        .get('/rest/app', { params: { app_id: this.appId } })
        .then(resp => {
          this.info = resp.data.data
        })
        .catch(err => {
          this.$refs.error.showAxios(err)
        })
    },
    open () {
      window.location.href = this.info.app.url
    },
    install () {
      this.action = 'Install'
      this.actionUrl = '/rest/app/install'
      this.appActionConfirmationVisible = true
    },
    upgrade () {
      this.action = 'Upgrade'
      this.actionUrl = '/rest/app/upgrade'
      this.appActionConfirmationVisible = true
    },
    remove () {
      this.action = 'Remove'
      this.actionUrl = '/rest/app/remove'
      this.appActionConfirmationVisible = true
    },
    backupConfirm () {
      this.backupConfirmationVisible = true
    },
    backup () {
      this.backupConfirmationVisible = false
      this.progressShow()
      this.progressSummary = 'Creating a backup'

      const error = this.$refs.error
      const onError = (err) => {
        this.progressHide()
        error.showAxios(err)
      }
      axios.post('/rest/backup/create', { app: this.appId })
        .then((resp) => {
          Common.checkForServiceError(resp.data, () => {
            Common.runAfterJobIsComplete(
              setTimeout,
              // this.loadApp,
              () => { this.loadApp().then(() => this.progressHide()) },
              onError,
              Common.JOB_STATUS_URL,
              Common.JOB_STATUS_PREDICATE)
          }, onError)
        })
        .catch(onError)
    },
    status () {
      const error = this.$refs.error
      const onError = (err) => {
        this.progressHide()
        error.showAxios(err)
      }
      Common.runAfterJobIsComplete(
        setTimeout,
        () => { this.loadApp().then(() => this.progressHide()) },
        onError,
        Common.INSTALLER_STATUS_URL,
        (response) => {
          if (!response.data.data.progress) {
            return false
          }
          const progress = response.data.data.progress[this.appId]
          if (progress) {
            this.progressShow()
            this.progressSummary = progress.summary
            if (progress.indeterminate) {
              this.progressIndeterminate = true
              this.progressPercentage = 20
            } else {
              this.progressIndeterminate = false
              this.progressPercentage = progress.percentage
            }
            return true
          } else {
            return false
          }
        }
      )
    },
    confirm () {
      this.appActionConfirmationVisible = false
      this.progressShow()
      const error = this.$refs.error
      const onError = (err) => {
        this.progressHide()
        error.showAxios(err)
      }
      axios.post(this.actionUrl, { app_id: this.appId })
        .then(resp => {
          Common.checkForServiceError(resp.data, this.status, onError)
        })
        .catch(onError)
    }
  }
}
</script>
<style>
@import '../style/site.css';
@import 'material-icons/iconfont/material-icons.css';
</style>
