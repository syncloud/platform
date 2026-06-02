<template>
  <label class="s-choice" :class="{ 'is-selected': isChecked }" v-bind="$attrs">
    <input type="checkbox" class="s-choice-input" :checked="isChecked" @change="toggle">
    <span class="s-choice-body"><slot/></span>
  </label>
</template>

<script>
export default {
  name: 'SCheckbox',
  inheritAttrs: false,
  props: {
    label: { type: [String, Number], default: '' }
  },
  inject: ['sCheckboxGroup'],
  computed: {
    isChecked () { return this.sCheckboxGroup.has(this.label) }
  },
  methods: {
    toggle () { this.sCheckboxGroup.toggle(this.label) }
  }
}
</script>
