<template>
  <div class="min-h-screen bg-gray-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <div class="flex justify-between items-center mb-6">
        <h1 class="text-2xl sm:text-3xl font-bold text-gray-900">Mon profil</h1>
        <router-link 
          to="/user/dashboard" 
          class="text-blue-600 hover:text-blue-800"
        >
          ← Retour au tableau de bord
        </router-link>
      </div>

      <!-- User Information Card -->
      <div class="bg-white rounded-lg shadow mb-6">
        <div class="px-6 py-4 border-b border-gray-200">
          <h2 class="text-xl font-semibold text-gray-900">Informations personnelles</h2>
        </div>
        <div class="p-6">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Prénom -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Prénom</label>
              <input 
                v-model="profileData.firstname"
                type="text"
                :disabled="!isEditing"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
              />
            </div>

            <!-- Nom -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Nom</label>
              <input 
                v-model="profileData.lastname"
                type="text"
                :disabled="!isEditing"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
              />
            </div>

            <!-- Email -->
            <div class="md:col-span-2">
              <label class="block text-sm font-medium text-gray-700 mb-2">Email</label>
              <input 
                v-model="profileData.email"
                type="email"
                :disabled="!isEditing"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
              />
            </div>

            <!-- Rôle (lecture seule) -->
            <div class="md:col-span-2">
              <label class="block text-sm font-medium text-gray-700 mb-2">Rôle</label>
              <input 
                :value="getRoleLabel()"
                type="text"
                disabled
                class="w-full px-4 py-2 border border-gray-300 rounded-lg bg-gray-100 cursor-not-allowed"
              />
            </div>
          </div>

          <!-- Action Buttons -->
          <div class="mt-6 flex justify-end space-x-4">
            <button 
              v-if="!isEditing"
              @click="startEditing"
              class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
            >
              Modifier
            </button>
            <template v-else>
              <button 
                @click="cancelEditing"
                :disabled="loading"
                class="px-6 py-2 bg-gray-300 text-gray-700 rounded-lg hover:bg-gray-400 transition disabled:opacity-50"
              >
                Annuler
              </button>
              <button 
                @click="saveChanges"
                :disabled="loading"
                class="px-6 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition disabled:opacity-50"
              >
                {{ loading ? 'Enregistrement...' : 'Enregistrer' }}
              </button>
            </template>
          </div>

          <!-- Success/Error Messages -->
          <div v-if="successMessage" class="mt-4 p-4 bg-green-100 text-green-700 rounded-lg">
            {{ successMessage }}
          </div>
          <div v-if="errorMessage" class="mt-4 p-4 bg-red-100 text-red-700 rounded-lg">
            {{ errorMessage }}
          </div>
        </div>
      </div>

      <!-- Change Password Card -->
      <div class="bg-white rounded-lg shadow mb-6">
        <div class="px-6 py-4 border-b border-gray-200">
          <h2 class="text-xl font-semibold text-gray-900">Changer le mot de passe</h2>
        </div>
        <div class="p-6">
          <div class="space-y-4">
            <!-- Current Password -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Mot de passe actuel</label>
              <div class="relative">
                <input 
                  v-model="passwordData.current"
                  :type="showCurrentPassword ? 'text' : 'password'"
                  class="w-full px-4 py-2 pr-12 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                />
                <button
                  type="button"
                  @click="showCurrentPassword = !showCurrentPassword"
                  class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none"
                  :aria-label="showCurrentPassword ? 'Masquer le mot de passe' : 'Afficher le mot de passe'"
                >
                  <svg v-if="!showCurrentPassword" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                  </svg>
                  <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"></path>
                  </svg>
                </button>
              </div>
            </div>

            <!-- New Password -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Nouveau mot de passe</label>
              <div class="relative">
                <input 
                  v-model="passwordData.new"
                  :type="showNewPassword ? 'text' : 'password'"
                  class="w-full px-4 py-2 pr-12 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                />
                <button
                  type="button"
                  @click="showNewPassword = !showNewPassword"
                  class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none"
                  :aria-label="showNewPassword ? 'Masquer le mot de passe' : 'Afficher le mot de passe'"
                >
                  <svg v-if="!showNewPassword" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                  </svg>
                  <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"></path>
                  </svg>
                </button>
              </div>
            </div>

            <!-- Confirm Password -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Confirmer le mot de passe</label>
              <div class="relative">
                <input 
                  v-model="passwordData.confirm"
                  :type="showConfirmPassword ? 'text' : 'password'"
                  class="w-full px-4 py-2 pr-12 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                />
                <button
                  type="button"
                  @click="showConfirmPassword = !showConfirmPassword"
                  class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none"
                  :aria-label="showConfirmPassword ? 'Masquer le mot de passe' : 'Afficher le mot de passe'"
                >
                  <svg v-if="!showConfirmPassword" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                  </svg>
                  <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"></path>
                  </svg>
                </button>
              </div>
            </div>
          </div>

          <div class="mt-6">
            <button 
              @click="changePasswordHandler"
              :disabled="passwordLoading"
              class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition disabled:opacity-50"
            >
              {{ passwordLoading ? 'Modification...' : 'Changer le mot de passe' }}
            </button>
          </div>

          <!-- Password Messages -->
          <div v-if="passwordSuccess" class="mt-4 p-4 bg-green-100 text-green-700 rounded-lg">
            {{ passwordSuccess }}
          </div>
          <div v-if="passwordError" class="mt-4 p-4 bg-red-100 text-red-700 rounded-lg">
            {{ passwordError }}
          </div>
        </div>
      </div>

      <!-- Account Statistics -->
      <div class="bg-white rounded-lg shadow">
        <div class="px-6 py-4 border-b border-gray-200">
          <h2 class="text-xl font-semibold text-gray-900">Statistiques du compte</h2>
        </div>
        <div class="p-6">
          <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div class="text-center">
              <p class="text-3xl font-bold text-blue-600">{{ databaseCount }}</p>
              <p class="text-gray-500 mt-2">Bases de données</p>
            </div>
            <div class="text-center">
              <p class="text-3xl font-bold text-green-600">{{ backupCount }}</p>
              <p class="text-gray-500 mt-2">Sauvegardes totales</p>
            </div>
            <div class="text-center">
              <p class="text-3xl font-bold text-orange-600">{{ completedBackups.length }}</p>
              <p class="text-gray-500 mt-2">Sauvegardes réussies</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import { useSafebaseStore } from '@/stores/safebase'
import { updateProfile, changePassword } from '@/api/profile_api'

const authStore = useAuthStore()
const safebaseStore = useSafebaseStore()

const { user } = storeToRefs(authStore)
const { databaseCount, backupCount, completedBackups } = storeToRefs(safebaseStore)

const isEditing = ref(false)
const loading = ref(false)
const successMessage = ref('')
const errorMessage = ref('')

const profileData = reactive({
  firstname: '',
  lastname: '',
  email: ''
})

const passwordLoading = ref(false)
const passwordSuccess = ref('')
const passwordError = ref('')

const showCurrentPassword = ref(false)
const showNewPassword = ref(false)
const showConfirmPassword = ref(false)

const passwordData = reactive({
  current: '',
  new: '',
  confirm: ''
})

const startEditing = () => {
  isEditing.value = true
  successMessage.value = ''
  errorMessage.value = ''
}

const cancelEditing = () => {
  isEditing.value = false
  // Restaurer les données originales
  if (user.value) {
    profileData.firstname = user.value.firstname
    profileData.lastname = user.value.lastname
    profileData.email = user.value.email
  }
  successMessage.value = ''
  errorMessage.value = ''
}

const saveChanges = async () => {
  loading.value = true
  successMessage.value = ''
  errorMessage.value = ''

  try {
    // Appeler l'API pour mettre à jour le profil
    const updatedUser = await updateProfile({
      firstname: profileData.firstname,
      lastname: profileData.lastname,
      email: profileData.email
    })
    
    successMessage.value = 'Profil mis à jour avec succès !'
    isEditing.value = false
    
    // Rafraîchir les données utilisateur dans le store
    await authStore.checkAuth()
  } catch (err: any) {
    errorMessage.value = err.message || 'Erreur lors de la mise à jour du profil'
  } finally {
    loading.value = false
  }
}

const changePasswordHandler = async () => {
  passwordLoading.value = true
  passwordSuccess.value = ''
  passwordError.value = ''

  // Validation
  if (!passwordData.current || !passwordData.new || !passwordData.confirm) {
    passwordError.value = 'Tous les champs sont requis'
    passwordLoading.value = false
    return
  }

  if (passwordData.new !== passwordData.confirm) {
    passwordError.value = 'Les mots de passe ne correspondent pas'
    passwordLoading.value = false
    return
  }

  if (passwordData.new.length < 8) {
    passwordError.value = 'Le mot de passe doit contenir au moins 8 caractères'
    passwordLoading.value = false
    return
  }

  try {
    // Appeler l'API pour changer le mot de passe
    await changePassword({
      current_password: passwordData.current,
      new_password: passwordData.new,
      confirm_password: passwordData.confirm
    })
    
    passwordSuccess.value = 'Mot de passe changé avec succès !'
    
    // Réinitialiser le formulaire
    passwordData.current = ''
    passwordData.new = ''
    passwordData.confirm = ''
  } catch (err: any) {
    passwordError.value = err.message || 'Erreur lors du changement de mot de passe'
  } finally {
    passwordLoading.value = false
  }
}

const getRoleLabel = (): string => {
  const roleName = user.value?.role?.name
  if (roleName === 'admin') return 'Administrateur'
  if (roleName === 'user') return 'Utilisateur'
  return roleName || 'Non défini'
}

onMounted(() => {
  // Initialiser les données du profil
  if (user.value) {
    profileData.firstname = user.value.firstname
    profileData.lastname = user.value.lastname
    profileData.email = user.value.email
  }

  // Charger les statistiques
  safebaseStore.fetchDatabases().catch(console.error)
  safebaseStore.fetchBackups().catch(console.error)
})
</script>
