<template>
  <div class="min-h-screen bg-gray-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <div class="flex justify-between items-center mb-6">
        <h1 class="text-2xl sm:text-3xl font-bold text-gray-900">Mes sauvegardes planifiées</h1>
        <button @click="showCreateModal = true" class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
          + Nouvelle planification
        </button>
      </div>

      <!-- Filters -->
      <div class="bg-white rounded-lg shadow p-4 mb-6">
        <div class="flex flex-wrap justify-between items-center gap-4">
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

      <div v-if="loading" class="text-center py-12">Chargement...</div>
      <div v-else-if="error" class="bg-red-100 text-red-700 p-4 rounded-lg">{{ error }}</div>
      <div v-else-if="schedules.length === 0" class="text-center py-12">
        <p class="text-gray-500 mb-4">Aucune sauvegarde planifiée</p>
        <button @click="showCreateModal = true" class="text-blue-600 hover:text-blue-800">
          Créer votre premi�re planification
        </button>
      </div>
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div v-for="schedule in filteredAndSortedSchedules" :key="schedule.id" class="bg-white rounded-lg shadow p-6">
          <div class="flex justify-between items-start mb-4">
            <div class="flex-1">
              <div class="flex items-center gap-2 mb-2">
                <h3 class="text-lg font-semibold text-gray-900">
                  {{ schedule.name }}
                </h3>
                <span
                  :class="schedule.active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
                  class="inline-block px-2 py-1 text-xs font-semibold rounded"
                >
                  {{ schedule.active ? 'Actif' : 'Inactif' }}
                </span>
              </div>
              <p class="text-sm text-gray-600 mb-2">
                Base de données : {{ schedule.database?.name || 'Base inconnue' }}
              </p>
              <span class="inline-block px-2 py-1 text-xs font-semibold rounded bg-blue-100 text-blue-800">
                {{ schedule.database?.type || 'N/A' }}
              </span>
            </div>
            <div class="flex gap-2">
              <button @click="editSchedule(schedule)" class="text-blue-600 hover:text-blue-800">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
                </svg>
              </button>
              <button @click="deleteSchedule(schedule.id)" class="text-red-600 hover:text-red-800">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                </svg>
              </button>
            </div>
          </div>

          <div class="text-sm text-gray-600 space-y-2">
            <div>
              <strong>Fréquence:</strong>
              <span class="text-gray-800">{{ CronUtils.getFrequencyDescription(schedule.cron_expression) }}</span>
            </div>
            <div>
              <strong>Prochaine exécution:</strong>
              <span class="text-gray-500">{{ CronUtils.getNextExecution(schedule.cron_expression) }}</span>
            </div>
            <div>
              <strong>Créé le:</strong>
              <span>{{ formatDate(schedule.created_at) }}</span>
            </div>
          </div>

          <div class="mt-4 flex gap-2">
            <button
              @click="toggleSchedule(schedule)"
              :class="schedule.active ? 'bg-red-600 hover:bg-red-700' : 'bg-green-600 hover:bg-green-700'"
              class="flex-1 px-3 py-2 text-white text-sm rounded-lg"
            >
              {{ schedule.active ? 'Désactiver' : 'Activer' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <div v-if="showCreateModal || showEditModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div class="bg-white rounded-lg max-w-2xl w-full max-h-[90vh] overflow-y-auto p-6">
        <h2 class="text-2xl font-bold mb-4">
          {{ showEditModal ? 'Modifier la planification' : 'Nouvelle planification' }}
        </h2>

        <form @submit.prevent="showEditModal ? updateSchedule() : createSchedule()" class="space-y-4">
          <div>
            <label class="block text-sm font-medium mb-2">Base de données</label>
            <select v-model="form.database_id" required class="w-full px-4 py-2 border rounded-lg">
              <option value="">Sélectionnez une base de données</option>
              <option v-for="db in databases" :key="db.id" :value="db.id">
                {{ db.name }} ({{ db.type }})
              </option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium mb-2">Nom de la planification</label>
            <input
              v-model="form.name"
              type="text"
              required
              placeholder="Ex: Sauvegarde quotidienne production"
              class="w-full px-4 py-2 border rounded-lg"
            />
            <p class="text-xs text-gray-500 mt-1">
              Donnez un nom descriptif � votre planification
            </p>
          </div>

          <div>
            <label class="block text-sm font-medium mb-2">Fréquence de sauvegarde</label>
            <div class="space-y-2">
              <select v-model="form.cron_expression" class="w-full px-4 py-2 border rounded-lg">
                <option value="">Personnalisé</option>
                <option v-for="preset in cronPresets" :key="preset.expression" :value="preset.expression">
                  {{ preset.label }} - {{ preset.description }}
                </option>
              </select>

              <div v-if="!isPresetSelected" class="mt-2">
                <label class="block text-xs text-gray-600 mb-1">Expression Cron personnalisée</label>
                <input
                  v-model="form.cron_expression"
                  placeholder="0 0 * * * (tous les jours � minuit)"
                  required
                  class="w-full px-4 py-2 border rounded-lg text-sm"
                />
                <p class="text-xs text-gray-500 mt-1">
                  Format: minute heure jour mois jour-semaine
                </p>
              </div>
            </div>
          </div>

          <div v-if="showEditModal" class="flex items-center">
            <input
              id="active"
              v-model="form.active"
              type="checkbox"
              class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
            />
            <label for="active" class="ml-2 block text-sm text-gray-900">
              Planification active
            </label>
          </div>

          <div class="flex gap-4 pt-4">
            <button
              type="button"
              @click="closeModal"
              class="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50"
            >
              Annuler
            </button>
            <button
              type="submit"
              :disabled="loading"
              class="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
            >
              {{ loading ? 'Enregistrement...' : (showEditModal ? 'Modifier' : 'Créer') }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useSafebaseStore } from '@/stores/safebase'
import { scheduleService } from '@/services/schedule_service'
import type { Schedule, ScheduleCreateRequest, ScheduleUpdateRequest, CronPreset } from '@/types/schedule'
import { CRON_PRESETS } from '@/types/schedule'
import { CronUtils } from '@/utils/cron-utils'

// Composables
const safebaseStore = useSafebaseStore()
const {
  databases,
  schedules,
  loading,
  error
} = storeToRefs(safebaseStore)

// �tat local
const showCreateModal = ref(false)
const showEditModal = ref(false)
const editingSchedule = ref<Schedule | null>(null)
const cronPresets = ref<CronPreset[]>(CRON_PRESETS)
const filterDatabaseId = ref<string>('')

const form = ref<ScheduleCreateRequest & { active?: boolean }>({
  database_id: 0,
  name: '',
  cron_expression: '',
  active: true
})

// Computed
const isPresetSelected = computed(() => {
  return cronPresets.value.some(preset => preset.expression === form.value.cron_expression)
})

// Fonction pour calculer la fréquence d'une expression cron
const getFrequencyScore = (cronExpression: string): number => {
  const parts = cronExpression.split(' ')
  if (parts.length !== 5) return 0

  const [minute, hour, day, month, dayOfWeek] = parts

  // Score basé sur la fréquence (plus le score est élevé, plus la fréquence est élevée)
  let score = 0

  // Minutes spécifiques (ex: "0" = tr�s fréquent)
  if (minute !== '*') score += 10

  // Heures spécifiques (ex: "0" = fréquent)
  if (hour !== '*') score += 5

  // Jours spécifiques du mois (moins fréquent)
  if (day !== '*') score += 2

  // Mois spécifiques (encore moins fréquent)
  if (month !== '*') score += 1

  // Jours de la semaine spécifiques (fréquent pour les sauvegardes hebdomadaires)
  if (dayOfWeek !== '*') score += 3

  // Pénalités pour les expressions tr�s spécifiques
  if (minute.includes(',') || hour.includes(',') || day.includes(',') || month.includes(',') || dayOfWeek.includes(',')) {
    score -= 2 // Moins fréquent si plusieurs valeurs
  }

  // Bonus pour les expressions quotidiennes
  if (minute !== '*' && hour !== '*' && day === '*' && month === '*' && dayOfWeek === '*') {
    score += 15 // Tous les jours � heure fixe = tr�s fréquent
  }

  // Bonus pour les expressions horaires
  if (minute !== '*' && hour === '*' && day === '*' && month === '*' && dayOfWeek === '*') {
    score += 20 // Toutes les heures = tr�s tr�s fréquent
  }

  return score
}

const filteredAndSortedSchedules = computed(() => {
  let filtered = schedules.value

  // Appliquer le filtre par base de données
  if (filterDatabaseId.value) {
    filtered = filtered.filter(schedule => schedule.database_id === parseInt(filterDatabaseId.value))
  }

  // Trier par fréquence (plus fréquentes en premier)
  return filtered.sort((a, b) => {
    const scoreA = getFrequencyScore(a.cron_expression)
    const scoreB = getFrequencyScore(b.cron_expression)
    return scoreB - scoreA // Score décroissant = fréquence décroissante
  })
})

// Méthodes
const closeModal = () => {
  showCreateModal.value = false
  showEditModal.value = false
  editingSchedule.value = null
  form.value = { database_id: 0, name: '', cron_expression: '', active: true }
}

const createSchedule = async () => {
  try {
    const newSchedule = await scheduleService.createSchedule({
      database_id: form.value.database_id,
      name: form.value.name,
      cron_expression: form.value.cron_expression
    })
    safebaseStore.addSchedule(newSchedule)
    closeModal()
  } catch (err: any) {
    error.value = err.message
  }
}

const editSchedule = (schedule: Schedule) => {
  editingSchedule.value = schedule
  form.value = {
    database_id: schedule.database_id,
    name: schedule.name,
    cron_expression: schedule.cron_expression,
    active: schedule.active
  }
  showEditModal.value = true
}

const updateSchedule = async () => {
  if (!editingSchedule.value) return

  try {
    const updateData: ScheduleUpdateRequest = {}
    if (form.value.name !== editingSchedule.value.name) {
      updateData.name = form.value.name
    }
    if (form.value.cron_expression !== editingSchedule.value.cron_expression) {
      updateData.cron_expression = form.value.cron_expression
    }
    if (form.value.active !== editingSchedule.value.active) {
      updateData.active = form.value.active
    }

    const updatedSchedule = await scheduleService.updateSchedule(editingSchedule.value.id, updateData)
    safebaseStore.updateSchedule(updatedSchedule)
    closeModal()
  } catch (err: any) {
    error.value = err.message
  }
}

const toggleSchedule = async (schedule: Schedule) => {
  try {
    const updatedSchedule = await scheduleService.toggleSchedule(schedule.id, !schedule.active)
    safebaseStore.updateSchedule(updatedSchedule)
  } catch (err: any) {
    error.value = err.message
  }
}

const deleteSchedule = async (id: number) => {
  if (!confirm('�tes-vous sûr de vouloir supprimer cette planification ?')) return

  try {
    await scheduleService.deleteSchedule(id)
    safebaseStore.removeSchedule(id)
  } catch (err: any) {
    error.value = err.message
  }
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('fr-FR', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// Lifecycle
onMounted(async () => {
  await safebaseStore.fetchDatabases()
  await safebaseStore.fetchSchedules()
})
onMounted(async () => {
  await safebaseStore.fetchDatabases()
  await safebaseStore.fetchSchedules()
})
</script>