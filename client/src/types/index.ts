import type { Ref } from 'vue'

export type MaybeRef<T> = Ref<T> | T

export type LoginForm = {
  login: string
  password: string
  repeatPassword: string
}
