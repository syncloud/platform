<template>
  <label>
    <input :id="id"
           ref="switch"
           type="checkbox"
           :data-on-text="onLabel"
           :data-off-text="offLabel"
           data-label-width="8"
    />
  </label>
</template>
<script>

export default {
  name: 'Switch',
  emits: ['toggle'],
  props: {
    id: String,
    checked: { type: Boolean, default: false },
    onLabel: String,
    offLabel: String
  },
  watch: {
    checked (val) {
      this.init()
    }
  },
  methods: {
    init () {
      const input = $(this.$refs.switch)
      input.off('switchChange.bootstrapSwitch')
      input.bootstrapSwitch('state', this.checked)
      const that = this
      input.on('switchChange.bootstrapSwitch', function (e, s) {
        that.$emit('toggle')
      })
    }
  },
  mounted () {
    this.init()
  }
}
</script>
