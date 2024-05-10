<script setup lang="ts">
import { computed, ref } from 'vue'
import BaseInput from '@/components/UI/BaseInput.vue'
import BaseButton from '@/components/UI/BaseButton.vue'
import BaseInputPassword from '@/components/UI/BaseInputPassword.vue'
import AuthLayout from '@/layouts/AuthLayout.vue'
import { type RegistrationForm } from '@/types'
import {
  containsNumber,
  containsSpecialSymbol,
  email,
  required,
  upperCaseAndLowerCase,
  withoutSpaces,
  maxLength,
  minLength,
  sameAsPassword
} from '@/utils/validators'
import useVuelidate from '@vuelidate/core'
import { auth } from '@/services/api/auth'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const router = useRouter()

const form = ref<RegistrationForm>({
  email: '',
  username: '',
  password: '',
  repeatPassword: ''
})

const rules = {
  email: { required, email },
  username: { required },
  password: {
    required,
    minLength: minLength(8),
    maxLength: maxLength(40),
    upperCaseAndLowerCase,
    containsNumber,
    withoutSpaces,
    containsSpecialSymbol
  },
  repeatPassword: {
    required,
    sameAs: sameAsPassword(computed(() => form.value.password))
  }
}

const v = useVuelidate<RegistrationForm>(rules, form)

const register = async () => {
  const isValid = await v.value.$validate()

  if (!isValid) {
    return
  }

  const { setAccessToken, setUser } = useAuthStore()
  const { repeatPassword, ...data } = form.value

  try {
    const token = await auth.register(data)
    setAccessToken(token.data.token)

    const user = await auth.user()
    setUser(user.data)

    await router.push({ name: 'home' })
  } catch (e) {
    console.error(e)
  }
}
</script>

<template>
  <AuthLayout>
    <div class="registration-form">
      <div class="registration-form__inner">
        <h1 class="registration-form__welcome">Welcome to Gmessenger</h1>
        <h1 class="registration-form__title">Please fill the fields below to register</h1>

        <div class="registration-form__inputs">
          <BaseInput
            v-model="form.email"
            label="Email address"
            placeholder="Enter email address"
            class="registration-form__input"
            :error-messages="v.email.$errors"
          />
          <BaseInput
            v-model="form.username"
            label="Username"
            placeholder="Enter username"
            class="registration-form__input"
            :error-messages="v.username.$errors"
          />
          <BaseInputPassword
            v-model="form.password"
            label="Password"
            placeholder="Enter password"
            type="password"
            class="registration-form__input"
            :error-messages="v.password.$errors"
          />
          <BaseInputPassword
            v-model="form.repeatPassword"
            label="Repeat password"
            placeholder="Enter password again"
            type="password"
            class="registration-form__input"
            :error-messages="v.repeatPassword.$errors"
          />
          <BaseButton class="registration-form__button" @click="register">Sign up</BaseButton>
        </div>

        <div class="registration-form__controls">
          Already have an account ?
          <RouterLink :to="{ name: 'login' }" class="registration-form__link">Sign in</RouterLink>
        </div>

        <div class="registration-form__social">
          <button class="registration-form__link registration-form__link_underlined">Login with Google</button>
        </div>
      </div>
    </div>
  </AuthLayout>
</template>

<style scoped lang="scss">
.registration-form {
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
