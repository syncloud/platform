<template>
  <el-dialog class="syncloud-dialog" :modelValue="visible" style="min-width: 300px; max-width: 500px" :before-close="close" align-center>
    <template #header>
      <h4 class="modal-title">
        <slot name="title"></slot>
      </h4>
    </template>
    <div style="font-size: 18px">
      <slot name="text" ></slot>
    </div>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="close">{{ cancelText || $t('common.cancel') }}</el-button>
        <el-button id="btn_confirm" type="primary" @click="yes" v-if="confirmEnabled">{{ $t('common.confirm') }}</el-button>
      </span>
    </template>
  </el-dialog>
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
<style>
@import '../style/site.css';
</style>
