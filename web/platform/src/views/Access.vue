<template>
  <div class="sc-page">
    <div class="sc-card" id="block1">
      <h1 class="sc-title">{{ $t('access.title') }}</h1>
      <div :style="{ visibility: visibility }">

        <div class="access-block" data-testid="ipv4-section">
          <div class="access-head">
            <div class="access-heading">
              <span class="access-name">{{ $t('access.ipv4') }}</span>
              <button type="button" @click="showIpv4Info" class="access-help">
                <i class='fa fa-question-circle fa-lg'></i>
              </button>
            </div>
          </div>

          <div class="seg" data-testid="ipv4-modes">
            <button type="button" class="seg-btn" :class="{ active: ipv4Mode === 'relay' }"
                    data-testid="ipv4-mode-relay" @click="ipv4Mode = 'relay'">{{ $t('access.ipv4Relay') }}</button>
            <button type="button" class="seg-btn" :class="{ active: ipv4Mode === 'public' }"
                    data-testid="ipv4-mode-public" @click="ipv4Mode = 'public'">{{ $t('access.ipv4Public') }}</button>
            <button type="button" class="seg-btn" :class="{ active: ipv4Mode === 'local' }"
                    data-testid="ipv4-mode-local" @click="ipv4Mode = 'local'">{{ $t('access.ipv4Local') }}</button>
            <button type="button" class="seg-btn" :class="{ active: ipv4Mode === 'off' }"
                    data-testid="ipv4-mode-off" @click="ipv4Mode = 'off'">{{ $t('access.ipv4Off') }}</button>
          </div>

          <div class="reveal" :class="{ open: ipv4Mode === 'relay' }">
            <div class="reveal-inner">
              <div class="detail" data-testid="relay-status">
                <p class="detail-desc">{{ $t('access.relayDescription') }}</p>
              </div>
            </div>
          </div>

          <div class="reveal" :class="{ open: ipv4Mode === 'public' }">
            <div class="reveal-inner">
              <div class="detail" data-testid="ipv4-public">
                <p class="detail-desc">{{ $t('access.ipv4PublicText') }}</p>
                <div class="row">
                  <span class="row-label">{{ $t('access.detectIp') }}</span>
                  <s-switch id="tgl_ip_autodetect" size="large" v-model="ipAutoDetect" style="--el-switch-on-color: #2faa5d" />
                </div>

                <div class="reveal" :class="{ open: !ipAutoDetect }">
                  <div class="reveal-inner">
                    <div class="row">
                      <label class="row-label" for="ipv4">{{ $t('access.publicIp') }}</label>
                      <input class="sc-input row-input" id="ipv4" data-testid="ipv4-input" type="text" v-model="ipv4">
                    </div>
                  </div>
                </div>

                <div class="row">
                  <span class="row-label">{{ $t('access.publicPort') }}
                    <button type="button" @click="showPortInfo" class="access-help"><i class='fa fa-question-circle fa-lg'></i></button>
                  </span>
                  <div class="row-right">
                    <button id="access_port_warning" type="button" @click="showAccessPortWarning" class="access-help warn" v-show="accessPort!==443"><i class='fa fa-exclamation-circle fa-lg'></i></button>
                    <input class="sc-input row-input" id="access_port" type="number" v-model.number="accessPort" />
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="reveal" :class="{ open: ipv4Mode === 'local' }">
            <div class="reveal-inner">
              <div class="detail" data-testid="ipv4-local">
                <p class="detail-desc">{{ $t('access.ipv4LocalText') }}</p>
              </div>
            </div>
          </div>

          <div class="reveal" :class="{ open: ipv4Mode === 'off' }">
            <div class="reveal-inner">
              <div class="detail" data-testid="ipv4-off">
                <p class="detail-desc">{{ $t('access.ipv4OffText') }}</p>
              </div>
            </div>
          </div>
        </div>

        <div class="access-block" data-testid="ipv6-section">
          <div class="access-head">
            <div class="access-heading">
              <span class="access-name">{{ $t('access.ipv6') }}</span>
              <button type="button" @click="showIpv6Info" class="access-help">
                <i class='fa fa-question-circle fa-lg'></i>
              </button>
            </div>
            <s-switch id="tgl_ipv6_enabled" data-testid="ipv6-toggle" size="large" v-model="ipv6Enabled"
                      style="--el-switch-on-color: #2faa5d" />
          </div>
          <p class="detail-desc">{{ $t('access.ipv6Description') }}</p>
        </div>

        <div class="sc-actions">
          <button class="sc-btn sc-btn-success" id="btn_save" data-testid="access-save" type="submit"
                  data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Working..."
                  @click="save">{{ $t('access.save') }}
          </button>
        </div>
      </div>
    </div>
  </div>

  <Dialog :visible="relayInfoVisible" @cancel="relayInfoVisible = false" :confirm-enabled="false" :cancel-text="$t('common.close')">
    <template v-slot:title>{{ $t('access.relayInfoTitle') }}</template>
    <template v-slot:text>
      {{ $t('access.relayInfoText') }}
    </template>
  </Dialog>
  <Dialog :visible="accessPortInfoVisible" @cancel="accessPortInfoVisible = false" :confirm-enabled="false" :cancel-text="$t('common.close')">
    <template v-slot:title>{{ $t('access.accessPortTitle') }}</template>
    <template v-slot:text>
      {{ $t('access.accessPortText') }}
    </template>
  </Dialog>
  <Dialog :visible="accessPortWarningVisible" @cancel="accessPortWarningVisible = false" :confirm-enabled="false" :cancel-text="$t('common.close')">
    <template v-slot:title>{{ $t('access.accessPortWarningTitle') }}</template>
    <template v-slot:text>
      {{ $t('access.accessPortWarningText') }}
    </template>
  </Dialog>
  <Dialog :visible="ipv4InfoVisible" @cancel="ipv4InfoVisible = false" :confirm-enabled="false" :cancel-text="$t('common.close')">
    <template v-slot:title>{{ $t('access.ipv4Title') }}</template>
    <template v-slot:text>
      {{ $t('access.ipv4Text') }}
    </template>
  </Dialog>
  <Dialog :visible="ipv6InfoVisible" @cancel="ipv6InfoVisible = false" :confirm-enabled="false" :cancel-text="$t('common.close')">
    <template v-slot:title>{{ $t('access.ipv6Title') }}</template>
    <template v-slot:text>
      {{ $t('access.ipv6Text') }}
    </template>
  </Dialog>

  <Error ref="error"/>

</template>

<script>
import Error from '../components/Error.vue'
import Dialog from '../components/Dialog.vue'
import * as Common from '../js/common.js'
import axios from 'axios'
import Loading from '../util/loading'

function isValidPort (port) {
  return !(Number.isNaN(port) || port < 1 || port > 65535)
}

function error (message) {
  return {
    response: {
      status: 200,
      data: {
        message: message
      }
    }
  }
}

export default {
  name: 'Access',
  data () {
    return {
      ipAutoDetect: undefined,
      ipv4: '',
      accessPort: 443,
      visibility: 'hidden',
      ipv4Mode: 'off',
      ipv6Enabled: undefined,
      loading: undefined,
      relayInfoVisible: false,
      accessPortInfoVisible: false,
      accessPortWarningVisible: false,
      ipv4InfoVisible: false,
      ipv6InfoVisible: false
    }
  },
  components: {
    Error,
    Dialog
  },
  mounted () {
    this.progressShow()
    this.reload()
  },
  methods: {
    progressShow (text) {
      this.loading = Loading.service({ lock: true, text: text || this.$t('common.loading'), background: 'rgba(0, 0, 0, 0.7)' })
    },
    waitReachable () {
      const attempt = (left) => axios.get('/rest/access', { timeout: 4000 })
        .then(() => true)
        .catch(() => {
          if (left <= 0) {
            return false
          }
          return new Promise(resolve => setTimeout(resolve, 2000)).then(() => attempt(left - 1))
        })
      return attempt(45)
    },
    progressHide () {
      this.visibility = 'visible'
      this.loading.close()
    },
    showRelayInfo () {
      this.relayInfoVisible = true
    },
    showAccessPortWarning () {
      this.accessPortWarningVisible = true
    },
    showIpv4Info () {
      this.ipv4InfoVisible = true
    },
    showIpv6Info () {
      this.ipv6InfoVisible = true
    },
    showPortInfo () {
      this.accessPortInfoVisible = true
    },
    reload () {
      const that = this
      const err = this.$refs.error

      const onError = (e) => {
        err.showAxios(e)
        this.progressHide()
      }
      const onComplete = (data) => {
        const accessData = data
        if (accessData.ipv4) {
          that.ipAutoDetect = false
          that.ipv4 = accessData.ipv4
        } else {
          that.ipAutoDetect = true
        }
        if (accessData.access_port !== undefined) {
          that.accessPort = accessData.access_port
        }
        that.ipv6Enabled = accessData.ipv6_enabled
        if (accessData.relay_enabled === true) {
          that.ipv4Mode = 'relay'
        } else if (accessData.ipv4_enabled === true) {
          that.ipv4Mode = accessData.ipv4_public === true ? 'public' : 'local'
        } else {
          that.ipv4Mode = 'off'
        }
        this.progressHide()
      }
      axios.get('/rest/access')
        .then(resp => Common.checkForServiceError(resp.data.data, () => onComplete(resp.data.data), onError))
        .catch(onError)
    },
    save (event) {
      event.preventDefault()
      const mode = this.ipv4Mode
      const requestData = {
        relay_enabled: mode === 'relay',
        access_port: this.accessPort,
        ipv4_enabled: mode === 'public' || mode === 'local',
        ipv4_public: mode === 'public',
        ipv6_enabled: this.ipv6Enabled
      }
      if (mode === 'public') {
        if (!isValidPort(this.accessPort)) {
          this.$refs.error.showAxios(error(this.$t('access.errorPortRange', { port: this.accessPort })))
          return
        }
        if (!this.ipAutoDetect) {
          if (this.ipv4.trim() === '') {
            this.$refs.error.showAxios(error(this.$t('access.errorEmptyIp')))
            return
          }
          requestData.ipv4 = this.ipv4
        }
      }

      this.progressShow(this.$t('access.applying'))
      axios.post('/rest/access', requestData, { timeout: 90000 })
        .then(response => {
          if (response.data && 'success' in response.data && !response.data.success) {
            this.$refs.error.showAxios({ response: { status: 200, data: response.data } })
            this.progressHide()
            return
          }
          this.reload()
        })
        .catch(err => {
          if (err.response) {
            this.$refs.error.showAxios(err)
            this.progressHide()
            return
          }
          this.waitReachable().then(back => {
            if (back) {
              this.reload()
            } else {
              this.$refs.error.showAxios(error(this.$t('access.applyTimeout')))
              this.progressHide()
            }
          })
        })
    }
  }
}
</script>
<style scoped>
.access-block {
  border: 1px solid var(--sc-border);
  border-radius: 16px;
  padding: 20px 22px;
  margin-bottom: 16px;
  background: var(--sc-surface);
}
.access-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}
.access-heading {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}
.access-name {
  font-size: 17px;
  font-weight: 700;
  color: var(--sc-ink);
}
.access-help {
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 0;
  color: var(--sc-faint);
  line-height: 1;
  display: inline-flex;
  align-items: center;
}
.access-help.warn {
  color: #f56c6c;
}

.seg {
  display: grid;
  grid-auto-flow: column;
  grid-auto-columns: 1fr;
  gap: 4px;
  margin-top: 16px;
  padding: 4px;
  background: var(--sc-surface-3);
  border-radius: 12px;
}
.seg-btn {
  padding: 9px 12px;
  border: none;
  border-radius: 9px;
  background: transparent;
  color: var(--sc-ink-2);
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.15s ease, color 0.15s ease, box-shadow 0.15s ease;
}
.seg-btn:hover:not(.active) {
  background: rgba(127, 127, 127, 0.08);
}
.seg-btn.active {
  background: var(--sc-surface);
  color: var(--sc-primary);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.12);
}

.reveal {
  display: grid;
  grid-template-rows: 0fr;
  transition: grid-template-rows 0.25s ease;
}
.reveal.open {
  grid-template-rows: 1fr;
}
.reveal > .reveal-inner {
  overflow: hidden;
  min-height: 0;
}

.detail {
  padding-top: 8px;
}
.detail-desc {
  color: var(--sc-muted);
  font-size: 14px;
  margin: 8px 0 0;
  max-width: 460px;
  line-height: 1.5;
}
.row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  min-height: 48px;
  border-top: 1px solid var(--sc-border-soft);
}
.detail > .row:first-child {
  border-top: none;
}
.row-label {
  color: var(--sc-ink-2);
  font-size: 14px;
  font-weight: 400;
}
.row-right {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}
.row-input {
  width: 130px;
  height: 38px;
}
</style>
