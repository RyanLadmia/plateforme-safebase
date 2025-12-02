// API pour la gestion des sauvegardes - Appels réseau purs avec Axios
import { apiClient } from './axios'
import type { Backup, BackupResponse, BackupListResponse } from '@/types/backup'

/**
 * Récupère toutes les sauvegardes de l'utilisateur
 */
export async function getBackups(): Promise<Backup[]> {
  const { data } = await apiClient.get<BackupListResponse>('/api/backups')
  return data.backups || []
}

/**
 * Récupère les sauvegardes d'une base de données spécifique
 */
export async function getBackupsByDatabase(databaseId: number): Promise<Backup[]> {
  const { data } = await apiClient.get<BackupListResponse>(`/api/backups/database/${databaseId}`)
  return data.backups || []
}

/**
 * Récupère une sauvegarde par son ID
 */
export async function getBackupById(id: number): Promise<Backup> {
  const { data } = await apiClient.get<BackupResponse>(`/api/backups/${id}`)
  return data.backup
}

/**
 * Crée une nouvelle sauvegarde pour une base de données
 */
export async function createBackup(databaseId: number): Promise<Backup> {
  const { data } = await apiClient.post<BackupResponse>(`/api/backups/database/${databaseId}`)
  return data.backup
}

/**
 * Supprime une sauvegarde
 */
export async function deleteBackup(id: number): Promise<void> {
  await apiClient.delete(`/api/backups/${id}`)
}

/**
 * Télécharge une sauvegarde
 */
export async function downloadBackup(id: number, filename: string): Promise<void> {
  const { data } = await apiClient.get(`/api/backups/${id}/download`, {
    responseType: 'blob', // Important pour télécharger des fichiers
  })

  // Créer un blob et déclencher le téléchargement
  const url = window.URL.createObjectURL(new Blob([data]))
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  
  // Nettoyage
  window.URL.revokeObjectURL(url)
  document.body.removeChild(link)
}

/**
 * Lance une restauration d'une sauvegarde vers une base de données
 */
export async function restoreBackup(backupId: number, databaseId: number): Promise<any> {
  const { data } = await apiClient.post(`/api/restores/backup/${backupId}/database/${databaseId}`)
  return data.restore
}

/**
 * Récupère toutes les restaurations de l'utilisateur
 */
export async function getRestores(): Promise<any[]> {
  const { data } = await apiClient.get('/api/restores')
  return data.restores || []
}

/**
 * Récupère une restauration par son ID
 */
export async function getRestoreById(id: number): Promise<any> {
  const { data } = await apiClient.get(`/api/restores/${id}`)
  return data.restore
}
