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
                    Custom domain
                  </a>
                </li>
                <li>
                  <a id="domain_type_managed" data-toggle="tab" href="#domain_managed" @click="domainType = 'managed'">
                    Managed domain
                  </a>
                </li>
              </ul>

            </div>

            <div class="tab-content">
              <div id="domain_syncloud" class="tab-pane fade in active">
                <div style="text-align: center">
                  <h2 style="display: inline-block">Syncloud Account</h2>
                  <button @click="showSyncloudAccountHelp" type=button
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
                <button @click="showCustomDomainHelp" type=button
                        style="vertical-align: super; background:transparent;">
                  <i class='fa fa-question-circle fa-lg'></i>
                </button>
                <input placeholder="Top level domain like example.com"
                       class="domain" id="full_domain" type="text" style="width:100% !important;" v-model="domain"
                       >
                <div class="alert alert-danger alert90" id="full_domain_alert" style="display: none;"></div>
              </div>
              <div id="domain_managed" class="tab-pane fade" style="text-align: center">
                <h2 style="display: inline-block">Device Name</h2>
                <button @click="showCustomDomainHelp" type=button
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
              <button @click="showDeviceCredentialHelp" type=button
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

  <Dialog ref="help_managed_domain">
    <template v-slot:title>Managed domain</template>
    <template v-slot:text>
      <div class="btext">If you have a domain you own, we can manage DNS records for you (requires Premium Account).
      </div>
      <span><br></span>
      <div class="btext" style="padding-left: 10px">A [Device IP] example.com</div>
      <div class="btext" style="padding-left: 10px">CNAME *.example.com example.com</div>

      <span><br></span>

      <div class="btext">If you do not have a DNS server and want to try on your LAN, edit your hosts file:</div>
      <span><br></span>
      <div class="btext" style="padding-left: 10px">[Device IP] example.com (device itself)</div>
      <div class="btext" style="padding-left: 10px">[Device IP] [app].example.com (line per app)</div>
    </template>
  </Dialog>

  <Dialog ref="help_custom_domain">
    <template v-slot:title>Custom domain</template>
    <template v-slot:text>
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
    </template>
  </Dialog>

  <Dialog ref="help_syncloud_account">
    <template v-slot:title>Syncloud account</template>
    <template v-slot:text>
      You can use free Syncloud name service (DNS) to get a device name at syncloud.it:
      <br>
      You need to register at <a href="https://syncloud.it" role="button">syncloud.it</a> to control one or more
      device names.
      <br>
      Syncloud account is also used for notifications about new releases.
      <br>
      It is only used to assign a dns name to IP of your device and update IP when it changes.
      Data transfer happens directly between your apps and device.
    </template>
  </Dialog>

  <Dialog ref="help_device_credential">
    <template v-slot:title>Device credentials</template>
    <template v-slot:text>
      Device credentials are used to access your device and all the apps (as admin user).
      They are stored on device and no one knows them. If you forget them you will need to reactivate your
      device.
    </template>
  </Dialog>

  <Error ref="error" :enable-logs="false"/>

</template>

<script>
import axios from 'axios'
import $ from 'jquery'
import 'bootstrap'
import 'bootstrap-switch'
import Error from '@/components/Error'
import 'gasparesganga-jquery-loading-overlay'
import Dialog from '@/components/Dialog'

export default {
  name: 'Activate',
  props: {
    onLogin: Function,
    onLogout: Function
  },
  components: {
    Dialog,
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
    },
    showDeviceCredentialHelp () {
      this.$refs.help_device_credential.show()
    },
    showSyncloudAccountHelp () {
      this.$refs.help_syncloud_account.show()
    },
    showCustomDomainHelp () {
      this.$refs.help_custom_domain.show()
    },
    showManagedDomainHelp () {
      this.$refs.help_custom_domain.show()
    }
  }
}
</script>
<style>
@import '../style/site.css';
@import '../style/material-icons.css';
</style>
