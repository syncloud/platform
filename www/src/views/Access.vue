<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12">
        <h1>Access</h1>
        <div class="row-no-gutters settingsblock">
          <div class="col2" :style="{ visibility: visibility }">
            <div class="setline">
              <h3>IP v4</h3>
            </div>
            <div class="setline" style='display: flex'>
              <span class="span name-alignment">Support:</span>
              <div class="value-alignment">
                <el-switch id="tgl_ipv4_enabled" size="large" v-model="ipv4Enabled" style="--el-switch-on-color: #36ad40; float: right" />
              </div>
              <button type=button @click="showIpv4Info" class="control" style="order: 3; background:transparent;">
                <i class='fa fa-question-circle fa-lg'></i>
              </button>
            </div>

            <div id="ipv4_mode_block">
              <div class="setline" style='display: flex'>
                  <span class="span name-alignment">Public:</span>
                  <div class="value-alignment">
                    <el-switch id="tgl_ipv4_public" size="large" v-model="ipv4Public" style="--el-switch-on-color: #36ad40; float: right" />
                  </div>
              </div>
              <div id="ipv4_public_block">
                <div class="setline" style='display: flex'>
                    <span class="span name-alignment">Detect IP:</span>
                    <div class="value-alignment">
                      <el-switch id="tgl_ip_autodetect" size="large" v-model="ipAutoDetect" style="--el-switch-on-color: #36ad40; float: right" />
                    </div>
                </div>

                <div class="setline" id="ipv4_block" style='display: flex'>
                  <label class="span name-alignment" for="ipv4" style="font-weight: 300">Public IP:</label>
                  <input class="value-alignment" id="ipv4" type="text"
                         style="width: 130px; height: 30px; padding: 0 10px 0 10px"
                         :disabled="ipAutoDetect" v-model="ipv4">
                </div>

                <div class="setline" style='display: flex'>
                    <label for="access_port" class="span name-alignment" style="font-weight: 300">Public port:</label>
                    <input class="value-alignment" id="access_port" type="number"
                           style="width: 100px; height: 30px; padding: 0 10px 0 10px"
                           v-model.number="accessPort"
                    />
                    <button type=button @click="showPortInfo" class="control" style="order: 3; background:transparent;">
                      <i class='fa fa-question-circle fa-lg'></i>
                    </button>
                    <button id="access_port_warning" type=button @click="showAccessPortWarning"
                            class="control" style="order: 4; background:transparent;" v-show="false">
                      <i class='fa fa-exclamation-circle fa-lg' style='color: red;'></i>
                    </button>

                </div>
              </div>
            </div>

            <div class="setline">
              <h3>IP v6</h3>
            </div>

            <div class="setline" style='display: flex'>
              <span class="span name-alignment">Support:</span>
              <div class="value-alignment">
                <el-switch id="tgl_ipv6_enabled" size="large" v-model="ipv6Enabled" style="--el-switch-on-color: #36ad40; float: right" />
              </div>
              <button type=button @click="showIpv6Info" class="control" style="order: 3; background:transparent;">
                <i class='fa fa-question-circle fa-lg'></i>
              </button>
            </div>

            <div class="setline">
              <div class="spandiv">
                <button class="submit buttongreen control" id="btn_save" type="submit"
                        data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Working..."
                        style="width: 150px" @click="save">Save
                </button>
              </div>
            </div>

          </div>

        </div>
      </div>
    </div>
  </div>

  <Dialog ref="access_port_info">
    <template v-slot:title>Access port</template>
    <template v-slot:text>
      If your Syncloud device is not visible directly from the Internet you will need to create a port mapping on your
      router.
      Ideally port 443 on your router should be mapped to port 443 on your device.
    </template>
  </Dialog>
  <Dialog ref="access_port_warning">
    <template v-slot:title>Access port warning</template>
    <template v-slot:text>
      Access port is not default 443.
      You may not be able to access your device from networks with strict firewalls allowing only port 443.
    </template>
  </Dialog>
  <Dialog ref="ipv4_info">
    <template v-slot:title>IP v4</template>
    <template v-slot:text>
      Enables IP v4 DNS record and allows you to control which IP v4 address (public/private) is used for DNS.
      <br><br>
      Syncloud DNS service verifies open public ip/port (internet accessibility) for convenience on save.
    </template>
  </Dialog>
  <Dialog ref="ipv6_info">
    <template v-slot:title>IP v6</template>
    <template v-slot:text>
      Enables IP v6 DNS record.
      <br><br>
      Syncloud DNS service verifies device connection (internet accessibility) for convenience on save.
    </template>
  </Dialog>

  <Error ref="error"/>

</template>

<script>
import $ from 'jquery'
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
  props: {
    checkUserSession: Function,
    activated: Boolean
  },
  data () {
    return {
      interfaces: undefined,
      ipAutoDetect: false,
      ipv4: '',
      accessPort: 443,
      visibility: 'hidden',
      ipv4Enabled: true,
      ipv4Public: false,
      ipv6Enabled: true,
      loading: undefined
    }
  },
  components: {
    Error,
    Dialog
  },
  watch: {
    ipv4Enabled (val) {
      this.displayIpv4Mode(val)
    },
    ipv4Public (val) {
      if (val) {
        $('#ipv4_public_block').show('slow')
      } else {
        $('#ipv4_public_block').hide('slow')
        this.accessPort = 443
      }
    },
    ipAutoDetect (val) {
      this.displayIpv4Manual(val)
    },
    accessPort (val) {
      if (val !== 443) {
        $('#access_port_warning').show('slow')
      } else {
        $('#access_port_warning').hide('slow')
      }
    }
  },
  mounted () {
    this.progressShow()
    this.reload()
  },
  methods: {
    progressShow () {
      this.loading = ElLoading.service({ lock: true, text: 'Loading', background: 'rgba(0, 0, 0, 0.7)' })
    },
    progressHide () {
      this.visibility = 'visible'
      this.loading.close()
    },
    displayIpv4Manual (val) {
      if (val) {
        $('#ipv4_block').hide('slow')
      } else {
        $('#ipv4_block').show('slow')
      }
    },
    displayIpv4Mode (val) {
      if (val) {
        $('#ipv4_mode_block').show('slow')
      } else {
        $('#ipv4_mode_block').hide('slow')
      }
    },
    showAccessPortWarning () {
      this.$refs.access_port_warning.show()
    },
    showIpv4Info () {
      this.$refs.ipv4_info.show()
    },
    showIpv6Info () {
      this.$refs.ipv6_info.show()
    },
    showPortInfo () {
      this.$refs.access_port_info.show()
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
          this.$refs.error.showAxios(error('Access port (' + this.accessPort + ') has to be between 1 and 65535'))
          this.progressHide()
          return
        }
        if (!this.ipAutoDetect) {
          if (this.ipv4.trim() === '') {
            this.$refs.error.showAxios(error('Empty IP'))
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
    },
    toggleIpAutoDetect () {
      this.ipAutoDetect = !this.ipAutoDetect
    },
    toggleIpv4 () {
      this.ipv4Enabled = !this.ipv4Enabled
    },
    toggleIpv4Public () {
      this.ipv4Public = !this.ipv4Public
    },
    toggleIpv6 () {
      this.ipv6Enabled = !this.ipv6Enabled
    }
  }
}
</script>
<style>
@import '../style/site.css';
@import '../style/material-icons.css';

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
</style>
