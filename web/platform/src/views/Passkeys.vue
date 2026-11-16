<template>
  <div class="sc-page">
    <div class="sc-card" id="block1">
      <h1 class="sc-title" data-testid="passkeys-title">{{ $t('passkeys.title') }}</h1>
      <p class="sc-lead passkeys-hint">{{ $t('passkeys.hint') }}</p>
      <PasskeyList/>
    </div>
  </div>
</template>

<script>
import PasskeyList from '../components/PasskeyList.vue'

export default {
  name: 'Passkeys',
  components: {
    PasskeyList
  }
}
</script>

<style scoped>
.passkeys-hint { margin-bottom: 16px; }
</style>
