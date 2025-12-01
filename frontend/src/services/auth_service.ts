// Service d'authentification - Logique métier
import * as authApi from '@/api/auth_api'
import type { User, LoginRequest, RegisterRequest } from '@/types/auth'

/**
 * Service d'authentification qui encapsule la logique métier
 */
export class AuthService {
  /**
   * Vérifie si l'utilisateur est authentifié
   */
  async checkAuthentication(): Promise<User | null> {
    return await authApi.checkAuth()
  }

  /**
   * Inscrit un nouvel utilisateur
   */
  async registerUser(userData: RegisterRequest): Promise<{ success: boolean; message: string }> {
    try {
      const response = await authApi.register(userData)
      return {
        success: true,
        message: response.message || 'Inscription réussie'
      }
    } catch (error: any) {
      return {
        success: false,
        message: error.message || 'Erreur lors de l\'inscription'
      }
    }
  }

  /**
   * Connecte un utilisateur
   * Le token JWT est automatiquement géré via cookie HTTP-only sécurisé
   */
  async loginUser(credentials: LoginRequest): Promise<User> {
    return await authApi.login(credentials)
  }

  /**
   * Déconnecte l'utilisateur
   */
  async logoutUser(): Promise<void> {
    await authApi.logout()
  }
}

// Export d'une instance unique du service
export const authService = new AuthService()
