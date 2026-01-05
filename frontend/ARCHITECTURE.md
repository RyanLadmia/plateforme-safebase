# Architecture Frontend - SafeBase

## Structure du projet

```
src/
├── api/                    → Appels réseau purs (Axios)
│   ├── axios.ts           → Configuration d'Axios
│   ├── auth_api.ts        → API d'authentification
│   ├── database_api.ts    → API des bases de données
│   ├── backup_api.ts      → API des sauvegardes
│   └── user_api.ts        → API des utilisateurs (Admin)
│
├── services/               → Logique métier
│   ├── auth_service.ts    → Service d'authentification
│   ├── database_service.ts → Service des bases de données
│   └── backup_service.ts  → Service des sauvegardes
│
├── stores/                 → État global (Pinia)
│   ├── auth.ts            → Store d'authentification
│   └── safebase.ts        → Store des données SafeBase
│
├── types/                  → Types TypeScript
│   ├── auth.ts            → Types d'authentification
│   ├── database.ts        → Types des bases de données
│   ├── backup.ts          → Types des sauvegardes
│   └── ui.ts              → Types UI
│
├── components/             → Composants réutilisables
│   ├── auth/              → Composants d'authentification
│   └── ui/                → Composants UI génériques
│
├── views/                  → Pages de l'application
│   ├── [PUBLIC]           → Pages publiques (racine)
│   ├── users/             → Pages utilisateurs authentifiés
│   └── admins/            → Pages administrateurs
│
├── router/                 → Configuration des routes
│   └── index.ts           → Router Vue avec guards
│
└── layout/                 → Layouts de l'application
    └── Header.vue         → En-tête global
```

## Séparation des responsabilités

### 1. **API Layer** (`src/api/`)
- **Responsabilité** : Appels HTTP purs avec Axios
- **Ne contient pas** : Logique métier, gestion d'état
- **Exports** : Fonctions asynchrones qui retournent des promesses
- **Exemple** :
```typescript
export async function getDatabases(): Promise<Database[]> {
  const { data } = await apiClient.get<DatabaseListResponse>('/api/databases')
  return data.databases || []
}
```

### 2. **Services Layer** (`src/services/`)
- **Responsabilité** : Logique métier, validation, transformations
- **Utilise** : Les APIs pour les appels réseau
- **Ne contient pas** : État réactif, gestion Vue
- **Exemple** :
```typescript
export class DatabaseService {
  async createDatabase(data: DatabaseCreateRequest): Promise<Database> {
    this.validateDatabaseData(data) // Validation métier
    return await databaseApi.createDatabase(data)
  }
  
  private validateDatabaseData(data: DatabaseCreateRequest): void {
    // Logique de validation
  }
}
```

### 3. **Stores Layer** (`src/stores/`)
- **Responsabilité** : Gestion d'état global réactif
- **Utilise** : Les services pour les opérations
- **Contient** : État réactif, getters computed, actions
- **Exemple** :
```typescript
export const useSafebaseStore = defineStore('safebase', () => {
  const databases = ref<Database[]>([])
  
  const fetchDatabases = async () => {
    databases.value = await databaseService.fetchDatabases()
  }
  
  return { databases, fetchDatabases }
})
```

## Organisation des vues par rôle

### Pages publiques (`src/views/`)
- **Accessible par** : Tous (authentifiés ou non)
- **Exemples** :
  - `HomeView.vue` - Page d'accueil
  - `LoginView.vue` - Page de connexion
  - `AboutView.vue` - À propos

### Pages utilisateurs (`src/views/users/`)
- **Accessible par** : Utilisateurs authentifiés (role: `user` ou `admin`)
- **Exemples à créer** :
  - `DashboardView.vue` - Tableau de bord
  - `DatabasesView.vue` - Gestion des bases de données
  - `BackupsView.vue` - Gestion des sauvegardes

### Pages administrateurs (`src/views/admins/`)
- **Accessible par** : Administrateurs uniquement (role: `admin`)
- **Exemples à créer** :
  - `AdminDashboardView.vue` - Tableau de bord admin
  - `UsersManagementView.vue` - Gestion des utilisateurs

## Configuration du Router

Le router utilise des guards de navigation pour contrôler l'accès :

```typescript
// Méta des routes
interface RouteMeta {
  requiresAuth?: boolean      // Nécessite une authentification
  requiresAdmin?: boolean     // Nécessite le rôle admin
  requiresGuest?: boolean     // Pour les pages publiques uniquement
  title?: string              // Titre de la page
}
```

**Exemples de routes** :
```typescript
// Route publique
{ path: '/', component: HomeView }

// Route utilisateur
{ 
  path: '/user/dashboard', 
  component: DashboardView,
  meta: { requiresAuth: true }
}

// Route admin
{
  path: '/admin/users',
  component: UsersManagementView,
  meta: { requiresAuth: true, requiresAdmin: true }
}
```

## Flux de données

```
Composant Vue
    ↓ (appel action)
Store Pinia
    ↓ (appel méthode)
Service
    ↓ (appel fonction)
API
    ↓ (requête HTTP)
Backend
```

**Exemple concret** :
```typescript
// 1. Composant appelle le store
const safebaseStore = useSafebaseStore()
await safebaseStore.fetchDatabases()

// 2. Store appelle le service
const fetchDatabases = async () => {
  databases.value = await databaseService.fetchDatabases()
}

// 3. Service appelle l'API
async fetchDatabases(): Promise<Database[]> {
  return await databaseApi.getDatabases()
}

// 4. API fait la requête HTTP
export async function getDatabases(): Promise<Database[]> {
  const { data } = await apiClient.get('/api/databases')
  return data.databases
}
```

## Bonnes pratiques

### À faire
- **API** : Fonctions pures qui font uniquement des appels HTTP
- **Services** : Classes avec méthodes de validation et transformation
- **Stores** : État réactif avec actions qui orchestrent
- **Composants** : Logique UI uniquement, délèguer au store

### À éviter
- Appels HTTP directs dans les composants
- Logique métier dans les stores
- État réactif dans les services
- Duplication de logique entre couches

## Technologies utilisées

- **Vue 3** : Framework frontend
- **TypeScript** : Typage statique
- **Pinia** : Gestion d'état
- **Axios** : Client HTTP
- **Vue Router** : Routing
- **Tailwind CSS** : Styling (si configuré)

## Prochaines étapes

1. Architecture de base mise en place
2. API Layer avec Axios
3. Services Layer avec logique métier
4. Stores refactorisés
5. Créer les vues utilisateurs
6. Créer les vues admin
7. Mettre à jour le router avec les guards

---

Cette architecture permet une séparation claire des responsabilités, facilite les tests et la maintenance, et garantit une évolutivité du code.
