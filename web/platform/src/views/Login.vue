<template>
  <div class="sc-page">
    <div class="sc-card sc-card-narrow" id="block1">
      <h1 class="sc-title">{{ $t('login.title') }}</h1>
      <p class="sc-lead">{{ $t('login.redirecting') }}</p>
    </div>
  </div>

  <Error ref="error"/>

</template>

<script>
import Error from '../components/Error.vue'

export default {
  name: 'Login',
  components: {
    Error
  },
  mounted () {
    window.location.href = '/rest/oidc/login'
  }
}
</script>
