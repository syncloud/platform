<template>
  <button class="sc-btn" :class="typeClass" :disabled="disabled || loading" v-bind="$attrs">
    <span v-if="loading" class="sc-btn-spinner" aria-hidden="true"></span>
    <slot/>
  </button>
</template>

<script>
export default {
  name: 'SButton',
  inheritAttrs: false,
  props: {
    type: { type: String, default: '' },
    loading: Boolean,
    disabled: Boolean,
    size: { type: String, default: '' }
  },
  computed: {
    typeClass () {
      return {
        'sc-btn-primary': this.type === 'primary',
        'sc-btn-success': this.type === 'success',
        'sc-btn-danger': this.type === 'danger',
        'sc-btn-ghost': this.type === '' || this.type === 'default' || this.type === 'info',
        'sc-btn-sm': this.size === 'small'
      }
    }
  }
}
</script>
