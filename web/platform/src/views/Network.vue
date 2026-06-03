<template>
  <div class="sc-page">
    <div class="sc-card" id="block1">
      <h1 class="sc-title">{{ $t('network.title') }}</h1>
      <span v-if="interfaces === undefined || interfaces.length === 0" class="sc-lead">{{ $t('network.noNetworks') }}</span>
      <div v-else class="net-grid">
        <div class="net-card" v-for="(iface, index) in interfaces" :key="index">
          <div class="net-name">{{ $t('network.interface') }} {{ iface.name }}</div>
          <div class="net-addr" v-for="(name, idx) in iface.addresses" :key="idx">{{ name }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'Network',
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

<style scoped>
.net-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 16px;
}
.net-card {
  background: var(--sc-field-bg);
  border: 1px solid var(--sc-border);
  border-radius: 14px;
  padding: 16px 18px;
  text-align: left;
}
.net-name {
  font-weight: 700;
  color: var(--sc-ink);
  margin-bottom: 8px;
}
.net-addr {
  font-size: 14px;
  color: var(--sc-muted);
  font-variant-numeric: tabular-nums;
  word-break: break-all;
}
</style>
