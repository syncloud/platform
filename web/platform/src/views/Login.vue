<template>
  <div class="sc-page">
    <div class="sc-card sc-card-narrow" id="block1">
      <h1 class="sc-title">{{ $t('login.title') }}</h1>
      <template v-if="reason">
        <p class="sc-lead" data-testid="login-reason">{{ $t('login.' + reason) }}</p>
        <button id="btn_login_retry" class="act-btn act-btn-primary" @click="signIn">
          {{ $t('login.retry') }}
        </button>
      </template>
      <p v-else class="sc-lead">{{ $t('login.redirecting') }}</p>
    </div>
  </div>

  <Error ref="error"/>

</template>

<script>
import Error from '../components/Error.vue'

const REASONS = ['session', 'signin']

export default {
  name: 'Login',
  components: {
    Error
  },
  data () {
    return {
      reason: ''
    }
  },
  mounted () {
    const reason = this.$route.query.error
    if (REASONS.includes(reason)) {
      this.reason = reason
      return
    }
    this.signIn()
  },
  methods: {
    signIn () {
      window.location.href = '/rest/oidc/login'
    }
  }
}
</script>
