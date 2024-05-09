<script setup lang="ts">
import { ref } from 'vue'
import BaseInput from '@/components/UI/BaseInput.vue'
import { type LoginForm } from '@/types'
import BaseButton from '@/components/UI/BaseButton.vue'

const form = ref<LoginForm>({
  login: '',
  password: ''
})
</script>

<template>
  <div class="login-form">
    <h1 class="login-form__welcome">Welcome back</h1>
    <h1 class="login-form__title">Please enter your details to login</h1>

    <div class="login-form__inputs">
      <BaseInput v-model="form.login" label="Email address" placeholder="Enter email address" />
      <BaseInput v-model="form.password" label="Password" placeholder="Enter password" />
      <BaseButton>Login</BaseButton>
    </div>

    <div class="login-form__controls">
      Don`t have an account ? <RouterLink to="/">Create now</RouterLink>
    </div>

    <div class="login-form__social">
      <button>Login with Google</button>
    </div>
  </div>
</template>

<style scoped></style>
