import * as validators from '@vuelidate/validators'
import type { MaybeRef } from 'vue'

export const required = validators.helpers.withMessage('This field is required', validators.required)
export const email = validators.helpers.withMessage('Enter login as example@gmail.com', validators.email)

export const requiredWithMessage = (message: string) => {
  return validators.helpers.withMessage(message, validators.required)
}

export const minLength = (count: number) => {
  return validators.helpers.withMessage(`Enter at least ${count} symbols`, validators.minLength(count))
}

export const maxLength = (count: number) => {
  return validators.helpers.withMessage(`Enter at most ${count} symbols`, validators.maxLength(count))
}

export const upperCaseAndLowerCase = validators.helpers.withMessage(
  'Enter at least 1 capital and 1 lowercase latter',
  (value: string) => {
    // eslint-disable-next-line @typescript-eslint/no-extra-parens
    return !validators.helpers.req(value) || (value !== value.toLowerCase() && value !== value.toUpperCase())
  }
)

export const containsNumber = validators.helpers.withMessage('Enter at least 1 number', (value: string) => {
  return !validators.helpers.req(value) || /\d/g.test(value)
})

export const withoutSpaces = validators.helpers.withMessage('Enter symbols without spaces', (value: string) => {
  return !validators.helpers.req(value) || !/\s/g.test(value)
})

const specialSymbols = [
  '~',
  '!',
  '?',
  '@',
  '#',
  '$',
  '%',
  '^',
  '&',
  '*',
  '_',
  '-',
  '+',
  '(',
  ')',
  '[',
  ']',
  '{',
  '}',
  '>',
  '<',
  '/',
  '\\',
  '|',
  '"',
  "'",
  '.',
  ',',
  ':'
]
export const containsSpecialSymbol = validators.helpers.withMessage(
  `Should contain at least 1 symbol: ${specialSymbols.join(' ')}`,
  (value: string) => {
    return specialSymbols.some((symbol) => !validators.helpers.req(value) || value.includes(symbol))
  }
)

export const lettersWithSpaces = validators.helpers.withMessage(
  'This filed should contain only letters with spaces',
  (value: string) => {
    return !validators.helpers.req(value) || !/[^\sA-zЁА-яё]+/.test(value)
  }
)

export const sameAsPassword = (value: MaybeRef<string>) => {
  return validators.helpers.withMessage('Passwords don`t match', validators.sameAs(value))
}
