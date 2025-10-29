<template>
  <header class="bg-gradient-to-r from-blue-500 to-purple-600 text-white p-6 shadow-lg sticky top-0 z-50">
    <div class="container mx-auto">
      <div class="flex flex-col md:flex-row items-center justify-between">
        <!-- Logo et titre -->
        <div class="text-center md:text-left mb-4 md:mb-0">
          <RouterLink to="/" class="block">
            <h1 class="text-3xl font-bold hover:text-blue-200 transition-colors">
              🔐 SafeBase
            </h1>
            <p class="text-blue-100 text-sm">
              Plateforme de gestion sécurisée
            </p>
          </RouterLink>
        </div>
        
        <!-- Navigation -->
        <nav class="flex items-center space-x-4">
          <RouterLink 
            to="/" 
            class="bg-white/10 hover:bg-white/20 px-4 py-2 rounded-lg transition-all duration-200 backdrop-blur-sm"
            active-class="bg-white/30"
          >
            Accueil
          </RouterLink>
          <RouterLink 
            to="/about" 
            class="bg-white/10 hover:bg-white/20 px-4 py-2 rounded-lg transition-all duration-200 backdrop-blur-sm"
            active-class="bg-white/30"
          >
            À propos
          </RouterLink>
          
          <!-- Bouton conditionnel : Connexion ou Déconnexion -->
          <div v-if="isAuthenticated" class="flex items-center space-x-4">
            <!-- Informations utilisateur (mobile hidden) -->
            <div class="hidden md:flex flex-col items-end text-sm">
              <span class="font-medium">
                {{ currentUser?.firstname }} {{ currentUser?.lastname }}
              </span>
              <span 
                class="text-xs px-2 py-1 rounded-full font-semibold"
                :class="currentUser?.role_id === 1 
                  ? 'bg-yellow-400 text-yellow-900' 
                  : 'bg-green-400 text-green-900'"
              >
                {{ currentUser?.role_id === 1 ? 'Admin' : 'Utilisateur' }}
              </span>
            </div>
            
            <!-- Bouton de déconnexion -->
            <button 
              @click="handleLogout" 
              :disabled="loading"
              class="bg-red-500 hover:bg-red-600 px-4 py-2 rounded-lg transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center space-x-2"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l-4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
              </svg>
              <span>{{ loading ? 'Déconnexion...' : 'Déconnexion' }}</span>
            </button>
          </div>
          
          <!-- Bouton de connexion si non connecté -->
          <RouterLink 
            v-else
            to="/login" 
            class="bg-white text-blue-600 hover:bg-blue-50 px-4 py-2 rounded-lg transition-all duration-200 font-medium flex items-center space-x-2"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1" />
            </svg>
            <span>Connexion</span>
          </RouterLink>
        </nav>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth'

// Composables
const router = useRouter()
const authStore = useAuthStore()

// State
const loading = ref<boolean>(false)

// Computed réactifs depuis le store (avec storeToRefs pour préserver la réactivité)
const { isAuthenticated, user: currentUser } = storeToRefs(authStore)

// Methods
const handleLogout = async (): Promise<void> => {
  loading.value = true
  try {
    await authStore.logout()
    await router.push('/login')
  } catch (error) {
    console.error('Erreur de déconnexion:', error)
  } finally {
    loading.value = false
  }
}

// Watchers pour la réactivité (production ready)
watch(() => isAuthenticated, (newValue) => {
  // Logique additionnelle si nécessaire lors du changement d'état
}, { immediate: true })
</script>
