<template>
  <div class="wrapper">
    <div class="content">
      <div class="block1 wd12" id="block1">
        <h1>{{ $t('login.title') }}</h1>
        <div class="formblock">
          <p>{{ $t('login.redirecting') }}</p>
        </div>
      </div>
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
<style>
@import '../style/site.css';
@import 'material-icons/iconfont/material-icons.css';
</style>
