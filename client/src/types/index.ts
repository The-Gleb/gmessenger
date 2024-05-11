import type { Ref } from 'vue'

export type ISODateString = string
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

export type ChatLastMessage = {
  id: number
  receiver_id: number
  sender_id: number
  sender_name: string
  status: string
  text: string
  time: ISODateString
}

export type ChatListItem = {
  reciever_id: number
  receiver_name: string
  type: 'PERSONAL_CHAT' | 'DIALOG'
  unread: number
  last_message: ChatLastMessage
}
