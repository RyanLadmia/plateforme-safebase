// Service de gestion des bases de données - Logique métier
import * as databaseApi from '@/api/database_api'
import type { Database, DatabaseCreateRequest, DatabaseUpdateRequest } from '@/types/database'

/**
 * Service de gestion des bases de données
 */
export class DatabaseService {
  /**
   * Récupère toutes les bases de données
   */
  async fetchDatabases(): Promise<Database[]> {
    return await databaseApi.getDatabases()
  }

  /**
   * Récupère une base de données par son ID
   */
  async fetchDatabaseById(id: number): Promise<Database> {
    return await databaseApi.getDatabaseById(id)
  }

  /**
   * Crée une nouvelle base de données avec validation
   */
  async createDatabase(databaseData: DatabaseCreateRequest): Promise<Database> {
    // Validation des données
    this.validateDatabaseData(databaseData)
    
    return await databaseApi.createDatabase(databaseData)
  }

  /**
   * Met à jour une base de données
   */
  async updateDatabase(id: number, databaseData: DatabaseUpdateRequest): Promise<Database> {
    // Validation des données
    this.validateDatabaseData(databaseData)
    
    return await databaseApi.updateDatabase(id, databaseData)
  }

  /**
   * Met à jour partiellement une base de données (seulement le nom pour la sécurité)
   */
  async updateDatabasePartial(id: number, updates: { name: string }): Promise<Database> {
    // Validation simple du nom
    if (!updates.name || updates.name.trim() === '') {
      throw new Error('Le nom de la base de données est requis')
    }

    return await databaseApi.updateDatabasePartial(id, { name: updates.name.trim() })
  }

  /**
   * Valide les données de la base de données
   */
  private validateDatabaseData(data: DatabaseCreateRequest | DatabaseUpdateRequest): void {
    if (!data.name || data.name.trim() === '') {
      throw new Error('Le nom de la base de données est requis')
    }

    if (!data.type || !['mysql', 'postgresql'].includes(data.type)) {
      throw new Error('Le type de base de données doit être mysql ou postgresql')
    }

    // Si une URL est fournie, les champs individuels ne sont pas requis
    if (data.url && data.url.trim() !== '') {
      // Validation basique de l'URL
      if (!data.url.includes('://')) {
        throw new Error('L\'URL doit être au format: mysql://user:pass@host:port/db ou postgresql://user:pass@host:port/db')
      }
      return // Pas besoin de valider les champs individuels
    }

    // Si pas d'URL, tous les champs individuels sont requis
    if (!data.host || data.host.trim() === '') {
      throw new Error('L\'hôte est requis')
    }

    if (!data.port || data.port.trim() === '') {
      throw new Error('Le port est requis')
    }

    if (!data.username || data.username.trim() === '') {
      throw new Error('Le nom d\'utilisateur est requis')
    }

    if (!data.db_name || data.db_name.trim() === '') {
      throw new Error('Le nom de la base de données est requis')
    }

    // Validation du port (doit être un nombre)
    const portNum = parseInt(data.port)
    if (isNaN(portNum) || portNum < 1 || portNum > 65535) {
      throw new Error('Le port doit être un nombre entre 1 et 65535')
    }
  }

  /**
   * Obtient le port par défaut selon le type de base de données
   */
  getDefaultPort(type: 'mysql' | 'postgresql'): string {
    return type === 'mysql' ? '3306' : '5432'
  }

  /**
   * Teste la connexion à une base de données
   * TODO: À implémenter côté backend
   */
  async testConnection(databaseData: DatabaseCreateRequest): Promise<boolean> {
    // Cette fonctionnalité nécessite un endpoint backend dédié
    console.warn('Test de connexion non implémenté côté backend')
    return true
  }

  /**
   * Supprime une base de données
   */
  async deleteDatabase(id: number): Promise<void> {
    await databaseApi.deleteDatabase(id)
  }
}

// Export d'une instance unique du service
export const databaseService = new DatabaseService()
