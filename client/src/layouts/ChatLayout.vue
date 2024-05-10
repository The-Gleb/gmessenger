<script setup lang="ts"></script>

<template>
  <main class="chat-layout">
    <Transition name="chat-layout" :appear="true">
      <slot />
    </Transition>
  </main>
</template>

<style scoped lang="scss">
.chat-layout {
  flex-grow: 1;
  width: 100%;
  display: grid;
  place-items: center;

  &-enter-active,
  &-leave-active {
    transition: var(--transition) opacity;
  }

  &-enter-from,
  &-leave-to {
    opacity: 0;
  }
}
</style>
