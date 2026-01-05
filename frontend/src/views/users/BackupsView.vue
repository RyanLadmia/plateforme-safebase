<template>
  <div class="min-h-screen bg-gray-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <div class="mb-6">
        <h1 class="text-2xl sm:text-3xl font-bold text-gray-900">Mes sauvegardes</h1>
      </div>

      <!-- Filters -->
      <div class="bg-white rounded-lg shadow p-4 mb-6">
        <div class="flex flex-wrap justify-between items-center gap-4">
          <!-- Status filters -->
          <div class="flex gap-2">
            <button 
              @click="filterStatus = null" 
              :class="filterStatus === null ? 'bg-blue-600 text-white' : 'bg-gray-200'"
              class="px-3 py-2 sm:px-4 rounded-lg text-sm font-medium transition-colors duration-200"
            >
              Toutes ({{ backups.length }})
            </button>
            <button 
              @click="filterStatus = 'completed'" 
              :class="filterStatus === 'completed' ? 'bg-green-600 text-white' : 'bg-gray-200'"
              class="px-3 py-2 sm:px-4 rounded-lg text-sm font-medium transition-colors duration-200"
            >
              Terminées ({{ completedBackups.length }})
            </button>
            <button 
              @click="filterStatus = 'pending'" 
              :class="filterStatus === 'pending' ? 'bg-orange-600 text-white' : 'bg-gray-200'"
              class="px-3 py-2 sm:px-4 rounded-lg text-sm font-medium transition-colors duration-200"
            >
              En cours ({{ pendingBackups.length }})
            </button>
            <button 
              @click="filterStatus = 'failed'" 
              :class="filterStatus === 'failed' ? 'bg-red-600 text-white' : 'bg-gray-200'"
              class="px-3 py-2 sm:px-4 rounded-lg text-sm font-medium transition-colors duration-200"
            >
              Échouées ({{ failedBackups.length }})
            </button>
          </div>

          <!-- Database filter -->
          <div class="flex items-center gap-2">
            <label class="text-sm font-medium text-gray-700">Base de données:</label>
            <select 
              v-model="filterDatabaseId" 
              class="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="">Toutes les bases</option>
              <option 
                v-for="db in databases" 
                :key="db.id" 
                :value="db.id"
              >
                {{ db.name }} ({{ db.type }})
              </option>
            </select>
          </div>
        </div>
      </div>

      <!-- Statistics -->
      <div class="my-6 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="bg-white rounded-lg shadow p-4">
          <p class="text-gray-500 text-sm">Taille totale</p>
          <p class="text-2xl font-bold text-gray-900">{{ formatTotalSize() }}</p>
        </div>
        <div class="bg-white rounded-lg shadow p-4">
          <p class="text-gray-500 text-sm">Sauvegardes réussies</p>
          <p class="text-2xl font-bold text-green-600">{{ filteredCompletedBackups.length }}</p>
        </div>
        <div class="bg-white rounded-lg shadow p-4">
          <p class="text-gray-500 text-sm">En cours</p>
          <p class="text-2xl font-bold text-orange-600">{{ filteredPendingBackups.length }}</p>
        </div>
        <div class="bg-white rounded-lg shadow p-4">
          <p class="text-gray-500 text-sm">Échouées</p>
          <p class="text-2xl font-bold text-red-600">{{ filteredFailedBackups.length }}</p>
        </div>
      </div>

      <!-- Backups List -->
      <div v-if="loading" class="text-center py-12">Chargement...</div>
      <div v-else-if="error" class="bg-red-100 text-red-700 p-4 rounded-lg">{{ error }}</div>
      <div v-else-if="filteredBackups.length === 0" class="text-center py-12">
        <p class="text-gray-500">Aucune sauvegarde trouvée</p>
      </div>
      <div v-else class="bg-white rounded-lg shadow overflow-hidden">
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Fichier</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Taille</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Type</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Statut</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="backup in filteredBackups" :key="backup.id">
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  <div class="max-w-xs truncate" :title="backup.filename">
                    {{ backup.filename }}
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {{ formatDate(backup.created_at) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {{ formatSize(backup.size) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span 
                    :class="getBackupTypeClass(backup)"
                    class="px-2 py-1 text-xs font-semibold rounded-full inline-block"
                  >
                    {{ getBackupTypeLabel(backup) }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex flex-col space-y-1">
                    <span 
                      :class="getStatusClass(backup.status)"
                      class="px-2 py-1 text-xs font-semibold rounded-full inline-block w-fit"
                    >
                      {{ getStatusLabel(backup.status) }}
                    </span>
                    <div v-if="backup.error_msg" class="relative">
                      <p 
                        class="text-xs text-red-600 max-w-xs truncate cursor-help" 
                        :title="backup.error_msg"
                        @click="showFullError(backup.error_msg)"
                      >
                         {{ backup.error_msg }}
                      </p>
                    </div>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <div class="flex space-x-2">
                    <button 
                      v-if="backup.status === 'completed'"
                      @click="downloadBackup(backup)"
                      class="text-blue-600 hover:text-blue-800 transition-colors duration-200"
                    >
                      Télécharger
                    </button>
                    <button 
                      v-if="backup.status === 'completed'"
                      @click="restoreBackup(backup)"
                      class="text-green-600 hover:text-green-800 transition-colors duration-200"
                    >
                      Restaurer
                    </button>
                    <button 
                      @click="deleteBackup(backup.id)"
                      class="text-red-600 hover:text-red-800 transition-colors duration-200"
                    >
                      Supprimer
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useSafebaseStore } from '@/stores/safebase'
import { backupService } from '@/services/backup_service'
import type { Backup } from '@/types/backup'

const safebaseStore = useSafebaseStore()
const { backups, loading, error, completedBackups, pendingBackups, failedBackups, databases } = storeToRefs(safebaseStore)

const filterStatus = ref<string | null>(null)
const filterDatabaseId = ref<string>('')

const filteredBackups = computed(() => {
  let filtered = backups.value

  // Appliquer le filtre par statut
  if (filterStatus.value) {
    filtered = backupService.filterByStatus(filtered, filterStatus.value)
  }

  // Appliquer le filtre par base de données
  if (filterDatabaseId.value) {
    filtered = filtered.filter(backup => backup.database_id === parseInt(filterDatabaseId.value))
  }

  return backupService.sortByDate(filtered)
})

const filteredCompletedBackups = computed(() => {
  return filteredBackups.value.filter(backup => backup.status === 'completed')
})

const filteredPendingBackups = computed(() => {
  return filteredBackups.value.filter(backup => backup.status === 'pending')
})

const filteredFailedBackups = computed(() => {
  return filteredBackups.value.filter(backup => backup.status === 'failed')
})

const formatDate = (dateString: string): string => {
  return new Date(dateString).toLocaleDateString('fr-FR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatSize = (bytes: number): string => {
  return backupService.formatFileSize(bytes)
}

const formatTotalSize = (): string => {
  const total = backupService.getTotalSize(backups.value)
  return backupService.formatFileSize(total)
}

const getStatusLabel = (status: string): string => {
  return backupService.getStatusLabel(status)
}

const getStatusClass = (status: string): string => {
  const colors: Record<string, string> = {
    'pending': 'bg-orange-100 text-orange-800',
    'completed': 'bg-green-100 text-green-800',
    'failed': 'bg-red-100 text-red-800'
  }
  return colors[status] || 'bg-gray-100 text-gray-800'
}

const getBackupTypeLabel = (backup: Backup): string => {
  return backupService.getBackupTypeLabel(backup)
}

const getBackupTypeClass = (backup: Backup): string => {
  return backupService.getBackupTypeClass(backup)
}

const downloadBackup = async (backup: Backup) => {
  try {
    await backupService.downloadBackup(backup)
  } catch (err: any) {
    alert(err.message)
  }
}

const deleteBackup = async (id: number) => {
  if (!confirm('Êtes-vous sûr de vouloir supprimer cette sauvegarde ?')) return
  try {
    await backupService.deleteBackup(id)
    safebaseStore.removeBackup(id)
  } catch (err: any) {
    alert(err.message)
  }
}

const restoreBackup = async (backup: Backup) => {
  // Confirmation simple avant restauration
  if (!confirm(`Êtes-vous sûr de vouloir restaurer la sauvegarde "${backup.filename}" ?\n\n Attention : Cette opération va remplacer les données actuelles de la base de données d'origine.`)) {
    return
  }
  
  try {
    await backupService.restoreBackup(backup)
    alert('Restauration lancée avec succès ! Le processus peut prendre quelques minutes.')
  } catch (err: any) {
    alert(`Erreur lors de la restauration : ${err.message}`)
  }
}

const showFullError = (errorMsg: string) => {
  alert(`Erreur complète :\n\n${errorMsg}`)
}

onMounted(() => {
  safebaseStore.fetchBackups()
  safebaseStore.fetchDatabases()
})
</script>
