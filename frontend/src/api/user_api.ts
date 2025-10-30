// API pour la gestion des utilisateurs (Admin uniquement) - Appels réseau purs avec Axios
import { apiClient } from './axios'
import type { User } from '@/types/auth'

export interface UserListResponse {
  users: User[]
}

/**
 * Récupère tous les utilisateurs (Admin uniquement)
 * TODO: À implémenter côté backend
 */
export async function getAllUsers(): Promise<User[]> {
  const { data } = await apiClient.get<UserListResponse>('/api/admin/users')
  return data.users || []
}

/**
 * Supprime un utilisateur (Admin uniquement)
 * TODO: À implémenter côté backend
 */
export async function deleteUser(userId: number): Promise<void> {
  await apiClient.delete(`/api/admin/users/${userId}`)
}
