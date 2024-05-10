<script setup lang="ts"></script>

<template>
  <main class="auth-layout">
    <Transition name="auth-layout" :appear="true">
      <slot />
    </Transition>
  </main>
</template>

<style scoped lang="scss">
.auth-layout {
  flex-grow: 1;
  width: 100%;
  display: grid;
  place-items: center;

  &-enter-active,
  &-leave-active {
    transition:
      var(--transition) opacity,
      var(--transition) transform;
  }

  &-enter-from,
  &-leave-to {
    opacity: 0;
    transform: translateY(30px);
  }
}
</style>
