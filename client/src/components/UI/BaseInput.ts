import type { ErrorObject } from '@vuelidate/core'

export type InputProps = {
  modelValue: string
  prependIcon?: string
  prependIconColor?: string
  prependIconClickable?: boolean
  appendIcon?: string
  appendIconColor?: string
  appendIconClickable?: boolean
  label?: string
  hint?: string
  autofocus?: boolean
  type?: string
  disabled?: boolean
  placeholder?: string
  autocomplete?: string
  errorMessages?: ErrorObject[]
  required?: boolean
}

export type InputEmits = {
  (e: 'update:modelValue', value: string): void
  (e: 'click:prepend-icon', event: MouseEvent): void
  (e: 'click:append-icon', event: MouseEvent): void
  (e: 'focus', event: FocusEvent): void
  (e: 'blur', event: FocusEvent): void
}
