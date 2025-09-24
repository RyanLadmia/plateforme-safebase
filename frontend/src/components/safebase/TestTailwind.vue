<template>
  <div class="bg-gradient-to-br from-blue-50 to-indigo-100 p-8 rounded-xl shadow-lg max-w-md mx-auto">
    <h2 class="text-2xl font-bold text-gray-800 mb-4 text-center">
      SafeBase
    </h2>
    <p class="text-gray-600 text-center mb-6">
      Plateforme de sauvegarde et restauration de bases de données
    </p>
    
    <div class="space-y-3">
      <button class="w-full bg-blue-600 hover:bg-blue-700 text-white font-semibold py-3 px-4 rounded-lg transition-colors duration-200 shadow-md">
        Nouvelle sauvegarde
      </button>
      
      <button class="w-full bg-green-600 hover:bg-green-700 text-white font-semibold py-3 px-4 rounded-lg transition-colors duration-200 shadow-md">
        Restaurer
      </button>
      
      <button class="w-full bg-gray-600 hover:bg-gray-700 text-white font-semibold py-3 px-4 rounded-lg transition-colors duration-200 shadow-md">
        Historique
      </button>
    </div>
    
    <div class="mt-6 p-4 bg-yellow-50 border-l-4 border-yellow-400 rounded-r-lg">
      <p class="text-sm text-yellow-800">
        <strong>Tailwind CSS installé avec succès!</strong><br>
        Vue.js + TypeScript + Vite + Tailwind CSS sont prêts à l'usage.
      </p>
    </div>
  </div>
</template>
