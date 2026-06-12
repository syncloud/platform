<template>
  <div class="sc-page">
    <div class="sc-card" :style="{ visibility: visibility }">
      <h1 class="sc-title" data-testid="users-title">{{ $t('users.title') }}</h1>

      <h3>{{ $t('users.addUser') }}</h3>
      <div class="setline user-field">
        <label class="span user-label" for="user_username">{{ $t('users.username') }}</label>
        <input class="user-input sc-input" id="user_username" type="text" v-model="newUsername" :placeholder="$t('users.usernamePlaceholder')">
      </div>
      <div class="setline user-field">
        <label class="span user-label" for="user_password">{{ $t('users.password') }}</label>
        <input class="user-input sc-input" id="user_password" type="password" v-model="newPassword">
      </div>
      <div class="setline user-field">
        <label class="span user-label" for="user_email">{{ $t('users.email') }}</label>
        <input class="user-input sc-input" id="user_email" type="text" v-model="newEmail" :placeholder="(newUsername || $t('users.usernameFallback')) + '@' + domain">
      </div>
      <div class="setline user-toggle">
        <label class="span user-label" for="user_admin">{{ $t('users.admin') }}</label>
        <s-switch id="user_admin" data-testid="user-admin-new" v-model="newAdmin"/>
      </div>
      <div class="sc-actions">
        <button class="sc-btn sc-btn-success" id="btn_add_user" type="button" @click="addUser">{{ $t('users.add') }}</button>
      </div>

      <h3>{{ $t('users.users') }}</h3>
      <div v-if="users.length === 0" class="setline">
        <span>{{ $t('users.noUsers') }}</span>
      </div>
      <div v-for="user in users" :key="user.username" class="user-entry" :data-testid="'user-row-' + user.username">
        <div class="user-row-main">
          <span class="user-name" :data-testid="'user-name-' + user.username">{{ user.username }}</span>
          <div class="user-email-edit">
            <input class="sc-input user-email-input" type="text" :id="'user_email_' + user.username"
                   :data-testid="'user-email-' + user.username" v-model="user.email"
                   @keyup.enter="saveEmail(user)">
            <button class="sc-btn sc-btn-primary user-email-save" type="button" :id="'btn_email_' + user.username"
                    @click="saveEmail(user)">{{ $t('users.save') }}</button>
          </div>
          <div class="user-admin-toggle">
            <span class="user-admin-label">{{ $t('users.admin') }}</span>
            <s-switch :id="'user_admin_' + user.username" :data-testid="'user-admin-' + user.username"
                      :modelValue="user.admin" @update:modelValue="value => setAdmin(user, value)"/>
          </div>
          <button class="sc-btn sc-btn-danger user-remove" type="button" :id="'btn_remove_user_' + user.username"
                  @click="removeUser(user.username)">{{ $t('users.remove') }}
          </button>
        </div>
        <div v-if="customGroups.length > 0" class="user-groups">
          <span class="user-groups-label">{{ $t('users.groups') }}</span>
          <label v-for="group in customGroups" :key="group.name" class="user-group-check">
            <input type="checkbox" :id="'user_group_' + user.username + '_' + group.name"
                   :data-testid="'user-group-' + user.username + '-' + group.name"
                   :checked="user.groups.includes(group.name)"
                   @change="toggleGroup(user, group.name, $event.target.checked)">
            <span>{{ group.name }}</span>
          </label>
        </div>
      </div>

      <h3>{{ $t('users.groupsManage') }}</h3>
      <div class="setline user-field">
        <label class="span user-label" for="group_name">{{ $t('users.groupName') }}</label>
        <input class="user-input sc-input" id="group_name" type="text" v-model="newGroup" :placeholder="$t('users.groupPlaceholder')">
        <button class="sc-btn sc-btn-success group-add" id="btn_add_group" type="button" @click="addGroup">{{ $t('users.add') }}</button>
      </div>
      <div v-if="customGroups.length === 0" class="setline">
        <span>{{ $t('users.noGroups') }}</span>
      </div>
      <div v-for="group in customGroups" :key="group.name" class="group-entry" :data-testid="'group-row-' + group.name">
        <span class="group-name">{{ group.name }}</span>
        <span class="group-members">{{ group.members.join(', ') }}</span>
        <button class="sc-btn sc-btn-danger group-remove" type="button" :id="'btn_remove_group_' + group.name"
                @click="removeGroup(group.name)">{{ $t('users.remove') }}
        </button>
      </div>
    </div>
  </div>

  <Error ref="error"/>

</template>

<script>
import Error from '../components/Error.vue'
import * as Common from '../js/common.js'
import axios from 'axios'
import Loading from '../util/loading'

export default {
  name: 'Users',
  props: {
    checkUserSession: Function,
    activated: Boolean
  },
  data () {
    return {
      users: [],
      groups: [],
      newUsername: '',
      newPassword: '',
      newEmail: '',
      newAdmin: false,
      newGroup: '',
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
    },
    customGroups () {
      return this.groups.filter(group => group.name !== 'syncloud')
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
      this.loading.close()
    },
    onError (err) {
      this.$refs.error.showAxios(err)
      this.progressHide()
    },
    reload () {
      const onError = (err) => this.onError(err)
      axios.get('/rest/groups')
        .then(resp => Common.checkForServiceError(resp.data, () => {
          this.groups = resp.data.data || []
          axios.get('/rest/users')
            .then(usersResp => Common.checkForServiceError(usersResp.data, () => {
              this.users = usersResp.data.data || []
              this.progressHide()
            }, onError))
            .catch(onError)
        }, onError))
        .catch(onError)
    },
    addUser () {
      this.progressShow()
      const onError = (err) => this.onError(err)
      axios.post('/rest/users/add', {
        username: this.newUsername,
        password: this.newPassword,
        email: this.newEmail,
        admin: this.newAdmin
      })
        .then(resp => Common.checkForServiceError(resp.data, () => {
          this.newUsername = ''
          this.newPassword = ''
          this.newEmail = ''
          this.newAdmin = false
          this.reload()
        }, onError))
        .catch(onError)
    },
    removeUser (username) {
      this.progressShow()
      const onError = (err) => this.onError(err)
      axios.post('/rest/users/remove', { username: username })
        .then(resp => Common.checkForServiceError(resp.data, () => this.reload(), onError))
        .catch(onError)
    },
    saveEmail (user) {
      this.progressShow()
      const onError = (err) => this.onError(err)
      axios.post('/rest/users/email', { username: user.username, email: user.email })
        .then(resp => Common.checkForServiceError(resp.data, () => this.reload(), onError))
        .catch(onError)
    },
    setAdmin (user, admin) {
      this.progressShow()
      const onError = (err) => this.onError(err)
      axios.post('/rest/users/admin', { username: user.username, admin: admin })
        .then(resp => Common.checkForServiceError(resp.data, () => this.reload(), onError))
        .catch(onError)
    },
    toggleGroup (user, group, member) {
      this.progressShow()
      const onError = (err) => this.onError(err)
      axios.post('/rest/groups/member', { group: group, username: user.username, member: member })
        .then(resp => Common.checkForServiceError(resp.data, () => this.reload(), onError))
        .catch(onError)
    },
    addGroup () {
      this.progressShow()
      const onError = (err) => this.onError(err)
      axios.post('/rest/groups/add', { name: this.newGroup })
        .then(resp => Common.checkForServiceError(resp.data, () => {
          this.newGroup = ''
          this.reload()
        }, onError))
        .catch(onError)
    },
    removeGroup (name) {
      this.progressShow()
      const onError = (err) => this.onError(err)
      axios.post('/rest/groups/remove', { name: name })
        .then(resp => Common.checkForServiceError(resp.data, () => this.reload(), onError))
        .catch(onError)
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

.user-field {
  display: flex;
  align-items: center;
}

.user-toggle {
  display: flex;
  align-items: center;
}

.user-input {
  width: 220px;
  height: 38px;
}

.group-add {
  margin-left: 12px;
}

.user-entry {
  padding: 12px 0;
  border-bottom: 1px solid #eef3f9;
}

.user-row-main {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.user-name {
  font-weight: 600;
  flex: 0 0 120px;
}

.user-email-edit {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.user-email-input {
  flex: 1;
  min-width: 160px;
  height: 38px;
}

.user-admin-toggle {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 0 0 auto;
}

.user-admin-label {
  color: var(--sc-muted);
}

.user-remove {
  min-width: 80px;
  padding: 8px 14px;
  flex: 0 0 auto;
}

.user-groups {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-top: 8px;
  flex-wrap: wrap;
}

.user-groups-label {
  color: var(--sc-faint);
}

.user-group-check {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.group-entry {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 0;
  border-bottom: 1px solid #eef3f9;
}

.group-name {
  font-weight: 600;
  flex: 0 0 120px;
}

.group-members {
  flex: 1;
  color: var(--sc-muted);
  word-break: break-all;
}

.group-remove {
  min-width: 80px;
  padding: 8px 14px;
  flex: 0 0 auto;
}

@media (max-width: 600px) {
  .user-field {
    flex-direction: column;
    align-items: stretch;
    gap: 6px;
  }
  .user-field .user-label { width: auto; min-width: 0; font-weight: 600; }
  .user-input { width: 100%; }

  .user-row-main {
    align-items: stretch;
    flex-direction: column;
  }
  .user-name { flex: 1 1 auto; }
  .user-email-edit { flex-wrap: wrap; }
}
</style>
