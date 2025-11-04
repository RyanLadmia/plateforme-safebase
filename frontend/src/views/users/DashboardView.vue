<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <header class="bg-white shadow">
      <div class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center">
          <h1 class="text-3xl font-bold text-gray-900">Tableau de bord</h1>
          <div class="flex items-center space-x-4">
            <span class="text-gray-600">{{ user?.email }}</span>
            <router-link 
              to="/user/profile" 
              class="text-blue-600 hover:text-blue-800"
            >
              Mon profil
            </router-link>
          </div>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <!-- Stats Cards -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <!-- Databases Count -->
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="flex-1">
              <p class="text-gray-500 text-sm">Bases de données</p>
              <p class="text-3xl font-bold text-gray-900">{{ databaseCount }}</p>
            </div>
            <div class="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
              <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4"></path>
              </svg>
            </div>
          </div>
        </div>

        <!-- Backups Count -->
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="flex-1">
              <p class="text-gray-500 text-sm">Sauvegardes</p>
              <p class="text-3xl font-bold text-gray-900">{{ backupCount }}</p>
            </div>
            <div class="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center">
              <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
              </svg>
            </div>
          </div>
        </div>

        <!-- Schedules Count -->
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="flex-1">
              <p class="text-gray-500 text-sm">Planifications</p>
              <p class="text-3xl font-bold text-gray-900">{{ scheduleCount }}</p>
            </div>
            <div class="w-12 h-12 bg-purple-100 rounded-full flex items-center justify-center">
              <svg class="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
            </div>
          </div>
        </div>

        <!-- Pending Backups -->
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="flex-1">
              <p class="text-gray-500 text-sm">En cours</p>
              <p class="text-3xl font-bold text-gray-900">{{ pendingBackups.length }}</p>
            </div>
            <div class="w-12 h-12 bg-orange-100 rounded-full flex items-center justify-center">
              <svg class="w-6 h-6 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <router-link 
          to="/user/databases" 
          class="bg-white rounded-lg shadow p-6 hover:shadow-lg transition"
        >
          <h3 class="text-xl font-semibold text-gray-900 mb-2">Gérer mes bases de données</h3>
          <p class="text-gray-600">Ajoutez, modifiez ou supprimez vos configurations de bases de données</p>
        </router-link>

        <router-link 
          to="/user/backups" 
          class="bg-white rounded-lg shadow p-6 hover:shadow-lg transition"
        >
          <h3 class="text-xl font-semibold text-gray-900 mb-2">Gérer mes sauvegardes</h3>
          <p class="text-gray-600">Créez et téléchargez vos sauvegardes</p>
        </router-link>

        <router-link 
          to="/user/schedules" 
          class="bg-white rounded-lg shadow p-6 hover:shadow-lg transition"
        >
          <h3 class="text-xl font-semibold text-gray-900 mb-2">Gérer mes planifications</h3>
          <p class="text-gray-600">Planifiez des sauvegardes automatiques</p>
        </router-link>
      </div>

      <!-- Recent Backups -->
      <div class="bg-white rounded-lg shadow">
        <div class="px-6 py-4 border-b border-gray-200">
          <h2 class="text-xl font-semibold text-gray-900">Dernières sauvegardes</h2>
        </div>
        <div class="p-6">
          <div v-if="loading" class="text-center py-8">
            <p class="text-gray-500">Chargement...</p>
          </div>
          <div v-else-if="error" class="text-center py-8">
            <p class="text-red-600">{{ error }}</p>
          </div>
          <div v-else-if="recentBackups.length === 0" class="text-center py-8">
            <p class="text-gray-500">Aucune sauvegarde pour le moment</p>
            <router-link 
              to="/user/databases" 
              class="mt-4 inline-block text-blue-600 hover:text-blue-800"
            >
              Commencer par ajouter une base de données
            </router-link>
          </div>
          <table v-else class="min-w-full divide-y divide-gray-200">
            <thead>
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Nom</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Date</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Taille</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Statut</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="backup in recentBackups" :key="backup.id">
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ backup.filename }}</td>
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
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useSafebaseStore } from '@/stores/safebase'
import { backupService } from '@/services/backup_service'

const router = useRouter()
const authStore = useAuthStore()
const safebaseStore = useSafebaseStore()

const { user } = storeToRefs(authStore)
const { databases, backups, loading, error, databaseCount, backupCount, scheduleCount, pendingBackups } = storeToRefs(safebaseStore)

const recentBackups = computed(() => {
  return backupService.sortByDate(backups.value).slice(0, 5)
})

const handleLogout = async () => {
  await authStore.logout()
  router.push('/login')
}

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

onMounted(async () => {
  // Charger les données
  try {
    await Promise.all([
      safebaseStore.fetchDatabases(),
      safebaseStore.fetchBackups()
    ])
  } catch (err) {
    console.error('Erreur lors du chargement des données:', err)
  }
})
</script>
