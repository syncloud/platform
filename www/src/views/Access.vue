<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12">
        <h1>Access</h1>
        <div class="row-no-gutters settingsblock">
          <div class="col2" id="wrapper" :style="{ visibility: visibility }">
            <h3>Domain name</h3>
            <div class="setline">
              <div class="spandiv" id="external_mode">
                <span class="span">External Access:</span>
                <Switch
                  id="tgl_external"
                  :checked="externalAccess"
                  @toggle="toggleExternalAccess"
                  on-label="ON"
                  off-label="OFF"
                />
              </div>
              <button type=button @click="showExternalAccessInfo" class="control" style=" background:transparent;">
                <i class='fa fa-question-circle fa-lg'></i>
              </button>
            </div>
            <div id="external_block">
              <div class="setline">
                <div class="spandiv">
                  <span class="span">Auto detect IP:</span>
                  <Switch
                    id="tgl_ip_autodetect"
                    :checked="ipAutoDetect"
                    @toggle="toggleIpAutoDetect"
                    on-label="ON"
                    off-label="OFF"
                  />
                </div>
              </div>

              <div class="setline" id="public_ip_block">
                <label class="span" for="public_ip" style="font-weight: 300">Public IP:</label>
                <input id="public_ip" type="text"
                       style="width: 150px; height: 30px; padding: 0 10px 0 10px"
                       :disabled="ipAutoDetect" v-model="publicIp">
              </div>

              <div class="setline">
                <h3>Router external ports</h3>
              </div>

              <div class="setline">
                <div class="spandiv">
                  <span class="span">Auto (UPnP):</span>
                  <Switch
                    id="tgl_upnp"
                    :checked="upnp"
                    @toggle="toggleUpnp"
                    on-label="ON"
                    off-label="OFF"
                  />
                  <button id="upnp_warning" type=button @click="showUpnpDisabledWarning"
                          class="control" style="background:transparent;">
                    <i class='fa fa-exclamation-circle fa-lg' style='color: red;'></i>
                  </button>
                </div>
              </div>
              <div id="ports_block">
                <div class="setline">
                  <span class="span">External certificate</span>
                  <span style='white-space: nowrap;'>
                      <label for="certificate_port" class="span" style="font-weight: 300">HTTP port 80:</label>
                      <input class="span" id="certificate_port" type="number"
                             style="width: 100px; height: 30px; padding: 0 10px 0 10px"
                             v-model.number="certificatePort"
                      >
                        <button id="certificate_port_warning" type=button @click="showCertificatePortWarning"
                                class="control" style="background:transparent;">
                          <i class='fa fa-exclamation-circle fa-lg' style='color: red;'></i>
                        </button>
                  </span>
                </div>

                <div class="setline">
                  <div class="spandiv">
                    <span class="span">External access</span>
                    <span style='white-space: nowrap;'>
                          <label for="access_port" class="span" style="font-weight: 300">HTTPS port 443:</label>
                          <input class="span" id="access_port" type="number"
                                 style="width: 100px; height: 30px; padding: 0 10px 0 10px"
                                 v-model.number="accessPort"
                          />
                            <button id="access_port_warning" type=button @click="showAccessPortWarning"
                                    class="control" style="background:transparent;">
                                <i class='fa fa-exclamation-circle fa-lg' style='color: red;'></i>
                            </button>
                    </span>
                  </div>
                </div>
              </div>
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

  <Dialog ref="certificate_port_warning">
    <template v-slot:title>Certificate warning</template>
    <template v-slot:text>
      Certificate can only be obtained if external port is 80.
      It is only used for certificate validation purposes.
    </template>
  </Dialog>

  <Dialog ref="access_port_warning">
    <template v-slot:title>Access port warning</template>
    <template v-slot:text>
      Access port is not default 443.
      You may not be able to access your device from networks with strict firewalls allowing only port 443.
    </template>
  </Dialog>

  <Dialog ref="upnp_disabled_warning">
    <template v-slot:title>UPnP is not available</template>
    <template v-slot:text>Your router does not have port mapping feature enabled.</template>
  </Dialog>

  <Dialog ref="external_access_info">
    <template v-slot:title>External access</template>
    <template v-slot:text>
      External access is used for two things:
      <br><br>
      1. Allow syncloud.it device DNS to use public IP. By default it uses internal IP to allow DNS based app
      access inside local network.
      <br>
      2. Maintain UPnP (auto) port mappings.
      <br><br>
      External access will not change device itself ports, they are always 80 and 443.
      <br><br>
      Syncloud.it DNS service verifies open ports (internet accessibility) before enabling external access for
      convenience.
      <br><br>
      In case of custom domain name external access page should not be used.
    </template>
  </Dialog>

  <Error ref="error"/>

</template>

<script>
import $ from 'jquery'
import Error from '@/components/Error'
import Dialog from '@/components/Dialog'
import 'bootstrap'
import 'bootstrap-switch'
import * as Common from '../js/common.js'
import axios from 'axios'
import Switch from '@/components/Switch'
import 'gasparesganga-jquery-loading-overlay'

function isValidPort (port) {
  return Number.isNaN(port) || port < 1 || port > 65535
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
    onLogin: Function,
    onLogout: Function
  },
  data () {
    return {
      interfaces: undefined,
      externalAccess: false,
      ipAutoDetect: false,
      publicIp: 0,
      upnp: false,
      upnpAvailable: false,
      accessPort: 0,
      certificatePort: 0,
      visibility: 'hidden'
    }
  },
  components: {
    Error,
    Dialog,
    Switch
  },
  watch: {
    externalAccess (val) {
      this.initExternalAccess(val)
    },
    ipAutoDetect (val) {
      this.initIpAutoDetect(val)
    },
    upnp (val) {
      this.initUpnp(val)
    },
    upnpAvailable (val) {
      if (val) {
        $('#upnp_warning').hide('slow')
      } else {
        $('#upnp_warning').show('slow')
      }
    },
    certificatePort (val) {
      if (val !== 80) {
        $('#certificate_port_warning').show('slow')
      } else {
        $('#certificate_port_warning').hide('slow')
      }
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
      $('#wrapper').LoadingOverlay('show', { background: 'rgb(0,0,0,0)' })
    },
    progressHide () {
      this.visibility = 'visible'
      $('#wrapper').LoadingOverlay('hide')
    },
    initExternalAccess (val) {
      if (val) {
        $('#external_block').show('slow')
      } else {
        $('#external_block').hide('slow')
      }
    },
    initIpAutoDetect (val) {
      if (val) {
        $('#public_ip_block').hide('slow')
      } else {
        $('#public_ip_block').show('slow')
      }
    },
    initUpnp (val) {
      if (val) {
        $('#ports_block').hide('slow')
      } else {
        $('#ports_block').show('slow')
      }
    },
    showUpnpDisabledWarning () {
      this.$refs.upnp_disabled_warning.show()
    },
    showAccessPortWarning () {
      this.$refs.access_port_warning.show()
    },
    showCertificatePortWarning () {
      this.$refs.certificate_port_warning.show()
    },
    showExternalAccessInfo () {
      this.$refs.external_access_info.show()
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
        that.externalAccess = accessData.external_access
        that.initExternalAccess(that.externalAccess)
        if ('public_ip' in accessData) {
          that.ipAutoDetect = false
          that.publicIp = accessData.public_ip
        } else {
          that.ipAutoDetect = true
        }
        this.initIpAutoDetect(that.ipAutoDetect)
        that.upnp = accessData.upnp_enabled
        this.initUpnp(this.upnp)
        that.upnpAvailable = accessData.upnp_available
        this.reloadPortMappings()
      }
      axios.get('/rest/access/access')
        .then(resp => Common.checkForServiceError(resp.data.data, () => onComplete(resp.data.data), onError))
        .catch(onError)
    },
    reloadPortMappings () {
      axios.get('/rest/access/port_mappings')
        .then(resp => {
          const certificatePortMapping = resp.data.port_mappings.find(function (mapping) {
            return mapping.local_port === 80
          })
          if (certificatePortMapping) {
            this.certificatePort = certificatePortMapping.external_port
          }
          const accessPortMapping = resp.data.port_mappings.find(function (mapping) {
            return mapping.local_port === 443
          })
          if (accessPortMapping) {
            this.accessPort = accessPortMapping.external_port
          }
          this.progressHide()
        })
        .catch(err => {
          error.showAxios(err)
          this.progressHide()
        })
    },
    save (event) {
      this.progressShow()

      event.preventDefault()
      const that = this
      const requestData = {
        external_access: this.externalAccess,
        upnp_enabled: false,
        certificate_port: 0,
        access_port: 0
      }
      if (this.externalAccess) {
        requestData.upnp_enabled = this.upnp
        if (!this.upnp) {
          if (isValidPort(this.certificatePort)) {
            this.$refs.error.showAxios(error('certificate port (' + this.certificatePort + ') has to be between 1 and 65535'))
            this.progressHide()
            return
          }
          requestData.certificate_port = this.certificatePort
          if (isValidPort(this.accessPort)) {
            this.$refs.error.showAxios(error('access port (' + this.accessPort + ') has to be between 1 and 65535'))
            this.progressHide()
            return
          }
          requestData.access_port = this.accessPort
        }
        if (!this.ipAutoDetect) {
          requestData.public_ip = this.publicIp
        }
      }

      const onError = (err) => {
        that.$refs.error.showAxios(err)
        this.progressHide()
      }
      axios.post('/rest/access/set_access', requestData)
        .then(response => Common.checkForServiceError(response.data, this.reload, onError))
        .catch(onError)
    },
    toggleExternalAccess () {
      this.externalAccess = !this.externalAccess
    },
    toggleIpAutoDetect () {
      this.ipAutoDetect = !this.ipAutoDetect
    },
    toggleUpnp () {
      this.upnp = !this.upnp
    }
  }
}
</script>
<style>
@import '../style/site.css';
@import '../style/material-icons.css';
</style>
