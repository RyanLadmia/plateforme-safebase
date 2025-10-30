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
      <!-- Info Message -->
      <div class="bg-blue-100 border-l-4 border-blue-500 text-blue-700 p-4 mb-6 rounded">
        <p class="font-medium">Note:</p>
        <p class="text-sm">La gestion des utilisateurs nécessite l'implémentation des endpoints backend correspondants (/api/admin/users).</p>
      </div>

      <!-- Users List (TODO: Backend implementation needed) -->
      <div class="bg-white rounded-lg shadow">
        <div class="px-6 py-4 border-b border-gray-200">
          <h2 class="text-xl font-semibold text-gray-900">Liste des utilisateurs</h2>
        </div>
        <div class="p-6">
          <div class="text-center py-8">
            <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"></path>
            </svg>
            <h3 class="mt-2 text-sm font-medium text-gray-900">Fonctionnalité à venir</h3>
            <p class="mt-1 text-sm text-gray-500">
              Cette fonctionnalité sera disponible une fois les endpoints backend implémentés.
            </p>
            <div class="mt-6">
              <h4 class="text-sm font-medium text-gray-700 mb-2">Endpoints requis :</h4>
              <ul class="text-sm text-gray-600 space-y-1">
                <li>• GET /api/admin/users - Liste des utilisateurs</li>
                <li>• DELETE /api/admin/users/:id - Supprimer un utilisateur</li>
                <li>• PUT /api/admin/users/:id/role - Modifier le rôle</li>
              </ul>
            </div>
          </div>
        </div>
      </div>

      <!-- Sample User Card (for demonstration) -->
      <div class="mt-6 bg-white rounded-lg shadow p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Aperçu (exemple)</h3>
        <div class="border border-gray-200 rounded-lg p-4">
          <div class="flex items-center justify-between">
            <div>
              <p class="font-medium text-gray-900">{{ user?.email }}</p>
              <p class="text-sm text-gray-500">ID: {{ user?.id }}</p>
            </div>
            <div class="flex items-center space-x-4">
              <span class="px-3 py-1 bg-red-100 text-red-800 rounded-full text-xs font-semibold">
                {{ user?.role || 'ADMIN' }}
              </span>
              <button 
                disabled
                class="px-4 py-2 bg-gray-300 text-gray-500 rounded-lg cursor-not-allowed"
              >
                Actions
              </button>
            </div>
          </div>
        </div>
        <p class="text-xs text-gray-500 mt-4">
          * Ceci est votre propre compte administrateur. Les fonctionnalités de gestion apparaîtront ici une fois l'API implémentée.
        </p>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const { user } = authStore

// TODO: Implémenter les fonctionnalités suivantes une fois le backend prêt:
// - Charger la liste des utilisateurs
// - Supprimer un utilisateur
// - Modifier le rôle d'un utilisateur
// - Créer un nouvel utilisateur
</script>
