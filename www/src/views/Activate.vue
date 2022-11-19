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
                  <ul class="plan">
                    <li class="header">Premium</li>
                    <li class="description">
                      Syncloud will manage DNS records for your domain (like example.com)
                      <br><br>
                      Personal support for your device
                    </li>
                    <!--                    <li class="description">Personal support for your device</li>-->
                    <li>
                      <el-button id="btn_premium_domain" class="buttongreen" type="success"
                                 @click="selectPremiumDomain">Select
                      </el-button>
                    </li>
                  </ul>
                </div>
                <div class="columns">
                  <ul class="plan">
                    <li class="header">Free</li>
                    <li class="description">Syncloud will manage DNS records for [name].{{ redirect_domain }} domain
                    </li>
                    <li>
                      <el-button id="btn_free_domain" class="buttongreen" type="success" @click="selectFreeDomain">
                        Select
                      </el-button>
                    </li>
                  </ul>
                </div>
                <div class="columns">
                  <ul class="plan">
                    <li class="header">Custom</li>
                    <li class="description">You will manage DNS records for your domain (like example.com)</li>
                    <li>
                      <el-button id="btn_custom_domain" class="buttongreen" type="success"
                                 @click="selectCustomDomain">
                        Select
                      </el-button>
                    </li>
                  </ul>
                </div>
              </div>
              <div id="domain-account-part" class="content formblock" role="tabpanel"
                   aria-labelledby="domain-account-part-trigger">
                <div v-if="domainType === 'free'">
                  <div style="text-align: center">
                    <h2 style="display: inline-block">Domain Account</h2>
                    <button @click="showFreeAccountHelp" type=button class="control"
                            style="vertical-align: super; background:transparent;">
                      <i class='fa fa-question-circle fa-lg'></i>
                    </button>
                  </div>

                  <input :placeholder="redirect_domain + ' email'" class="emailinput" id="email"
                         type="text" v-model="redirectEmail">
                  <div class="alert alert-danger alert90" id="email_alert" style="display: none;"></div>
                  <input :placeholder="redirect_domain + ' password'" class="passinput"
                         id="redirect_password" type="password" v-model="redirectPassword">
                  <div class="alert alert-danger alert90" id="redirect_password_alert" style="display: none;"></div>
                  <div style=" display: flow-root">
                    <div style="padding-right:10px; float: right">
                      Do not have an account?
                      <a :href="'https://' + redirect_domain" class="btn btn-info" role="button"
                         style="line-height: 10px"
                         target="_blank">register</a>
                    </div>
                  </div>
                  <div style="text-align: center">
                    <h2 style="display: inline-block">Device Name</h2>
                    <button @click="showManagedDomainHelp" type=button class="control"
                            style="vertical-align: super; background:transparent;">
                      <i class='fa fa-question-circle fa-lg'></i>
                    </button>
                  </div>

                  <div id="domain">
                    <input placeholder="Name" class="domain" id="domain_input" type="text" v-model="domain">
                    <span>.{{ redirect_domain }}</span>
                  </div>
                  <div class="alert alert-danger alert90" id="domain_alert" style="display: none;"></div>

                </div>

                <div v-if=" domainType === 'custom' ">
                  <div style="text-align: center">
                    <h2 style="display: inline-block">Device Name</h2>
                    <button @click="showCustomDomainHelp" type=button class="control"
                            style="vertical-align: super; background:transparent;">
                      <i class='fa fa-question-circle fa-lg'></i>
                    </button>
                  </div>
                  <input placeholder="Top level domain like example.com"
                         class="domain" id="domain" type="text" style="width:100% !important;" v-model="domain">
                  <div class="alert alert-danger alert90" id="domain_alert" style="display: none;"></div>

                </div>

                <div v-if=" domainType === 'premium' ">
                  <div style="text-align: center">
                    <h2 style="display: inline-block">Syncloud Account</h2>
                    <button @click="showPremiumAccountHelp" type="button" class="control"
                            style="vertical-align: super; background:transparent;">
                      <i class='fa fa-question-circle fa-lg'></i>
                    </button>
                  </div>

                  <input :placeholder="redirect_domain + ' email'" class="emailinput" id="email"
                         type="text" v-model="redirectEmail">
                  <div class="alert alert-danger alert90" id="alert" style="display: none;"></div>
                  <input :placeholder="redirect_domain + ' password'" class="passinput"
                         id="redirect_password" type="password" v-model="redirectPassword">
                  <div class="alert alert-danger alert90" id="redirect_password_alert"
                       style="display: none;"></div>

                  <div style="text-align: center">
                    <h2 style="display: inline-block">Device Name</h2>
                    <button @click="showManagedDomainHelp" type=button class="control"
                            style="vertical-align: super; background:transparent;">
                      <i class='fa fa-question-circle fa-lg'></i>
                    </button>
                  </div>
                  <div id="domain">
                    <input placeholder="Top level domain like example.com"
                           class="domain" id="domain_premium" type="text" style="width:100% !important;"
                           v-model="domain">
                  </div>
                  <div class="alert alert-danger alert90" id="domain_alert" style="display: none;"></div>

                </div>

                <div style="padding: 10px; float: left;">
                  <el-button class="buttonblue" type="primary" @click="stepper.previous()">
                    Previous
                  </el-button>
                </div>
                <div style="padding: 10px; float: right;">
                  <el-button id="btn_next" type="primary" class="buttonblue" @click="selectDeviceName">
                    Next
                  </el-button>
                </div>
              </div>
              <div id="device-credentials-part" class="content formblock" role="tabpanel"
                   aria-labelledby="device-credentials-part-trigger">

                <div style="text-align: center">
                  <h2 style="display: inline-block">Device Credentials</h2>
                  <button @click="showDeviceCredentialHelp" type=button class="control"
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
                  <el-button class="buttonblue" type="primary" @click="stepper.previous()">
                    Previous
                  </el-button>
                </div>
                <div style="padding: 10px; float: right;">
                  <el-button id="btn_activate" class="buttonblue" type="primary" @click="activate"
                             data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Activating...">
                    Finish
                  </el-button>
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

  <Dialog ref="help_free_account">
    <template v-slot:title>Domain account</template>
    <template v-slot:text>
      Free Syncloud account name service (DNS) for device names at <b>{{ redirect_domain }}</b>.
      <br>
      You need to <a :href="'https://' + redirect_domain" class="btn btn-info" role="button"
                     target="_blank">register</a> an
      account to control one or more
      device names.
      <br>
      Domain account is also used for notifications about new releases.
      <br>
      It is only used to assign a dns name to IP of your device and update IP when it changes.
      Data transfer happens directly between your apps and device.
    </template>
  </Dialog>

  <Dialog ref="help_premium_account">
    <template v-slot:title>Domain account</template>
    <template v-slot:text>
      Premium Syncloud account name service (DNS) for personal domain name management (like example.com).
      <br>
      You need to <a :href="'https://' + redirect_domain" class="btn btn-info" role="button"
                     target="_blank">register</a> an
      account to control one or more
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
import Stepper from 'bs-stepper'
import Error from '../components/Error.vue'
import Dialog from '../components/Dialog.vue'
import { ElLoading } from 'element-plus'

export default {
  name: 'Activate',
  props: {
    checkUserSession: Function,
    activated: Boolean
  },
  components: {
    Dialog,
    Error
  },
  data () {
    return {
      domainType: 'free',
      loading: undefined,
      redirectEmail: '',
      redirectPassword: '',
      domain: '',
      redirect_domain: 'syncloud.it',
      deviceUsername: '',
      devicePassword: '',
      stepper: Stepper
    }
  },
  mounted () {
    this.progressShow()
    this.stepper = new Stepper(document.querySelector('.bs-stepper'))
    axios
      .get('/rest/redirect_info')
      .then(response => {
        this.redirect_domain = response.data.data.domain
        this.progressHide()
      })
      .catch(err => {
        this.progressHide()
        if (err.response.status !== 502) {
          this.$refs.error.showAxios(err)
        }
      })
  },
  methods: {
    progressShow () {
      this.loading = ElLoading.service({ lock: true, text: 'Loading', background: 'rgba(0, 0, 0, 0.7)' })
    },
    progressHide () {
      if (this.loading) {
        this.loading.close()
      }
    },
    activate (event) {
      event.preventDefault()
      this.progressShow()
      $('#form_activate .alert').remove()
      switch (this.domainType) {
        case 'premium':
          this.activatePremiumDomain()
          break
        case 'custom':
          this.activateCustomDomain()
          break
        default:
          this.activateFreeDomain()
      }
    },
    forceCertificateRecheck () {
      window.location = '/?t=' + (new Date()).getTime()
    },
    activateFreeDomain () {
      axios
        .post('/rest/activate/managed', {
          redirect_email: this.redirectEmail,
          redirect_password: this.redirectPassword,
          domain: this.fullDomain(),
          device_username: this.deviceUsername,
          device_password: this.devicePassword
        })
        .then(this.forceCertificateRecheck)
        .catch(err => {
          this.progressHide()
          this.$refs.error.showAxios(err)
        })
    },
    activatePremiumDomain () {
      axios
        .post('/rest/activate/managed', {
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
        .post('/rest/activate/custom', {
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
    showPremiumAccountHelp () {
      this.$refs.help_premium_account.show()
    },
    showCustomDomainHelp () {
      this.$refs.help_custom_domain.show()
    },
    showManagedDomainHelp () {
      this.$refs.help_managed_domain.show()
    },
    selectPremiumDomain () {
      this.domainType = 'premium'
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
    selectDeviceName () {
      if (this.domainType === 'custom') {
        this.stepper.next()
      } else {
        this.domainAvailability()
      }
    },
    fullDomain () {
      if (this.domainType === 'free') {
        return this.domain + '.' + this.redirect_domain
      }
      return this.domain
    },
    domainAvailability () {
      this.progressShow()
      axios
        .post('/rest/redirect/domain/availability',
          {
            email: this.redirectEmail,
            password: this.redirectPassword,
            domain: this.fullDomain()
          })
        .then(_ => {
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
@import 'material-icons/iconfont/material-icons.css';
@import 'bs-stepper/dist/css/bs-stepper.css';
@import 'font-awesome/css/font-awesome.css';

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

.plan {
  list-style-type: none;
  border: 1px solid #eee;
  margin: 0;
  padding: 0;
  -webkit-transition: 0.3s;
  transition: 0.3s;
}

.plan:hover {
  box-shadow: 0 8px 12px 0 rgba(0, 0, 0, 0.2)
}

.plan li {
  border-bottom: 1px solid #eee;
  padding: 20px;
  text-align: center;
}

.plan .header {
  background-color: #00aeef;
  font-size: 20px;
}

.plan .description {
  min-height: 125px;
  padding: 10px;
}

@media only screen and (max-width: 600px) {
  .columns {
    width: 100%;
  }

  .plan .description {
    min-height: 0;
  }
}
</style>
