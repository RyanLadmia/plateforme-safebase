// Types pour l'authentification

export interface Role {
  id: number
  name: string
  created_at: string
  updated_at: string
}

export interface User {
  id: number
  firstname: string
  lastname: string
  email: string
  role_id: number
  role?: Role // Objet role complet du backend
  created_at: string
  updated_at: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  firstname: string
  lastname: string
  email: string
  password: string
}

export interface AuthResponse {
  token: string
}

export interface RegisterResponse {
  id: number
  firstname: string
  lastname: string
  email: string
  role_id: number
}

export type MessageType = 'success' | 'error'
export type FormType = 'login' | 'register'
