<template>
  <div class="sc-page">
    <div class="sc-card sc-card-wide" id="block1">
      <router-link to="/storage" class="sc-link back-link" data-testid="storage-back">{{ $t('storage.back') }}</router-link>
      <h1 class="sc-title" data-testid="disk-title">{{ disk ? disk.name + ' - ' + disk.size : device }}</h1>

      <p v-if="loaded && !disk" class="sc-lead" data-testid="disk-missing">{{ $t('storage.diskMissing') }}</p>

      <dl v-if="disk" class="disk-details">
        <dt>{{ $t('storage.device') }}</dt>
        <dd data-testid="disk-device">{{ disk.device }}</dd>
        <template v-if="disk.serial">
          <dt>{{ $t('storage.serial') }}</dt>
          <dd data-testid="disk-serial">{{ disk.serial }}</dd>
        </template>
        <template v-if="port">
          <dt>{{ $t('storage.port') }}</dt>
          <dd data-testid="disk-port">{{ port }}</dd>
        </template>
        <template v-if="disk.devid">
          <dt>{{ $t('storage.devid') }}</dt>
          <dd data-testid="disk-devid">{{ disk.devid }}</dd>
        </template>
        <template v-if="disk.raid">
          <dt>{{ $t('storage.profile') }}</dt>
          <dd data-testid="disk-profile">{{ $t('storage.profileValue', { data: disk.raid, metadata: disk.raid_metadata || '-' }) }}</dd>
        </template>
        <template v-if="disk.errors">
          <dt>{{ $t('storage.errorCounters') }}</dt>
          <dd data-testid="disk-errors" :class="{ 'disk-bad': disk.has_errors }">{{ $t('storage.counters', disk.errors) }}</dd>
        </template>
        <template v-if="health">
          <dt>{{ $t('storage.smart') }}</dt>
          <dd data-testid="disk-smart" :class="{ 'disk-bad': smartBad(health.smart) }">{{ smartText(t, health.smart) }}</dd>
          <template v-if="health.smart.available">
            <dt>{{ $t('storage.selfTest') }}</dt>
            <dd class="disk-selftest">
              <button class="sc-btn sc-btn-ghost sc-btn-small" type="button" data-testid="disk-smart-test"
                      :disabled="busy || health.smart.self_test.running"
                      @click="startSmartTest">{{ $t('storage.smartTestStart') }}</button>
              <span data-testid="disk-selftest">{{ selfTestText(t, health.smart.self_test) }}</span>
            </dd>
          </template>
          <template v-if="health.smart.log && health.smart.log.length > 0">
            <dt>{{ $t('storage.selfTestLog') }}</dt>
            <dd>
              <ul class="selftest-log" data-testid="disk-selftest-log">
                <li v-for="(entry, index) in health.smart.log" :key="index" :class="{ 'disk-bad': entry.passed === false }">{{ logEntryText(t, entry) }}</li>
              </ul>
            </dd>
          </template>
          <dt>{{ $t('storage.linkEvents') }}</dt>
          <dd data-testid="disk-link" :class="{ 'disk-bad': health.link.count > 0 }">{{ eventsText(t, health.link) }}</dd>
          <dt>{{ $t('storage.ioErrors') }}</dt>
          <dd data-testid="disk-io" :class="{ 'disk-bad': health.io.count > 0 }">{{ eventsText(t, health.io) }}</dd>
        </template>
      </dl>
      <p v-if="busy" class="sc-muted" data-testid="checks-busy">{{ $t('storage.checksBusy') }}</p>
    </div>
  </div>
  <Error ref="error"/>
</template>

<script>
import axios from 'axios'
import Error from '../components/Error.vue'
import { eventsText, healthOf, logEntryText, portOf, selfTestText, smartBad, smartText } from '../js/storage'

export default {
  name: 'StorageDisk',
  components: {
    Error
  },
  data () {
    return {
      disk: undefined,
      health: undefined,
      busy: false,
      loaded: false,
      timer: undefined
    }
  },
  computed: {
    device () {
      return this.$route.query.device
    },
    port () {
      return this.disk ? portOf(this.disk, this.health) : ''
    }
  },
  mounted () {
    this.reload()
  },
  unmounted () {
    clearTimeout(this.timer)
  },
  methods: {
    t (...args) {
      return this.$t(...args)
    },
    smartBad,
    smartText,
    selfTestText,
    eventsText,
    logEntryText,
    reload () {
      clearTimeout(this.timer)
      return Promise.all([axios.get('/rest/storage/disks'), axios.get('/rest/storage/health')])
        .then(([disks, health]) => {
          this.disk = (disks.data.data || []).find(d => d.device === this.device)
          this.health = this.disk ? healthOf(health.data.data, this.disk) : undefined
          this.busy = health.data.data.busy
          this.loaded = true
          if (this.health && this.health.smart.self_test && this.health.smart.self_test.running) {
            this.timer = setTimeout(this.reload, 10000)
          }
        })
        .catch(err => this.$refs.error.showAxios(err))
    },
    startSmartTest () {
      return axios.post('/rest/storage/check/smart', { device: this.disk.device })
        .then(() => this.reload())
        .catch(err => this.$refs.error.showAxios(err))
    }
  }
}
</script>
<style scoped>
.back-link {
  display: inline-block;
  margin-bottom: 8px;
}
.disk-details {
  display: grid;
  grid-template-columns: max-content minmax(0, 1fr);
  gap: 8px 16px;
  margin: 0;
}
.disk-details dt {
  color: var(--sc-muted);
}
.disk-details dd {
  margin: 0;
  overflow-wrap: anywhere;
}
.disk-bad {
  color: var(--sc-danger);
}
.selftest-log {
  margin: 0;
  padding: 0;
  list-style: none;
}
.disk-selftest {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}
.disk-selftest button {
  flex-shrink: 0;
}
.sc-btn-small {
  padding: 4px 10px;
  font-size: 13px;
}
</style>
