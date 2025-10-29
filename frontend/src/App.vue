<script setup lang="ts">
import { RouterView } from 'vue-router'
import AppHeader from '@/layout/AppHeader.vue'
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

