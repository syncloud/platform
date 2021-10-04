<template>

  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12" id="block1">
        <h1>Updates</h1>
        <div class="row-no-gutters settingsblock" id="block_updates">
          <div class="col2">
            <div class="setline">
              <button
                @click="check"
                class="buttongreen bwidth smbutton btn-lg" id="btn_check_updates"
                data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Checking..."
              >Check for updates
              </button>
            </div>
            <div class="setline">
              <span class="span">Note: upgrade System first if available before upgrading Installer.</span>
            </div>
            <div class="setline">
              <span class="span">System: </span>
              <span id="txt_platform_version" style="padding-right: 10px">{{ platformVersion }}</span>
              <button
                v-if="platformVersion !== platformVersionAvailable"
                id="btn_platform_upgrade"
                @click="upgradePlatform"
                class="buttongreen bwidth smbutton btn-lg"
                :data-loading-text="'<i class=\'fa fa-circle-o-notch fa-spin\'></i> Upgrading to ' + platformVersionAvailable + ' ...'"
              >
                Upgrade to {{ platformVersionAvailable }}
              </button>
            </div>
            <div class="setline">
              <span class="span">Installer: </span>
              <span id="txt_installer_version" style="padding-right: 10px">{{ installerVersion }}</span>
              <button
                v-if="installerVersion !== installerVersionAvailable"
                id="btn_installer_upgrade"
                @click="upgradeInstaller"
                class="buttongreen bwidth smbutton btn-lg"
                :data-loading-text="'<i class=\'fa fa-circle-o-notch fa-spin\'></i> Upgrading to ' + installerVersionAvailable + ' ...'"
              >
                Upgrade to {{ installerVersionAvailable }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <Error ref="error"/>

</template>

<script>
import $ from 'jquery'
import axios from 'axios'
import 'bootstrap'
import 'bootstrap-switch'
import * as Common from '../js/common.js'
import Error from '@/components/Error'
import 'gasparesganga-jquery-loading-overlay'

export default {
  name: 'Updates',
  components: {
    Error
  },
  props: {
    checkUserSession: Function,
    activated: Boolean
  },
  data () {
    return {
      platformVersion: undefined,
      platformVersionAvailable: undefined,
      installerVersion: undefined,
      installerVersionAvailable: undefined
    }
  },
  mounted () {
    this.progressShow()
    this.versions()
  },
  methods: {
    progressShow () {
      $('#block_updates').LoadingOverlay('show', { background: 'rgb(0,0,0,0)' })
    },
    progressHide () {
      $('#block_updates').LoadingOverlay('hide')
    },
    upgradePlatform () {
      this.progressShow()
      const that = this
      const onError = (err) => {
        that.$refs.error.showAxios(err)
        this.progressHide()
      }

      axios.post('/rest/upgrade', { app_id: 'platform' })
        .then((resp) => {
          Common.checkForServiceError(resp.data, () => {
            Common.runAfterJobIsComplete(
              setTimeout,
              this.versions,
              onError,
              Common.INSTALLER_STATUS_URL,
              Common.DEFAULT_STATUS_PREDICATE)
          }, onError)
        })
        .catch(onError)
    },
    upgradeInstaller () {
      this.progressShow()
      const that = this
      const onError = err => {
        that.$refs.error.showAxios(err)
        this.progressHide()
      }

      axios.post('/rest/installer/upgrade')
        .then((resp) => {
          Common.checkForServiceError(resp.data, () => {
            Common.runAfterJobIsComplete(
              setTimeout,
              this.versions,
              onError,
              Common.JOB_STATUS_URL,
              Common.JOB_STATUS_PREDICATE)
          }, onError)
        })
        .catch(onError)
    },
    findApp (appsData, appId) {
      for (const data of appsData) {
        if (data.app.id === appId) {
          return data
        }
      }
      return null
    },
    check () {
      this.progressShow()
      this.versions()
    },
    versions () {
      axios.get('/rest/settings/versions')
        .then((resp) => {
          const platformData = this.findApp(resp.data.data, 'platform')
          this.platformVersion = platformData.installed_version
          this.platformVersionAvailable = platformData.current_version
          const installerData = this.findApp(resp.data.data, 'installer')
          this.installerVersion = installerData.installed_version
          this.installerVersionAvailable = installerData.current_version
          this.progressHide()
        })
        .catch(err => {
          this.progressHide()
          this.$refs.error.showAxios(err)
        })
    }
  }
}
</script>
<style>
@import '../style/site.css';
@import '../style/material-icons.css';
</style>
