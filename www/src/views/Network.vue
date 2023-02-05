<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12" id="block1">
        <h1>Network</h1>
        <div class="row-no-gutters settingsblock">
          <div>
            <div class="setline" style="display: flex; flex-wrap: wrap; justify-content: center;">
              <span v-if="interfaces === undefined || interfaces.length === 0" class="span">No networks found</span>
              <div style="font-size: 18px; margin: 10px; min-width: 300px;" v-for="(iface, index) in interfaces" :key="index">
                <div>
                <span style="font-weight: bold">Interface: {{ iface.name }}</span>
                </div>
                <div v-for="(name, index) in iface.addresses" :key="index">
                  {{ name }}
                </div>

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
    axios.get('/rest/network/interfaces')
      .then(response => {
        this.interfaces = response.data.data
      })
      .catch(err => this.$refs.error.showAxios(err))
  }
}
</script>
<style>
@import '../style/site.css';
@import 'material-icons/iconfont/material-icons.css';
</style>
