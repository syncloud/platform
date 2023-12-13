<template>
  <el-dialog :modelValue="visible" style="min-width: 300px; max-width: 500px" :before-close="close" align-center>
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
        <el-button @click="close">{{ cancelText }}</el-button>
        <el-button type="primary" @click="yes" v-if="confirmEnabled">Confirm</el-button>
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
      default: 'Cancel'
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
