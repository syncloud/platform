<template>
  <input
    class="sc-input"
    :class="{ 'sc-input-sm': size === 'small' }"
    :value="modelValue"
    @input="$emit('update:modelValue', $event.target.value)"
    v-bind="$attrs">
</template>

<script>
export default {
  name: 'SInput',
  inheritAttrs: false,
  props: {
    modelValue: { type: [String, Number], default: '' },
    size: { type: String, default: '' }
  },
  emits: ['update:modelValue']
}
</script>
