import { httpClient } from '@/plugins/httpClient'
import { BASE_URL } from '@/utils/constants'

export const chat = {
  list: () => {
    return httpClient.get<any>('/chats', { baseURL: BASE_URL })
  }
}
