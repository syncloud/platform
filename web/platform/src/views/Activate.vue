<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12" id="block_activate">
        <h1>{{ $t('activate.title') }}</h1>

        <div>
          <el-steps :active="step" simple finish-status="success" style="max-width: 500px; margin: 0 auto">
            <el-step :title="$t('activate.stepType')" />
            <el-step :title="$t('activate.stepName')" />
            <el-step :title="$t('activate.stepUser')" />
          </el-steps>

          <div v-if="step === 0" id="domain-type-part"  >
            <div class="language-picker">
              <span class="language-label">{{ $t('language.select') }}:</span>
              <el-select id="activate_language" v-model="locale" @change="onLocaleChange" size="default" style="width: 180px">
                <el-option v-for="l in locales" :key="l.code" :label="l.name" :value="l.code" :id="'lang_' + l.code"/>
              </el-select>
            </div>
            <div style="text-align: center; max-width: 800px; margin: 0 auto; display: flex; flex-wrap: wrap; justify-content: center;">
              <div class="columns">
                <ul class="plan">
                  <li class="description">
                    {{ $t('activate.premiumDescription') }}
                  </li>
                  <li>
                    <el-button id="btn_premium_domain" style="width: 100%;height: 40px;" type="primary" @click="selectPremiumDomain">
                      {{ $t('activate.premiumButton') }}
                    </el-button>
                  </li>
                </ul>
              </div>
              <div class="columns">
                <ul class="plan">
                  <li class="description">
                    {{ $t('activate.freeDescription', { domain: redirect_domain }) }}
                  </li>
                  <li>
                    <el-button id="btn_free_domain" style="width: 100%;height: 40px;" type="primary" @click="selectFreeDomain">
                      {{ $t('activate.freeButton') }}
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
                  <h2 style="display: inline-block">{{ $t('activate.syncloudAccount') }}</h2>
                  <button @click="showFreeAccountHelp" type=button class="control"
                          style="vertical-align: super; background:transparent;">
                    <i class='fa fa-question-circle fa-lg'></i>
                  </button>
                </div>

                <input :placeholder="$t('activate.emailPlaceholder', { domain: redirect_domain })" class="input" id="email"
                       type="text" v-model="redirectEmail">
                <div id="email_alert" class="alert alert-danger alert90" v-show="redirectEmailAlertVisible ">{{ redirectEmailAlert }}</div>

                <Password
                  id="redirect_password"
                  v-model="redirectPassword"
                  :placeholder="$t('activate.passwordPlaceholder', { domain: redirect_domain })"
                  :show-error="redirectPasswordAlertVisible"
                  :error="redirectPasswordAlert"
                />

                <div style=" display: flow-root">
                  <div style="padding-right:10px; float: right">
                    {{ $t('activate.noAccount') }}
                    <a :href="'https://' + redirect_domain" class="register"
                       target="_blank">{{ $t('activate.register') }}</a>
                  </div>
                </div>
                <div style="text-align: center">
                  <h2 style="display: inline-block">{{ $t('activate.deviceName') }}</h2>
                  <button @click="showManagedDomainHelp" type=button class="control"
                          style="vertical-align: super; background:transparent;">
                    <i class='fa fa-question-circle fa-lg'></i>
                  </button>
                </div>

                <div id="domain" style="text-align: left">
                  <input :placeholder="$t('activate.domainNamePlaceholder')" class="domain input" id="domain_input" type="text" v-model="domain">
                  <span>.{{ redirect_domain }}</span>
                </div>
                <div id="domain_alert"  class="alert alert-danger alert90" v-show="domainAlertVisible" >{{ domainAlert }}</div>

              </div>

              <div v-if=" domainType === 'premium' ">
                <div style="text-align: center">
                  <h2 style="display: inline-block">{{ $t('activate.syncloudAccount') }}</h2>
                  <button @click="showPremiumAccountHelp" type="button" class="control"
                          style="vertical-align: super; background:transparent;">
                    <i class='fa fa-question-circle fa-lg'></i>
                  </button>
                </div>

                <input :placeholder="$t('activate.emailPlaceholder', { domain: redirect_domain })" class="input" id="email"
                       type="text" v-model="redirectEmail">
                <div class="alert alert-danger alert90" id="alert" style="display: none;"></div>

                <Password
                  id="redirect_password"
                  v-model="redirectPassword"
                  :placeholder="$t('activate.passwordPlaceholder', { domain: redirect_domain })"
                  :show-error="redirectPasswordAlertVisible"
                  :error="redirectPasswordAlert"
                />

                <div style="text-align: center">
                  <h2 style="display: inline-block">{{ $t('activate.deviceName') }}</h2>
                  <button @click="showManagedDomainHelp" type=button class="control"
                          style="vertical-align: super; background:transparent;">
                    <i class='fa fa-question-circle fa-lg'></i>
                  </button>
                </div>
                <div id="domain">
                  <input :placeholder="$t('activate.premiumDomainPlaceholder')"
                         class="domain input" id="domain_premium" type="text" style="width:100% !important;"
                         v-model="domain">
                </div>
                <div id="domain_alert" class="alert alert-danger alert90" v-show="domainAlertVisible" >{{ domainAlert }}</div>
              </div>

              <div style="padding: 10px; float: left;">
                <el-button type="primary" @click="step--">
                  {{ $t('common.previous') }}
                </el-button>
              </div>
              <div style="padding: 10px; float: right;">
                <el-button id="btn_next" type="primary"  @click="selectDeviceName">
                  {{ $t('common.next') }}
                </el-button>
              </div>
            </div>
          </div>

          <div v-if="step === 2" id="device-credentials-part">
            <div style="text-align: center; max-width: 400px; margin: 0 auto;">

              <div style="text-align: center">
                <h2 style="display: inline-block">{{ $t('activate.deviceCredentials') }}</h2>
                <button @click="showDeviceCredentialHelp" type=button class="control"
                        style="vertical-align: super; background:transparent;">
                  <i class='fa fa-question-circle fa-lg'></i>
                </button>
              </div>

              <input :placeholder="$t('activate.loginPlaceholder')" id="device_username" type="text" v-model="deviceUsername"
                     v-on:keyup.enter="activate"
                     class="input">
              <div class="alert alert-danger alert90" v-show="deviceUsernameAlertVisible">{{
                  deviceUsernameAlert
                }}
              </div>

              <Password
                id="device_password"
                v-model="devicePassword"
                :placeholder="$t('activate.passwordInputPlaceholder')"
                :show-error="devicePasswordAlertVisible"
                :error="devicePasswordAlert"
                @trigger="activate"
              />

              <Password
                id="device_password_confirm"
                v-model="devicePasswordConfirm"
                :placeholder="$t('activate.confirmPasswordPlaceholder')"
                :show-error="devicePassword !== devicePasswordConfirm"
                :error="$t('activate.passwordsMismatch')"
                @trigger="activate"
              />

              <div style="padding: 10px; float: left;">
                <el-button type="primary" @click="step--">
                  {{ $t('common.previous') }}
                </el-button>
              </div>
              <div style="padding: 10px; float: right;">
                <el-button id="btn_activate" type="primary" @click="activate" :disabled="!validDeviceCredentials()">
                {{ $t('activate.finish') }}
                </el-button>
              </div>
            </div>
          </div>

        </div>
      </div>
    </div>
  </div>

  <Dialog :visible="helpManagedDomainVisible" @cancel="helpManagedDomainVisible = false" :confirm-enabled="false" :cancel-text="$t('common.close')">
    <template v-slot:title>{{ $t('activate.help.managedDomainTitle') }}</template>
    <template v-slot:text>
      <div class="btext">{{ $t('activate.help.managedDomainText') }}</div>
    </template>
  </Dialog>

  <Dialog :visible="helpFreeAccountVisible" @cancel="helpFreeAccountVisible = false" :confirm-enabled="false" :cancel-text="$t('common.close')">
    <template v-slot:title>{{ $t('activate.help.syncloudAccountTitle') }}</template>
    <template v-slot:text>
      {{ $t('activate.help.freeAccountText') }}
    </template>
  </Dialog>

  <Dialog :visible="helpPremiumAccountVisible" @cancel="helpPremiumAccountVisible = false" :confirm-enabled="false" :cancel-text="$t('common.close')">
    <template v-slot:title>{{ $t('activate.help.syncloudAccountTitle') }}</template>
    <template v-slot:text>
      {{ $t('activate.help.premiumAccountText') }}
    </template>
  </Dialog>

  <Dialog :visible="helpDeviceCredentialVisible" @cancel="helpDeviceCredentialVisible = false" :confirm-enabled="false" :cancel-text="$t('common.close')">
    <template v-slot:title>{{ $t('activate.help.deviceCredentialsTitle') }}</template>
    <template v-slot:text>
      {{ $t('activate.help.deviceCredentialsText') }}
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
import { SUPPORTED_LOCALES, setLocale } from '../i18n'

export default {
  name: 'Activate',
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
      step: 0,
      locale: this.$i18n ? this.$i18n.locale : 'en',
      locales: SUPPORTED_LOCALES
    }
  },
  mounted () {
    this.progressShow()
    this.locale = this.$i18n.locale
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
    onLocaleChange (code) {
      setLocale(code)
    },
    progressShow () {
      this.loading = ElLoading.service({ lock: true, text: this.$t('common.loading'), background: 'rgba(0, 0, 0, 0.7)' })
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
      if (event && event.preventDefault) {
        event.preventDefault()
      }
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
        .then(() => {
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

.language-picker {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 10px;
  margin: 20px 0 10px 0;
}

.language-label {
  font-size: 14px;
  color: #555;
}

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
