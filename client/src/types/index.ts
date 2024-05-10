import type { Ref } from 'vue'

export type MaybeRef<T> = Ref<T> | T

export type AccessToken = string
export type RefreshToken = string

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
