<template>
  <div class="s-checkbox-group" v-bind="$attrs"><slot/></div>
</template>

<script>
export default {
  name: 'SCheckboxGroup',
  inheritAttrs: false,
  props: {
    modelValue: { type: Array, default: () => [] }
  },
  emits: ['update:modelValue'],
  provide () {
    return { sCheckboxGroup: this }
  },
  methods: {
    toggle (value) {
      const arr = Array.isArray(this.modelValue) ? this.modelValue.slice() : []
      const i = arr.indexOf(value)
      if (i >= 0) { arr.splice(i, 1) } else { arr.push(value) }
      this.$emit('update:modelValue', arr)
    },
    has (value) {
      return Array.isArray(this.modelValue) && this.modelValue.indexOf(value) >= 0
    }
  }
}
</script>
