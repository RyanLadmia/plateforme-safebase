// API pour la gestion des utilisateurs (Admin uniquement) - Appels réseau purs avec Axios
import { apiClient } from './axios'
import type { User } from '@/types/auth'
import type { UserUpdateRequest, UserRoleUpdateRequest, UserListResponse, UserResponse, MessageResponse } from '@/types/user'

/**
 * Récupère tous les utilisateurs (Admin uniquement)
 */
export async function getAllUsers(): Promise<User[]> {
  const { data } = await apiClient.get<UserListResponse>('/api/admin/users')
  return data.users || []
}

/**
 * Récupère tous les utilisateurs actifs (Admin uniquement)
 */
export async function getAllActiveUsers(): Promise<User[]> {
  const { data } = await apiClient.get<UserListResponse>('/api/admin/users/active')
  return data.users || []
}

/**
 * Récupère un utilisateur par ID (Admin uniquement)
 */
export async function getUserById(userId: number): Promise<User> {
  const { data } = await apiClient.get<UserResponse>(`/api/admin/users/${userId}`)
  return data.user
}

/**
 * Met à jour un utilisateur (Admin uniquement)
 */
export async function updateUser(userId: number, userData: UserUpdateRequest): Promise<MessageResponse> {
  const { data } = await apiClient.put<MessageResponse>(`/api/admin/users/${userId}`, userData)
  return data
}

/**
 * Change le rôle d'un utilisateur (Admin uniquement)
 */
export async function changeUserRole(userId: number, roleData: UserRoleUpdateRequest): Promise<MessageResponse> {
  const { data } = await apiClient.put<MessageResponse>(`/api/admin/users/${userId}/role`, roleData)
  return data
}

/**
 * Désactive un utilisateur (Admin uniquement)
 */
export async function deactivateUser(userId: number): Promise<MessageResponse> {
  const { data } = await apiClient.put<MessageResponse>(`/api/admin/users/${userId}/deactivate`)
  return data
}

/**
 * Active un utilisateur (Admin uniquement)
 */
export async function activateUser(userId: number): Promise<MessageResponse> {
  const { data } = await apiClient.put<MessageResponse>(`/api/admin/users/${userId}/activate`)
  return data
}
