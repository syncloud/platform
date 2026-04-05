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
  <div class="alert alert-danger alert90" v-show="showError">
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
  height: 40px;
  padding: 0 50px 0 20px;
  border-radius: 3px;
  border: 1px solid #dcdee0;
  background-color: #fff!important;
  background-size: 14px 14px;
  transition: all .3s ease-out;
}
</style>
