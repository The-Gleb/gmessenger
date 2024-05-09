<script setup lang="ts"></script>

<template>
  <main class="auth-layout">
    <slot />
  </main>
</template>

<style scoped lang="scss">
.auth-layout {
  flex-grow: 1;
  width: 100%;
  display: grid;
  place-items: center;
}
</style>
