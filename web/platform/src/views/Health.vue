<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12" id="block1" data-testid="health-page">
        <h1>{{ $t('health.title') }}</h1>

        <div class="row-no-gutters settingsblock">
          <div class="col2">
            <div class="setline">
              <h3>{{ $t('health.cpu') }} — {{ cpuPct.toFixed(0) }}%</h3>
              <el-progress :percentage="cpuPct" :stroke-width="14" :show-text="false" :status="pctStatus(cpuPct)" data-testid="health-cpu-bar" />
            </div>

            <div class="setline">
              <h3>{{ $t('health.memory') }} — {{ memUsedMb }} / {{ memTotalMb }} MB</h3>
              <el-progress :percentage="memPct" :stroke-width="14" :show-text="false" :status="pctStatus(memPct)" data-testid="health-mem-bar" />
              <div class="muted">{{ memAvailMb }} MB {{ $t('health.available') }}</div>
            </div>

            <div class="setline" v-if="swapTotalMb > 0">
              <h3>{{ $t('health.swap') }} — {{ swapUsedMb }} / {{ swapTotalMb }} MB</h3>
              <el-progress :percentage="swapPct" :stroke-width="14" :show-text="false" :status="pctStatus(swapPct)" data-testid="health-swap-bar" />
            </div>
          </div>

          <div class="col2">
            <div class="setline">
              <h3>{{ $t('health.disks') }}</h3>
              <div v-for="m in metrics.mounts" :key="m.path" class="setline-sub" :data-testid="'health-mount-' + m.path">
                <span class="span">{{ m.path }}</span>
                <span class="muted">{{ mb(m.used_kb) }} / {{ mb(m.total_kb) }} MB</span>
                <el-progress :percentage="mountPct(m)" :stroke-width="8" :show-text="false" :status="pctStatus(mountPct(m))" />
              </div>
              <div v-for="d in diskRates" :key="d.name" class="setline-sub" :data-testid="'health-disk-' + d.name">
                <span class="span">{{ d.name }}</span>
                <span class="muted">{{ $t('health.ioRead') }} {{ d.readKBs }} KB/s · {{ $t('health.ioWrite') }} {{ d.writeKBs }} KB/s</span>
              </div>
            </div>

            <div class="setline">
              <h3>{{ $t('health.network') }}</h3>
              <div v-for="n in netRates" :key="n.name" class="setline-sub" :data-testid="'health-net-' + n.name">
                <span class="span">{{ n.name }}</span>
                <span class="muted">{{ $t('health.netRx') }} {{ n.rxKBs }} KB/s · {{ $t('health.netTx') }} {{ n.txKBs }} KB/s</span>
              </div>
            </div>
          </div>
        </div>

        <h2>{{ $t('health.events') }}</h2>
        <div v-if="events.length === 0" class="muted" data-testid="health-events-empty">{{ $t('health.noEvents') }}</div>
        <el-table v-else :data="events" data-testid="health-events-table">
          <el-table-column prop="time" :label="$t('health.colTime')" width="200">
            <template #default="scope">{{ fmtTime(scope.row.time) }}</template>
          </el-table-column>
          <el-table-column prop="kind" :label="$t('health.colKind')" width="240">
            <template #default="scope">{{ $t('health.kind' + kindCamel(scope.row.kind)) }}</template>
          </el-table-column>
          <el-table-column :label="$t('health.colDetails')">
            <template #default="scope">{{ fmtDetails(scope.row) }}</template>
          </el-table-column>
        </el-table>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'

const METRICS_INTERVAL_MS = 2000
const EVENTS_INTERVAL_MS = 10000

export default {
  name: 'Health',
  data () {
    return {
      metrics: { cpu: {}, memory: {}, disks: [], mounts: [], net: [] },
      prevMetrics: null,
      events: [],
      cpuPct: 0,
      diskRates: [],
      netRates: [],
      metricsTimer: null,
      eventsTimer: null
    }
  },
  computed: {
    memTotalMb () { return this.kbToMb(this.metrics.memory.total_kb) },
    memUsedMb () { return this.kbToMb((this.metrics.memory.total_kb || 0) - (this.metrics.memory.available_kb || 0)) },
    memAvailMb () { return this.kbToMb(this.metrics.memory.available_kb) },
    memPct () {
      const t = this.metrics.memory.total_kb || 0
      if (!t) return 0
      return ((t - (this.metrics.memory.available_kb || 0)) / t) * 100
    },
    swapTotalMb () { return this.kbToMb(this.metrics.memory.swap_total_kb) },
    swapUsedMb () { return this.kbToMb((this.metrics.memory.swap_total_kb || 0) - (this.metrics.memory.swap_free_kb || 0)) },
    swapPct () {
      const t = this.metrics.memory.swap_total_kb || 0
      if (!t) return 0
      return ((t - (this.metrics.memory.swap_free_kb || 0)) / t) * 100
    }
  },
  methods: {
    kbToMb (kb) { return Math.round((kb || 0) / 1024) },
    mb (kb) { return Math.round((kb || 0) / 1024) },
    mountPct (m) {
      if (!m.total_kb) return 0
      return (m.used_kb / m.total_kb) * 100
    },
    pctStatus (pct) {
      if (pct >= 90) return 'exception'
      if (pct >= 75) return 'warning'
      return 'success'
    },
    kindCamel (kind) {
      return kind.split('_').map(s => s.charAt(0).toUpperCase() + s.slice(1)).join('')
    },
    fmtTime (iso) {
      try { return new Date(iso).toLocaleString() } catch { return iso }
    },
    fmtDetails (e) {
      const parts = []
      if (e.comm) parts.push(e.comm)
      if (e.pid) parts.push('pid=' + e.pid)
      if (e.rss_kb) parts.push('rss=' + this.kbToMb(e.rss_kb) + 'MB')
      if (e.avail_ratio) parts.push('avail=' + (e.avail_ratio * 100).toFixed(1) + '%')
      if (e.psi_avg10) parts.push('psi=' + e.psi_avg10.toFixed(1))
      if (e.path) parts.push(e.path)
      if (e.size_bytes) parts.push(Math.round(e.size_bytes / (1024 * 1024)) + ' MB')
      return parts.join(' · ')
    },
    fetchMetrics () {
      axios.get('/rest/settings/health/metrics')
        .then(resp => {
          const next = resp.data.data
          this.computeDeltas(this.prevMetrics, next)
          this.prevMetrics = this.metrics.cpu && this.metrics.cpu.user ? this.metrics : null
          this.metrics = next
        })
        .catch(() => {})
    },
    computeDeltas (prev, next) {
      if (!prev) {
        this.cpuPct = 0
        this.diskRates = (next.disks || []).map(d => ({ name: d.name, readKBs: 0, writeKBs: 0 }))
        this.netRates = (next.net || []).map(n => ({ name: n.name, rxKBs: 0, txKBs: 0 }))
        return
      }
      const totalDelta = (next.cpu.user + next.cpu.nice + next.cpu.system + next.cpu.idle + next.cpu.iowait + next.cpu.irq + next.cpu.softirq + next.cpu.steal) -
                        (prev.cpu.user + prev.cpu.nice + prev.cpu.system + prev.cpu.idle + prev.cpu.iowait + prev.cpu.irq + prev.cpu.softirq + prev.cpu.steal)
      const idleDelta = (next.cpu.idle + next.cpu.iowait) - (prev.cpu.idle + prev.cpu.iowait)
      this.cpuPct = totalDelta > 0 ? Math.max(0, Math.min(100, ((totalDelta - idleDelta) / totalDelta) * 100)) : 0

      const secs = METRICS_INTERVAL_MS / 1000
      const prevDisks = {}
      ;(prev.disks || []).forEach(d => { prevDisks[d.name] = d })
      this.diskRates = (next.disks || []).map(d => {
        const p = prevDisks[d.name]
        if (!p) return { name: d.name, readKBs: 0, writeKBs: 0 }
        return {
          name: d.name,
          readKBs: Math.round(((d.sectors_read - p.sectors_read) * 512 / 1024) / secs),
          writeKBs: Math.round(((d.sectors_written - p.sectors_written) * 512 / 1024) / secs)
        }
      })

      const prevNet = {}
      ;(prev.net || []).forEach(n => { prevNet[n.name] = n })
      this.netRates = (next.net || []).map(n => {
        const p = prevNet[n.name]
        if (!p) return { name: n.name, rxKBs: 0, txKBs: 0 }
        return {
          name: n.name,
          rxKBs: Math.round((n.rx_bytes - p.rx_bytes) / 1024 / secs),
          txKBs: Math.round((n.tx_bytes - p.tx_bytes) / 1024 / secs)
        }
      })
    },
    fetchEvents () {
      axios.get('/rest/settings/health/events?limit=100')
        .then(resp => { this.events = resp.data.data || [] })
        .catch(() => {})
    }
  },
  mounted () {
    this.fetchMetrics()
    this.fetchEvents()
    this.metricsTimer = setInterval(this.fetchMetrics, METRICS_INTERVAL_MS)
    this.eventsTimer = setInterval(this.fetchEvents, EVENTS_INTERVAL_MS)
  },
  beforeUnmount () {
    if (this.metricsTimer) clearInterval(this.metricsTimer)
    if (this.eventsTimer) clearInterval(this.eventsTimer)
  }
}
</script>

<style scoped>
.setline-sub {
  display: block;
  padding: 4px 0;
}
.setline-sub .span {
  display: inline-block;
  min-width: 110px;
  font-weight: 500;
}
.muted {
  color: #888;
  font-size: 0.9em;
  margin-left: 8px;
}
h2 {
  margin-top: 24px;
}
</style>
