// API pour l'authentification - Appels réseau purs avec Axios
import { apiClient } from './axios'
import type { User, LoginRequest, RegisterRequest } from '@/types/auth'

export interface AuthResponse {
  user: User
  token?: string
  message?: string
}

export interface RegisterResponse {
  message: string
  user?: User
}

/**
 * Vérifie l'authentification de l'utilisateur actuel
 */
export async function checkAuth(): Promise<User | null> {
  try {
    const { data } = await apiClient.get<AuthResponse>('/auth/me')
    return data.user
  } catch (error: any) {
    // 401 est normal quand l'utilisateur n'est pas connecté
    if (error.status === 401) {
      return null
    }
    // Pour les autres erreurs, logger et retourner null
    console.warn('Erreur lors de la vérification de l\'authentification:', error.message)
    return null
  }
}

/**
 * Inscription d'un nouvel utilisateur
 */
export async function register(userData: RegisterRequest): Promise<RegisterResponse> {
  const { data } = await apiClient.post<RegisterResponse>('/auth/register', userData)
  return data
}

/**
 * Connexion d'un utilisateur
 * Le token JWT est automatiquement stocké dans un cookie HTTP-only sécurisé par le backend
 */
export async function login(credentials: LoginRequest): Promise<User> {
  const { data } = await apiClient.post<AuthResponse>('/auth/login', credentials)
  
  // Le token est déjà dans le cookie HTTP-only (géré par le backend)
  // Pas besoin de le stocker côté frontend (plus sécurisé)
  return data.user
}

/**
 * Déconnexion de l'utilisateur
 */
export async function logout(): Promise<void> {
  // Le cookie sera supprimé par le backend
  await apiClient.post('/auth/logout')
}
