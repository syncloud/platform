<template>
  <div class="s-radio-group" role="radiogroup" v-bind="$attrs"><slot/></div>
</template>

<script>
export default {
  name: 'SRadioGroup',
  inheritAttrs: false,
  props: {
    modelValue: { type: [String, Number], default: '' }
  },
  emits: ['update:modelValue'],
  provide () {
    return { sRadioGroup: this }
  },
  methods: {
    pick (value) { this.$emit('update:modelValue', value) }
  }
}
</script>
