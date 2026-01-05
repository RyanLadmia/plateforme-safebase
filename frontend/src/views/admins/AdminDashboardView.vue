<template>
  <div class="min-h-screen bg-gray-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <div class="flex justify-between items-center mb-6">
        <h1 class="text-2xl sm:text-3xl font-bold text-gray-900">Administration - Tableau de bord</h1>
        <div class="flex items-center space-x-4">
          <span class="px-3 py-1 bg-red-100 text-red-800 rounded-full text-sm font-semibold">ADMIN</span>
          <span class="text-gray-600">{{ user?.email }}</span>
        </div>
      </div>

      <!-- Admin Navigation -->
      <div class="bg-white rounded-lg shadow mb-6 p-4">
        <nav class="flex space-x-4">
          <router-link 
            to="/admin/dashboard" 
            class="px-4 py-2 bg-blue-600 text-white rounded-lg"
          >
            Tableau de bord
          </router-link>
          <router-link 
            to="/admin/users" 
            class="px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300"
          >
            Gestion des utilisateurs
          </router-link>
        </nav>
      </div>

      <!-- Global Statistics -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="flex-1">
              <p class="text-gray-500 text-sm">Utilisateurs total</p>
              <p class="text-3xl font-bold text-gray-900">-</p>
            </div>
            <div class="w-12 h-12 bg-purple-100 rounded-full flex items-center justify-center">
              <svg class="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"></path>
              </svg>
            </div>
          </div>
        </div>

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

      <!-- System Information -->
      <div class="bg-white rounded-lg shadow mb-6">
        <div class="px-6 py-4 border-b border-gray-200">
          <h2 class="text-xl font-semibold text-gray-900">Informations système</h2>
        </div>
        <div class="p-6">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <h3 class="text-sm font-medium text-gray-500 mb-2">Backend</h3>
              <p class="text-gray-900">API REST Go - Port 8080</p>
              <p class="text-sm text-gray-500">Version: 1.0.0</p>
            </div>
            <div>
              <h3 class="text-sm font-medium text-gray-500 mb-2">Frontend</h3>
              <p class="text-gray-900">Vue 3 + TypeScript</p>
              <p class="text-sm text-gray-500">Version: 1.0.0</p>
            </div>
            <div>
              <h3 class="text-sm font-medium text-gray-500 mb-2">Base de données</h3>
              <p class="text-gray-900">PostgreSQL</p>
              <p class="text-sm text-gray-500">Gestion des utilisateurs et métadonnées</p>
            </div>
            <div>
              <h3 class="text-sm font-medium text-gray-500 mb-2">Sauvegardes</h3>
              <p class="text-gray-900">MySQL & PostgreSQL</p>
              <p class="text-sm text-gray-500">Format: ZIP</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Recent Activity -->
      <div class="bg-white rounded-lg shadow">
        <div class="px-6 py-4 border-b border-gray-200">
          <h2 class="text-xl font-semibold text-gray-900">Activité récente</h2>
        </div>
        <div class="p-6">
          <div v-if="loading" class="text-center py-8">Chargement...</div>
          <div v-else-if="recentBackups.length === 0" class="text-center py-8">
            <p class="text-gray-500">Aucune activité récente</p>
          </div>
          <div v-else class="space-y-4">
            <div 
              v-for="backup in recentBackups" 
              :key="backup.id"
              class="flex items-center justify-between p-4 bg-gray-50 rounded-lg"
            >
              <div class="flex-1">
                <p class="font-medium text-gray-900">{{ backup.filename }}</p>
                <p class="text-sm text-gray-500">{{ formatDate(backup.created_at) }}</p>
              </div>
              <span 
                :class="getStatusClass(backup.status)"
                class="px-3 py-1 text-xs font-semibold rounded-full"
              >
                {{ getStatusLabel(backup.status) }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useSafebaseStore } from '@/stores/safebase'
import { backupService } from '@/services/backup_service'

const router = useRouter()
const authStore = useAuthStore()
const safebaseStore = useSafebaseStore()

const { user } = storeToRefs(authStore)
const { databaseCount, backupCount, backups, pendingBackups, loading } = storeToRefs(safebaseStore)

const recentBackups = computed(() => {
  return backupService.sortByDate(backups.value).slice(0, 10)
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
  await Promise.all([
    safebaseStore.fetchDatabases(),
    safebaseStore.fetchBackups()
  ])
})
</script>
