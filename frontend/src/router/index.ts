import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth'

// Types pour les meta des routes
declare module 'vue-router' {
  interface RouteMeta {
    requiresAuth?: boolean
    requiresGuest?: boolean
    title?: string
  }
}

// Configuration des routes
const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'home',
    component: HomeView,
    meta: { 
      title: 'Accueil - SafeBase' 
    }
  },
  {
    path: '/login',
    name: 'login',
    component: LoginView,
    meta: { 
      requiresGuest: true,
      title: 'Connexion - SafeBase'
    }
  },
  {
    path: '/about',
    name: 'about',
    // Route avec code-splitting
    // Génère un chunk séparé (About.[hash].js) pour cette route
    // qui est chargé de manière lazy quand la route est visitée
    component: () => import('../views/AboutView.vue'),
    meta: {
      title: 'À propos - SafeBase'
    }
  },
]

// Création du router
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

// Guard de navigation pour l'authentification
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  const { isAuthenticated } = storeToRefs(authStore)
  
  // Mise à jour du titre de la page
  if (to.meta.title) {
    document.title = to.meta.title
  }
  
  // Si la route nécessite une authentification et l'utilisateur n'est pas connecté
  if (to.meta.requiresAuth && !isAuthenticated.value) {
    next('/login')
    return
  }
  
  // Si la route est pour les invités et l'utilisateur est connecté
  if (to.meta.requiresGuest && isAuthenticated.value) {
    next('/')
    return
  }
  
  next()
})

export default router
