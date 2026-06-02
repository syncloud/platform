<template>
  <label class="s-choice" :class="{ 'is-selected': isSelected }" v-bind="$attrs">
    <input type="radio" class="s-choice-input" :checked="isSelected" @change="pick">
    <span class="s-choice-body"><slot/></span>
  </label>
</template>

<script>
export default {
  name: 'SRadio',
  inheritAttrs: false,
  props: {
    label: { type: [String, Number], default: '' }
  },
  inject: ['sRadioGroup'],
  computed: {
    isSelected () { return this.sRadioGroup.modelValue === this.label }
  },
  methods: {
    pick () { this.sRadioGroup.pick(this.label) }
  }
}
</script>
