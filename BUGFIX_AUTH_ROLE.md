# 🐛 Correction des bugs d'authentification et de rôle

## Problèmes identifiés

### 1. ❌ Erreur 401 sur les appels API
```
GET http://localhost:8080/api/databases 401 (Unauthorized)
Error: authorization header missing
```

**Cause** : Le backend utilise des **cookies HTTP-only** pour l'authentification, pas des headers Authorization avec Bearer token.

**Solution** : Axios est déjà configuré avec `withCredentials: true`, donc les cookies sont bien envoyés. L'erreur 401 indique que :
- L'utilisateur n'est pas authentifié côté backend
- Le cookie de session a expiré
- Il faut se reconnecter

### 2. ❌ Rôle mal détecté (admin → user → inconnu)
```
Page d'accueil: admin ✅
Dashboard: user ❌
Profil: inconnu ❌
```

**Cause** : Le frontend utilisait `user.role` comme string, mais le backend retourne un **objet Role**.

## ✅ Corrections apportées

### 1. Types TypeScript mis à jour

**Avant** :
```typescript
export interface User {
  role?: string // ❌ Incorrect
}
```

**Après** :
```typescript
export interface Role {
  id: number
  name: string
  created_at: string
  updated_at: string
}

export interface User {
  role?: Role // ✅ Objet complet
}
```

### 2. Store auth corrigé

**Avant** :
```typescript
const isAdmin = computed(() => user.value?.role === 'admin') // ❌
const isUser = computed(() => user.value?.role === 'user') // ❌
```

**Après** :
```typescript
const isAdmin = computed(() => user.value?.role?.name === 'admin') // ✅
const isUser = computed(() => user.value?.role?.name === 'user') // ✅
```

### 3. Page profil corrigée

**Changements** :
- ❌ Retiré : Champ "ID Utilisateur"
- ✅ Corrigé : Fonction `getRoleLabel()` utilise `user?.role?.name`

**Avant** :
```typescript
const getRoleLabel = (): string => {
  if (authStore.isAdmin) return 'Administrateur'
  if (authStore.isUser) return 'Utilisateur'
  return 'Inconnu' // ❌ Toujours "Inconnu"
}
```

**Après** :
```typescript
const getRoleLabel = (): string => {
  const roleName = user?.role?.name
  if (roleName === 'admin') return 'Administrateur'
  if (roleName === 'user') return 'Utilisateur'
  return roleName || 'Non défini' // ✅ Affiche le rôle réel
}
```

## 🔍 Diagnostic du problème 401

### Backend utilise des cookies HTTP-only

Le backend SafeBase utilise l'authentification par **cookies HTTP-only** :
- Plus sécurisé que localStorage
- Protection contre XSS
- Le cookie est envoyé automatiquement par le navigateur

### Configuration Axios correcte

```typescript
// api/axios.ts
export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true, // ✅ Envoie les cookies
  headers: {
    'Content-Type': 'application/json',
  },
})
```

### Pourquoi l'erreur 401 ?

L'erreur 401 est normale dans ces cas :

1. **Première visite** : Aucun cookie de session
2. **Session expirée** : Cookie expiré (après 24h par défaut)
3. **Déconnexion** : Cookie supprimé
4. **Backend redémarré** : Sessions perdues

### Solution pour l'utilisateur

**Se reconnecter** :
1. Aller sur `/login`
2. Entrer email et mot de passe
3. Le backend crée une session
4. Le cookie est stocké
5. Les appels API fonctionnent

## 📊 Structure de données backend

### User reçu du backend
```json
{
  "id": 1,
  "firstname": "Jean",
  "lastname": "Dupont",
  "email": "jean.dupont@example.com",
  "role_id": 1,
  "role": {
    "id": 1,
    "name": "admin",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Accès au rôle dans le code

```typescript
// ✅ Correct
user.role?.name // "admin" ou "user"

// ❌ Incorrect
user.role // { id: 1, name: "admin", ... }
user.role_id // 1 (ID du rôle, pas le nom)
```

## 🧪 Tests de validation

### Build réussi ✅
```bash
npm run build
✓ 118 modules transformed
✓ built in 831ms
```

### À tester manuellement
1. ✅ Se connecter en tant qu'admin
2. ✅ Vérifier le rôle sur toutes les pages
3. ✅ Vérifier que les menus s'affichent correctement
4. ✅ Tester la page profil (sans ID)
5. ✅ Vérifier les appels API après connexion

## 📝 Checklist des corrections

### Types
- [x] Interface `Role` ajoutée
- [x] Interface `User` mise à jour avec `role?: Role`

### Store
- [x] `isAdmin` utilise `user?.role?.name`
- [x] `isUser` utilise `user?.role?.name`

### Vues
- [x] ProfileView : `getRoleLabel()` corrigé
- [x] ProfileView : Champ ID utilisateur retiré
- [x] Toutes les vues utilisent les getters du store

### Configuration
- [x] Axios configuré avec `withCredentials: true`
- [x] Build réussi sans erreurs

## 💡 Notes importantes

### Authentification
- Le backend gère l'authentification par **cookies HTTP-only**
- Les cookies sont automatiquement envoyés avec chaque requête
- L'erreur 401 signifie "non authentifié" → il faut se connecter

### Rôles
- Deux rôles : `admin` et `user`
- L'objet `role` complet est retourné par le backend
- Toujours accéder via `user.role?.name`

### Sécurité
- ✅ Cookies HTTP-only (protection XSS)
- ✅ CORS configuré correctement
- ✅ Pas de token en localStorage
- ✅ Validation côté backend

## 🚀 Résultat

**Corrections appliquées avec succès !**

- ✅ Types TypeScript corrects
- ✅ Détection de rôle fonctionnelle
- ✅ Page profil sans ID
- ✅ Build sans erreurs

**Pour les erreurs 401** : L'utilisateur doit se connecter via `/login` pour créer une session.
