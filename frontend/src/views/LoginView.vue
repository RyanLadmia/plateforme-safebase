<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
    <div class="w-full max-w-md">
      <!-- Header de la page -->
      <div class="text-center mb-8">
        <div class="flex justify-center mb-4">
          <div class="bg-gradient-to-r from-blue-600 to-purple-600 p-4 rounded-full shadow-xl">
            <div class="text-4xl">🔐</div>
          </div>
        </div>
        <h1 class="text-3xl md:text-4xl font-bold text-gray-800 mb-2">
          Connexion
        </h1>
        <p class="text-gray-600">
          Accédez à votre espace sécurisé
        </p>
      </div>
      
      <!-- Formulaire d'authentification -->
      <AuthComponent @login-success="handleLoginSuccess" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import AuthComponent from '@/components/auth/AuthComponent.vue'

// Composables
const router = useRouter()

// Methods
const handleLoginSuccess = async (): Promise<void> => {
  // Redirection vers la page d'accueil
  await router.push('/')
}
</script>
