<template>
  <header v-if="initialized" class="bg-gradient-to-r from-blue-500 to-purple-600 text-white p-6 shadow-lg sticky top-0 z-50">
    <div class="px-4 sm:px-6 lg:px-8">
      <div class="flex flex-col md:flex-row items-center justify-between">
        <!-- Logo et titre -->
        <div class="text-center md:text-left mb-4 md:mb-0">
          <RouterLink :to="'/'" class="block">
            <h1 class="text-3xl font-bold hover:text-blue-200 transition-colors">
            SafeBase
            </h1>
          </RouterLink>
        </div>
        
        <!-- Navigation -->
        <nav class="flex flex-wrap items-center justify-center md:justify-end gap-2 md:gap-4">
          <!-- Navigation pour utilisateurs authentifiés -->
          <template v-if="isAuthenticated">
            <!-- Lien Dashboard -->
            <RouterLink 
              :to="dashboardLink"
              class="bg-white/10 hover:bg-white/20 px-2 md:px-3 lg:px-4 py-2 rounded-lg transition-all duration-200 backdrop-blur-sm flex items-center space-x-1 md:space-x-2"
              active-class="bg-white/30"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"></path>
              </svg>
              <span class="hidden lg:inline">{{ isAdmin ? 'Admin' : 'Dashboard' }}</span>
            </RouterLink>

            <!-- Lien Gestion utilisateurs (admin uniquement) -->
            <RouterLink 
              v-if="isAdmin"
              to="/admin/users"
              class="bg-white/10 hover:bg-white/20 px-2 md:px-3 lg:px-4 py-2 rounded-lg transition-all duration-200 backdrop-blur-sm flex items-center space-x-1 md:space-x-2"
              active-class="bg-white/30"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"></path>
              </svg>
              <span class="hidden lg:inline">Utilisateurs</span>
            </RouterLink>

            <!-- Lien Bases de données (utilisateurs uniquement) -->
            <RouterLink 
              to="/user/databases"
              class="bg-white/10 hover:bg-white/20 px-2 md:px-3 lg:px-4 py-2 rounded-lg transition-all duration-200 backdrop-blur-sm flex items-center space-x-1 md:space-x-2"
              active-class="bg-white/30"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4"></path>
              </svg>
              <span class="hidden lg:inline">Bases de données</span>
            </RouterLink>

            <!-- Lien Sauvegardes (utilisateurs uniquement) -->
            <RouterLink 
              to="/user/backups"
              class="bg-white/10 hover:bg-white/20 px-2 md:px-3 lg:px-4 py-2 rounded-lg transition-all duration-200 backdrop-blur-sm flex items-center space-x-1 md:space-x-2"
              active-class="bg-white/30"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4"></path>
              </svg>
              <span class="hidden lg:inline">Sauvegardes</span>
            </RouterLink>

            <!-- Lien Planifications (utilisateurs uniquement) -->
            <RouterLink 
              to="/user/schedules"
              class="bg-white/10 hover:bg-white/20 px-2 md:px-3 lg:px-4 py-2 rounded-lg transition-all duration-200 backdrop-blur-sm flex items-center space-x-1 md:space-x-2"
              active-class="bg-white/30"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
              <span class="hidden lg:inline">Planifications</span>
            </RouterLink>

            <!-- Lien Historique (utilisateurs uniquement) -->
            <RouterLink 
              to="/user/history"
              class="bg-white/10 hover:bg-white/20 px-2 md:px-3 lg:px-4 py-2 rounded-lg transition-all duration-200 backdrop-blur-sm flex items-center space-x-1 md:space-x-2"
              active-class="bg-white/30"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v3m0 0v3m0-3h3m-3 0H9"></path>
              </svg>
              <span class="hidden lg:inline">Historique</span>
            </RouterLink>

            <!-- Dropdown utilisateur -->
            <div class="relative">
              <button 
                @click="showUserMenu = !showUserMenu"
                class="bg-white/10 hover:bg-white/20 px-3 md:px-4 py-2 rounded-lg transition-all duration-200 backdrop-blur-sm flex items-center space-x-2"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
                </svg>
                <span class="hidden md:inline">{{ currentUser?.firstname || 'Utilisateur' }}</span>
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
                </svg>
              </button>

              <!-- Menu dropdown -->
              <div 
                v-if="showUserMenu"
                @click.self="showUserMenu = false"
                class="absolute right-0 mt-2 w-56 bg-white rounded-lg shadow-xl overflow-hidden z-50"
              >
                <!-- User info -->
                <div class="px-4 py-3 border-b border-gray-200 bg-gray-50">
                  <p class="text-sm font-medium text-gray-900">
                    {{ currentUser?.firstname || 'Prénom' }} {{ currentUser?.lastname || 'Nom' }}
                  </p>
                  <p class="text-xs text-gray-500">{{ currentUser?.email || 'email@example.com' }}</p>
                  <span 
                    class="inline-block mt-2 text-xs px-2 py-1 rounded-full font-semibold"
                    :class="isAdmin ? 'bg-red-100 text-red-800' : 'bg-blue-100 text-blue-800'"
                  >
                    {{ isAdmin ? 'Administrateur' : 'Utilisateur' }}
                  </span>
                </div>

                <!-- Menu items -->
                <RouterLink 
                  to="/user/profile"
                  @click="showUserMenu = false"
                  class="block px-4 py-3 text-sm text-gray-700 hover:bg-gray-100 transition flex items-center space-x-2"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
                  </svg>
                  <span>Mon profil</span>
                </RouterLink>
                <button 
                  @click="handleLogout"
                  :disabled="loading"
                  class="w-full text-left px-4 py-3 text-sm text-red-600 hover:bg-red-50 transition flex items-center space-x-2 disabled:opacity-50"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"></path>
                  </svg>
                  <span>{{ loading ? 'Déconnexion...' : 'Déconnexion' }}</span>
                </button>
              </div>
            </div>
          </template>
          
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
import { ref, computed } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth'

// Composables
const router = useRouter()
const authStore = useAuthStore()

// State
const loading = ref<boolean>(false)
const showUserMenu = ref<boolean>(false)

// Computed réactifs depuis le store
const { isAuthenticated, isAdmin, user: currentUser, initialized } = storeToRefs(authStore)

// Computed pour le lien du dashboard selon le rôle
const dashboardLink = computed(() => {
  return isAdmin.value ? '/admin/dashboard' : '/user/dashboard'
})

// Methods
const handleLogout = async (): Promise<void> => {
  loading.value = true
  showUserMenu.value = false
  try {
    await authStore.logout()
    await router.push('/login')
  } catch (error) {
    console.error('Erreur de déconnexion:', error)
  } finally {
    loading.value = false
  }
}

// Fermer le menu si on clique ailleurs
if (typeof window !== 'undefined') {
  window.addEventListener('click', (e) => {
    const target = e.target as HTMLElement
    if (!target.closest('.relative')) {
      showUserMenu.value = false
    }
  })
}
</script>
