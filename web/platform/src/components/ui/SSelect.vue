<template>
  <div class="s-select" :class="{ 'is-open': open, 'is-disabled': disabled }" :disabled="disabled || undefined" v-bind="$attrs">
    <div class="s-select__control" @click="onControlClick">
      <input
        v-if="filterable"
        class="s-select__search"
        :value="open ? query : displayLabel"
        :placeholder="placeholder || displayLabel"
        :disabled="disabled"
        autocomplete="off"
        @input="query = $event.target.value; open = true"
        @focus="openDropdown"
        @click.stop="openDropdown">
      <span v-else class="s-select__value" :class="{ 'is-placeholder': displayLabel === '' }">{{ displayLabel || placeholder }}</span>
      <span class="s-select__caret">▾</span>
    </div>
    <div v-if="open" class="s-select__backdrop" @click="close"></div>
    <div v-show="open" class="s-select__dropdown">
      <slot/>
    </div>
  </div>
</template>

<script>
export default {
  name: 'SSelect',
  inheritAttrs: false,
  props: {
    modelValue: { type: [String, Number], default: '' },
    placeholder: { type: String, default: '' },
    filterable: Boolean,
    disabled: Boolean
  },
  emits: ['update:modelValue'],
  provide () {
    return { sSelect: this }
  },
  data () {
    return {
      open: false,
      query: '',
      labels: {}
    }
  },
  computed: {
    displayLabel () {
      const l = this.labels[this.modelValue]
      return l !== undefined ? l : (this.modelValue === undefined || this.modelValue === null ? '' : String(this.modelValue))
    }
  },
  methods: {
    registerOption (value, label) { this.labels[value] = label },
    unregisterOption (value) { delete this.labels[value] },
    onControlClick () {
      if (this.disabled) return
      if (this.filterable) { this.openDropdown() } else { this.open = !this.open }
    },
    openDropdown () { if (!this.disabled) this.open = true },
    close () { this.open = false; this.query = '' },
    choose (value) {
      this.$emit('update:modelValue', value)
      this.close()
    }
  }
}
</script>
