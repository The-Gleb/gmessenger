import { httpClient } from '@/plugins/httpClient'
import { type ChatMessage, type ChatListItem } from '@/types'
import { BASE_URL } from '@/utils/constants'

export const chat = {
  list: () => {
    return httpClient.get<ChatListItem[]>('/chats', { baseURL: BASE_URL })
  },
  lastMessages: (receiverId: number) => {
    return httpClient.get<ChatMessage[]>(`/dialog/${receiverId}`, { baseURL: BASE_URL })
  }
}
