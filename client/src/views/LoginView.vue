<script setup lang="ts">
import { ref } from 'vue'
import BaseInput from '@/components/UI/BaseInput.vue'
import BaseButton from '@/components/UI/BaseButton.vue'
import BaseInputPassword from '@/components/UI/BaseInputPassword.vue'
import AuthLayout from '@/layouts/AuthLayout.vue'
import { type LoginForm } from '@/types'
import {
  containsNumber,
  containsSpecialSymbol,
  email,
  required,
  upperCaseAndLowerCase,
  withoutSpaces,
  maxLength,
  minLength
} from '@/utils/validators'
import useVuelidate from '@vuelidate/core'
import { auth } from '@/services/api/auth'

const form = ref<LoginForm>({
  email: '',
  password: ''
})

const rules = {
  email: { required, email },
  password: {
    required,
    minLength: minLength(8),
    maxLength: maxLength(40),
    upperCaseAndLowerCase,
    containsNumber,
    withoutSpaces,
    containsSpecialSymbol
  }
}

const v = useVuelidate<LoginForm>(rules, form)

const login = async () => {
  const isValid = await v.value.$validate()

  if (!isValid) {
    return
  }

  const test = await auth.login(form.value)
  console.log(test)
}
</script>

<template>
  <AuthLayout>
    <div class="login-form">
      <div class="login-form__inner">
        <h1 class="login-form__welcome">Welcome back</h1>
        <h1 class="login-form__title">Please enter your details to login</h1>

        <div class="login-form__inputs">
          <BaseInput
            v-model="form.email"
            label="Email address"
            placeholder="Enter email address"
            class="login-form__input"
            :error-messages="v.email.$errors"
          />
          <BaseInputPassword
            v-model="form.password"
            label="Password"
            placeholder="Enter password"
            type="password"
            class="login-form__input"
            :error-messages="v.password.$errors"
          />
          <BaseButton class="login-form__button" @click="login">Sign in</BaseButton>
        </div>

        <div class="login-form__controls">
          Don`t have an account ?
          <RouterLink :to="{ name: 'register' }" class="login-form__link">Create now</RouterLink>
        </div>

        <div class="login-form__social">
          <button class="login-form__link login-form__link_underlined">Login with Google</button>
        </div>
      </div>
    </div>
  </AuthLayout>
</template>

<style scoped lang="scss">
.login-form {
  width: 400px;
  display: flex;
  background-color: var(--vt-c-black);
  border-radius: 6px;
  padding: 2rem;
  box-shadow: var(--shadow-black);

  &__inner {
    flex: 1 1 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
  }

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
  }

  &__input {
    margin-bottom: 16px;
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
