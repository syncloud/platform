<template>
  <div v-show="visible" class="s-option" :class="{ 'is-selected': isSelected }" @click="choose" v-bind="$attrs">
    <slot>{{ label }}</slot>
  </div>
</template>

<script>
export default {
  name: 'SOption',
  inheritAttrs: false,
  props: {
    value: { type: [String, Number], default: '' },
    label: { type: [String, Number], default: '' }
  },
  inject: ['sSelect'],
  computed: {
    isSelected () { return this.sSelect.modelValue === this.value },
    visible () {
      if (!this.sSelect.filterable) return true
      const q = (this.sSelect.query || '').trim().toLowerCase()
      if (!q) return true
      return String(this.label).toLowerCase().includes(q)
    }
  },
  mounted () { this.sSelect.registerOption(this.value, this.label) },
  unmounted () { this.sSelect.unregisterOption(this.value) },
  methods: {
    choose () { this.sSelect.choose(this.value) }
  }
}
</script>
