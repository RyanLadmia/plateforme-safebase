import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth'

// Types pour les meta des routes
declare module 'vue-router' {
  interface RouteMeta {
    requiresAuth?: boolean
    requiresAdmin?: boolean
    requiresGuest?: boolean
    title?: string
  }
}

// Configuration des routes
const routes: RouteRecordRaw[] = [
  // Routes publiques
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
    component: () => import('../views/AboutView.vue'),
    meta: {
      title: 'À propos - SafeBase'
    }
  },

  // Routes utilisateurs (authentification requise)
  {
    path: '/user/dashboard',
    name: 'user-dashboard',
    component: () => import('@/views/users/DashboardView.vue'),
    meta: {
      requiresAuth: true,
      title: 'Tableau de bord - SafeBase'
    }
  },
  {
    path: '/user/profile',
    name: 'user-profile',
    component: () => import('@/views/users/ProfileView.vue'),
    meta: {
      requiresAuth: true,
      title: 'Mon profil - SafeBase'
    }
  },
  {
    path: '/user/databases',
    name: 'user-databases',
    component: () => import('@/views/users/DatabasesView.vue'),
    meta: {
      requiresAuth: true,
      title: 'Mes bases de données - SafeBase'
    }
  },
  {
    path: '/user/backups',
    name: 'user-backups',
    component: () => import('@/views/users/BackupsView.vue'),
    meta: {
      requiresAuth: true,
      title: 'Mes sauvegardes - SafeBase'
    }
  },
  {
    path: '/user/schedules',
    name: 'user-schedules',
    component: () => import('@/views/users/SchedulesView.vue'),
    meta: {
      requiresAuth: true,
      title: 'Mes sauvegardes planifiées - SafeBase'
    }
  },

  // Routes administrateur (rôle admin requis)
  {
    path: '/admin/dashboard',
    name: 'admin-dashboard',
    component: () => import('@/views/admins/AdminDashboardView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
      title: 'Administration - SafeBase'
    }
  },
  {
    path: '/admin/users',
    name: 'admin-users',
    component: () => import('@/views/admins/UsersManagementView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: true,
      title: 'Gestion des utilisateurs - SafeBase'
    }
  },

  // Route 404
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('@/views/NotFoundView.vue'),
    meta: {
      title: 'Page introuvable - SafeBase'
    }
  }
]

// Création du router
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

// Guard de navigation pour l'authentification
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  const { isAuthenticated, isAdmin, initialized } = storeToRefs(authStore)
  
  // Attendre l'initialisation du store
  if (!initialized.value) {
    await new Promise(resolve => {
      const unwatch = authStore.$subscribe((mutation, state) => {
        if (state.initialized) {
          unwatch()
          resolve(true)
        }
      })
    })
  }
  
  // Mise à jour du titre de la page
  if (to.meta.title) {
    document.title = to.meta.title
  }
  
  // Si la route nécessite une authentification et l'utilisateur n'est pas connecté
  if (to.meta.requiresAuth && !isAuthenticated.value) {
    next({ name: 'login', query: { redirect: to.fullPath } })
    return
  }
  
  // Si la route nécessite le rôle admin et l'utilisateur n'est pas admin
  if (to.meta.requiresAdmin && !isAdmin.value) {
    next({ name: 'home' })
    return
  }
  
  // Si la route est pour les invités et l'utilisateur est connecté
  if (to.meta.requiresGuest && isAuthenticated.value) {
    // Rediriger selon le rôle
    if (isAdmin.value) {
      next({ name: 'admin-dashboard' })
    } else {
      next({ name: 'user-dashboard' })
    }
    return
  }
  
  next()
})

export default router
