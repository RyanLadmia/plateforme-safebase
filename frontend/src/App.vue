<script setup lang="ts">
import { onMounted } from 'vue'
import { RouterView } from 'vue-router'
import AppHeader from './layout/Header.vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

// Initialiser l'authentification au montage de l'app
onMounted(async () => {
  await authStore.checkAuth()
})
</script>

<template>
  <div class="min-h-screen flex flex-col">
    <AppHeader />
    <main class="flex-grow">
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
/* Styles spécifiques au composant App si nécessaire */
</style>

