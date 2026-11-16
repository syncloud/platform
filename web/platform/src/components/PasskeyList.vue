<template>
  <div class="passkeys" data-testid="passkeys">
    <div v-if="loaded && passkeys.length === 0" class="passkeys-empty" data-testid="passkeys-empty">{{ $t('passkeys.none') }}</div>
    <div v-for="passkey in passkeys" :key="passkey.id" class="passkey-row" :data-testid="'passkey-' + passkey.id">
      <i class="material-icons passkey-icon">key</i>
      <div class="passkey-info">
        <div class="passkey-name" :data-testid="'passkey-name-' + passkey.id">{{ passkey.description }}</div>
        <div class="passkey-meta">
          {{ $t('passkeys.added', { date: format(passkey.created) }) }} ·
          {{ passkey.last_used ? $t('passkeys.lastUsed', { date: format(passkey.last_used) }) : $t('passkeys.neverUsed') }}
        </div>
      </div>
      <button class="sc-btn sc-btn-danger" type="button" :data-testid="'passkey-remove-' + passkey.id" @click="confirmRemove(passkey)">
        {{ $t('passkeys.remove') }}
      </button>
    </div>
    <div v-if="passkeys.length > 1" class="passkeys-actions">
      <button class="sc-btn sc-btn-danger" type="button" data-testid="passkeys-remove-all" @click="confirmRemoveAll">
        {{ $t('passkeys.removeAll') }}
      </button>
    </div>
  </div>

  <Dialog :visible="pending !== undefined" @cancel="pending = undefined" @confirm="remove">
    <template v-slot:title>{{ pending === 'all' ? $t('passkeys.removeAll') : $t('passkeys.remove') }}</template>
    <template v-slot:text>{{ pending === 'all' ? $t('passkeys.confirmRemoveAll') : $t('passkeys.confirmRemove', { name: pending && pending.description }) }}</template>
  </Dialog>

  <Error ref="error"/>
</template>

<script>
import axios from 'axios'
import Error from './Error.vue'
import Dialog from './Dialog.vue'
import * as Common from '../js/common.js'

export default {
  name: 'PasskeyList',
  props: {
    username: {
      type: String,
      default: ''
    }
  },
  components: {
    Error,
    Dialog
  },
  data () {
    return {
      passkeys: [],
      loaded: false,
      pending: undefined
    }
  },
  computed: {
    base () {
      return this.username ? '/rest/users/passkeys' : '/rest/passkeys'
    }
  },
  mounted () {
    this.reload()
  },
  methods: {
    format (date) {
      return new Date(date).toLocaleDateString()
    },
    async reload () {
      try {
        const params = this.username ? { username: this.username } : {}
        const resp = await axios.get(this.base, { params })
        this.passkeys = resp.data.data || []
        this.loaded = true
      } catch (err) {
        this.$refs.error.showAxios(err)
      }
    },
    confirmRemove (passkey) {
      this.pending = passkey
    },
    confirmRemoveAll () {
      this.pending = 'all'
    },
    async remove () {
      const pending = this.pending
      this.pending = undefined
      try {
        if (pending === 'all') {
          await Common.post(this.base + '/reset', this.request())
        } else {
          await Common.post(this.base + '/remove', { ...this.request(), id: pending.id })
        }
        await this.reload()
      } catch (err) {
        this.$refs.error.showAxios(err)
      }
    },
    request () {
      return this.username ? { username: this.username } : {}
    }
  }
}
</script>

<style scoped>
.passkeys { display: flex; flex-direction: column; gap: 10px; }
.passkeys-empty { color: var(--sc-faint); }
.passkey-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  border: 1px solid var(--sc-border);
  border-radius: 10px;
  background: var(--sc-field-bg);
}
.passkey-icon { color: var(--sc-muted); }
.passkey-info { flex: 1; min-width: 0; }
.passkey-name { font-weight: 600; color: var(--sc-text); overflow-wrap: anywhere; }
.passkey-meta { font-size: 13px; color: var(--sc-faint); }
.passkeys-actions { display: flex; justify-content: flex-end; }
</style>
