import type { Ref } from 'vue'

export type MaybeRef<T> = Ref<T> | T

export type AccessToken = {
  token: string
}

export type RefreshToken = {
  token: string
}

export type LoginForm = {
  email: string
  password: string
}

export type RegistrationForm = {
  email: string
  password: string
  repeatPassword: string
  username: string
}

export type User = {
  id: number
  username: string
  email: string
}
