// Store Pinia pour SafeBase - Gestion d'état centralisée
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { databaseService } from '@/services/database_service'
import { backupService } from '@/services/backup_service'
import type { Database } from '@/types/database'
import type { Backup } from '@/types/backup'

export const useSafebaseStore = defineStore('safebase', () => {
  // État réactif
  const databases = ref<Database[]>([])
  const backups = ref<Backup[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Getters calculés
  const databaseCount = computed(() => databases.value.length)
  const backupCount = computed(() => backups.value.length)
  const completedBackups = computed(() => 
    backups.value.filter(b => b.status === 'completed')
  )
  const pendingBackups = computed(() => 
    backups.value.filter(b => b.status === 'pending')
  )
  const failedBackups = computed(() => 
    backups.value.filter(b => b.status === 'failed')
  )

  // Actions pour les bases de données
  const fetchDatabases = async (): Promise<void> => {
    loading.value = true
    error.value = null
    try {
      databases.value = await databaseService.fetchDatabases()
    } catch (err: any) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const addDatabase = (database: Database): void => {
    databases.value.push(database)
  }

  const updateDatabase = (database: Database): void => {
    const index = databases.value.findIndex(d => d.id === database.id)
    if (index !== -1) {
      databases.value[index] = database
    }
  }

  const removeDatabase = (id: number): void => {
    databases.value = databases.value.filter(d => d.id !== id)
  }

  // Actions pour les sauvegardes
  const fetchBackups = async (): Promise<void> => {
    loading.value = true
    error.value = null
    try {
      backups.value = await backupService.fetchBackups()
    } catch (err: any) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const fetchBackupsByDatabase = async (databaseId: number): Promise<void> => {
    loading.value = true
    error.value = null
    try {
      const databaseBackups = await backupService.fetchBackupsByDatabase(databaseId)
      // Mettre à jour uniquement les sauvegardes de cette base de données
      backups.value = backups.value.filter(b => b.database_id !== databaseId)
      backups.value.push(...databaseBackups)
    } catch (err: any) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const addBackup = (backup: Backup): void => {
    backups.value.push(backup)
  }

  const updateBackup = (backup: Backup): void => {
    const index = backups.value.findIndex(b => b.id === backup.id)
    if (index !== -1) {
      backups.value[index] = backup
    }
  }

  const removeBackup = (id: number): void => {
    backups.value = backups.value.filter(b => b.id !== id)
  }

  // Utilitaires
  const clearError = (): void => {
    error.value = null
  }

  const reset = (): void => {
    databases.value = []
    backups.value = []
    loading.value = false
    error.value = null
  }

  return {
    // État
    databases,
    backups,
    loading,
    error,
    
    // Getters
    databaseCount,
    backupCount,
    completedBackups,
    pendingBackups,
    failedBackups,
    
    // Actions - Databases
    fetchDatabases,
    addDatabase,
    updateDatabase,
    removeDatabase,
    
    // Actions - Backups
    fetchBackups,
    fetchBackupsByDatabase,
    addBackup,
    updateBackup,
    removeBackup,
    
    // Utilitaires
    clearError,
    reset
  }
})
