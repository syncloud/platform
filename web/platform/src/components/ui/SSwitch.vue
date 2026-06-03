<template>
  <label class="s-switch" :class="{ 'is-disabled': disabled }">
    <input
      type="checkbox"
      class="s-switch-input"
      role="switch"
      :aria-checked="modelValue ? 'true' : 'false'"
      :checked="!!modelValue"
      :disabled="disabled"
      @change="$emit('update:modelValue', $event.target.checked)"
      v-bind="$attrs">
    <span class="s-switch-track"><span class="s-switch-thumb"></span></span>
    <span v-if="stateText" class="s-switch-text">{{ stateText }}</span>
  </label>
</template>

<script>
export default {
  name: 'SSwitch',
  inheritAttrs: false,
  props: {
    modelValue: Boolean,
    disabled: Boolean,
    activeText: { type: String, default: '' },
    inactiveText: { type: String, default: '' },
    size: { type: String, default: '' }
  },
  emits: ['update:modelValue'],
  computed: {
    stateText () {
      return this.modelValue ? this.activeText : this.inactiveText
    }
  }
}
</script>
