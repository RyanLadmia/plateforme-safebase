// Service de gestion des sauvegardes - Logique métier
import * as backupApi from '@/api/backup_api'
import type { Backup } from '@/types/backup'

/**
 * Service de gestion des sauvegardes
 */
export class BackupService {
  /**
   * Récupère toutes les sauvegardes
   */
  async fetchBackups(): Promise<Backup[]> {
    return await backupApi.getBackups()
  }

  /**
   * Récupère les sauvegardes d'une base de données
   */
  async fetchBackupsByDatabase(databaseId: number): Promise<Backup[]> {
    return await backupApi.getBackupsByDatabase(databaseId)
  }

  /**
   * Récupère une sauvegarde par son ID
   */
  async fetchBackupById(id: number): Promise<Backup> {
    return await backupApi.getBackupById(id)
  }

  /**
   * Crée une nouvelle sauvegarde
   */
  async createBackup(databaseId: number): Promise<Backup> {
    if (!databaseId || databaseId <= 0) {
      throw new Error('ID de base de données invalide')
    }
    
    return await backupApi.createBackup(databaseId)
  }

  /**
   * Supprime une sauvegarde
   */
  async deleteBackup(id: number): Promise<void> {
    await backupApi.deleteBackup(id)
  }

  /**
   * Télécharge une sauvegarde
   */
  async downloadBackup(backup: Backup): Promise<void> {
    if (backup.status !== 'completed') {
      throw new Error('La sauvegarde n\'est pas encore terminée')
    }
    
    await backupApi.downloadBackup(backup.id, backup.filename)
  }

  /**
   * Formate la taille d'un fichier en format lisible
   */
  formatFileSize(bytes: number): string {
    if (bytes === 0) return '0 B'
    
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
  }

  /**
   * Obtient le libellé du statut en français
   */
  getStatusLabel(status: string): string {
    const labels: Record<string, string> = {
      'pending': 'En cours',
      'completed': 'Terminée',
      'failed': 'Échouée'
    }
    
    return labels[status] || status
  }

  /**
   * Obtient la couleur du statut pour l'affichage
   */
  getStatusColor(status: string): string {
    const colors: Record<string, string> = {
      'pending': 'orange',
      'completed': 'green',
      'failed': 'red'
    }
    
    return colors[status] || 'gray'
  }

  /**
   * Vérifie si une sauvegarde peut être téléchargée
   */
  canDownload(backup: Backup): boolean {
    return backup.status === 'completed' && backup.size > 0
  }

  /**
   * Filtre les sauvegardes par statut
   */
  filterByStatus(backups: Backup[], status: string): Backup[] {
    return backups.filter(b => b.status === status)
  }

  /**
   * Trie les sauvegardes par date (plus récentes en premier)
   */
  sortByDate(backups: Backup[], descending: boolean = true): Backup[] {
    return [...backups].sort((a, b) => {
      const dateA = new Date(a.created_at).getTime()
      const dateB = new Date(b.created_at).getTime()
      return descending ? dateB - dateA : dateA - dateB
    })
  }

  /**
   * Calcule la taille totale des sauvegardes
   */
  getTotalSize(backups: Backup[]): number {
    return backups.reduce((total, backup) => total + backup.size, 0)
  }

  /**
   * Restaure une sauvegarde vers sa base de données d'origine
   */
  async restoreBackup(backup: Backup): Promise<any> {
    if (backup.status !== 'completed') {
      throw new Error('La sauvegarde n\'est pas encore terminée')
    }

    return await backupApi.restoreBackup(backup.id, backup.database_id)
  }

  /**
   * Détermine si une sauvegarde est manuelle ou automatique
   */
  getBackupType(backup: Backup): 'manual' | 'automatic' {
    // Si user_agent est "Scheduled Task", c'est automatique, sinon manuel
    return backup.user_agent === 'Scheduled Task' ? 'automatic' : 'manual'
  }

  /**
   * Obtient le libellé du type de sauvegarde
   */
  getBackupTypeLabel(backup: Backup): string {
    const type = this.getBackupType(backup)
    return type === 'automatic' ? 'Automatique' : 'Manuelle'
  }

  /**
   * Obtient la classe CSS pour le type de sauvegarde
   */
  getBackupTypeClass(backup: Backup): string {
    const type = this.getBackupType(backup)
    return type === 'automatic' ? 'bg-blue-100 text-blue-800' : 'bg-purple-100 text-purple-800'
  }
}

// Export d'une instance unique du service
export const backupService = new BackupService()
