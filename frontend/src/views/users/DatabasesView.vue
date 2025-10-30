<template>
  <div class="min-h-screen bg-gray-50">
    <header class="bg-white shadow">
      <div class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8 flex justify-between items-center">
        <h1 class="text-3xl font-bold text-gray-900">Mes bases de données</h1>
        <button @click="showCreateModal = true" class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
          + Nouvelle base de données
        </button>
      </div>
    </header>

    <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <div v-if="loading" class="text-center py-12">Chargement...</div>
      <div v-else-if="error" class="bg-red-100 text-red-700 p-4 rounded-lg">{{ error }}</div>
      <div v-else-if="databases.length === 0" class="text-center py-12">
        <p class="text-gray-500 mb-4">Aucune base de données configurée</p>
        <button @click="showCreateModal = true" class="text-blue-600 hover:text-blue-800">
          Ajouter votre première base de données
        </button>
      </div>
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div v-for="db in databases" :key="db.id" class="bg-white rounded-lg shadow p-6">
          <div class="flex justify-between items-start mb-4">
            <div>
              <h3 class="text-lg font-semibold text-gray-900">{{ db.name }}</h3>
              <span class="inline-block px-2 py-1 text-xs font-semibold rounded bg-blue-100 text-blue-800 mt-2">
                {{ db.type }}
              </span>
            </div>
            <button @click="deleteDatabase(db.id)" class="text-red-600 hover:text-red-800">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
              </svg>
            </button>
          </div>
          <div class="text-sm text-gray-600 space-y-1">
            <p><strong>Hôte:</strong> {{ db.host }}:{{ db.port }}</p>
            <p><strong>Base:</strong> {{ db.db_name }}</p>
            <p><strong>Utilisateur:</strong> {{ db.username }}</p>
          </div>
          <button 
            @click="createBackupForDb(db.id)" 
            class="mt-4 w-full px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700"
          >
            Créer une sauvegarde
          </button>
        </div>
      </div>
    </main>

    <!-- Create Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div class="bg-white rounded-lg max-w-2xl w-full max-h-[90vh] overflow-y-auto p-6">
        <h2 class="text-2xl font-bold mb-4">Nouvelle base de données</h2>
        <form @submit.prevent="createDatabase" class="space-y-4">
          <div>
            <label class="block text-sm font-medium mb-2">Nom</label>
            <input v-model="form.name" required class="w-full px-4 py-2 border rounded-lg" />
          </div>
          <div>
            <label class="block text-sm font-medium mb-2">Type</label>
            <select v-model="form.type" class="w-full px-4 py-2 border rounded-lg">
              <option value="postgresql">PostgreSQL</option>
              <option value="mysql">MySQL</option>
            </select>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium mb-2">Hôte</label>
              <input v-model="form.host" required class="w-full px-4 py-2 border rounded-lg" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-2">Port</label>
              <input v-model="form.port" required class="w-full px-4 py-2 border rounded-lg" />
            </div>
          </div>
          <div>
            <label class="block text-sm font-medium mb-2">Nom de la base</label>
            <input v-model="form.db_name" required class="w-full px-4 py-2 border rounded-lg" />
          </div>
          <div>
            <label class="block text-sm font-medium mb-2">Utilisateur</label>
            <input v-model="form.username" required class="w-full px-4 py-2 border rounded-lg" />
          </div>
          <div>
            <label class="block text-sm font-medium mb-2">Mot de passe</label>
            <div class="relative">
              <input 
                v-model="form.password" 
                :type="showPassword ? 'text' : 'password'" 
                required 
                class="w-full px-4 py-2 pr-12 border rounded-lg" 
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none"
                :aria-label="showPassword ? 'Masquer le mot de passe' : 'Afficher le mot de passe'"
              >
                <svg v-if="!showPassword" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                </svg>
                <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"></path>
                </svg>
              </button>
            </div>
          </div>
          <div class="flex justify-end space-x-4">
            <button type="button" @click="closeModal" class="px-6 py-2 bg-gray-300 rounded-lg hover:bg-gray-400">
              Annuler
            </button>
            <button type="submit" :disabled="formLoading" class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
              {{ formLoading ? 'Création...' : 'Créer' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useSafebaseStore } from '@/stores/safebase'
import { databaseService } from '@/services/database_service'
import { backupService } from '@/services/backup_service'
import type { DatabaseCreateRequest } from '@/types/database'

const safebaseStore = useSafebaseStore()
const { databases, loading, error } = storeToRefs(safebaseStore)

const showCreateModal = ref(false)
const formLoading = ref(false)
const showPassword = ref(false)

const form = reactive<DatabaseCreateRequest>({
  name: '',
  type: 'postgresql',
  host: 'localhost',
  port: '5432',
  username: '',
  password: '',
  db_name: ''
})

const createDatabase = async () => {
  formLoading.value = true
  try {
    const newDb = await databaseService.createDatabase(form)
    safebaseStore.addDatabase(newDb)
    closeModal()
  } catch (err: any) {
    alert(err.message)
  } finally {
    formLoading.value = false
  }
}

const deleteDatabase = async (id: number) => {
  if (!confirm('Êtes-vous sûr de vouloir supprimer cette base de données ?')) return
  try {
    await databaseService.deleteDatabase(id)
    safebaseStore.removeDatabase(id)
  } catch (err: any) {
    alert(err.message)
  }
}

const createBackupForDb = async (dbId: number) => {
  try {
    const backup = await backupService.createBackup(dbId)
    safebaseStore.addBackup(backup)
    alert('Sauvegarde lancée avec succès !')
  } catch (err: any) {
    alert(err.message)
  }
}

const closeModal = () => {
  showCreateModal.value = false
  Object.assign(form, {
    name: '',
    type: 'postgresql',
    host: 'localhost',
    port: '5432',
    username: '',
    password: '',
    db_name: ''
  })
}

onMounted(() => {
  safebaseStore.fetchDatabases()
})
</script>
