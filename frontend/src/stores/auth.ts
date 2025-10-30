// Store Pinia pour l'authentification - Gestion d'état uniquement
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authService } from '@/services/auth_service'
import type { User, LoginRequest, RegisterRequest } from '@/types/auth'

export const useAuthStore = defineStore('auth', () => {
  // État réactif
  const user = ref<User | null>(null)
  const loading = ref<boolean>(false)
  const initialized = ref<boolean>(false)

  // Getters calculés
  const isAuthenticated = computed(() => user.value !== null)
  const isAdmin = computed(() => user.value?.role?.name === 'admin')
  const isUser = computed(() => user.value?.role?.name === 'user' || user.value?.role?.name === 'admin')

  // Actions utilisant les services
  const checkAuth = async (): Promise<void> => {
    try {
      const authenticatedUser = await authService.checkAuthentication()
      user.value = authenticatedUser
    } catch (error) {
      console.error('Erreur lors de la vérification de l\'authentification:', error)
      user.value = null
    } finally {
      initialized.value = true
    }
  }

  const register = async (userData: RegisterRequest): Promise<{ success: boolean; message: string }> => {
    loading.value = true
    try {
      return await authService.registerUser(userData)
    } finally {
      loading.value = false
    }
  }

  const login = async (credentials: LoginRequest): Promise<void> => {
    loading.value = true
    try {
      const authenticatedUser = await authService.loginUser(credentials)
      user.value = authenticatedUser
    } finally {
      loading.value = false
    }
  }

  const logout = async (): Promise<void> => {
    loading.value = true
    try {
      await authService.logoutUser()
      user.value = null
    } finally {
      loading.value = false
    }
  }

  // Initialisation : Vérifier l'authentification au chargement
  checkAuth()

  return {
    // État
    user,
    loading,
    initialized,
    
    // Getters
    isAuthenticated,
    isAdmin,
    isUser,
    
    // Actions
    checkAuth,
    register,
    login,
    logout
  }
})
