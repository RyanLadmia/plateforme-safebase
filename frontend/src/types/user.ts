// Types pour la gestion des utilisateurs (Admin)
import type { User } from './auth'

// Réexporter User pour faciliter les imports
export type { User }

export interface UserUpdateRequest {
  firstname: string
  lastname: string
  email: string
  role_id?: number
}

export interface UserRoleUpdateRequest {
  role_id: number
}

export interface UserListResponse {
  users: User[]
}

export interface UserResponse {
  user: User
}

export interface MessageResponse {
  message: string
}