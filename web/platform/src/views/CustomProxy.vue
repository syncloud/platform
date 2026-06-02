<template>
  <div class="sc-page">
    <div class="sc-card" :style="{ visibility: visibility }">
      <h1 class="sc-title" data-testid="customproxy-title">{{ $t('customProxy.title') }}</h1>

      <h3>{{ $t('customProxy.addProxy') }}</h3>
      <div class="setline">
        <span class="proxy-warning">{{ $t('customProxy.warning') }}</span>
      </div>

      <div class="setline proxy-field">
        <label class="span proxy-label" for="proxy_name">{{ $t('customProxy.name') }}</label>
        <input class="proxy-input sc-input" id="proxy_name" type="text" v-model="newName" :placeholder="$t('customProxy.namePlaceholder')">
      </div>
      <div class="setline">
        <span class="proxy-hint">{{ $t('customProxy.urlPreview') }} {{ newName || $t('customProxy.nameFallback') }}.{{ domain }}</span>
      </div>

      <div class="setline proxy-field">
        <label class="span proxy-label" for="proxy_host">{{ $t('customProxy.host') }}</label>
        <input class="proxy-input sc-input" id="proxy_host" type="text" v-model="newHost" :placeholder="$t('customProxy.hostPlaceholder')">
      </div>

      <div class="setline proxy-field">
        <label class="span proxy-label" for="proxy_port">{{ $t('customProxy.port') }}</label>
        <input class="proxy-input sc-input" id="proxy_port" type="number" v-model.number="newPort" :placeholder="$t('customProxy.portPlaceholder')">
      </div>

      <div class="setline proxy-toggle">
        <label class="span proxy-label" for="proxy_https">{{ $t('customProxy.https') }}</label>
        <input id="proxy_https" type="checkbox" v-model="newHttps" class="proxy-check">
      </div>

      <div class="setline proxy-toggle">
        <label class="span proxy-label" for="proxy_authelia">{{ $t('customProxy.authelia') }}:</label>
        <input id="proxy_authelia" data-testid="proxy-authelia" type="checkbox" v-model="newAuthelia" class="proxy-check">
        <button type="button" class="proxy-help" @click="autheliaInfoVisible = true" aria-label="?">
          <i class="fa fa-question-circle fa-lg"></i>
        </button>
      </div>

      <div class="sc-actions">
        <button class="sc-btn sc-btn-success" id="btn_add" type="button" @click="add">{{ $t('customProxy.add') }}</button>
      </div>

      <h3>{{ $t('customProxy.proxies') }}</h3>

      <div v-if="proxies.length === 0" class="setline">
        <span>{{ $t('customProxy.noProxies') }}</span>
      </div>

      <div v-for="proxy in proxies" :key="proxy.name" class="proxy-entry">
        <a class="proxy-link sc-link" :href="'https://' + proxy.name + '.' + domain" target="_blank" :data-testid="'proxy-link-' + proxy.name">{{ proxy.name }}</a>
        <span class="proxy-url">{{ proxy.https ? 'https' : 'http' }}://{{ proxy.host }}:{{ proxy.port }}</span>
        <span v-if="proxy.authelia" class="proxy-badge"
              :data-testid="'proxy-row-' + proxy.name + '-authelia'">{{ $t('customProxy.authelia') }}</span>
        <button class="btn_remove sc-btn sc-btn-danger proxy-remove" type="button" :id="'btn_remove_' + proxy.name"
                @click="remove(proxy.name)">{{ $t('customProxy.remove') }}
        </button>
      </div>
    </div>
  </div>

  <Dialog :visible="autheliaInfoVisible" @cancel="autheliaInfoVisible = false" :confirm-enabled="false" :cancel-text="$t('common.close')">
    <template v-slot:title>{{ $t('customProxy.authelia') }}</template>
    <template v-slot:text>{{ $t('customProxy.autheliaHint') }}</template>
  </Dialog>

  <Error ref="error"/>

</template>

<script>
import Error from '../components/Error.vue'
import Dialog from '../components/Dialog.vue'
import * as Common from '../js/common.js'
import axios from 'axios'
import { ElLoading } from 'element-plus'

export default {
  name: 'CustomProxy',
  props: {
    checkUserSession: Function,
    activated: Boolean
  },
  data () {
    return {
      proxies: [],
      newName: '',
      newHost: '',
      newPort: null,
      newHttps: false,
      newAuthelia: false,
      visibility: 'hidden',
      loading: undefined,
      autheliaInfoVisible: false
    }
  },
  components: {
    Error,
    Dialog
  },
  computed: {
    domain () {
      return window.location.hostname
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
    reload () {
      const onError = (err) => {
        this.$refs.error.showAxios(err)
        this.progressHide()
      }
      axios.get('/rest/proxy_custom/list')
        .then(resp => Common.checkForServiceError(resp.data, () => {
          this.proxies = resp.data.data || []
          this.progressHide()
        }, onError))
        .catch(onError)
    },
    add () {
      this.progressShow()
      const onError = (err) => {
        this.$refs.error.showAxios(err)
        this.progressHide()
      }
      axios.post('/rest/proxy_custom/add', {
        name: this.newName,
        host: this.newHost,
        port: this.newPort,
        https: this.newHttps,
        authelia: this.newAuthelia
      })
        .then(resp => Common.checkForServiceError(resp.data, () => {
          this.newName = ''
          this.newHost = ''
          this.newPort = null
          this.newHttps = false
          this.newAuthelia = false
          this.reload()
        }, onError))
        .catch(onError)
    },
    remove (name) {
      this.progressShow()
      const onError = (err) => {
        this.$refs.error.showAxios(err)
        this.progressHide()
      }
      axios.post('/rest/proxy_custom/remove', { name: name })
        .then(resp => Common.checkForServiceError(resp.data, () => {
          this.reload()
        }, onError))
        .catch(onError)
    }
  }
}
</script>
<style scoped>
.proxy-label {
  width: 132px;
  min-width: 132px;
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
}

.proxy-field {
  display: flex;
  align-items: center;
}

.proxy-toggle {
  display: flex;
  align-items: center;
}

.proxy-input {
  width: 200px;
  height: 38px;
}

.proxy-entry {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 0;
  border-bottom: 1px solid #eef3f9;
}
.proxy-link { font-weight: 600; flex: 0 0 auto; }
.proxy-url {
  flex: 1;
  color: var(--sc-muted);
  word-break: break-all;
  font-variant-numeric: tabular-nums;
}
.proxy-remove { min-width: 80px; padding: 8px 14px; flex: 0 0 auto; }

.proxy-check {
  margin: 0;
  width: 18px;
  height: 18px;
}

.proxy-help {
  margin-left: 8px;
  background: transparent;
  border: none;
  cursor: pointer;
  color: var(--sc-faint);
  padding: 0;
  line-height: 1;
}
.proxy-help:hover { color: var(--sc-primary); }

.proxy-warning {
  color: var(--sc-danger);
  display: block;
  max-width: 500px;
  word-wrap: break-word;
}

.proxy-hint {
  margin-left: 10px;
  color: var(--sc-faint);
  display: inline-flex;
  align-items: center;
  white-space: nowrap;
}

.proxy-badge {
  background-color: var(--sc-success);
  color: white;
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 10px;
}

@media (max-width: 600px) {
  .proxy-field {
    flex-direction: column;
    align-items: stretch;
    gap: 6px;
  }
  .proxy-field .proxy-label { width: auto; min-width: 0; font-weight: 600; }
  .proxy-input { width: 100%; }

  .proxy-entry {
    flex-wrap: wrap;
    align-items: center;
    gap: 8px;
    padding: 12px;
    border: 1px solid var(--sc-border);
    border-radius: 12px;
    background: var(--sc-field-bg);
    margin-bottom: 10px;
  }
  .proxy-link { width: 100%; }
  .proxy-url { flex: 1 1 100%; width: 100%; }
  .proxy-remove { margin-left: auto; }
}
</style>
