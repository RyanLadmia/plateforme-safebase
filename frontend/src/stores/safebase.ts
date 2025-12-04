// Store Pinia pour SafeBase - Gestion d'état centralisée
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { databaseService } from '@/services/database_service'
import { backupService } from '@/services/backup_service'
import { scheduleService } from '@/services/schedule_service'
import { userService } from '@/services/user_service'
import type { Database, DatabaseUpdateRequest } from '@/types/database'
import type { Backup } from '@/types/backup'
import type { Schedule } from '@/types/schedule'
import type { User } from '@/types/auth'
import type { UserUpdateRequest, UserRoleUpdateRequest } from '@/types/user'

export const useSafebaseStore = defineStore('safebase', () => {
  // État réactif
  const databases = ref<Database[]>([])
  const backups = ref<Backup[]>([])
  const schedules = ref<Schedule[]>([])
  const users = ref<User[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Getters calculés
  const databaseCount = computed(() => databases.value.length)
  const backupCount = computed(() => backups.value.length)
  const scheduleCount = computed(() => schedules.value.length)
  const userCount = computed(() => users.value.length)
  const activeUsers = computed(() => users.value.filter(u => u.active))
  const inactiveUsers = computed(() => users.value.filter(u => !u.active))
  const adminUsers = computed(() => users.value.filter(u => u.role?.name === 'admin'))
  const regularUsers = computed(() => users.value.filter(u => u.role?.name === 'user'))
  const activeSchedules = computed(() => 
    schedules.value.filter(s => s.active)
  )
  const inactiveSchedules = computed(() => 
    schedules.value.filter(s => !s.active)
  )
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

  const updateDatabaseAsync = async (id: number, databaseData: DatabaseUpdateRequest): Promise<void> => {
    loading.value = true
    error.value = null
    try {
      const updatedDatabase = await databaseService.updateDatabase(id, databaseData)
      updateDatabase(updatedDatabase)
    } catch (err: any) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const updateDatabasePartialAsync = async (id: number, updates: { name: string }): Promise<void> => {
    loading.value = true
    error.value = null
    try {
      const updatedDatabase = await databaseService.updateDatabasePartial(id, updates)
      updateDatabase(updatedDatabase)
    } catch (err: any) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const deleteDatabaseAsync = async (id: number): Promise<void> => {
    loading.value = true
    error.value = null
    try {
      await databaseService.deleteDatabase(id)
      removeDatabase(id)
    } catch (err: any) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const getDatabaseWithBackupCountAsync = async (id: number): Promise<{ database: Database; backup_count: number }> => {
    loading.value = true
    error.value = null
    try {
      return await databaseService.getDatabaseWithBackupCount(id)
    } catch (err: any) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
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

  // Actions pour les schedules
  const fetchSchedules = async (): Promise<void> => {
    loading.value = true
    error.value = null
    try {
      schedules.value = await scheduleService.fetchSchedules()
    } catch (err: any) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const addSchedule = (schedule: Schedule): void => {
    schedules.value.push(schedule)
  }

  const updateSchedule = (schedule: Schedule): void => {
    const index = schedules.value.findIndex(s => s.id === schedule.id)
    if (index !== -1) {
      schedules.value[index] = schedule
    }
  }

  const removeSchedule = (id: number): void => {
    schedules.value = schedules.value.filter(s => s.id !== id)
  }

  // Actions pour les utilisateurs (Admin)
  const fetchUsers = async (): Promise<void> => {
    loading.value = true
    error.value = null
    try {
      users.value = await userService.fetchAllUsers()
    } catch (err: any) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const fetchActiveUsers = async (): Promise<void> => {
    loading.value = true
    error.value = null
    try {
      users.value = await userService.fetchAllActiveUsers()
    } catch (err: any) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const fetchUserById = async (id: number): Promise<User> => {
    loading.value = true
    error.value = null
    try {
      return await userService.fetchUserById(id)
    } catch (err: any) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const updateUser = async (id: number, userData: UserUpdateRequest): Promise<void> => {
    loading.value = true
    error.value = null
    try {
      await userService.updateUser(id, userData)
      // Mettre à jour l'utilisateur dans la liste
      const index = users.value.findIndex(u => u.id === id)
      if (index !== -1) {
        const updatedUser = await userService.fetchUserById(id)
        users.value[index] = updatedUser
      }
    } catch (err: any) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const changeUserRole = async (id: number, roleData: UserRoleUpdateRequest): Promise<void> => {
    loading.value = true
    error.value = null
    try {
      await userService.changeUserRole(id, roleData)
      // Mettre à jour l'utilisateur dans la liste
      const index = users.value.findIndex(u => u.id === id)
      if (index !== -1) {
        const updatedUser = await userService.fetchUserById(id)
        users.value[index] = updatedUser
      }
    } catch (err: any) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const deactivateUser = async (id: number): Promise<void> => {
    loading.value = true
    error.value = null
    try {
      await userService.deactivateUser(id)
      // Mettre à jour le statut dans la liste
      const user = users.value.find(u => u.id === id)
      if (user) {
        user.active = false
      }
    } catch (err: any) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const activateUser = async (id: number): Promise<void> => {
    loading.value = true
    error.value = null
    try {
      await userService.activateUser(id)
      // Mettre à jour le statut dans la liste
      const user = users.value.find(u => u.id === id)
      if (user) {
        user.active = true
      }
    } catch (err: any) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
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
    schedules,
    users,
    loading,
    error,
    
    // Getters
    databaseCount,
    backupCount,
    scheduleCount,
    activeSchedules,
    inactiveSchedules,
    completedBackups,
    pendingBackups,
    failedBackups,
    userCount,
    activeUsers,
    inactiveUsers,
    adminUsers,
    
    // Actions - Databases
    fetchDatabases,
    addDatabase,
    updateDatabase,
    updateDatabaseAsync,
    updateDatabasePartialAsync,
    getDatabaseWithBackupCountAsync,
    deleteDatabaseAsync,
    removeDatabase,
    
    // Actions - Backups
    fetchBackups,
    fetchBackupsByDatabase,
    addBackup,
    updateBackup,
    removeBackup,
    
    // Actions - Schedules
    fetchSchedules,
    addSchedule,
    updateSchedule,
    removeSchedule,
    
    // Actions - Users (Admin)
    fetchUsers,
    fetchActiveUsers,
    fetchUserById,
    updateUser,
    changeUserRole,
    deactivateUser,
    activateUser,
    
    // Utilitaires
    clearError,
    reset
  }
})
