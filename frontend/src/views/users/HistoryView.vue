<template>
  <div class="min-h-screen bg-gray-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <div class="mb-6">
        <h1 class="text-2xl sm:text-3xl font-bold text-gray-900">Historique des activités</h1>
        <p class="text-gray-600 mt-2">Suivez toutes vos actions sur la plateforme</p>
      </div>

      <!-- Filters -->
      <div class="bg-white rounded-lg shadow p-4 mb-6">
        <div class="flex flex-wrap justify-between items-center gap-4">
          <!-- Activity type filters -->
          <div class="flex flex-wrap gap-2 sm:gap-4">
            <button
              @click="activeFilter = 'all'"
              :class="activeFilter === 'all' ? 'bg-blue-600 text-white' : 'bg-gray-200 text-gray-700'"
              class="px-3 py-2 sm:px-4 rounded-lg text-sm font-medium transition-colors duration-200"
            >
              Toutes les activités ({{ totalActivities }})
            </button>
            <button
              @click="activeFilter = 'database'"
              :class="activeFilter === 'database' ? 'bg-purple-600 text-white' : 'bg-gray-200 text-gray-700'"
              class="px-3 py-2 sm:px-4 rounded-lg text-sm font-medium transition-colors duration-200"
            >
              Bases de données ({{ databaseActivities }})
            </button>
            <button
              @click="activeFilter = 'backup'"
              :class="activeFilter === 'backup' ? 'bg-green-600 text-white' : 'bg-gray-200 text-gray-700'"
              class="px-3 py-2 sm:px-4 rounded-lg text-sm font-medium transition-colors duration-200"
            >
              Sauvegardes ({{ backupActivities }})
            </button>
            <button
              @click="activeFilter = 'schedule'"
              :class="activeFilter === 'schedule' ? 'bg-orange-600 text-white' : 'bg-gray-200 text-gray-700'"
              class="px-3 py-2 sm:px-4 rounded-lg text-sm font-medium transition-colors duration-200"
            >
              Planifications ({{ scheduleActivities }})
            </button>
            <button
              @click="activeFilter = 'restore'"
              :class="activeFilter === 'restore' ? 'bg-indigo-600 text-white' : 'bg-gray-200 text-gray-700'"
              class="px-3 py-2 sm:px-4 rounded-lg text-sm font-medium transition-colors duration-200"
            >
              Restaurations ({{ restoreActivities }})
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
      <div class="mt-6 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-4 mb-6">
        <div class="bg-white rounded-lg shadow p-4">
          <p class="text-gray-500 text-sm">Total activités</p>
          <p class="text-2xl font-bold text-gray-900">{{ totalActivities }}</p>
        </div>
        <div class="bg-white rounded-lg shadow p-4">
          <p class="text-gray-500 text-sm">Bases de données</p>
          <p class="text-2xl font-bold text-purple-600">{{ databaseActivities }}</p>
        </div>
        <div class="bg-white rounded-lg shadow p-4">
          <p class="text-gray-500 text-sm">Sauvegardes</p>
          <p class="text-2xl font-bold text-green-600">{{ backupActivities }}</p>
        </div>
        <div class="bg-white rounded-lg shadow p-4">
          <p class="text-gray-500 text-sm">Planifications</p>
          <p class="text-2xl font-bold text-orange-600">{{ scheduleActivities }}</p>
        </div>
        <div class="bg-white rounded-lg shadow p-4">
          <p class="text-gray-500 text-sm">Restaurations</p>
          <p class="text-2xl font-bold text-indigo-600">{{ restoreActivities }}</p>
        </div>
      </div>

      <!-- History List -->
      <div v-if="loading" class="text-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
        <p class="text-gray-500 mt-4">Chargement de l'historique...</p>
      </div>
      <div v-else-if="error" class="bg-red-100 text-red-700 p-4 rounded-lg">{{ error }}</div>
      <div v-else-if="filteredHistory.length === 0" class="text-center py-12">
        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v3m0 0v3m0-3h3m-3 0H9"></path>
        </svg>
        <p class="text-gray-500 mt-4">Aucune activité trouvée pour ce filtre</p>
      </div>
      <div v-else class="bg-white rounded-lg shadow overflow-hidden">
        <div class="divide-y divide-gray-200">
          <div
            v-for="item in filteredHistory"
            :key="item.id"
            class="p-6 hover:bg-gray-50 transition-colors duration-200"
          >
            <div class="flex items-start space-x-4">
              <!-- Icon -->
              <div class="flex-shrink-0">
                <div
                  :class="getActionIconClass(item)"
                  class="w-10 h-10 rounded-full flex items-center justify-center"
                >
                  <svg v-if="item.action === 'created' || item.action === 'create' || item.action === 'restored'" class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v3m0 0v3m0-3h3m-3 0H9m12 0a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                  </svg>
                  <svg v-else-if="item.action === 'updated' || item.action === 'update'" class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
                  </svg>
                  <svg v-else-if="item.action === 'deleted' || item.action === 'delete'" class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                  </svg>
                  <svg v-else-if="item.action === 'completed'" class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                  </svg>
                  <svg v-else-if="item.action === 'failed'" class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                  </svg>
                  <svg v-else-if="item.action === 'executed'" class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h1.586a1 1 0 01.707.293l.707.707A1 1 0 0012.414 11H15m-3-3h3a1 1 0 011 1v3a1 1 0 01-1 1h-3m-3-3h-3a1 1 0 00-1 1v3a1 1 0 001 1h3m-3-3v-3a1 1 0 011-1h3z"></path>
                  </svg>
                  <svg v-else-if="item.action === 'download'" class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                  </svg>
                  <svg v-else class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                  </svg>
                </div>
              </div>

              <!-- Content -->
              <div class="flex-1 min-w-0">
                <!-- Affichage spécial pour les bases de données -->
                <div v-if="isDatabaseItem(item)" class="space-y-1">
                  <div class="flex items-center justify-between">
                    <p class="text-sm font-medium text-gray-900">
                      {{ getDatabaseContent(item).title }}
                    </p>
                    <p class="text-sm text-gray-500">
                      {{ formatDate(item.created_at) }}
                    </p>
                  </div>
                  <p class="text-sm text-gray-600">
                    {{ getDatabaseContent(item).subtitle }}
                  </p>
                  <p v-for="detail in getDatabaseContent(item).details" :key="detail" class="text-sm text-gray-500">
                    {{ detail }}
                  </p>
                  <div class="flex items-center space-x-4 mt-2">
                    <span
                      :class="getResourceTypeClass(item.resource_type)"
                      class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                    >
                      {{ getResourceTypeLabel(item.resource_type) }}
                    </span>
                  </div>
                </div>
                <!-- Affichage spécial pour les sauvegardes -->
                <div v-else-if="isBackupItem(item)" class="space-y-1">
                  <div class="flex items-center justify-between">
                    <p class="text-sm font-medium text-gray-900">
                      {{ getBackupContent(item).title }}
                    </p>
                    <p class="text-sm text-gray-500">
                      {{ formatDate(item.created_at) }}
                    </p>
                  </div>
                  <p class="text-sm text-gray-600" v-html="getBackupContent(item).subtitle">
                  </p>
                  <p v-for="detail in getBackupContent(item).details" :key="detail" class="text-sm text-gray-500">
                    {{ detail }}
                  </p>
                  <div class="flex items-center space-x-4 mt-2">
                    <span
                      :class="getResourceTypeClass(item.resource_type)"
                      class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                    >
                      {{ getResourceTypeLabel(item.resource_type) }}
                    </span>
                  </div>
                </div>
                <!-- Affichage spécial pour les restaurations -->
                <div v-else-if="isRestoreItem(item)" class="space-y-1">
                  <div class="flex items-center justify-between">
                    <p class="text-sm font-medium text-gray-900">
                      {{ getRestoreContent(item).title }}
                    </p>
                    <p class="text-sm text-gray-500">
                      {{ formatDate(item.created_at) }}
                    </p>
                  </div>
                  <p class="text-sm text-gray-600">
                    {{ getRestoreContent(item).subtitle }}
                  </p>
                  <p v-for="detail in getRestoreContent(item).details" :key="detail" class="text-sm text-gray-500">
                    {{ detail }}
                  </p>
                  <div class="flex items-center space-x-4 mt-2">
                    <span
                      :class="getResourceTypeClass(item.resource_type)"
                      class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                    >
                      {{ getResourceTypeLabel(item.resource_type) }}
                    </span>
                  </div>
                </div>
                <!-- Affichage spécial pour les planifications -->
                <div v-else-if="isScheduleItem(item)" class="space-y-1">
                  <div class="flex items-center justify-between">
                    <p class="text-sm font-medium text-gray-900">
                      {{ getScheduleContent(item).title }}
                    </p>
                    <p class="text-sm text-gray-500">
                      {{ formatDate(item.created_at) }}
                    </p>
                  </div>
                  <p class="text-sm text-gray-600">
                    {{ getScheduleContent(item).subtitle }}
                  </p>
                  <p v-for="detail in getScheduleContent(item).details" :key="detail" class="text-sm text-gray-500">
                    {{ detail }}
                  </p>
                  <div class="flex items-center space-x-4 mt-2">
                    <span
                      :class="getResourceTypeClass(item.resource_type)"
                      class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                    >
                      {{ getResourceTypeLabel(item.resource_type) }}
                    </span>
                  </div>
                </div>
                <!-- Affichage normal pour les autres types -->
                <div v-else>
                  <div class="flex items-center justify-between">
                    <p class="text-sm font-medium text-gray-900">
                      {{ getActionText(item) }}
                    </p>
                    <p class="text-sm text-gray-500">
                      {{ formatDate(item.created_at) }}
                    </p>
                  </div>
                  <p class="text-sm text-gray-600 mt-1">
                    {{ item.description }}
                  </p>
                  <div class="flex items-center space-x-4 mt-2">
                    <span
                      :class="getResourceTypeClass(item.resource_type)"
                      class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                    >
                      {{ getResourceTypeLabel(item.resource_type) }}
                    </span>
                    <span class="text-xs text-gray-500">
                      ID: {{ item.resource_id }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="mt-6 flex justify-center">
        <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
          <button
            @click="currentPage = Math.max(1, currentPage - 1)"
            :disabled="currentPage === 1"
            class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
            </svg>
          </button>

          <button
            v-for="page in visiblePages"
            :key="page"
            @click="currentPage = typeof page === 'number' ? page : currentPage"
            :class="page === currentPage ? 'bg-blue-50 border-blue-500 text-blue-600' : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50'"
            class="relative inline-flex items-center px-4 py-2 border text-sm font-medium"
            :disabled="typeof page !== 'number'"
          >
            {{ page }}
          </button>

          <button
            @click="currentPage = Math.min(totalPages, currentPage + 1)"
            :disabled="currentPage === totalPages"
            class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
            </svg>
          </button>
        </nav>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useSafebaseStore } from '@/stores/safebase'
import type { ActivityType, HistoryItem, HistoryResponse } from '@/types/history'
import { historyService } from '@/services/history_service'

// Composables
const safebaseStore = useSafebaseStore()
const { databases } = storeToRefs(safebaseStore)

// State
const loading = ref(false)
const error = ref('')
const activeFilter = ref<ActivityType>('all')
const currentPage = ref(1)
const itemsPerPage = 10
const history = ref<HistoryItem[]>([])
const totalItems = ref(0)
const totalByType = ref<Record<ActivityType, number>>({
  all: 0,
  database: 0,
  backup: 0,
  schedule: 0,
  restore: 0
})
const filterDatabaseId = ref<string>('')

// Computed
const filteredHistory = computed(() => {
  let filtered = history.value

  // Appliquer le filtre par base de données
  if (filterDatabaseId.value) {
    const selectedDbId = parseInt(filterDatabaseId.value)
    filtered = filtered.filter(item => {
      // Pour les éléments de base de données, le resource_id correspond à l'ID de la base
      if (item.resource_type === 'database') {
        return item.resource_id === selectedDbId
      }

      // Pour les autres types, vérifier dans les métadonnées
      if (item.metadata?.database_id) {
        return item.metadata.database_id === selectedDbId
      }

      // Essayer de trouver l'ID de base de données dans d'autres champs de métadonnées
      if (item.metadata && item.metadata.database_name) {
        // Chercher la base de données correspondante par nom
        const db = databases.value.find(d => d.name === item.metadata!.database_name)
        return db?.id === selectedDbId
      }

      return false
    })
  }

  return filtered
})

const totalActivities = computed(() => totalByType.value.all)

const databaseActivities = computed(() => totalByType.value.database)

const backupActivities = computed(() => totalByType.value.backup)

const scheduleActivities = computed(() => totalByType.value.schedule)

const restoreActivities = computed(() => totalByType.value.restore)

const totalPages = computed(() => {
  const total = activeFilter.value === 'all' ? totalByType.value.all : totalByType.value[activeFilter.value]
  return Math.ceil(total / itemsPerPage)
})

const visiblePages = computed(() => {
  const pages = []
  const total = totalPages.value
  const current = currentPage.value

  if (total <= 7) {
    for (let i = 1; i <= total; i++) {
      pages.push(i)
    }
  } else {
    if (current <= 4) {
      pages.push(1, 2, 3, 4, 5, '...', total)
    } else if (current >= total - 3) {
      pages.push(1, '...', total - 4, total - 3, total - 2, total - 1, total)
    } else {
      pages.push(1, '...', current - 1, current, current + 1, '...', total)
    }
  }

  return pages.filter(p => p !== '...' || typeof p === 'number')
})

// Methods
const formatDate = (dateString: string): string => {
  return historyService.formatDate(dateString)
}

const getActionIconClass = (item: HistoryItem): string => {
  return historyService.getActionIconClass(item.action)
}

const getActionText = (item: HistoryItem): string => {
  return historyService.getActionText(item.action)
}

const getResourceTypeLabel = (type: string): string => {
  return historyService.getResourceTypeLabel(type)
}

const getResourceTypeClass = (type: string): string => {
  return historyService.getResourceTypeClass(type)
}

const getDatabaseContent = (item: HistoryItem) => {
  return historyService.getDatabaseContent(item)
}

const getBackupContent = (item: HistoryItem) => {
  return historyService.getBackupContent(item)
}

const getRestoreContent = (item: HistoryItem) => {
  return historyService.getRestoreContent(item)
}

const getScheduleContent = (item: HistoryItem) => {
  return historyService.getScheduleContent(item)
}

const isDatabaseItem = (item: HistoryItem): boolean => {
  return historyService.isDatabaseItem(item)
}

const isBackupItem = (item: HistoryItem): boolean => {
  return historyService.isBackupItem(item)
}

const isRestoreItem = (item: HistoryItem): boolean => {
  return historyService.isRestoreItem(item)
}

const isScheduleItem = (item: HistoryItem): boolean => {
  return historyService.isScheduleItem(item)
}

const loadHistory = async () => {
  loading.value = true
  error.value = ''

  try {
    const response: HistoryResponse = await historyService.fetchHistoryByType(activeFilter.value, currentPage.value, itemsPerPage)
    history.value = response.history
    totalItems.value = response.total
  } catch (err: any) {
    error.value = err.message || 'Erreur lors du chargement de l\'historique'
    console.error('Erreur chargement historique:', err)
  } finally {
    loading.value = false
  }
}

const loadTotals = async () => {
  try {
    // Charger les totaux pour chaque type
    const types: ActivityType[] = ['all', 'database', 'backup', 'schedule', 'restore']
    const promises = types.map(async (type) => {
      const response = await historyService.fetchHistoryByType(type, 1, 1)
      return { type, total: response.total }
    })

    const results = await Promise.all(promises)
    results.forEach(({ type, total }) => {
      totalByType.value[type] = total
    })
  } catch (err: any) {
    console.error('Erreur chargement totaux:', err)
  }
}

// Watch for filter changes
watch(activeFilter, () => {
  currentPage.value = 1
  loadHistory()
})

watch(currentPage, () => {
  loadHistory()
})

watch(filterDatabaseId, () => {
  currentPage.value = 1
  loadHistory()
})

// Lifecycle
onMounted(async () => {
  await safebaseStore.fetchDatabases()
  await loadTotals()
  await loadHistory()
})
</script>