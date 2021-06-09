<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12" id="block_activate">
        <h1>Activate</h1>

        <div>
          <div class="bs-stepper">
            <div class="bs-stepper-header" role="tablist" style="max-width: 500px; margin: 0 auto">
              <div class="step" data-target="#domain-type-part">
                <button type="button" class="step-trigger" role="tab" aria-controls="domain-type-part"
                        id="domain-type-part-trigger">
                  <span class="bs-stepper-circle">1</span>
                  <span class="bs-stepper-label">Type</span>
                </button>
              </div>
              <div class="line"></div>
              <div class="step" data-target="#domain-account-part">
                <button type="button" class="step-trigger" role="tab" aria-controls="domain-account-part"
                        id="domain-account-part-trigger">
                  <span class="bs-stepper-circle">2</span>
                  <span class="bs-stepper-label">Name</span>
                </button>
              </div>
              <div class="line"></div>
              <div class="step" data-target="#device-credentials-part">
                <button type="button" class="step-trigger" role="tab" aria-controls="device-credentials-part"
                        id="device-credentials-part-trigger">
                  <span class="bs-stepper-circle">3</span>
                  <span class="bs-stepper-label">Credentials</span>
                </button>
              </div>
            </div>
            <div class="bs-stepper-content">
              <div id="domain-type-part" class="content" role="tabpanel" aria-labelledby="domain-type-part-trigger"
                   style="text-align: center; max-width: 800px; margin: 0 auto">
                <div class="columns">
                  <ul class="price">
                    <li class="header">Managed</li>
                    <li class="description">Syncloud will manage DNS records for your domain (like example.com), requires Premium account</li>
                    <li>
                      <button id="btn_managed_domain" class="buttongreen"
                              @click="selectManagedDomain">
                        Select
                      </button>
                    </li>
                  </ul>
                </div>
                <div class="columns">
                  <ul class="price">
                    <li class="header">Free</li>
                    <li class="description">Syncloud will manage DNS records for [name].syncloud.it domain</li>
                    <li>
                      <button id="btn_free_domain" class="buttongreen"
                              @click="selectFreeDomain">
                        Select
                      </button>
                    </li>
                  </ul>
                </div>
                <div class="columns">
                  <ul class="price">
                    <li class="header">Custom</li>
                    <li class="description">You will manage DNS records for your domain (like example.com)</li>
                    <li>
                      <button id="btn_custom_domain" class="buttongreen"
                              @click="selectCustomDomain">
                        Select
                      </button>
                    </li>
                  </ul>
                </div>
              </div>
              <div id="domain-account-part" class="content formblock" role="tabpanel"
                   aria-labelledby="domain-account-part-trigger">
                <div v-if="domainType === 'free'">
                  <div style="text-align: center">
                    <h2 style="display: inline-block">Domain Account</h2>
                    <button @click="showFreeAccountHelp" type=button
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
                  <div style=" display: flow-root">
                    <div style="padding-right:10px; float: right">
                      Do not have an account?
                      <a href="https://syncloud.it" class="btn btn-info" role="button" style="line-height: 10px" target="_blank">register</a>
                    </div>
                  </div>
                  <div style="text-align: center">
                    <h2 style="display: inline-block">Device Name</h2>
                    <button @click="showCustomDomainHelp" type=button
                            style="vertical-align: super; background:transparent;">
                      <i class='fa fa-question-circle fa-lg'></i>
                    </button>
                  </div>

                  <div id="domain">
                    <input placeholder="Name" class="domain" id="user_domain_input" type="text" v-model="domain">
                    <span>.syncloud.it</span>
                  </div>
                  <div class="alert alert-danger alert90" id="domain_alert" style="display: none;"></div>

                </div>

                <div v-if=" domainType === 'custom' ">
                  <div style="text-align: center">
                    <h2 style="display: inline-block">Device Name</h2>
                    <button @click="showCustomDomainHelp" type=button
                            style="vertical-align: super; background:transparent;">
                      <i class='fa fa-question-circle fa-lg'></i>
                    </button>
                  </div>
                  <input placeholder="Top level domain like example.com"
                         class="domain" id="domain" type="text" style="width:100% !important;" v-model="domain">
                  <div class="alert alert-danger alert90" id="domain_alert" style="display: none;"></div>

                </div>

                <div v-if=" domainType === 'managed' ">
                  <div style="text-align: center">
                    <h2 style="display: inline-block">Syncloud Account</h2>
                    <button @click="showManagedAccountHelp" type=button
                            style="vertical-align: super; background:transparent;">
                      <i class='fa fa-question-circle fa-lg'></i>
                    </button>
                  </div>

                  <input placeholder="syncloud.it email" class="emailinput" id="email"
                         type="text" v-model="redirectEmail">
                  <div class="alert alert-danger alert90" id="alert" style="display: none;"></div>
                  <input placeholder="syncloud.it password" class="passinput"
                         id="redirect_password" type="password" v-model="redirectPassword">
                  <div class="alert alert-danger alert90" id="redirect_password_alert"
                       style="display: none;"></div>

                  <div style="text-align: center">
                    <h2 style="display: inline-block">Device Name</h2>
                    <button @click="showManagedDomainHelp" type=button
                            style="vertical-align: super; background:transparent;">
                      <i class='fa fa-question-circle fa-lg'></i>
                    </button>
                  </div>
                  <input placeholder="Top level domain like example.com"
                         class="domain" id="domain" type="text" style="width:100% !important;" v-model="domain">
                  <div class="alert alert-danger alert90" id="domain_alert" style="display: none;"></div>

                </div>

                <div style="padding: 10px; float: left;">
                  <button class="buttonblue" @click="stepper.previous()">
                    Previous
                  </button>
                </div>
                <div style="padding: 10px; float: right;">
                  <button id="btn_next" class="buttonblue" @click="domainAvailability">
                    Next
                  </button>
                </div>
              </div>
              <div id="device-credentials-part" class="content formblock" role="tabpanel"
                   aria-labelledby="device-credentials-part-trigger">

                <div style="text-align: center">
                  <h2 style="display: inline-block">Device Credentials</h2>
                  <button @click="showDeviceCredentialHelp" type=button
                          style="vertical-align: super; background:transparent;">
                    <i class='fa fa-question-circle fa-lg'></i>
                  </button>
                </div>

                <input placeholder="Login" class="nameinput" id="device_username" type="text" v-model="deviceUsername">
                <div class="alert alert-danger alert90" id="device_username_alert" style="display: none;"></div>
                <input placeholder="Password" class="passinput" id="device_password" type="password"
                       v-model="devicePassword">
                <div class="alert alert-danger alert90" id="device_password_alert" style="display: none;"></div>

                <div style="padding: 10px; float: left;">
                  <button class="buttonblue" @click="stepper.previous()">
                    Previous
                  </button>
                </div>
                <div style="padding: 10px; float: right;">
                  <button id="btn_activate" class="buttonblue" @click="activate"
                          data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Activating...">
                    Finish
                  </button>
                </div>
              </div>

            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <Dialog ref="help_managed_domain">
    <template v-slot:title>Managed domain</template>
    <template v-slot:text>
      <div class="btext">Syncloud will manage DNS records for your personal domain name
      </div>
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
    <template v-slot:title>Domain account</template>
    <template v-slot:text>
      Free Syncloud account name service (DNS) for device names at <b>syncloud.it</b>.
      <br>
      You need to <a href="https://syncloud.it" class="btn btn-info" role="button" target="_blank">register</a> an account to control one or more
      device names.
      <br>
      Domain account is also used for notifications about new releases.
      <br>
      It is only used to assign a dns name to IP of your device and update IP when it changes.
      Data transfer happens directly between your apps and device.
    </template>
  </Dialog>

  <Dialog ref="help_managed_account">
    <template v-slot:title>Domain account</template>
    <template v-slot:text>
      Premium Syncloud account name service (DNS) for personal domain name management (like example.com).
      <br>
      You need to <a href="https://syncloud.it" class="btn btn-info" role="button" target="_blank">register</a> an account to control one or more
      device names. Then request a premium plan in your Account services.
      <br>
      Domain account is also used for notifications about new releases.
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
import Stepper from 'bs-stepper'
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
      domainType: 'free',
      loading: false,
      redirectEmail: '',
      redirectPassword: '',
      domain: '',
      deviceUsername: '',
      devicePassword: '',
      stepper: Stepper
    }
  },
  mounted () {
    this.stepper = new Stepper(document.querySelector('.bs-stepper'))
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
      if (this.domainType === 'free') {
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
        .post('/rest/activate/free', {
          redirect_email: this.redirectEmail,
          redirect_password: this.redirectPassword,
          domain: this.domain,
          device_username: this.deviceUsername,
          device_password: this.devicePassword
        })
        .then(this.forceCertificateRecheck)
        .catch(err => {
          this.progressHide()
          this.$refs.error.showAxios(err)
        })
    },
    activateManagedDomain () {
      axios
        .post('/rest/activate_managed_domain', {
          redirect_email: this.redirectEmail,
          redirect_password: this.redirectPassword,
          domain: this.domain,
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
          domain: this.domain,
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
    showFreeAccountHelp () {
      this.$refs.help_free_account.show()
    },
    showManagedAccountHelp () {
      this.$refs.help_managed_account.show()
    },
    showCustomDomainHelp () {
      this.$refs.help_custom_domain.show()
    },
    showManagedDomainHelp () {
      this.$refs.help_managed_domain.show()
    },
    selectManagedDomain () {
      this.domainType = 'managed'
      this.stepper.next()
    },
    selectCustomDomain () {
      this.domainType = 'custom'
      this.stepper.next()
    },
    selectFreeDomain () {
      this.domainType = 'free'
      this.stepper.next()
    },
    domainAvailability () {
      this.progressShow()
      axios
        .post('/rest/redirect/domain/availability',
          {
            email: this.redirectEmail,
            password: this.redirectPassword,
            user_domain: this.domain
          })
        .then(response => {
          this.stepper.next()
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
@import '~bs-stepper/dist/css/bs-stepper.css';

.active .bs-stepper-circle {
  background-color: #02a0dc;
}

* {
  box-sizing: border-box;
}

.columns {
  float: left;
  width: 33.3%;
  padding: 8px;
}

.price {
  list-style-type: none;
  border: 1px solid #eee;
  margin: 0;
  padding: 0;
  -webkit-transition: 0.3s;
  transition: 0.3s;
}

.price:hover {
  box-shadow: 0 8px 12px 0 rgba(0,0,0,0.2)
}

.price li {
  border-bottom: 1px solid #eee;
  padding: 20px;
  text-align: center;
}

.price .header {
  background-color: #00aeef;
  font-size: 20px;
}

.price .description {
  min-height: 125px;
}

@media only screen and (max-width: 600px) {
  .columns {
    width: 100%;
  }
  .price .description {
    min-height: 0px;
  }
}
</style>
