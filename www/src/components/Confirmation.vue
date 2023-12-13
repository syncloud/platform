<template>
  <el-dialog :modelValue="visible" style="min-width: 300px; max-width: 500px" :before-close="close" align-center>
    <template #header>
      <h4 class="modal-title">
        <slot name="title"></slot>
      </h4>
    </template>
    <slot name="text"></slot>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="close">Cancel</el-button>
        <el-button id="btn_confirm" type="primary" @click="yes">Confirm</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script>
export default {
  name: 'Confirmation',
  emits: ['confirm', 'cancel'],
  props: {
    visible: Boolean
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
