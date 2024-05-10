<script setup lang="ts">
import { ref } from 'vue'
import BaseInput from '@/components/UI/BaseInput.vue'
import BaseButton from '@/components/UI/BaseButton.vue'
import BaseInputPassword from '@/components/UI/BaseInputPassword.vue'
import AuthLayout from '@/layouts/AuthLayout.vue'
import { type LoginForm } from '@/types'

const form = ref<LoginForm>({
  login: '',
  password: ''
})
</script>

<template>
  <AuthLayout>
    <div class="login-form">
      <h1 class="login-form__welcome">Welcome back</h1>
      <h1 class="login-form__title">Please enter your details to login</h1>

      <div class="login-form__inputs">
        <BaseInput v-model="form.login" label="Email address" placeholder="Enter email address" />
        <BaseInputPassword
          v-model="form.password"
          label="Password"
          placeholder="Enter password"
          type="password"
        />
        <BaseButton class="login-form__button">Login</BaseButton>
      </div>

      <div class="login-form__controls">
        Don`t have an account ? <RouterLink to="/" class="login-form__link">Create now</RouterLink>
      </div>

      <div class="login-form__social">
        <button class="login-form__link login-form__link_underlined">Login with Google</button>
      </div>
    </div>
  </AuthLayout>
</template>

<style scoped lang="scss">
.login-form {
  display: flex;
  flex-direction: column;
  align-items: center;
  background-color: var(--vt-c-black);
  border-radius: 6px;
  padding: 2rem;
  box-shadow: var(--shadow-black);

  &__welcome {
    font-size: 1.25rem;
    font-weight: 600;
  }

  &__title {
    font-size: 1rem;
    margin-bottom: 3rem;
  }

  &__inputs {
    align-self: stretch;
    display: flex;
    flex-direction: column;
    row-gap: 16px;
  }

  &__controls {
    margin-bottom: 16px;
  }

  &__button {
    margin-bottom: 12px;
  }

  &__link {
    color: var(--vt-c-light-green);
    font-weight: 600;

    transition: color var(--transition);

    &:hover {
      color: var(--vt-c-green);
    }

    &_underlined {
      text-decoration: underline;
    }
  }
}
</style>
