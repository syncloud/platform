<template>
  <div class="s-alert" :class="'s-alert-' + type" role="alert">
    <i v-if="showIcon" class="material-icons s-alert-icon">{{ icon }}</i>
    <div class="s-alert-content"><slot/></div>
  </div>
</template>

<script>
export default {
  name: 'SAlert',
  props: {
    type: { type: String, default: 'info' },
    showIcon: Boolean,
    closable: { type: Boolean, default: true }
  },
  computed: {
    icon () {
      switch (this.type) {
        case 'success': return 'check_circle'
        case 'warning': return 'warning'
        case 'error': return 'error'
        default: return 'info'
      }
    }
  }
}
</script>
