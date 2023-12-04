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
            <div v-if="progress">
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

  <div id="app_action_confirmation" class="modal fade bs-are-use-sure" tabindex="-1" role="dialog"
       aria-labelledby="mySmallModalLabel">
    <div class="modal-dialog" role="document">
      <div class="modal-content">
        <div class="modal-header">
          <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span
            aria-hidden="true">&times;</span>
          </button>
          <h4 class="modal-title"><span id="confirm_caption">{{ action }}</span></h4>
        </div>
        <div class="modal-body">
          <div class="bodymod">
            <div class="btext">
              Are you sure?
            </div>

          </div>
          <div class="modal-footer">
            <button type="button" class="btn buttonlight bwidth smbutton" data-dismiss="modal">Close
            </button>
            <button type="button" id="btn_confirm" class="btn buttonlight bwidth smbutton"
                    data-dismiss="modal" @click="confirm">OK
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>

  <div id="backup_confirmation" class="modal fade bs-are-use-sure" tabindex="-1" role="dialog"
       aria-labelledby="mySmallModalLabel">
    <div class="modal-dialog" role="document">
      <div class="modal-content">
        <div class="modal-header">
          <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span
            aria-hidden="true">&times;</span>
          </button>
          <h4 class="modal-title">Backup</h4>
        </div>
        <div class="modal-body">
          <div class="bodymod">
            <div class="btext">
              This will backup app settings excluding files uploaded to the disk storage.<br>
              Later you can restore it from Settings - Backup<br>
              Are you sure?
            </div>

          </div>
          <div class="modal-footer">
            <button type="button" class="btn buttonlight bwidth smbutton" data-dismiss="modal">Close
            </button>
            <button type="button" id="btn_backup_confirm" class="btn buttonlight bwidth smbutton"
                    data-dismiss="modal" @click="backup">OK
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>

  <Error ref="error"/>

</template>

<script>
import axios from 'axios'
import $ from 'jquery'
import Error from '../components/Error.vue'
import 'bootstrap'
import * as Common from '../js/common.js'

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
      progressIndeterminate: true
    }
  },
  components: {
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
      $('#app_action_confirmation').modal('show')
    },
    upgrade () {
      this.action = 'Upgrade'
      this.actionUrl = '/rest/app/upgrade'
      $('#app_action_confirmation').modal('show')
    },
    remove () {
      this.action = 'Remove'
      this.actionUrl = '/rest/app/remove'
      $('#app_action_confirmation').modal('show')
    },
    backupConfirm () {
      $('#backup_confirmation').modal('show')
    },
    backup () {
      this.progressShow()

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
              this.loadApp,
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
