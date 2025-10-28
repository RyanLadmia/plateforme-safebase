<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 py-12">
    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- Header -->
      <div class="text-center mb-12">
        <h1 class="text-4xl md:text-5xl font-bold text-gray-800 mb-4">
          À propos de SafeBase
        </h1>
        <p class="text-xl text-gray-600 max-w-2xl mx-auto">
          Une solution moderne et sécurisée pour la gestion de vos données
        </p>
      </div>

      <!-- Contenu principal -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
        <!-- Description -->
        <div class="space-y-6">
          <div class="bg-white rounded-xl shadow-lg p-8">
            <h2 class="text-2xl font-bold text-gray-800 mb-4">
              🚀 Notre Mission
            </h2>
            <p class="text-gray-600 leading-relaxed">
              SafeBase est une plateforme développée avec les technologies les plus modernes 
              pour offrir une solution complète de gestion sécurisée des données. 
              Notre objectif est de simplifier l'authentification et la sécurité 
              tout en maintenant les plus hauts standards de protection.
            </p>
          </div>

          <div class="bg-white rounded-xl shadow-lg p-8">
            <h2 class="text-2xl font-bold text-gray-800 mb-4">
              🛡️ Sécurité
            </h2>
            <p class="text-gray-600 leading-relaxed">
              Nous utilisons les meilleures pratiques de sécurité : 
              chiffrement bcrypt pour les mots de passe, tokens JWT sécurisés, 
              validation stricte des données et architecture moderne 
              pour garantir la protection de vos informations.
            </p>
          </div>
        </div>

        <!-- Technologies -->
        <div class="bg-white rounded-xl shadow-lg p-8">
          <h2 class="text-2xl font-bold text-gray-800 mb-6">
            ⚡ Technologies utilisées
          </h2>
          
          <div class="space-y-4">
            <div class="flex items-center space-x-3">
              <div class="w-3 h-3 bg-blue-500 rounded-full"></div>
              <span class="font-semibold text-gray-700">Backend :</span>
              <span class="text-gray-600">Go (Gin Framework)</span>
            </div>
            
            <div class="flex items-center space-x-3">
              <div class="w-3 h-3 bg-green-500 rounded-full"></div>
              <span class="font-semibold text-gray-700">Frontend :</span>
              <span class="text-gray-600">Vue.js 3 + TypeScript</span>
            </div>
            
            <div class="flex items-center space-x-3">
              <div class="w-3 h-3 bg-purple-500 rounded-full"></div>
              <span class="font-semibold text-gray-700">Base de données :</span>
              <span class="text-gray-600">PostgreSQL + GORM</span>
            </div>
            
            <div class="flex items-center space-x-3">
              <div class="w-3 h-3 bg-yellow-500 rounded-full"></div>
              <span class="font-semibold text-gray-700">Styles :</span>
              <span class="text-gray-600">Tailwind CSS</span>
            </div>
            
            <div class="flex items-center space-x-3">
              <div class="w-3 h-3 bg-red-500 rounded-full"></div>
              <span class="font-semibold text-gray-700">Authentification :</span>
              <span class="text-gray-600">JWT + bcrypt</span>
            </div>
          </div>

          <div class="mt-8 p-4 bg-gray-50 rounded-lg">
            <p class="text-sm text-gray-600 text-center">
              Développé avec ❤️ pour l'apprentissage et la démonstration 
              des bonnes pratiques de développement web moderne.
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
