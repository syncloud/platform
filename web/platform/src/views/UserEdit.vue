<template>
  <div class="sc-page">
    <div class="sc-card" :style="{ visibility: visibility }">
      <h1 class="sc-title" data-testid="user-edit-title">{{ isNew ? $t('users.add') : $t('users.edit') }}</h1>

      <div class="setline user-field">
        <label class="span user-label" for="user_username">{{ $t('users.username') }}</label>
        <input class="user-input sc-input" id="user_username" type="text" v-model="username"
               :disabled="!isNew" :placeholder="$t('users.usernamePlaceholder')">
      </div>

      <div class="setline user-field">
        <label class="span user-label" for="user_password">{{ $t('users.password') }}</label>
        <input class="user-input sc-input" id="user_password" type="password" v-model="password"
               :placeholder="isNew ? '' : $t('users.passwordKeep')">
      </div>
      <div class="setline password-rules" v-if="isNew || password" data-testid="password-rules">
        <ul class="pw-rules">
          <li v-for="rule in passwordRules" :key="rule.key" :data-testid="'pwrule-' + rule.key"
              class="pw-rule" :class="{ 'pw-ok': rule.ok }">
            <i class="material-icons pw-rule-icon">{{ rule.ok ? 'check_circle' : 'radio_button_unchecked' }}</i>
            <span>{{ $t(rule.label) }}</span>
          </li>
        </ul>
      </div>

      <div class="setline user-field">
        <label class="span user-label" for="user_email">{{ $t('users.email') }}</label>
        <input class="user-input sc-input" id="user_email" type="text" v-model="email"
               :placeholder="(username || $t('users.usernameFallback')) + '@' + domain">
      </div>

      <div class="setline user-toggle">
        <label class="span user-label" for="user_admin">{{ $t('users.admin') }}</label>
        <s-switch id="user_admin" data-testid="user-admin" v-model="admin" :disabled="adminLocked"/>
      </div>
      <div v-if="adminLocked" class="setline">
        <span class="user-hint" data-testid="user-admin-last">{{ $t('users.adminLast') }}</span>
      </div>

      <div class="user-groups-section">
        <h3 class="user-groups-heading">{{ $t('users.groups') }}</h3>
        <div class="user-chips">
          <button v-for="group in customGroups" :key="group.name" type="button"
                  class="user-chip" :class="{ 'is-on': selectedGroups.includes(group.name) }"
                  :data-testid="'user-group-' + group.name" @click="toggleGroup(group.name)">
            {{ group.name }}
            <i v-if="selectedGroups.includes(group.name)" class="material-icons user-chip-check">check</i>
          </button>
          <span class="user-newgroup">
            <input class="sc-input user-newgroup-input" id="new_group" type="text" v-model="newGroup"
                   :placeholder="$t('users.newGroup')" @keyup.enter="createGroup">
            <button type="button" class="sc-btn sc-btn-primary" id="btn_add_group" data-testid="group-create"
                    @click="createGroup">{{ $t('users.addGroup') }}</button>
          </span>
        </div>
      </div>

      <div class="sc-actions user-edit-actions">
        <button class="sc-btn sc-btn-success" id="btn_save" type="button" :disabled="saveDisabled" @click="save">{{ $t('users.save') }}</button>
        <button class="sc-btn" id="btn_cancel" type="button" @click="cancel">{{ $t('users.cancel') }}</button>
        <button v-if="!isNew" class="sc-btn sc-btn-danger user-delete" id="btn_delete" type="button" @click="deleteVisible = true">
          {{ $t('users.remove') }}
        </button>
      </div>
    </div>
  </div>

  <Dialog :visible="deleteVisible" @cancel="deleteVisible = false" @confirm="remove" :cancel-text="$t('users.cancel')">
    <template v-slot:title>{{ $t('users.remove') }}</template>
    <template v-slot:text>{{ $t('users.confirmDelete') }}</template>
  </Dialog>

  <Error ref="error"/>

</template>

<script>
import Error from '../components/Error.vue'
import Dialog from '../components/Dialog.vue'
import * as Common from '../js/common.js'
import axios from 'axios'
import Loading from '../util/loading'

export default {
  name: 'UserEdit',
  data () {
    return {
      username: '',
      password: '',
      email: '',
      admin: false,
      groups: [],
      selectedGroups: [],
      newGroup: '',
      originalEmail: '',
      originalAdmin: false,
      originalGroups: [],
      adminCount: 0,
      deleteVisible: false,
      visibility: 'hidden',
      loading: undefined
    }
  },
  components: {
    Error,
    Dialog
  },
  computed: {
    isNew () {
      return !this.$route.query.username
    },
    domain () {
      return window.location.hostname
    },
    customGroups () {
      return this.groups.filter(group => group.name !== 'syncloud')
    },
    adminLocked () {
      return !this.isNew && this.originalAdmin && this.adminCount <= 1
    },
    passwordRules () {
      const p = this.password
      return [
        { key: 'length', label: 'users.ruleLength', ok: p.length >= 8 },
        { key: 'letter', label: 'users.ruleLetter', ok: /[a-zA-Z]/.test(p) },
        { key: 'number', label: 'users.ruleNumber', ok: /[0-9]/.test(p) }
      ]
    },
    passwordValid () {
      return this.passwordRules.every(rule => rule.ok)
    },
    saveDisabled () {
      if (this.isNew) {
        return this.username.trim() === '' || !this.passwordValid
      }
      return this.password !== '' && !this.passwordValid
    }
  },
  mounted () {
    this.progressShow()
    this.reload()
  },
  methods: {
    progressShow () {
      this.loading = Loading.service({ lock: true, text: this.$t('common.loading'), background: 'rgba(0, 0, 0, 0.7)' })
    },
    progressHide () {
      this.visibility = 'visible'
      if (this.loading) this.loading.close()
    },
    onError (err) {
      this.$refs.error.showAxios(err)
      this.progressHide()
    },
    async reload () {
      try {
        const groupsResp = await axios.get('/rest/groups')
        this.groups = (groupsResp.data && groupsResp.data.data) || []
        if (!this.isNew) {
          const usersResp = await axios.get('/rest/users')
          const users = (usersResp.data && usersResp.data.data) || []
          this.adminCount = users.filter(u => u.admin).length
          const user = users.find(u => u.username === this.$route.query.username)
          if (user) {
            this.username = user.username
            this.email = user.email
            this.admin = user.admin
            this.selectedGroups = [...(user.groups || [])]
            this.originalEmail = user.email
            this.originalAdmin = user.admin
            this.originalGroups = [...(user.groups || [])]
          }
        }
        this.progressHide()
      } catch (err) {
        this.onError(err)
      }
    },
    toggleGroup (name) {
      if (this.selectedGroups.includes(name)) {
        this.selectedGroups = this.selectedGroups.filter(g => g !== name)
      } else {
        this.selectedGroups.push(name)
      }
    },
    async createGroup () {
      const name = this.newGroup.trim()
      if (!name) return
      this.progressShow()
      try {
        await Common.post('/rest/groups/add', { name: name })
        const groupsResp = await axios.get('/rest/groups')
        this.groups = (groupsResp.data && groupsResp.data.data) || []
        if (!this.selectedGroups.includes(name)) this.selectedGroups.push(name)
        this.newGroup = ''
        this.progressHide()
      } catch (err) {
        this.onError(err)
      }
    },
    async save () {
      this.progressShow()
      try {
        if (this.isNew) {
          await Common.post('/rest/users/add', {
            username: this.username,
            password: this.password,
            email: this.email,
            admin: this.admin
          })
          for (const group of this.selectedGroups) {
            await Common.post('/rest/groups/member/add', { group: group, username: this.username })
          }
        } else {
          if (this.password) {
            await Common.post('/rest/users/password', { username: this.username, password: this.password })
          }
          if (this.email !== this.originalEmail) {
            await Common.post('/rest/users/email', { username: this.username, email: this.email })
          }
          if (this.admin !== this.originalAdmin) {
            await Common.post('/rest/users/admin', { username: this.username, admin: this.admin })
          }
          for (const group of this.selectedGroups.filter(g => !this.originalGroups.includes(g))) {
            await Common.post('/rest/groups/member/add', { group: group, username: this.username })
          }
          for (const group of this.originalGroups.filter(g => !this.selectedGroups.includes(g))) {
            await Common.post('/rest/groups/member/remove', { group: group, username: this.username })
          }
        }
        this.progressHide()
        this.$router.push('/users')
      } catch (err) {
        this.onError(err)
      }
    },
    async remove () {
      this.deleteVisible = false
      this.progressShow()
      try {
        await Common.post('/rest/users/remove', { username: this.username })
        this.progressHide()
        this.$router.push('/users')
      } catch (err) {
        this.onError(err)
      }
    },
    cancel () {
      this.$router.push('/users')
    }
  }
}
</script>
<style scoped>
.user-label {
  width: 132px;
  min-width: 132px;
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
}
.user-field { display: flex; align-items: center; }
.user-toggle { display: flex; align-items: center; }
.user-input { width: 260px; height: 38px; }
.user-hint { color: var(--sc-faint); margin-left: 132px; }

.password-rules { margin-left: 132px; }
.pw-rules { list-style: none; margin: 0; padding: 0; }
.pw-rule {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--sc-faint);
  font-size: 13px;
  line-height: 1.8;
}
.pw-rule.pw-ok { color: var(--sc-success); }
.pw-rule-icon { font-size: 16px; }

.user-groups-section {
  margin-top: 18px;
  padding: 16px;
  border: 1px solid var(--sc-border);
  border-radius: 12px;
  background: var(--sc-field-bg);
}
.user-groups-heading {
  margin: 0 0 12px;
  font-size: 15px;
  color: var(--sc-text);
}
.user-chips { display: flex; flex-wrap: wrap; gap: 8px; align-items: center; }
.user-chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  border: 1px solid var(--sc-border);
  background: var(--sc-field-bg);
  color: var(--sc-muted);
  border-radius: 16px;
  padding: 5px 12px;
  cursor: pointer;
}
.user-chip.is-on { background: var(--sc-success); border-color: var(--sc-success); color: white; }
.user-chip-check { font-size: 16px; }
.user-newgroup { display: inline-flex; align-items: center; gap: 6px; }
.user-newgroup-input { width: 130px; height: 34px; }

.user-edit-actions { display: flex; gap: 10px; align-items: center; }
.user-delete { margin-left: auto; }

@media (max-width: 600px) {
  .user-field { flex-direction: column; align-items: stretch; gap: 6px; }
  .user-field .user-label { width: auto; min-width: 0; font-weight: 600; }
  .user-input { width: 100%; }
  .user-hint { margin-left: 0; }
  .password-rules { margin-left: 0; }
}
</style>
