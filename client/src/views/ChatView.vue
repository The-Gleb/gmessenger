<script setup lang="ts">
import ChatLayout from '@/layouts/ChatLayout.vue'
import ChatMessages from '@/modules/chat/components/ChatMessages.vue'
import ChatUserList from '@/modules/chat/components/ChatUserList.vue'
import { chat } from '@/services/api/chat'
import { onMounted, ref } from 'vue'
import { type ChatListItem } from '@/types'

const chatList = ref<ChatListItem[]>([])

onMounted(async () => {
  const { data } = await chat.list()
  chatList.value = data
})
</script>

<template>
  <ChatLayout>
    <main class="chat">
      <ChatUserList v-if="chatList.length" :chat-list="chatList" class="chat__user-list" />
      <ChatMessages class="chat__messages" />
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
