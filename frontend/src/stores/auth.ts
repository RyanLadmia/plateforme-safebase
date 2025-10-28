// Store Pinia sécurisé pour l'authentification avec cookies HTTP-only
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, LoginRequest, RegisterRequest } from '@/types/auth'

// Configuration sécurisée via variables d'environnement
const API_BASE_URL: string = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

export const useAuthStore = defineStore('auth', () => {
  // State réactif
  const user = ref<User | null>(null)
  const loading = ref<boolean>(false)
  const initialized = ref<boolean>(false)

  // Getters
  const isAuthenticated = computed(() => user.value !== null)

  // Actions
  const checkAuth = async (): Promise<void> => {
    try {
      const response = await fetch(`${API_BASE_URL}/auth/me`, {
        method: 'GET',
        credentials: 'include', // Important: inclut les cookies
        headers: {
          'Content-Type': 'application/json',
        },
      })

      if (response.status === 200) {
        const data = await response.json()
        user.value = data.user
      } else if (response.status === 401) {
        // 401 Unauthorized est normal quand l'utilisateur n'est pas connecté
        user.value = null
      } else {
        // Autres erreurs (500, 403, etc.)
        console.warn(`⚠️ Réponse inattendue du serveur: ${response.status}`)
        user.value = null
      }
    } catch (error) {
      // Erreur réseau (serveur inaccessible, etc.)
      if (error instanceof TypeError && error.message.includes('Failed to fetch')) {
        // Serveur inaccessible - message informatif
        console.info('ℹ️ Serveur d\'authentification inaccessible')
      } else {
        console.error('❌ Erreur lors de la vérification de l\'authentification:', error)
      }
      user.value = null
    } finally {
      initialized.value = true
    }
  }

  const register = async (userData: RegisterRequest): Promise<void> => {
    loading.value = true
    try {
      const response = await fetch(`${API_BASE_URL}/auth/register`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(userData),
      })

      if (!response.ok) {
        const error = await response.json()
        throw new Error(error.error || 'Erreur lors de l\'inscription')
      }

      // L'inscription ne connecte pas automatiquement l'utilisateur
      const data = await response.json()
      return data
    } finally {
      loading.value = false
    }
  }

  const login = async (credentials: LoginRequest): Promise<void> => {
    loading.value = true
    try {
      const response = await fetch(`${API_BASE_URL}/auth/login`, {
        method: 'POST',
        credentials: 'include', // Important: inclut les cookies
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(credentials),
      })

      if (!response.ok) {
        const error = await response.json()
        throw new Error(error.error || 'Erreur lors de la connexion')
      }

      const data = await response.json()
      user.value = data.user // Le serveur retourne les infos utilisateur
    } finally {
      loading.value = false
    }
  }

  const logout = async (): Promise<void> => {
    loading.value = true
    try {
      const response = await fetch(`${API_BASE_URL}/auth/logout`, {
        method: 'POST',
        credentials: 'include', // Important: inclut les cookies
        headers: {
          'Content-Type': 'application/json',
        },
      })

      if (!response.ok) {
        const error = await response.json()
        throw new Error(error.error || 'Erreur lors de la déconnexion')
      }

      user.value = null
    } finally {
      loading.value = false
    }
  }

  // Nettoyage sécuritaire : supprimer tout token résiduel du localStorage
  const clearInsecureStorage = () => {
    // Supprimer tous les tokens potentiels du localStorage (sécurité)
    const keysToRemove = ['auth_token', 'token', 'jwt_token', 'access_token', 'user_token']
    keysToRemove.forEach(key => {
      if (localStorage.getItem(key)) {
        console.warn(`🚨 Token non sécurisé trouvé dans localStorage (${key}) - suppression automatique`)
        localStorage.removeItem(key)
      }
    })
    
    // Idem pour sessionStorage
    keysToRemove.forEach(key => {
      if (sessionStorage.getItem(key)) {
        console.warn(`🚨 Token non sécurisé trouvé dans sessionStorage (${key}) - suppression automatique`)
        sessionStorage.removeItem(key)
      }
    })
  }

  // Initialisation - nettoyage puis vérification de l'authentification
  clearInsecureStorage()
  checkAuth()

  return {
    // State
    user,
    loading,
    initialized,
    
    // Getters
    isAuthenticated,
    
    // Actions
    checkAuth,
    register,
    login,
    logout
  }
})
