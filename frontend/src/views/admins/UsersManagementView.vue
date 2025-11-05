<template>
  <div class="min-h-screen bg-gray-50">
    <header class="bg-white shadow">
      <div class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center">
          <h1 class="text-3xl font-bold text-gray-900">Gestion des utilisateurs</h1>
          <router-link
            to="/admin/dashboard"
            class="text-blue-600 hover:text-blue-800"
          >
            ← Retour au tableau de bord
          </router-link>
        </div>
      </div>
    </header>

    <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <!-- Statistiques -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-6">
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <svg class="h-8 w-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"></path>
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-500">Total utilisateurs</p>
              <p class="text-2xl font-semibold text-gray-900">{{ userCount }}</p>
            </div>
          </div>
        </div>

        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <svg class="h-8 w-8 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-500">Utilisateurs actifs</p>
              <p class="text-2xl font-semibold text-gray-900">{{ activeUsersCount }}</p>
            </div>
          </div>
        </div>

        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <svg class="h-8 w-8 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-500">Utilisateurs inactifs</p>
              <p class="text-2xl font-semibold text-gray-900">{{ inactiveUsersCount }}</p>
            </div>
          </div>
        </div>

        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <svg class="h-8 w-8 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"></path>
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-500">Administrateurs</p>
              <p class="text-2xl font-semibold text-gray-900">{{ adminUsersCount }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Filtres et actions -->
      <div class="bg-white rounded-lg shadow mb-6">
        <div class="px-6 py-4 border-b border-gray-200">
          <div class="flex justify-between items-center">
            <h2 class="text-xl font-semibold text-gray-900">Liste des utilisateurs</h2>
            <div class="flex space-x-4">
              <button
                @click="showActiveOnly = !showActiveOnly"
                :class="showActiveOnly ? 'bg-blue-600 text-white' : 'bg-gray-200 text-gray-700'"
                class="px-4 py-2 rounded-lg text-sm font-medium"
              >
                {{ showActiveOnly ? 'Tous les utilisateurs' : 'Utilisateurs actifs uniquement' }}
              </button>
              <button
                @click="fetchUsers"
                class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 text-sm font-medium"
              >
                Actualiser
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Message d'erreur -->
      <div v-if="error" class="bg-red-100 text-red-700 p-4 rounded-lg mb-6">
        {{ error }}
      </div>

      <!-- Liste des utilisateurs -->
      <div v-if="loading" class="text-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
        <p class="mt-4 text-gray-500">Chargement des utilisateurs...</p>
      </div>

      <div v-else-if="filteredUsers.length === 0" class="text-center py-12">
        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"></path>
        </svg>
        <h3 class="mt-2 text-sm font-medium text-gray-900">Aucun utilisateur trouvé</h3>
        <p class="mt-1 text-sm text-gray-500">
          {{ showActiveOnly ? 'Aucun utilisateur actif.' : 'Aucun utilisateur.' }}
        </p>
      </div>

      <div v-else class="grid grid-cols-1 gap-6">
        <div
          v-for="user in filteredUsers"
          :key="user.id"
          class="bg-white rounded-lg shadow p-6"
        >
          <div class="flex items-center justify-between">
            <div class="flex items-center space-x-4">
              <div class="flex-shrink-0">
                <div class="h-12 w-12 rounded-full bg-gray-300 flex items-center justify-center">
                  <span class="text-lg font-medium text-gray-700">
                    {{ user.firstname.charAt(0).toUpperCase() }}{{ user.lastname.charAt(0).toUpperCase() }}
                  </span>
                </div>
              </div>
              <div>
                <h3 class="text-lg font-medium text-gray-900">
                  {{ user.firstname }} {{ user.lastname }}
                </h3>
                <p class="text-sm text-gray-500">{{ user.email }}</p>
                <p class="text-xs text-gray-400">
                  Créé le {{ formatDate(user.created_at) }}
                </p>
              </div>
            </div>

            <div class="flex items-center space-x-4">
              <span
                :class="user.active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
                class="inline-flex px-3 py-1 text-xs font-semibold rounded-full"
              >
                {{ user.active ? 'Actif' : 'Inactif' }}
              </span>

              <span class="inline-flex px-3 py-1 text-xs font-semibold rounded-full bg-blue-100 text-blue-800">
                {{ user.role?.name || 'N/A' }}
              </span>

              <div class="flex space-x-2">
                <button
                  @click="editUser(user)"
                  class="text-blue-600 hover:text-blue-800 p-2"
                  title="Modifier"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
                  </svg>
                </button>

                <button
                  v-if="user.active"
                  @click="confirmDeactivateUser(user)"
                  class="text-red-600 hover:text-red-800 p-2"
                  title="Désactiver"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 5.636l-12.728 12.728m0 0L5.636 18.364m12.728-12.728L18.364 18.364M12 2a10 10 0 100 20 10 10 0 000-20z"></path>
                  </svg>
                </button>

                <button
                  v-else
                  @click="confirmActivateUser(user)"
                  class="text-green-600 hover:text-green-800 p-2"
                  title="Activer"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Modal d'édition -->
      <div
        v-if="showEditModal"
        class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50"
        @click="closeModal"
      >
        <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white" @click.stop>
          <div class="mt-3">
            <h3 class="text-lg font-medium text-gray-900 mb-4">
              Modifier l'utilisateur
            </h3>

            <form @submit.prevent="saveUser" class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700">Prénom</label>
                <input
                  v-model="editForm.firstname"
                  type="text"
                  required
                  class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700">Nom</label>
                <input
                  v-model="editForm.lastname"
                  type="text"
                  required
                  class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700">Email</label>
                <input
                  v-model="editForm.email"
                  type="email"
                  required
                  class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700">Rôle</label>
                <select
                  v-model="editForm.role_id"
                  class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                >
                  <option v-for="role in availableRoles" :key="role.id" :value="role.id">
                    {{ role.name }}
                  </option>
                </select>
              </div>

              <div class="flex justify-end space-x-3 pt-4">
                <button
                  type="button"
                  @click="closeModal"
                  class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 border border-gray-300 rounded-md hover:bg-gray-200"
                >
                  Annuler
                </button>
                <button
                  type="submit"
                  :disabled="saving"
                  class="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 disabled:opacity-50"
                >
                  {{ saving ? 'Enregistrement...' : 'Enregistrer' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>

      <!-- Modal de confirmation -->
      <div
        v-if="showConfirmModal"
        class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50"
        @click="closeConfirmModal"
      >
        <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white" @click.stop>
          <div class="mt-3">
            <h3 class="text-lg font-medium text-gray-900 mb-4">
              {{ confirmAction === 'deactivate' ? 'Désactiver' : 'Activer' }} l'utilisateur
            </h3>
            <p class="text-sm text-gray-500 mb-4">
              Êtes-vous sûr de vouloir {{ confirmAction === 'deactivate' ? 'désactiver' : 'activer' }}
              l'utilisateur <strong>{{ selectedUser?.firstname }} {{ selectedUser?.lastname }}</strong> ?
            </p>

            <div class="flex justify-end space-x-3">
              <button
                @click="closeConfirmModal"
                class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 border border-gray-300 rounded-md hover:bg-gray-200"
              >
                Annuler
              </button>
              <button
                @click="executeAction"
                :disabled="saving"
                :class="confirmAction === 'deactivate' ? 'bg-red-600 hover:bg-red-700' : 'bg-green-600 hover:bg-green-700'"
                class="px-4 py-2 text-sm font-medium text-white border border-transparent rounded-md disabled:opacity-50"
              >
                {{ saving ? 'Traitement...' : (confirmAction === 'deactivate' ? 'Désactiver' : 'Activer') }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useSafebaseStore } from '@/stores/safebase'
import type { User } from '@/types/auth'
import type { UserUpdateRequest, UserRoleUpdateRequest } from '@/types/user'

// Store
const safebaseStore = useSafebaseStore()
const {
  users,
  loading,
  error,
  userCount,
  activeUsers: activeUsersGetter,
  inactiveUsers: inactiveUsersGetter,
  adminUsers: adminUsersGetter
} = storeToRefs(safebaseStore)

// État local
const showActiveOnly = ref(false)
const showEditModal = ref(false)
const showConfirmModal = ref(false)
const saving = ref(false)
const selectedUser = ref<User | null>(null)
const confirmAction = ref<'activate' | 'deactivate'>('activate')

// État du formulaire d'édition
const editForm = ref<UserUpdateRequest & { role_id?: number }>({
  firstname: '',
  lastname: '',
  email: '',
  role_id: undefined
})

// Rôles disponibles (hardcodés pour l'instant)
const availableRoles = ref([
  { id: 1, name: 'admin' },
  { id: 2, name: 'user' }
])

// Computed
const activeUsersCount = computed(() => activeUsersGetter.value.length)
const inactiveUsersCount = computed(() => inactiveUsersGetter.value.length)
const adminUsersCount = computed(() => adminUsersGetter.value.length)

const filteredUsers = computed(() => {
  if (showActiveOnly.value) {
    return users.value.filter(user => user.active)
  }
  return users.value
})

// Méthodes
const fetchUsers = async () => {
  try {
    await safebaseStore.fetchUsers()
  } catch (err) {
    console.error('Erreur lors du chargement des utilisateurs:', err)
  }
}

const formatDate = (dateString: string): string => {
  return new Date(dateString).toLocaleDateString('fr-FR')
}

const editUser = (user: User) => {
  selectedUser.value = user
  editForm.value = {
    firstname: user.firstname,
    lastname: user.lastname,
    email: user.email,
    role_id: user.role_id
  }
  showEditModal.value = true
}

const closeModal = () => {
  showEditModal.value = false
  selectedUser.value = null
  editForm.value = {
    firstname: '',
    lastname: '',
    email: '',
    role_id: undefined
  }
}

const saveUser = async () => {
  if (!selectedUser.value) return

  saving.value = true
  try {
    await safebaseStore.updateUser(selectedUser.value.id, editForm.value)
    closeModal()
  } catch (err) {
    console.error('Erreur lors de la sauvegarde:', err)
  } finally {
    saving.value = false
  }
}

const confirmDeactivateUser = (user: User) => {
  selectedUser.value = user
  confirmAction.value = 'deactivate'
  showConfirmModal.value = true
}

const confirmActivateUser = (user: User) => {
  selectedUser.value = user
  confirmAction.value = 'activate'
  showConfirmModal.value = true
}

const closeConfirmModal = () => {
  showConfirmModal.value = false
  selectedUser.value = null
  confirmAction.value = 'activate'
}

const executeAction = async () => {
  if (!selectedUser.value) return

  saving.value = true
  try {
    if (confirmAction.value === 'deactivate') {
      await safebaseStore.deactivateUser(selectedUser.value.id)
    } else {
      await safebaseStore.activateUser(selectedUser.value.id)
    }
    closeConfirmModal()
  } catch (err) {
    console.error('Erreur lors de l\'action:', err)
  } finally {
    saving.value = false
  }
}

// Lifecycle
onMounted(async () => {
  await fetchUsers()
})
</script>
