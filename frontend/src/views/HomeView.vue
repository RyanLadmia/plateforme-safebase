<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth'

// Composables
const authStore = useAuthStore()

// Computed réactifs depuis le store (avec storeToRefs pour préserver la réactivité)
const { isAuthenticated, user: currentUser } = storeToRefs(authStore)
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
    <!-- Hero Section -->
    <section class="relative overflow-hidden">
      <!-- Background Pattern -->
      <div class="absolute inset-0 bg-gradient-to-br from-blue-600 via-purple-600 to-indigo-800 opacity-10"></div>
      
      <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24">
        <div class="text-center">  
          <!-- Titre principal -->
          <h1 class="text-5xl md:text-7xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent mb-6">
            SafeBase
          </h1>
          
          <!-- Sous-titre -->
          <p class="text-xl md:text-2xl text-gray-600 mb-8 max-w-3xl mx-auto">
            Plateforme de gestion sécurisée avec authentification avancée
          </p>
          
          <!-- Description -->
          <p class="text-lg text-gray-500 mb-12 max-w-2xl mx-auto">
            Système d'authentification moderne développé avec Go et Vue.js, 
            offrant une sécurité de niveau entreprise pour vos données sensibles.
          </p>

          <!-- Informations utilisateur si connecté -->
          <div v-if="isAuthenticated" class="bg-white rounded-2xl shadow-xl p-8 max-w-2xl mx-auto mb-12">
            
            
            <h2 class="text-2xl font-bold text-gray-800 mb-4">
              Bienvenue, {{ currentUser?.firstname }} {{ currentUser?.lastname }} !
            </h2>
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-left">
              <div class="bg-gray-50 p-4 rounded-lg">
                <span class="text-sm font-medium text-gray-500">Email</span>
                <p class="text-gray-800 font-semibold">{{ currentUser?.email }}</p>
              </div>
              <div class="bg-gray-50 p-4 rounded-lg">
                <span class="text-sm font-medium text-gray-500">Rôle</span>
                <p class="font-semibold">
                  <span 
                    :class="currentUser?.role_id === 1 
                      ? 'text-yellow-600 bg-yellow-100 px-2 py-1 rounded-full text-sm' 
                      : 'text-green-600 bg-green-100 px-2 py-1 rounded-full text-sm'"
                  >
                    {{ currentUser?.role_id === 1 ? 'Administrateur' : 'Utilisateur' }}
                  </span>
                </p>
              </div>
            </div>
          </div>

          <!-- Call to action si non connecté -->
          <div v-else class="space-y-4">
            <router-link 
              to="/login"
              class="inline-flex items-center px-8 py-4 bg-gradient-to-r from-blue-600 to-purple-600 text-white font-semibold rounded-full hover:from-blue-700 hover:to-purple-700 transform hover:scale-105 transition-all duration-200 shadow-lg"
            >
              <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1" />
              </svg>
              Accéder à la plateforme
            </router-link>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="py-20 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-16">
          <h2 class="text-3xl md:text-4xl font-bold text-gray-800 mb-4">
            Fonctionnalités principales
          </h2>
          <p class="text-lg text-gray-600 max-w-2xl mx-auto">
            Une solution complète pour la gestion sécurisée de vos données
          </p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
          <!-- Feature 1 -->
          <div class="bg-gradient-to-br from-blue-50 to-indigo-100 p-6 rounded-xl hover:shadow-lg transition-shadow duration-200">
            <h3 class="text-xl font-semibold text-gray-800 mb-2">Authentification JWT</h3>
            <p class="text-gray-600">Système de tokens sécurisés avec expiration automatique</p>
          </div>

          <!-- Feature 2 -->
          <div class="bg-gradient-to-br from-purple-50 to-pink-100 p-6 rounded-xl hover:shadow-lg transition-shadow duration-200">
            <h3 class="text-xl font-semibold text-gray-800 mb-2">Gestion des rôles</h3>
            <p class="text-gray-600">Système de permissions granulaires par utilisateur</p>
          </div>

          <!-- Feature 3 -->
          <div class="bg-gradient-to-br from-green-50 to-emerald-100 p-6 rounded-xl hover:shadow-lg transition-shadow duration-200">
            <h3 class="text-xl font-semibold text-gray-800 mb-2">Sécurité avancée</h3>
            <p class="text-gray-600">Chiffrement bcrypt et validation stricte</p>
          </div>

          <!-- Feature 4 -->
          <div class="bg-gradient-to-br from-yellow-50 to-orange-100 p-6 rounded-xl hover:shadow-lg transition-shadow duration-200">
            <h3 class="text-xl font-semibold text-gray-800 mb-2">Performance</h3>
            <p class="text-gray-600">API Go haute performance avec PostgreSQL</p>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>
