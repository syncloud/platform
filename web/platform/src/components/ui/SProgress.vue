<template>
  <div class="s-progress" v-bind="$attrs">
    <div class="s-progress-track" :style="{ height: strokeWidth + 'px' }">
      <div
        class="s-progress-bar"
        :class="['s-progress-' + (status || 'normal'), { indeterminate: indeterminate }]"
        :style="indeterminate ? {} : { width: pct + '%' }"></div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'SProgress',
  inheritAttrs: false,
  props: {
    percentage: { type: Number, default: 0 },
    indeterminate: Boolean,
    status: { type: String, default: '' },
    strokeWidth: { type: Number, default: 6 },
    showText: { type: Boolean, default: true }
  },
  computed: {
    pct () {
      return Math.max(0, Math.min(100, this.percentage || 0))
    }
  }
}
</script>
