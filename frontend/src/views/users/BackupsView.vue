<template>
  <div class="min-h-screen bg-gray-50">
    <header class="bg-white shadow">
      <div class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        <h1 class="text-3xl font-bold text-gray-900">Mes sauvegardes</h1>
      </div>
    </header>

    <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <!-- Filters -->
      <div class="bg-white rounded-lg shadow p-4 mb-6 flex space-x-4">
        <button 
          @click="filterStatus = null" 
          :class="filterStatus === null ? 'bg-blue-600 text-white' : 'bg-gray-200'"
          class="px-4 py-2 rounded-lg"
        >
          Toutes ({{ backups.length }})
        </button>
        <button 
          @click="filterStatus = 'completed'" 
          :class="filterStatus === 'completed' ? 'bg-green-600 text-white' : 'bg-gray-200'"
          class="px-4 py-2 rounded-lg"
        >
          Terminées ({{ completedBackups.length }})
        </button>
        <button 
          @click="filterStatus = 'pending'" 
          :class="filterStatus === 'pending' ? 'bg-orange-600 text-white' : 'bg-gray-200'"
          class="px-4 py-2 rounded-lg"
        >
          En cours ({{ pendingBackups.length }})
        </button>
        <button 
          @click="filterStatus = 'failed'" 
          :class="filterStatus === 'failed' ? 'bg-red-600 text-white' : 'bg-gray-200'"
          class="px-4 py-2 rounded-lg"
        >
          Échouées ({{ failedBackups.length }})
        </button>
      </div>

      <!-- Backups List -->
      <div v-if="loading" class="text-center py-12">Chargement...</div>
      <div v-else-if="error" class="bg-red-100 text-red-700 p-4 rounded-lg">{{ error }}</div>
      <div v-else-if="filteredBackups.length === 0" class="text-center py-12">
        <p class="text-gray-500">Aucune sauvegarde trouvée</p>
      </div>
      <div v-else class="bg-white rounded-lg shadow overflow-hidden">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Fichier</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Date</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Taille</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Statut</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="backup in filteredBackups" :key="backup.id">
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {{ backup.filename }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ formatDate(backup.created_at) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ formatSize(backup.size) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span 
                  :class="getStatusClass(backup.status)"
                  class="px-2 py-1 text-xs font-semibold rounded-full"
                >
                  {{ getStatusLabel(backup.status) }}
                </span>
                <p v-if="backup.error_msg" class="text-xs text-red-600 mt-1">{{ backup.error_msg }}</p>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium space-x-2">
                <button 
                  v-if="backup.status === 'completed'"
                  @click="downloadBackup(backup)"
                  class="text-blue-600 hover:text-blue-800"
                >
                  Télécharger
                </button>
                <button 
                  @click="deleteBackup(backup.id)"
                  class="text-red-600 hover:text-red-800"
                >
                  Supprimer
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Statistics -->
      <div class="mt-6 grid grid-cols-1 md:grid-cols-4 gap-4">
        <div class="bg-white rounded-lg shadow p-4">
          <p class="text-gray-500 text-sm">Taille totale</p>
          <p class="text-2xl font-bold text-gray-900">{{ formatTotalSize() }}</p>
        </div>
        <div class="bg-white rounded-lg shadow p-4">
          <p class="text-gray-500 text-sm">Sauvegardes réussies</p>
          <p class="text-2xl font-bold text-green-600">{{ completedBackups.length }}</p>
        </div>
        <div class="bg-white rounded-lg shadow p-4">
          <p class="text-gray-500 text-sm">En cours</p>
          <p class="text-2xl font-bold text-orange-600">{{ pendingBackups.length }}</p>
        </div>
        <div class="bg-white rounded-lg shadow p-4">
          <p class="text-gray-500 text-sm">Échouées</p>
          <p class="text-2xl font-bold text-red-600">{{ failedBackups.length }}</p>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useSafebaseStore } from '@/stores/safebase'
import { backupService } from '@/services/backup_service'
import type { Backup } from '@/types/backup'

const safebaseStore = useSafebaseStore()
const { backups, loading, error, completedBackups, pendingBackups, failedBackups } = safebaseStore

const filterStatus = ref<string | null>(null)

const filteredBackups = computed(() => {
  if (!filterStatus.value) return backupService.sortByDate(backups)
  return backupService.sortByDate(backupService.filterByStatus(backups, filterStatus.value))
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
  const total = backupService.getTotalSize(backups)
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

onMounted(() => {
  safebaseStore.fetchBackups()
})
</script>
