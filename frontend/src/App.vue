<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <header class="bg-gradient-to-r from-blue-500 to-purple-600 text-white p-6 rounded-lg shadow-lg mb-8">
    <div class="container mx-auto">
      <h1 class="text-3xl font-bold text-center mb-4">SafeBase</h1>
      <p class="text-center text-blue-100 mb-6">Plateforme de sauvegarde et restauration de bases de données</p>
      
      <nav class="text-center">
        <RouterLink to="/" class="bg-white text-blue-600 px-4 py-2 rounded-md mx-2 hover:bg-gray-100 transition-colors">Accueil</RouterLink>
        <RouterLink to="/about" class="bg-white text-blue-600 px-4 py-2 rounded-md mx-2 hover:bg-gray-100 transition-colors">À propos</RouterLink>
      </nav>
    </div>
  </header>

  <RouterView />
</template>

