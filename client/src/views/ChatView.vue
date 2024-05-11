<script setup lang="ts">
import ChatLayout from '@/layouts/ChatLayout.vue'
import ChatMessages from '@/modules/chat/components/ChatMessages.vue'
import ChatUserList from '@/modules/chat/components/ChatUserList.vue'
import { chat } from '@/services/api/chat'
import { onMounted, ref } from 'vue'
import { type ChatListItem, type ChatMessage } from '@/types'

const currentChatId = ref(0)

const chatList = ref<ChatListItem[]>([])
const chatMessages = ref<ChatMessage[]>([])

onMounted(async () => {
  const list = await chat.list()

  chatList.value = list.data
  currentChatId.value = chatList.value[0].receiver_id

  const lastMessages = await chat.lastMessages(currentChatId.value)
  chatMessages.value = lastMessages.data
})
</script>

<template>
  <ChatLayout>
    <main class="chat">
      <ChatUserList
        v-if="chatList.length"
        :chat-list="chatList"
        v-model:current-chat-id="currentChatId"
        class="chat__user-list"
      />
      <ChatMessages v-if="chatMessages" :chat-messages="chatMessages" class="chat__messages" />
    </main>
  </ChatLayout>
</template>

<style scoped lang="scss">
.chat {
  width: 80%;
  height: 80vh;
  background-color: var(--vt-c-black);
  padding: 16px 0;
  border-radius: 8px;
  display: flex;

  &__user-list {
    border-right: 1px solid var(--vt-c-black-mute);
  }
}
</style>
