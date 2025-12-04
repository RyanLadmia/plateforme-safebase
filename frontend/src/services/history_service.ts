// Service de gestion de l'historique des actions - Logique métier
import * as historyApi from '@/api/history_api'
import type { HistoryItem, HistoryResponse, ActivityType } from '@/types/history'

/**
 * Service de gestion de l'historique des actions
 */
export class HistoryService {
  /**
   * Récupère l'historique des actions de l'utilisateur
   */
  async fetchUserHistory(page: number = 1, limit: number = 20): Promise<HistoryResponse> {
    return await historyApi.getUserActionHistory(page, limit)
  }

  /**
   * Récupère l'historique par type d'activité
   */
  async fetchHistoryByType(resourceType: ActivityType, page: number = 1, limit: number = 20): Promise<HistoryResponse> {
    if (resourceType === 'all') {
      return await this.fetchUserHistory(page, limit)
    }
    return await historyApi.getActionHistoryByType(resourceType, page, limit)
  }

  /**
   * Récupère l'historique pour une ressource spécifique
   */
  async fetchResourceHistory(resourceType: string, resourceId: number): Promise<HistoryItem[]> {
    const response = await historyApi.getResourceActionHistory(resourceType, resourceId)
    return response.history
  }

  /**
   * Récupère l'historique récent (admin uniquement)
   */
  async fetchRecentHistory(page: number = 1, limit: number = 20): Promise<HistoryResponse> {
    return await historyApi.getRecentActionHistory(page, limit)
  }

  /**
   * Formate une date pour l'affichage
   */
  formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString('fr-FR', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  /**
   * Formate la taille d'un fichier
   */
  formatFileSize(bytes: number): string {
    if (bytes === 0) return '0 Bytes'

    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))

    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

  /**
   * Obtient le libellé d'un type de ressource
   */
  getResourceTypeLabel(type: string): string {
    const labels: Record<string, string> = {
      database: 'Base de données',
      backup: 'Sauvegarde',
      schedule: 'Planification',
      restore: 'Restauration'
    }

    return labels[type] || type
  }

  /**
   * Obtient la classe CSS pour un type de ressource
   */
  getResourceTypeClass(type: string): string {
    const classes: Record<string, string> = {
      database: 'bg-purple-100 text-purple-800',
      backup: 'bg-green-100 text-green-800',
      schedule: 'bg-orange-100 text-orange-800',
      restore: 'bg-indigo-100 text-indigo-800'
    }

    return classes[type] || 'bg-gray-100 text-gray-800'
  }

  /**
   * Obtient le texte d'une action
   */
  getActionText(action: string): string {
    const actionTexts: Record<string, string> = {
      // For BDD
      created: 'Ajout',
      updated: 'Modification',
      deleted: 'Suppression',
      // For Backups/Restores/Schedules
      create: 'Création',
      update: 'Modification',
      delete: 'Suppression',
      completed: 'Succès',
      failed: 'Échec',
      executed: 'Exécution',
      download: 'Téléchargement',
      restored: 'Restauration'
    }

    return actionTexts[action] || action
  }

  /**
   * Obtient le contenu formaté spécial pour les bases de données
   */
  getDatabaseContent(item: HistoryItem): { title: string, subtitle: string, details: string[] } {
    const actionText = this.getActionText(item.action)
    
    // Récupérer les informations depuis les métadonnées
    let databaseName = 'Inconnu'
    let databaseType = 'Inconnu'
    
    if (item.metadata) {
      databaseName = item.metadata.database_name || item.metadata.name || 'Inconnu'
      databaseType = item.metadata.database_type || item.metadata.type || 'Inconnu'
    }
    
    // Essayer d'extraire depuis la description si les métadonnées ne sont pas disponibles
    if (databaseName === 'Inconnu' && item.description) {
      // La description contient souvent le nom de la base
      const match = item.description.match(/Base de données '([^']+)'/)
      if (match) {
        databaseName = match[1]
      }
    }
    
    // Formater le type
    const typeLabel = databaseType === 'postgresql' ? 'PostgreSQL' : 
                     databaseType === 'mysql' ? 'MySQL' : 
                     databaseType === 'postgres' ? 'PostgreSQL' :
                     databaseType.toUpperCase()
    
    // Construire les détails
    const details: string[] = []
    
    // Toujours ajouter le type
    details.push(`Type : ${typeLabel}`)
    
    // Vérifier s'il y a des changements détaillés
    if (item.metadata?.changes) {
      const changes = item.metadata.changes as Record<string, any>
      
      // Traiter les changements de nom
      if (changes.name) {
        const nameChange = changes.name
        if (nameChange.from && nameChange.to) {
          details.push(`Changements : nom de '${nameChange.from}' à '${nameChange.to}'`)
        }
      }
      
      // Ici on pourrait ajouter d'autres types de changements si nécessaire
      // (type, host, etc.) mais pour l'instant on se limite au nom
    }
    
    return {
      title: actionText,
      subtitle: `Base de données : ${databaseName}`,
      details: details
    }
  }

  /**
   * Obtient le contenu formaté spécial pour les sauvegardes
   */
  getBackupContent(item: HistoryItem): { title: string, subtitle: string, details: string[] } {
    const actionText = this.getActionText(item.action)
    
    // Récupérer les informations depuis les métadonnées
    let databaseName = 'Inconnu'
    let fileName = 'Inconnu'
    
    if (item.metadata) {
      databaseName = item.metadata.database_name || 'Inconnu'
      // Pour les téléchargements, le nom du fichier peut être dans différents champs
      fileName = item.metadata.file_name || item.metadata.filename || 'Inconnu'
    }
    
    // Essayer d'extraire depuis la description si les métadonnées ne sont pas disponibles
    if (fileName === 'Inconnu' && item.description) {
      // La description contient souvent le nom du fichier entre guillemets
      const fileMatch = item.description.match(/Sauvegarde '([^']+)'/)
      if (fileMatch) {
        fileName = fileMatch[1]
      }
    }
    
    if (databaseName === 'Inconnu' && item.description) {
      // La description contient souvent le nom de la base entre parenthèses
      const dbMatch = item.description.match(/Base de données: ([^)]+)/)
      if (dbMatch) {
        databaseName = dbMatch[1]
      }
    }
    
    return {
      title: actionText,
      subtitle: `Fichier : ${fileName}`,
      details: [`Base de données : ${databaseName}`]
    }
  }

  /**
   * Obtient le contenu formaté spécial pour les restaurations
   */
  getRestoreContent(item: HistoryItem): { title: string, subtitle: string, details: string[] } {
    const actionText = this.getActionText(item.action)
    
    // Récupérer les informations depuis les métadonnées
    let backupName = 'Inconnu'
    let databaseName = 'Inconnu'
    
    if (item.metadata) {
      backupName = item.metadata.backup_name || 'Inconnu'
      databaseName = item.metadata.database_name || 'Inconnu'
    }
    
    return {
      title: actionText,
      subtitle: `Sauvegarde utilisée : ${backupName}`,
      details: [`Base de donnée : ${databaseName}`]
    }
  }

  /**
   * Obtient la classe CSS pour une action
   */
  getActionIconClass(action: string): string {
    const colorClasses: Record<string, string> = {
      created: 'bg-green-500',
      updated: 'bg-blue-500',
      deleted: 'bg-red-500',
      create: 'bg-green-500',
      update: 'bg-blue-500',
      delete: 'bg-red-500',
      completed: 'bg-green-500',
      failed: 'bg-red-500',
      executed: 'bg-blue-500',
      download: 'bg-purple-500',
      restored: 'bg-green-500'
    }

    return colorClasses[action] || 'bg-gray-500'
  }

  /**
   * Vérifie si un élément historique concerne une base de données
   */
  isDatabaseItem(item: HistoryItem): boolean {
    return item.resource_type === 'database'
  }

  /**
   * Vérifie si un élément historique concerne une sauvegarde
   */
  isBackupItem(item: HistoryItem): boolean {
    return item.resource_type === 'backup'
  }

  /**
   * Vérifie si un élément historique concerne une restauration
   */
  isRestoreItem(item: HistoryItem): boolean {
    return item.resource_type === 'restore'
  }
}

// Instance singleton
export const historyService = new HistoryService()