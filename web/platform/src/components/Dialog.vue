<template>
  <div v-if="visible" class="s-modal-overlay" @click.self="close">
    <div class="s-modal syncloud-dialog" role="dialog">
      <h4 class="modal-title"><slot name="title"></slot></h4>
      <div class="s-modal-body"><slot name="text"></slot></div>
      <div class="s-modal-footer">
        <button class="sc-btn sc-btn-ghost" type="button" @click="close">{{ cancelText || $t('common.cancel') }}</button>
        <button v-if="confirmEnabled" id="btn_confirm" data-testid="btn_confirm" class="sc-btn sc-btn-primary" type="button" @click="yes">{{ $t('common.confirm') }}</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Dialog',
  emits: ['confirm', 'cancel'],
  props: {
    visible: Boolean,
    confirmEnabled: {
      type: Boolean,
      default: true
    },
    cancelText: {
      type: String,
      default: ''
    }
  },
  methods: {
    yes () {
      this.$emit('confirm')
    },
    close () {
      this.$emit('cancel')
    }
  }
}
</script>
