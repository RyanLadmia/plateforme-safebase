// Types pour l'authentification

export interface User {
  id: number
  firstname: string
  lastname: string
  email: string
  role_id: number
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
