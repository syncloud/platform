<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12" id="block1">
        <h1>Network</h1>
        <div class="row-no-gutters settingsblock">
          <div>
            <div class="setline" style="display: flex; flex-wrap: wrap; justify-content: center;">
              <span v-if="interfaces === undefined || interfaces.length === 0" class="span">No networks found</span>
              <div v-for="(iface, index) in interfaces" :key="index">
                <table style="font-size: 18px; margin: 10px; min-width: 300px;">
                  <tr style="text-align: left; font-weight: bold">
                    <td style="padding: 5px" colspan="2">Interface: {{ iface.name }}</td>
                  </tr>
                  <tr v-for="(address, index) in iface.ipv4" :key="index">
                    <td style="min-width: 20%;text-align: right; vertical-align: top; padding-right: 5px"
                        :rowspan="iface.ipv4.length" v-if="index === 0">IPv4:
                    </td>
                    <td style="text-align: left">{{ address.addr }}</td>
                  </tr>
                  <tr v-for="(address, index) in iface.ipv6" :key="index">
                    <td style="min-width: 20%;text-align: right; vertical-align: top; padding-right: 5px" :rowspan="iface.ipv6.length"
                        v-if="index === 0">IPv6:
                    </td>
                    <td style="text-align: left">{{ address.addr }}</td>
                  </tr>
                </table>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
import 'bootstrap'

export default {
  name: 'Network',
  props: {
    checkUserSession: Function,
    activated: Boolean
  },
  data () {
    return {
      interfaces: undefined
    }
  },
  mounted () {
    axios.get('/rest/access/network_interfaces')
      .then(response => {
        this.interfaces = response.data.data.interfaces
      })
      .catch(err => this.$refs.error.showAxios(err))
  }
}
</script>
<style>
@import '../style/site.css';
@import 'material-icons/iconfont/material-icons.css';
</style>
