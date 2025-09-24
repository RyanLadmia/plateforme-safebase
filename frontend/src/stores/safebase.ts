import { defineStore } from 'pinia'
import { ref } from 'vue'

interface Database {
  id: string
  name: string
  type: string
}

interface Backup {
  id: string
  databaseId: string
  createdAt: Date
  size: number
}

export const useSafebaseStore = defineStore('safebase', () => {
  // État
  const databases = ref<Database[]>([])
  const backups = ref<Backup[]>([])
  const isLoading = ref(false)
  
  // Actions
  function setLoading(loading: boolean) {
    isLoading.value = loading
  }
  
  function addDatabase(database: Database) {
    databases.value.push(database)
  }
  
  function addBackup(backup: Backup) {
    backups.value.push(backup)
  }

  return {
    databases,
    backups,
    isLoading,
    setLoading,
    addDatabase,
    addBackup
  }
})
