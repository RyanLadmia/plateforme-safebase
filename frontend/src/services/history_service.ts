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
      created: 'Créé',
      updated: 'Modifié',
      deleted: 'Supprimé',
      create: 'Créé',
      update: 'Modifié',
      delete: 'Supprimé',
      completed: 'Terminé',
      failed: 'Échoué',
      executed: 'Exécuté',
      download: 'Téléchargé'
    }

    return actionTexts[action] || action
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
      download: 'bg-purple-500'
    }

    return colorClasses[action] || 'bg-gray-500'
  }
}

// Instance singleton
export const historyService = new HistoryService()