<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50">
    <div class="text-center">
      <h1 class="text-9xl font-bold text-gray-300">404</h1>
      <h2 class="text-3xl font-semibold text-gray-700 mt-4">Page introuvable</h2>
      <p class="text-gray-500 mt-2">La page que vous recherchez n'existe pas.</p>
      <router-link 
        to="/" 
        class="mt-6 inline-block px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
      >
        Retour à l'accueil
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
// Page 404
</script>
