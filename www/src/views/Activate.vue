<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12" id="block_activate">
        <h1>Activate</h1>

        <div>
          <el-steps :active="step" simple finish-status="success" style="max-width: 500px; margin: 0 auto">
            <el-step title="Type" />
            <el-step title="Name" />
            <el-step title="User" />
          </el-steps>

          <div v-if="step === 0" id="domain-type-part"  >
            <div style="text-align: center; max-width: 800px; margin: 0 auto; display: flex; flex-wrap: wrap; justify-content: center;">
              <div class="columns">
                <ul class="plan">
                  <li class="description">
                    Syncloud manages<br>your domain (like example.com)<br>
                    Subscription is required today.
                  </li>
                  <li>
                    <el-button id="btn_premium_domain" style="width: 100%;height: 40px;" type="primary" @click="selectPremiumDomain">
                      Your name
                    </el-button>
                  </li>
                </ul>
              </div>
              <div class="columns">
                <ul class="plan">
                  <li class="description">
                    Syncloud manages<br>[name].{{ redirect_domain }} domain<br>
                    Subscription is required in 30 days.
                  </li>
                  <li>
                    <el-button id="btn_free_domain" style="width: 100%;height: 40px;" type="primary" @click="selectFreeDomain">
                      Our name
                    </el-button>
                  </li>
                </ul>
              </div>
            </div>
          </div>

          <div v-if="step === 1" id="domain-account-part">
            <div style="text-align: center; max-width: 400px; margin: 0 auto;">
              <div v-if="domainType === 'free'">
                <div style="text-align: center">
                  <h2 style="display: inline-block">Syncloud Account</h2>
                  <button @click="showFreeAccountHelp" type=button class="control"
                          style="vertical-align: super; background:transparent;">
                    <i class='fa fa-question-circle fa-lg'></i>
                  </button>
                </div>

                <input :placeholder="redirect_domain + ' email'" class="input" id="email"
                       type="text" v-model="redirectEmail">
                <div id="email_alert" class="alert alert-danger alert90" v-show="redirectEmailAlertVisible ">{{ redirectEmailAlert }}</div>

                <Password
                  id="redirect_password"
                  v-model="redirectPassword"
                  :placeholder="redirect_domain + ' password'"
                  :show-error="redirectPasswordAlertVisible"
                  :error="redirectPasswordAlert"
                />

                <div style=" display: flow-root">
                  <div style="padding-right:10px; float: right">
                    Do not have an account?
                    <a :href="'https://' + redirect_domain" class="register"
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

                <div id="domain" style="text-align: left">
                  <input placeholder="Name" class="domain input" id="domain_input" type="text" v-model="domain">
                  <span>.{{ redirect_domain }}</span>
                </div>
                <div id="domain_alert"  class="alert alert-danger alert90" v-show="domainAlertVisible" >{{ domainAlert }}</div>

              </div>

              <div v-if=" domainType === 'premium' ">
                <div style="text-align: center">
                  <h2 style="display: inline-block">Syncloud Account</h2>
                  <button @click="showPremiumAccountHelp" type="button" class="control"
                          style="vertical-align: super; background:transparent;">
                    <i class='fa fa-question-circle fa-lg'></i>
                  </button>
                </div>

                <input :placeholder="redirect_domain + ' email'" class="input" id="email"
                       type="text" v-model="redirectEmail">
                <div class="alert alert-danger alert90" id="alert" style="display: none;"></div>

                <Password
                  id="redirect_password"
                  v-model="redirectPassword"
                  :placeholder="redirect_domain + ' password'"
                  :show-error="redirectPasswordAlertVisible"
                  :error="redirectPasswordAlert"
                />

                <div style="text-align: center">
                  <h2 style="display: inline-block">Device Name</h2>
                  <button @click="showManagedDomainHelp" type=button class="control"
                          style="vertical-align: super; background:transparent;">
                    <i class='fa fa-question-circle fa-lg'></i>
                  </button>
                </div>
                <div id="domain">
                  <input placeholder="Top level domain like example.com"
                         class="domain input" id="domain_premium" type="text" style="width:100% !important;"
                         v-model="domain">
                </div>
                <div id="domain_alert" class="alert alert-danger alert90" v-show="domainAlertVisible" >{{ domainAlert }}</div>
              </div>

              <div style="padding: 10px; float: left;">
                <el-button type="primary" @click="step--">
                  Previous
                </el-button>
              </div>
              <div style="padding: 10px; float: right;">
                <el-button id="btn_next" type="primary"  @click="selectDeviceName">
                  Next
                </el-button>
              </div>
            </div>
          </div>

          <div v-if="step === 2" id="device-credentials-part">
            <div style="text-align: center; max-width: 400px; margin: 0 auto;">

              <div style="text-align: center">
                <h2 style="display: inline-block">Device Credentials</h2>
                <button @click="showDeviceCredentialHelp" type=button class="control"
                        style="vertical-align: super; background:transparent;">
                  <i class='fa fa-question-circle fa-lg'></i>
                </button>
              </div>

              <input placeholder="Login" id="device_username" type="text" v-model="deviceUsername"
                     v-on:keyup.enter="activate"
                     class="input">
              <div class="alert alert-danger alert90" v-show="deviceUsernameAlertVisible">{{
                  deviceUsernameAlert
                }}
              </div>

              <Password
                id="device_password"
                v-model="devicePassword"
                placeholder="Password"
                :show-error="devicePasswordAlertVisible"
                :error="devicePasswordAlert"
                @trigger="activate"
              />

              <Password
                id="device_password_confirm"
                v-model="devicePasswordConfirm"
                placeholder="Confirm your password"
                :show-error="devicePassword !== devicePasswordConfirm"
                error="Passwords do not match"
                @trigger="activate"
              />

              <div style="padding: 10px; float: left;">
                <el-button type="primary" @click="step--">
                  Previous
                </el-button>
              </div>
              <div style="padding: 10px; float: right;">
                <el-button id="btn_activate" type="primary" @click="activate" :disabled="!validDeviceCredentials()">
                Finish
                </el-button>
              </div>
            </div>
          </div>

        </div>
      </div>
    </div>
  </div>

  <Dialog :visible="helpManagedDomainVisible" @cancel="helpManagedDomainVisible = false" :confirm-enabled="false" cancel-text="Close">
    <template v-slot:title>Managed domain</template>
    <template v-slot:text>
      <div class="btext">Syncloud will manage DNS records for your personal domain name
      </div>
    </template>
  </Dialog>

  <Dialog :visible="helpFreeAccountVisible" @cancel="helpFreeAccountVisible = false" :confirm-enabled="false" cancel-text="Close">
    <template v-slot:title>Syncloud account</template>
    <template v-slot:text>
      You need to <a :href="'https://' + redirect_domain" target="_blank" class="register">register</a> an
      account and have a valid Subscription no more than 30 days after the registration.
      <br>
      It is only used to maintain dns records of your device.
      Data transfer happens directly between your apps and device.
    </template>
  </Dialog>

  <Dialog :visible="helpPremiumAccountVisible" @cancel="helpPremiumAccountVisible = false" :confirm-enabled="false" cancel-text="Close">
    <template v-slot:title>Syncloud account</template>
    <template v-slot:text>
      You need to <a :href="'https://' + redirect_domain" target="_blank" class="register">register</a> an
      account and have a valid Subscription.
      <br>
      It is only used to maintain dns records of your device.
      Data transfer happens directly between your apps and device.
    </template>
  </Dialog>

  <Dialog :visible="helpDeviceCredentialVisible" @cancel="helpDeviceCredentialVisible = false" :confirm-enabled="false" cancel-text="Close">
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
import Error from '../components/Error.vue'
import Dialog from '../components/Dialog.vue'
import Password from '../components/Password.vue'
import { ElLoading } from 'element-plus'

export default {
  name: 'Activate',
  props: {
    checkUserSession: Function,
    activated: Boolean
  },
  components: {
    Dialog,
    Error,
    Password
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
      devicePasswordConfirm: '',
      helpManagedDomainVisible: false,
      helpFreeAccountVisible: false,
      helpPremiumAccountVisible: false,
      helpDeviceCredentialVisible: false,
      deviceUsernameAlertVisible: false,
      deviceUsernameAlert: '',
      devicePasswordAlertVisible: false,
      devicePasswordAlert: '',
      redirectPasswordAlertVisible: false,
      redirectPasswordAlert: '',
      domainAlertVisible: false,
      domainAlert: '',
      redirectEmailAlertVisible: false,
      redirectEmailAlert: '',
      step: 0
    }
  },
  mounted () {
    this.progressShow()
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
    validDeviceCredentials () {
      if (this.deviceUsername === '') {
        return false
      }
      if (this.devicePassword === '') {
        return false
      }
      if (this.devicePassword !== this.devicePasswordConfirm) {
        return false
      }
      return true
    },
    activate (event) {
      if (!this.validDeviceCredentials()) {
        return
      }
      event.preventDefault()
      this.progressShow()
      this.hideAlerts()
      switch (this.domainType) {
        case 'premium':
          this.activatePremiumDomain()
          break
        default:
          this.activateFreeDomain()
      }
    },
    forceCertificateRecheck () {
      window.location = '/?t=' + (new Date()).getTime()
    },
    hideAlerts () {
      this.deviceUsernameAlertVisible = false
      this.devicePasswordAlertVisible = false
      this.redirectEmailAlertVisible = false
      this.redirectPasswordAlertVisible = false
      this.domainAlertVisible = false
    },
    showRedirectAlert (err) {
      if (err.response) {
        const response = err.response
        if (response.data) {
          const data = response.data
          if (data.parameters_messages) {
            for (let i = 0; i < data.parameters_messages.length; i++) {
              const pm = data.parameters_messages[i]
              const message = pm.messages.join('\n')
              if (pm.parameter === 'redirect_password') {
                this.redirectPasswordAlertVisible = true
                this.redirectPasswordAlert = message
              }
              if (pm.parameter === 'email') {
                this.redirectEmailAlertVisible = true
                this.redirectEmailAlert = message
              }
              if (pm.parameter === 'domain') {
                this.domainAlertVisible = true
                this.domainAlert = message
              }
            }
          } else {
            this.$refs.error.showAxios(err)
          }
        }
      }
    },
    showActivateAlert (err) {
      if (err.response) {
        const response = err.response
        if (response.data) {
          const data = response.data
          if (data.parameters_messages) {
            for (let i = 0; i < data.parameters_messages.length; i++) {
              const pm = data.parameters_messages[i]
              const message = pm.messages.join(', ')
              if (pm.parameter === 'device_username') {
                this.deviceUsernameAlertVisible = true
                this.deviceUsernameAlert = message
              }
              if (pm.parameter === 'device_password') {
                this.devicePasswordAlertVisible = true
                this.devicePasswordAlert = message
              }
            }
          } else {
            this.$refs.error.showAxios(err)
          }
        }
      }
    },
    activateFreeDomain () {
      let domain = this.fullDomain();
      axios
        .post('/rest/activate/managed', {
          redirect_email: this.redirectEmail,
          redirect_password: this.redirectPassword,
          domain: domain,
          device_username: this.deviceUsername,
          device_password: this.devicePassword
        })
        .then(this.forceCertificateRecheck)
        .catch(err => {
          this.progressHide()
          this.showActivateAlert(err)
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
          this.showActivateAlert(err)
        })
    },
    showDeviceCredentialHelp () {
      this.helpDeviceCredentialVisible = true
    },
    showFreeAccountHelp () {
      this.helpFreeAccountVisible = true
    },
    showPremiumAccountHelp () {
      this.helpPremiumAccountVisible = true
    },
    showManagedDomainHelp () {
      this.helpManagedDomainVisible = true
    },
    selectPremiumDomain () {
      this.hideAlerts()
      this.domainType = 'premium'
      this.step++
    },
    selectFreeDomain () {
      this.hideAlerts()
      this.domainType = 'free'
      this.step++
    },
    selectDeviceName () {
      this.hideAlerts()
      this.domainAvailability()
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
          this.step++
          this.progressHide()
        })
        .catch(err => {
          this.progressHide()
          this.showRedirectAlert(err)
        })
    }
  }
}
</script>
<style>
@import '../style/site.css';
@import 'material-icons/iconfont/material-icons.css';
@import 'font-awesome/css/font-awesome.css';

.register {
  color: #00aeef;
  font-weight: bold;
}

.domain {
  width: 70%!important;
  margin-right: 2%
}

.input {
  width: 100%;
  height: 40px;
  padding: 0 50px 0 20px;
  border-radius: 3px;
  border: 1px solid #dcdee0;
  margin-bottom: 10px;
  background-color: #fff!important;
  background-size: 14px 14px;
  transition: all .3s ease-out;
}

* {
  box-sizing: border-box;
}

.columns {
  width: 33.3%;
  padding: 8px;
}

.plan {
  list-style-type: none;
  margin: 0;
  padding: 0;
  transition: 0.3s;
}

.plan:hover {
  box-shadow: 0 8px 12px 0 rgba(0, 0, 0, 0.2)
}

.plan li {
  padding: 0 20px 20px 20px ;
  text-align: center;
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
