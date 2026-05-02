<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12">
        <h1>{{ $t('customProxy.title') }}</h1>
        <div class="row-no-gutters settingsblock">
          <div class="col2" :style="{ visibility: visibility }">

            <div class="setline">
              <h3>{{ $t('customProxy.addProxy') }}</h3>
            </div>

            <div class="setline">
              <span class="proxy-warning">{{ $t('customProxy.warning') }}</span>
            </div>

            <div class="setline" style='display: flex'>
              <label class="span proxy-label" for="proxy_name">{{ $t('customProxy.name') }}</label>
              <input class="proxy-input" id="proxy_name" type="text" v-model="newName" :placeholder="$t('customProxy.namePlaceholder')">
            </div>
            <div class="setline">
              <span class="proxy-hint">{{ $t('customProxy.urlPreview') }} {{ newName || $t('customProxy.nameFallback') }}.{{ domain }}</span>
            </div>

            <div class="setline" style='display: flex'>
              <label class="span proxy-label" for="proxy_host">{{ $t('customProxy.host') }}</label>
              <input class="proxy-input" id="proxy_host" type="text" v-model="newHost" :placeholder="$t('customProxy.hostPlaceholder')">
            </div>

            <div class="setline" style='display: flex'>
              <label class="span proxy-label" for="proxy_port">{{ $t('customProxy.port') }}</label>
              <input class="proxy-input" id="proxy_port" type="number" v-model.number="newPort" :placeholder="$t('customProxy.portPlaceholder')">
            </div>

            <div class="setline" style='display: flex; align-items: center;'>
              <label class="span proxy-label" for="proxy_https">{{ $t('customProxy.https') }}</label>
              <input id="proxy_https" type="checkbox" v-model="newHttps">
            </div>

            <div class="setline" style='display: flex; align-items: center;'>
              <label class="span proxy-label" for="proxy_authelia">{{ $t('customProxy.authelia') }}</label>
              <input id="proxy_authelia" data-testid="proxy-authelia" type="checkbox" v-model="newAuthelia">
              <span class="proxy-hint">{{ $t('customProxy.autheliaHint') }}</span>
            </div>

            <div class="setline">
              <div class="spandiv">
                <button class="submit buttongreen control" id="btn_add" type="button"
                        style="width: 150px" @click="add">{{ $t('customProxy.add') }}
                </button>
              </div>
            </div>

            <div class="setline">
              <h3>{{ $t('customProxy.proxies') }}</h3>
            </div>

            <div v-if="proxies.length === 0" class="setline">
              <span>{{ $t('customProxy.noProxies') }}</span>
            </div>

            <div v-for="proxy in proxies" :key="proxy.name" class="setline proxy-entry" style='display: flex; align-items: center;'>
              <a class="proxy-label" :href="'https://' + proxy.name + '.' + domain" target="_blank">{{ proxy.name }}</a>
              <span style="flex: 1">{{ proxy.https ? 'https' : 'http' }}://{{ proxy.host }}:{{ proxy.port }}</span>
              <span v-if="proxy.authelia" class="proxy-badge"
                    :data-testid="'proxy-row-' + proxy.name + '-authelia'">{{ $t('customProxy.authelia') }}</span>
              <button class="submit control btn_remove" type="button" :id="'btn_remove_' + proxy.name"
                      style="width: 80px; background-color: #d9534f; color: white;" @click="remove(proxy.name)">{{ $t('customProxy.remove') }}
              </button>
            </div>

          </div>
        </div>
      </div>
    </div>
  </div>

  <Error ref="error"/>

</template>

<script>
import Error from '../components/Error.vue'
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
      loading: undefined
    }
  },
  components: {
    Error
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
<style>
@import '../style/site.css';

.proxy-label {
  min-width: 80px;
  display: inline-flex;
  align-items: center;
}

.proxy-input {
  width: 200px;
  height: 30px;
  padding: 0 10px;
}

.proxy-entry {
  padding: 5px 0;
}

.proxy-warning {
  margin-left: 10px;
  color: #d9534f;
  display: block;
  max-width: 500px;
  word-wrap: break-word;
}

.proxy-hint {
  margin-left: 10px;
  color: #999;
  display: inline-flex;
  align-items: center;
  white-space: nowrap;
}

.proxy-badge {
  background-color: #5cb85c;
  color: white;
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 10px;
  margin-right: 10px;
}
</style>
