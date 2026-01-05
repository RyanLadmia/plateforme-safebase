<template>
  <div class="w-full">
    <!-- Formulaire de connexion/inscription -->
    <div v-if="!isAuthenticated" class="bg-white/95 backdrop-blur-sm rounded-2xl shadow-2xl p-8">
      <!-- Toggle buttons -->
      <div class="flex bg-gray-100 rounded-lg p-1 mb-8">
        <button 
          @click="currentForm = 'login'" 
          :class="[
            'flex-1 py-2 px-4 rounded-md font-medium transition-all duration-200',
            currentForm === 'login' 
              ? 'bg-white text-blue-600 shadow-sm' 
              : 'text-gray-600 hover:text-gray-800'
          ]"
        >
          Connexion
        </button>
        <button 
          @click="currentForm = 'register'" 
          :class="[
            'flex-1 py-2 px-4 rounded-md font-medium transition-all duration-200',
            currentForm === 'register' 
              ? 'bg-white text-blue-600 shadow-sm' 
              : 'text-gray-600 hover:text-gray-800'
          ]"
        >
          Inscription
        </button>
      </div>

      <!-- Formulaire de connexion -->
      <form v-if="currentForm === 'login'" @submit.prevent="handleLogin" class="space-y-6">
        <h2 class="text-2xl font-bold text-gray-800 text-center mb-6">Connexion</h2>
        
        <div class="space-y-4">
          <div>
            <label for="login-email" class="block text-sm font-medium text-gray-700 mb-2">
              Email
            </label>
            <input
              id="login-email"
              v-model="loginForm.email"
              type="email"
              required
              placeholder="votre@email.com"
              class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200"
            />
          </div>

          <div>
            <label for="login-password" class="block text-sm font-medium text-gray-700 mb-2">
              Mot de passe
            </label>
            <div class="relative">
              <input
                id="login-password"
                v-model="loginForm.password"
                :type="showLoginPassword ? 'text' : 'password'"
                required
                placeholder="Votre mot de passe"
                class="w-full px-4 py-3 pr-12 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200"
              />
              <button
                type="button"
                @click="showLoginPassword = !showLoginPassword"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none"
                :aria-label="showLoginPassword ? 'Masquer le mot de passe' : 'Afficher le mot de passe'"
              >
                <svg v-if="!showLoginPassword" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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

        <button 
          type="submit" 
          :disabled="loading" 
          class="w-full bg-gradient-to-r from-blue-600 to-purple-600 text-white py-3 px-4 rounded-lg font-semibold hover:from-blue-700 hover:to-purple-700 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
        >
          {{ loading ? 'Connexion...' : 'Se connecter' }}
        </button>
      </form>

      <!-- Formulaire d'inscription -->
      <form v-if="currentForm === 'register'" @submit.prevent="handleRegister" class="space-y-6">
        <h2 class="text-2xl font-bold text-gray-800 text-center mb-6">Inscription</h2>
        
        <div class="space-y-4">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label for="register-firstname" class="block text-sm font-medium text-gray-700 mb-2">
                Prénom
              </label>
              <input
                id="register-firstname"
                v-model="registerForm.firstname"
                type="text"
                required
                placeholder="Votre prénom"
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200"
              />
            </div>

            <div>
              <label for="register-lastname" class="block text-sm font-medium text-gray-700 mb-2">
                Nom
              </label>
              <input
                id="register-lastname"
                v-model="registerForm.lastname"
                type="text"
                required
                placeholder="Votre nom"
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200"
              />
            </div>
          </div>

          <div>
            <label for="register-email" class="block text-sm font-medium text-gray-700 mb-2">
              Email
            </label>
            <input
              id="register-email"
              v-model="registerForm.email"
              type="email"
              required
              placeholder="votre@email.com"
              class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200"
            />
          </div>

          <div>
            <label for="register-password" class="block text-sm font-medium text-gray-700 mb-2">
              Mot de passe
            </label>
            <div class="relative">
              <input
                id="register-password"
                v-model="registerForm.password"
                :type="showRegisterPassword ? 'text' : 'password'"
                required
                placeholder="Min. 10 caract�res, maj, min, chiffre, spécial"
                class="w-full px-4 py-3 pr-12 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200"
              />
              <button
                type="button"
                @click="showRegisterPassword = !showRegisterPassword"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none"
                :aria-label="showRegisterPassword ? 'Masquer le mot de passe' : 'Afficher le mot de passe'"
              >
                <svg v-if="!showRegisterPassword" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                </svg>
                <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"></path>
                </svg>
              </button>
            </div>
            <p class="text-xs text-gray-500 mt-2">
              Le mot de passe doit contenir au moins 10 caractères avec majuscules, minuscules, chiffres et caractères spéciaux.
            </p>
          </div>

          <div>
            <label for="register-confirm-password" class="block text-sm font-medium text-gray-700 mb-2">
              Confirmer le mot de passe
            </label>
            <div class="relative">
              <input
                id="register-confirm-password"
                v-model="registerForm.confirm_password"
                :type="showRegisterConfirmPassword ? 'text' : 'password'"
                required
                placeholder="Confirmez votre mot de passe"
                class="w-full px-4 py-3 pr-12 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200"
              />
              <button
                type="button"
                @click="showRegisterConfirmPassword = !showRegisterConfirmPassword"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none"
                :aria-label="showRegisterConfirmPassword ? 'Masquer le mot de passe' : 'Afficher le mot de passe'"
              >
                <svg v-if="!showRegisterConfirmPassword" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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

        <button 
          type="submit" 
          :disabled="loading" 
          class="w-full bg-gradient-to-r from-blue-600 to-purple-600 text-white py-3 px-4 rounded-lg font-semibold hover:from-blue-700 hover:to-purple-700 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
        >
          {{ loading ? 'Inscription...' : 'S\'inscrire' }}
        </button>
      </form>
    </div>

    <!-- Messages d'erreur/succ�s -->
    <div v-if="message" :class="[
      'mt-4 p-4 rounded-lg text-center font-medium transition-all duration-200',
      messageType === 'success' 
        ? 'bg-green-100 text-green-800 border border-green-200' 
        : 'bg-red-100 text-red-800 border border-red-200'
    ]">
      {{ message }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import type { LoginRequest, RegisterRequest, FormType, MessageType } from '@/types/auth'

// Composables
const authStore = useAuthStore()

// Emits
const emit = defineEmits<{
  'login-success': []
}>()

// State
const currentForm = ref<FormType>('login')
const loading = ref<boolean>(false)
const message = ref<string>('')
const messageType = ref<MessageType>('success')
const showLoginPassword = ref<boolean>(false)
const showRegisterPassword = ref<boolean>(false)
const showRegisterConfirmPassword = ref<boolean>(false)

// Form data
const loginForm = ref<LoginRequest>({
  email: '',
  password: ''
})

const registerForm = ref<RegisterRequest>({
  firstname: '',
  lastname: '',
  email: '',
  password: '',
  confirm_password: ''
})

// Computed réactifs depuis le store (avec storeToRefs pour préserver la réactivité)
const { isAuthenticated, user: currentUser, initialized } = storeToRefs(authStore)

// Methods
const showMessage = (text: string, type: MessageType = 'success'): void => {
  message.value = text
  messageType.value = type
  setTimeout(() => {
    message.value = ''
  }, 5000)
}

const handleLogin = async (): Promise<void> => {
  loading.value = true
  try {
    await authStore.login(loginForm.value)
    showMessage('Connexion réussie !', 'success')
    
    // Réinitialise le formulaire
    loginForm.value = { email: '', password: '' }
    
    // �met l'événement de connexion réussie
    emit('login-success')
  } catch (error) {
    showMessage(error instanceof Error ? error.message : 'Erreur de connexion', 'error')
  } finally {
    loading.value = false
  }
}

const handleRegister = async (): Promise<void> => {
  loading.value = true
  try {
    await authStore.register(registerForm.value)
    showMessage('Inscription réussie ! Vous pouvez maintenant vous connecter.', 'success')
    
    // Réinitialise le formulaire et passe � la connexion
    registerForm.value = { firstname: '', lastname: '', email: '', password: '', confirm_password: '' }
    currentForm.value = 'login'
  } catch (error) {
    showMessage(error instanceof Error ? error.message : 'Erreur d\'inscription', 'error')
  } finally {
    loading.value = false
  }
}

const handleLogout = async (): Promise<void> => {
  loading.value = true
  try {
    await authStore.logout()
    showMessage('Déconnexion réussie !', 'success')
  } catch (error) {
    showMessage(error instanceof Error ? error.message : 'Erreur de déconnexion', 'error')
  } finally {
    loading.value = false
  }
}

// Lifecycle
onMounted((): void => {
  // Attendre que le store soit initialisé avant de vérifier l'authentification
  watch(initialized, (isInit) => {
    if (isInit && isAuthenticated.value) {
      showMessage('Vous �tes déj� connecté !', 'success')
    }
  }, { immediate: true })
})
</script>