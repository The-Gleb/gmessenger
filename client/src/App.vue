<script setup lang="ts">
import { RouterView } from 'vue-router'
import { useAuthStore } from './stores/auth'
import { onMounted } from 'vue'

const { init } = useAuthStore()

onMounted(() => {
  init()
})
</script>

<template>
  <RouterView />
</template>
