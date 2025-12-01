// Service de gestion des utilisateurs (Admin) - Logique métier
import * as userApi from '@/api/user_api'
import type { User } from '@/types/auth'
import type { UserUpdateRequest, UserRoleUpdateRequest } from '@/types/user'

/**
 * Service de gestion des utilisateurs pour les administrateurs
 */
export class UserService {
  /**
   * Récupère tous les utilisateurs
   */
  async fetchAllUsers(): Promise<User[]> {
    return await userApi.getAllUsers()
  }

  /**
   * Récupère tous les utilisateurs actifs
   */
  async fetchAllActiveUsers(): Promise<User[]> {
    return await userApi.getAllActiveUsers()
  }

  /**
   * Récupère un utilisateur par son ID
   */
  async fetchUserById(id: number): Promise<User> {
    return await userApi.getUserById(id)
  }

  /**
   * Met à jour un utilisateur avec validation
   */
  async updateUser(id: number, userData: UserUpdateRequest): Promise<void> {
    // Validation des données
    this.validateUserData(userData)

    await userApi.updateUser(id, userData)
  }

  /**
   * Change le rôle d'un utilisateur
   */
  async changeUserRole(id: number, roleData: UserRoleUpdateRequest): Promise<void> {
    // Validation du rôle
    this.validateRoleData(roleData)

    await userApi.changeUserRole(id, roleData)
  }

  /**
   * Désactive un utilisateur
   */
  async deactivateUser(id: number): Promise<void> {
    await userApi.deactivateUser(id)
  }

  /**
   * Active un utilisateur
   */
  async activateUser(id: number): Promise<void> {
    await userApi.activateUser(id)
  }

  /**
   * Valide les données d'un utilisateur
   */
  private validateUserData(userData: UserUpdateRequest): void {
    if (!userData.firstname?.trim()) {
      throw new Error('Le prénom est requis')
    }
    if (!userData.lastname?.trim()) {
      throw new Error('Le nom est requis')
    }
    if (!userData.email?.trim()) {
      throw new Error('L\'email est requis')
    }

    // Validation basique de l'email
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!emailRegex.test(userData.email)) {
      throw new Error('L\'email n\'est pas valide')
    }

    if (userData.role_id && userData.role_id <= 0) {
      throw new Error('L\'ID du rôle doit être positif')
    }
  }

  /**
   * Valide les données de changement de rôle
   */
  private validateRoleData(roleData: UserRoleUpdateRequest): void {
    if (!roleData.role_id || roleData.role_id <= 0) {
      throw new Error('L\'ID du rôle est requis et doit être positif')
    }
  }
}

// Instance singleton du service
export const userService = new UserService()