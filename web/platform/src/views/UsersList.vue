<template>
  <div class="sc-page">
    <div class="sc-card" :style="{ visibility: visibility }">
      <h1 class="sc-title" data-testid="users-title">{{ $t('users.title') }}</h1>

      <div class="sc-actions">
        <button class="sc-btn sc-btn-success" id="btn_add_user" data-testid="users-add" type="button" @click="add">
          {{ $t('users.add') }}
        </button>
      </div>

      <div v-if="users.length === 0" class="setline">
        <span>{{ $t('users.noUsers') }}</span>
      </div>

      <router-link v-for="user in users" :key="user.username" class="user-entry"
                   :to="'/useredit?username=' + encodeURIComponent(user.username)"
                   :id="'user_' + user.username" :data-testid="'user-row-' + user.username">
        <span class="user-name">{{ user.username }}</span>
        <span class="user-email">{{ user.email }}</span>
        <span class="user-tags">
          <span v-if="user.admin" class="user-badge user-badge-admin" :data-testid="'user-admin-badge-' + user.username">{{ $t('users.admin') }}</span>
          <span v-for="group in user.groups" :key="group" class="user-badge">{{ group }}</span>
        </span>
        <i class="material-icons user-chevron">chevron_right</i>
      </router-link>
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
  name: 'UsersList',
  data () {
    return {
      users: [],
      visibility: 'hidden',
      loading: undefined
    }
  },
  components: {
    Error
  },
  mounted () {
    this.progressShow()
    this.reload()
  },
  methods: {
    add () {
      this.$router.push('/useredit')
    },
    progressShow () {
      this.loading = Loading.service({ lock: true, text: this.$t('common.loading'), background: 'rgba(0, 0, 0, 0.7)' })
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
      axios.get('/rest/users')
        .then(resp => Common.checkForServiceError(resp.data, () => {
          this.users = resp.data.data || []
          this.progressHide()
        }, onError))
        .catch(onError)
    }
  }
}
</script>
<style scoped>
.sc-actions {
  margin-bottom: 16px;
}

.user-entry {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 8px;
  border-bottom: 1px solid #eef3f9;
  text-decoration: none;
  color: inherit;
  cursor: pointer;
}
.user-entry:hover { background: var(--sc-field-bg); }

.user-name { font-weight: 600; flex: 0 0 140px; word-break: break-all; }
.user-email { flex: 1; color: var(--sc-muted); word-break: break-all; }
.user-tags { display: flex; gap: 6px; flex-wrap: wrap; flex: 0 0 auto; }
.user-badge {
  background: var(--sc-field-bg);
  border: 1px solid var(--sc-border);
  color: var(--sc-muted);
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 10px;
}
.user-badge-admin {
  background: var(--sc-success);
  border-color: var(--sc-success);
  color: white;
}
.user-chevron { color: var(--sc-faint); flex: 0 0 auto; }

@media (max-width: 600px) {
  .user-entry {
    flex-wrap: wrap;
    gap: 6px 10px;
    padding: 14px;
    border: 1px solid var(--sc-border);
    border-radius: 12px;
    background: var(--sc-field-bg);
    margin-bottom: 10px;
  }
  .user-name { order: 1; flex: 1 1 auto; font-size: 16px; }
  .user-chevron { order: 2; }
  .user-email { order: 3; flex: 1 1 100%; }
  .user-tags { order: 4; flex: 1 1 100%; }
}
</style>
