import { httpClient } from '@/plugins/httpClient'
import type { AccessToken } from '@/types'
import { BASE_URL } from '@/utils/constants'

export interface LoginPayload {
  email: string
  password: string
}

export interface RegisterPayload {
  email: string
  password: string
  username: string
}

export const auth = {
  login: (data: LoginPayload) => {
    return httpClient.post<AccessToken>('/login', data, { baseURL: BASE_URL })
  },
  register: (data: RegisterPayload) => {
    return httpClient.post<AccessToken>('/register', data, { baseURL: BASE_URL })
  }
}
