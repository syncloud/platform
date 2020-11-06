<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12" id="block1">
        <h1>Activate</h1>
        <div class="formblock" id="block_activate">
          <form id="form_activate" @submit="activate">
            <div class="centered-pills">

              <ul class="nav nav-pills">
                <li class="active">
                  <a id="domain_type_syncloud" data-toggle="tab" href="#domain_syncloud"
                     @click="domainType = 'syncloud'">[name].syncloud.it</a>
                </li>
                <li>
                  <a id="domain_type_custom" data-toggle="tab" href="#domain_custom" @click="domainType = 'custom'">
                    My own domain
                  </a>
                </li>
              </ul>

            </div>

            <div class="tab-content">
              <div id="domain_syncloud" class="tab-pane fade in active">
                <div style="text-align: center">
                  <h2 style="display: inline-block">Syncloud Account</h2>
                  <button data-toggle="modal" data-target="#help_syncloud_account" type=button
                          style="vertical-align: super; background:transparent;">
                    <i class='fa fa-question-circle fa-lg'></i>
                  </button>
                </div>
                <input placeholder="syncloud.it email" class="emailinput" id="email"
                       type="text" v-model="redirectEmail">
                <div class="alert alert-danger alert90" id="email_alert" style="display: none;"></div>
                <input placeholder="syncloud.it password" class="passinput"
                       id="redirect_password" type="password" v-model="redirectPassword">
                <div class="alert alert-danger alert90" id="redirect_password_alert" style="display: none;"></div>
                <input placeholder="Name" class="domain" id="user_domain" type="text" v-model="domain"
                       ><span>.syncloud.it</span>
                <div class="alert alert-danger alert90" id="user_domain_alert" style="display: none;"></div>

              </div>
              <div id="domain_custom" class="tab-pane fade" style="text-align: center">
                <h2 style="display: inline-block">Device Name</h2>
                <button data-toggle="modal" data-target="#help_custom_domain" type=button
                        style="vertical-align: super; background:transparent;">
                  <i class='fa fa-question-circle fa-lg'></i>
                </button>
                <input placeholder="Top level domain like example.com"
                       class="domain" id="full_domain" type="text" style="width:100% !important;" v-model="domain"
                       >
                <div class="alert alert-danger alert90" id="full_domain_alert" style="display: none;"></div>
              </div>
            </div>

            <div style="text-align: center">
              <h2 style="display: inline-block">Device Credentials</h2>
              <button data-toggle="modal" data-target="#help_device_credential" type=button
                      style="vertical-align: super; background:transparent;">
                <i class='fa fa-question-circle fa-lg'></i>
              </button>
            </div>

            <input placeholder="Login" class="nameinput" id="device_username" type="text" v-model="deviceUsername"
                   >
            <div class="alert alert-danger alert90" id="device_username_alert" style="display: none;"></div>
            <input placeholder="Password" class="passinput" id="device_password" type="password"
                   v-model="devicePassword"
                   >
            <div class="alert alert-danger alert90" id="device_password_alert" style="display: none;"></div>
            <button id="btn_activate" class="submit buttonblue" type="submit"
                    data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Activating...">Activate
            </button>
          </form>
        </div>
      </div>
    </div>
  </div>

  <div id="help_custom_domain" class="modal fade bs-are-use-sure" tabindex="-1" role="dialog"
       aria-labelledby="mySmallModalLabel">
    <div class="modal-dialog" role="document">
      <div class="modal-content">
        <div class="modal-header">
          <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span
            aria-hidden="true">&times;</span></button>
          <h4 class="modal-title">Custom domain</h4>
        </div>
        <div class="modal-body">
          <div class="bodymod">
            <div class="btext">If you have a domain you own, make sure you have correct records on your DNS server:
            </div>
            <span><br></span>
            <div class="btext" style="padding-left: 10px">A [Device IP] example.com</div>
            <div class="btext" style="padding-left: 10px">CNAME *.example.com example.com</div>

            <span><br></span>

            <div class="btext">If you do not have a DNS server and want to try on your LAN, edit your hosts file:</div>
            <span><br></span>
            <div class="btext" style="padding-left: 10px">[Device IP] example.com (device itself)</div>
            <div class="btext" style="padding-left: 10px">[Device IP] [app].example.com (line per app)</div>

          </div>
          <div class="modal-footer">
            <button type="button" class="btn buttonlight bwidth smbutton" data-dismiss="modal">Close</button>
          </div>
        </div>
      </div>
    </div>
  </div>

  <div id="help_syncloud_account" class="modal fade bs-are-use-sure" tabindex="-1" role="dialog"
       aria-labelledby="mySmallModalLabel">
    <div class="modal-dialog" role="document">
      <div class="modal-content">
        <div class="modal-header">
          <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span
            aria-hidden="true">&times;</span></button>
          <h4 class="modal-title">Syncloud account</h4>
        </div>
        <div class="modal-body">
          <div class="bodymod">
            <div class="btext">
              You can use free Syncloud name service (DNS) to get a device name at syncloud.it:
              <br>
              You need to register at <a href="https://syncloud.it" role="button">syncloud.it</a> to control one or more
              device names.
              <br>
              Syncloud account is also used for notifications about new releases.
              <br>
              It is only used to assign a dns name to IP of your device and update IP when it changes.
              Data transfer happens directly between your apps and device.
            </div>

          </div>
          <div class="modal-footer">
            <button type="button" class="btn buttonlight bwidth smbutton" data-dismiss="modal">Close</button>
          </div>
        </div>
      </div>
    </div>
  </div>

  <div id="help_device_credential" class="modal fade bs-are-use-sure" tabindex="-1" role="dialog"
       aria-labelledby="mySmallModalLabel">
    <div class="modal-dialog" role="document">
      <div class="modal-content">
        <div class="modal-header">
          <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span
            aria-hidden="true">&times;</span></button>
          <h4 class="modal-title">Device credentials</h4>
        </div>
        <div class="modal-body">
          <div class="bodymod">
            <div class="btext">
              Device credentials are used to access your device and all the apps (as admin user).
              They are stored on device and no one knows them. If you forget them you will need to reactivate your
              device.
            </div>

          </div>
          <div class="modal-footer">
            <button type="button" class="btn buttonlight bwidth smbutton" data-dismiss="modal">Close</button>
          </div>
        </div>
      </div>
    </div>
  </div>

  <Error ref="error" :enable-logs="false"/>

</template>

<script>
import axios from 'axios'
import $ from 'jquery'
import 'bootstrap'
import 'bootstrap-switch'
import Error from '@/components/Error'
import 'gasparesganga-jquery-loading-overlay'

export default {
  name: 'Activate',
  props: {
    onLogin: Function,
    onLogout: Function
  },
  components: {
    Error
  },
  data () {
    return {
      domainType: 'syncloud',
      loading: false,
      redirectEmail: '',
      redirectPassword: '',
      domain: '',
      deviceUsername: '',
      devicePassword: ''
    }
  },
  methods: {
    progressShow () {
      $('#block_activate').LoadingOverlay('show', { background: 'rgb(0,0,0,0)' })
    },
    progressHide () {
      $('#block_activate').LoadingOverlay('hide')
    },
    activate (event) {
      event.preventDefault()
      this.progressShow()
      $('#form_activate .alert').remove()
      if (this.domainType === 'syncloud') {
        this.activateFreeDomain()
      } else {
        this.activateCustomDomain()
      }
    },
    forceCertificateRecheck () {
      window.location.reload(true)
    },
    activateFreeDomain () {
      axios
        .post('/rest/activate', {
          redirect_email: this.redirectEmail,
          redirect_password: this.redirectPassword,
          user_domain: this.domain,
          device_username: this.deviceUsername,
          device_password: this.devicePassword
        })
        .then(this.forceCertificateRecheck)
        .catch(err => {
          this.progressHide()
          this.$refs.error.showAxios(err)
        })
    },
    activateCustomDomain () {
      axios
        .post('/rest/activate_custom_domain', {
          full_domain: this.domain,
          device_username: this.deviceUsername,
          device_password: this.devicePassword
        })
        .then(this.forceCertificateRecheck)
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
