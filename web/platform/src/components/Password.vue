<template>
  <div style="position: relative; margin-bottom: 10px">
    <input :placeholder="placeholder" :type="show ? 'text' : 'password'" class="password"
           :id="id"
           :value="modelValue"
           @input="$emit('update:modelValue', $event.target.value)"
           v-on:keyup.enter="$emit('trigger')"
           >
    <i class="fa" :class="show ? 'fa-eye-slash' : 'fa-eye'"
       style="position: absolute; right: 10px; top: 50%; transform: translateY(-50%); cursor: pointer;"
       @click="toggleVisibility"></i>
  </div>
  <div class="password-error" v-show="showError">
    {{ error }}
  </div>
</template>

<script>
export default {
  name: 'Password',
  emits: [
    'trigger',
    'update:modelValue'
  ],
  props: {
    id: String,
    placeholder: String,
    modelValue: String,
    showError: Boolean,
    error: String
  },
  data () {
    return {
      show: false
    }
  },
  methods: {
    toggleVisibility () {
      this.show = !this.show
    }
  }
}
</script>
<style scoped>
.password {
  width: 100%;
  height: 46px;
  padding: 0 50px 0 16px;
  border-radius: var(--sc-control-radius, 12px);
  border: 1px solid var(--sc-border, #d5dde8);
  background: var(--sc-field-bg, #f6f9fd);
  font-size: 16px;
  color: var(--sc-ink, #1a2a3a);
  box-sizing: border-box;
  transition: border-color 0.2s ease, background 0.2s ease, box-shadow 0.2s ease;
}
.password:focus {
  outline: none;
  background: var(--sc-surface);
  border-color: var(--sc-primary, #2b7bd6);
  box-shadow: 0 0 0 4px rgba(43, 123, 214, 0.12);
}
.fa-eye, .fa-eye-slash { color: var(--sc-faint, #8796a8); }
.password-error {
  margin-top: 6px;
  color: var(--sc-danger, #d9363e);
  font-size: 13px;
  white-space: pre-line;
}
</style>
