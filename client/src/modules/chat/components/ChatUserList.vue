<script setup lang="ts">
import { ref } from 'vue'
import { type ChatListItem } from '@/types'

type Props = {
  chatList: ChatListItem[]
}

const props = defineProps<Props>()

const currentChatId = ref(props.chatList[0].reciever_id)
</script>

<template>
  <section class="chat-user-list">
    <div class="chat-user-list__header">
      <h2 class="chat-user-list__title">Gmessenger</h2>
    </div>
    <div class="chat-user-list__content">
      <button
        v-for="item in chatList"
        :key="item.reciever_id"
        :class="['chat-user-list-item', { 'chat-user-list-item_active': item.reciever_id === currentChatId }]"
      >
        <img class="chat-user-list-item__avatar" src="@/assets/img/avatar.png" />
        <div class="chat-user-list-item__body">
          <div class="chat-user-list-item__user-info">
            <p class="chat-user-list-item__username">{{ item.receiver_name }}</p>
          </div>
          <p class="chat-user-list-item__text">{{ item.last_message.text }}</p>
        </div>
      </button>
    </div>
  </section>
</template>

<style scoped lang="scss">
.chat-user-list {
  flex: 0 1 30%;
  display: flex;
  flex-direction: column;

  &__header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 8px 0 24px 16px;
    border-bottom: 1px solid var(--vt-c-black-mute);
  }

  &__title {
    font-size: 1rem;
    font-weight: 600;
  }

  &__content {
    flex-grow: 1;
    display: flex;
    flex-direction: column;
    overflow: auto;
    margin-top: 16px;
    -ms-overflow-style: none;
    scrollbar-width: none;

    &::-webkit-scrollbar {
      display: none;
      width: 0;
    }
  }
}

.chat-user-list-item {
  display: flex;
  align-items: center;
  column-gap: 16px;
  padding: 12px 16px;

  &:not(:last-child) {
    border-bottom: 1px solid var(--vt-c-black-mute);
  }

  &_active {
    background-color: var(--vt-c-indigo);
    border-left: 3px solid var(--vt-c-dark-green);
  }

  &__avatar {
    width: 50px;
    height: 50px;
    border-radius: 50%;
    object-fit: cover;
  }

  &__body {
    flex-grow: 1;
    display: flex;
    flex-direction: column;
    row-gap: 12px;
  }

  &__user-info {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  &__username {
    font-size: 1rem;
    font-weight: 600;
    color: var(--vt-c-white-mute);
  }

  &__date {
    @include text-14s;

    color: var(--vt-c-grey);
  }

  &__text {
    @include text-14r;
    @include line-clamp(1);

    color: var(--);
  }
}
</style>
