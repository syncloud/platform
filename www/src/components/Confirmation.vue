<template>
  <el-dialog :modelValue="visible" style="min-width: 300px; max-width: 500px" :before-close="close">
    <template #header>
      <h4 class="modal-title">
        <slot name="title"></slot>
      </h4>
    </template>
    <div class="bodymod">
      <div class="btext">
        <slot name="text"></slot>
      </div>
    </div>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="close">Cancel</el-button>
        <el-button type="primary" @click="yes">Confirm</el-button>
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
