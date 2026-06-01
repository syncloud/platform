<template>
  <div class="sc-page">
    <div class="sc-card" id="block1">
      <h1 class="sc-title">{{ $t('access.title') }}</h1>
      <div :style="{ visibility: visibility }">
        <h3>{{ $t('access.ipv4') }}</h3>
        <div class="setline" style='display: flex'>
          <span class="span name-alignment">{{ $t('access.support') }}</span>
          <div class="value-alignment">
            <el-switch id="tgl_ipv4_enabled" size="large" v-model="ipv4Enabled" style="--el-switch-on-color: #2faa5d; float: right" />
          </div>
          <button type=button @click="showIpv4Info" class="control" style="order: 3; background:transparent; border: none; cursor: pointer;">
            <i class='fa fa-question-circle fa-lg' style="color: var(--sc-faint)"></i>
          </button>
        </div>

        <Transition @after-enter="(el) => el.setAttribute('data-ready', 'true')">
        <div id="ipv4_mode_block" v-if="ipv4Enabled">
          <div class="setline" style='display: flex'>
              <span class="span name-alignment">{{ $t('access.public') }}</span>
              <div class="value-alignment">
                <el-switch id="tgl_ipv4_public" size="large" v-model="ipv4Public" style="--el-switch-on-color: #2faa5d; float: right" />
              </div>
          </div>

          <Transition @after-enter="(el) => el.setAttribute('data-ready', 'true')">
          <div id="ipv4_public_block" v-if="ipv4Public">
            <div class="setline" style='display: flex'>
                <span class="span name-alignment">{{ $t('access.detectIp') }}</span>
                <div class="value-alignment">
                  <el-switch id="tgl_ip_autodetect" size="large" v-model="ipAutoDetect" style="--el-switch-on-color: #2faa5d; float: right" />
                </div>
            </div>

            <Transition>
            <div class="setline" id="ipv4_block" style='display: flex' v-if="!ipAutoDetect">
              <label class="span name-alignment" for="ipv4" style="font-weight: 300">{{ $t('access.publicIp') }}</label>
              <input class="value-alignment sc-input" id="ipv4" type="text"
                     style="width: 130px; height: 38px;"
                     :disabled="ipAutoDetect" v-model="ipv4">
            </div>
            </Transition>

            <div class="setline" style='display: flex'>
                <label for="access_port" class="span name-alignment" style="font-weight: 300">{{ $t('access.publicPort') }}</label>
                <input class="value-alignment sc-input" id="access_port" type="number"
                       style="width: 110px; height: 38px;"
                       v-model.number="accessPort"
                />
                <button type=button @click="showPortInfo" class="control" style="order: 3; background:transparent; border: none; cursor: pointer;">
                  <i class='fa fa-question-circle fa-lg' style="color: var(--sc-faint)"></i>
                </button>
                <button id="access_port_warning" type=button @click="showAccessPortWarning"
                        class="control" style="order: 4; background:transparent; border: none; cursor: pointer;" v-show="accessPort!==443">
                  <i class='fa fa-exclamation-circle fa-lg' style='color: red;'></i>
                </button>

            </div>
          </div>
          </Transition>

        </div>
        </Transition>

        <h3>{{ $t('access.ipv6') }}</h3>

        <div class="setline" style='display: flex'>
          <span class="span name-alignment">{{ $t('access.support') }}</span>
          <div class="value-alignment">
            <el-switch id="tgl_ipv6_enabled" size="large" v-model="ipv6Enabled" style="--el-switch-on-color: #2faa5d; float: right" />
          </div>
          <button type=button @click="showIpv6Info" class="control" style="order: 3; background:transparent; border: none; cursor: pointer;">
            <i class='fa fa-question-circle fa-lg' style="color: var(--sc-faint)"></i>
          </button>
        </div>

        <div class="sc-actions">
          <button class="sc-btn sc-btn-success" id="btn_save" type="submit"
                  data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Working..."
                  @click="save">{{ $t('access.save') }}
          </button>
        </div>
      </div>
    </div>
  </div>

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
import { ElLoading } from 'element-plus'

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
      interfaces: undefined,
      ipAutoDetect: undefined,
      ipv4: '',
      accessPort: 443,
      visibility: 'hidden',
      ipv4Enabled: undefined,
      ipv4Public: undefined,
      ipv6Enabled: undefined,
      loading: undefined,
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
  watch: {
    ipv4Public (val) {
      if (!val) {
        this.accessPort = 443
      }
    }
  },
  mounted () {
    this.progressShow()
    this.reload()
  },
  methods: {
    progressShow () {
      this.loading = ElLoading.service({ lock: true, text: this.$t('common.loading'), background: 'rgba(0, 0, 0, 0.7)' })
    },
    progressHide () {
      this.visibility = 'visible'
      this.loading.close()
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
      const error = this.$refs.error

      const onError = (err) => {
        error.showAxios(err)
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
        that.ipv4Enabled = accessData.ipv4_enabled
        that.ipv4Public = accessData.ipv4_public
        that.ipv6Enabled = accessData.ipv6_enabled
        this.progressHide()
      }
      axios.get('/rest/access')
        .then(resp => Common.checkForServiceError(resp.data.data, () => onComplete(resp.data.data), onError))
        .catch(onError)
    },
    save (event) {
      this.progressShow()

      event.preventDefault()
      const that = this
      const requestData = {
        access_port: this.accessPort,
        ipv4_enabled: this.ipv4Enabled,
        ipv4_public: this.ipv4Public,
        ipv6_enabled: this.ipv6Enabled
      }
      if (this.ipv4Enabled) {
        if (!isValidPort(this.accessPort)) {
          this.$refs.error.showAxios(error(this.$t('access.errorPortRange', { port: this.accessPort })))
          this.progressHide()
          return
        }
        if (!this.ipAutoDetect) {
          if (this.ipv4.trim() === '') {
            this.$refs.error.showAxios(error(this.$t('access.errorEmptyIp')))
            this.progressHide()
            return
          }
          requestData.ipv4 = this.ipv4
        }
      }

      const onError = (err) => {
        that.$refs.error.showAxios(err)
        this.progressHide()
      }
      axios.post('/rest/access', requestData)
        .then(response => Common.checkForServiceError(response.data, this.reload, onError))
        .catch(onError)
    }
  }
}
</script>
<style scoped>
.name-alignment {
  min-width: 100px;
  display: inline-flex;
  align-items: center;
  order: 1;
}

.value-alignment {
  order: 2;
  min-width: 130px;
  margin-right: 5px;
}

.v-enter-active,
.v-leave-active {
  transition: opacity 0.5s ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}
</style>
