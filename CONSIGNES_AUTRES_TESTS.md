# 📝 CONSIGNES POUR ADAPTER LES AUTRES TESTS E2E

## 🎯 Résumé des adaptations nécessaires

Tous les tests E2E doivent être adaptés à votre architecture spécifique. Voici les règles à suivre :

---

## ✅ RÈGLES IMPORTANTES

### 1. **Routes correctes**

❌ **INCORRECT** :
```typescript
cy.visit('/databases')
cy.visit('/backups')
cy.visit('/schedules')
cy.visit('/dashboard')
```

✅ **CORRECT** :
```typescript
cy.visit('/user/databases')
cy.visit('/user/backups')
cy.visit('/user/schedules')
cy.visit('/user/dashboard')  // ou /admin/dashboard pour admin
```

### 2. **Pas de localStorage, utiliser les cookies**

❌ **INCORRECT** :
```typescript
cy.window().then((win) => {
  const token = win.localStorage.getItem('token')
  // ...
})
```

✅ **CORRECT** :
```typescript
// Le token est dans les cookies HTTP-only
// Pas besoin de le récupérer manuellement
// Il est automatiquement envoyé avec chaque requête

// Pour vérifier l'authentification :
cy.getCookies().should('have.length.at.least', 1)
```

### 3. **Optimisation : Créer l'utilisateur UNE SEULE FOIS**

❌ **INCORRECT** (LENT) :
```typescript
beforeEach(() => {
  // Crée un utilisateur à CHAQUE test
  cy.registerUser({ ... })
  cy.login(...)
})
```

✅ **CORRECT** (RAPIDE) :
```typescript
let testUser: any

before(() => {
  // Crée l'utilisateur UNE SEULE FOIS pour toute la suite
  const timestamp = Date.now()
  testUser = {
    firstname: 'John',
    lastname: 'Doe',
    email: `test.${timestamp}@e2e.com`,
    password: 'TestP@ssw0rd123'
  }
  
  // Inscription
  cy.visit('/login')
  cy.contains('button', 'Inscription').click()
  cy.get('input#register-firstname').type(testUser.firstname)
  cy.get('input#register-lastname').type(testUser.lastname)
  cy.get('input#register-email').type(testUser.email)
  cy.get('input#register-password').type(testUser.password)
  cy.get('input#register-confirm-password').type(testUser.password)
  cy.get('button[type="submit"]').click()
  cy.wait(2000)
})

beforeEach(() => {
  // Se connecter avant chaque test (réutilise testUser)
  cy.visit('/login')
  cy.contains('button', 'Connexion').click()
  cy.get('input#login-email').type(testUser.email)
  cy.get('input#login-password').type(testUser.password)
  cy.get('button[type="submit"]').click()
  cy.url({ timeout: 10000 }).should('match', /dashboard/)
})
```

### 4. **Ne PAS utiliser `cy.createDatabase()` via API**

❌ **INCORRECT** (401 Unauthorized) :
```typescript
cy.createDatabase({
  name: 'Test DB',
  type: 'mysql',
  ...
})
```

✅ **CORRECT** (Via l'UI) :
```typescript
// Utiliser l'interface utilisateur
cy.visit('/user/databases')
cy.contains(/nouvelle|ajouter/i).click()
cy.get('input[name="name"]').type('Test DB')
cy.get('select[name="type"]').select('mysql')
// ... remplir les autres champs
cy.get('button[type="submit"]').click()
```

### 5. **Vérifier la structure réelle de l'UI**

Avant d'écrire un test, VÉRIFIEZ l'application manuellement :

```bash
# Ouvrir dans le navigateur
open http://localhost:3000/user/databases
```

Inspectez les éléments pour trouver les bons sélecteurs :
- IDs des inputs
- Textes des boutons
- Classes CSS
- Structure du formulaire

---

## 📋 STRUCTURE TYPE D'UN TEST

```typescript
describe('Feature X', () => {
  let testUser: any
  
  // Créer l'utilisateur UNE FOIS
  before(() => {
    const timestamp = Date.now()
    testUser = {
      firstname: 'John',
      lastname: 'Doe',
      email: `feature-x.${timestamp}@e2e.com`,
      password: 'TestP@ssw0rd123'
    }
    
    // Inscription
    cy.visit('/login')
    cy.contains('button', 'Inscription').click()
    cy.get('input#register-firstname').type(testUser.firstname)
    cy.get('input#register-lastname').type(testUser.lastname)
    cy.get('input#register-email').type(testUser.email)
    cy.get('input#register-password').type(testUser.password)
    cy.get('input#register-confirm-password').type(testUser.password)
    cy.get('button[type="submit"]').click()
    cy.wait(2000)
  })
  
  // Se connecter avant chaque test
  beforeEach(() => {
    cy.visit('/login')
    cy.contains('button', 'Connexion').click()
    cy.get('input#login-email').type(testUser.email)
    cy.get('input#login-password').type(testUser.password)
    cy.get('button[type="submit"]').click()
    cy.url({ timeout: 10000 }).should('match', /dashboard/)
    
    // Naviguer vers la page à tester
    cy.visit('/user/feature-x')
  })
  
  it('should do something', () => {
    // Votre test ici
  })
})
```

---

## 🔍 INSPECTION DES ÉLÉMENTS

### Pour trouver les bons sélecteurs :

1. **Ouvrir la console du navigateur** (F12)
2. **Cliquer sur l'inspecteur** (icône flèche)
3. **Sélectionner l'élément** dans la page
4. **Noter** :
   - L'ID (ex: `id="database-name"`)
   - Le name (ex: `name="name"`)
   - Les classes (ex: `class="btn-primary"`)
   - Le texte (ex: `Nouvelle base de données`)

### Exemple DatabasesView.vue :

```vue
<button @click="showCreateModal = true" class="px-4 py-2 bg-blue-600...">
  + Nouvelle base de données
</button>
```

**Sélecteur Cypress :**
```typescript
cy.contains('button', 'Nouvelle base de données').click()
// ou
cy.contains(/nouvelle base/i).click()
```

---

## ⚠️ PROBLÈMES COURANTS

### Problème 1 : "Expected to find content: '/ajouter|créer|nouvelle/i' but never did"

**Cause** : Le texte du bouton est différent

**Solution** : Vérifier manuellement le texte exact
```typescript
// Au lieu de
cy.contains(/ajouter|créer|nouvelle/i)

// Utiliser le texte exact trouvé dans l'UI
cy.contains('+ Nouvelle base de données')
```

### Problème 2 : "401 Unauthorized" lors de cy.request()

**Cause** : Le token n'est pas envoyé correctement dans les requêtes API

**Solution** : NE PAS utiliser `cy.request()` pour créer des données, utiliser l'UI à la place
```typescript
// ❌ ÉVITER
cy.request({
  method: 'POST',
  url: `${Cypress.env('apiUrl')}/databases`,
  body: { ... }
})

// ✅ PRÉFÉRER
cy.visit('/user/databases')
cy.contains('Nouvelle base de données').click()
// Remplir le formulaire via l'UI
```

### Problème 3 : Tests très lents

**Cause** : Inscription + Connexion à chaque test

**Solution** : Utiliser `before()` au lieu de `beforeEach()`
- `before()` = exécuté UNE FOIS pour toute la suite
- `beforeEach()` = exécuté avant CHAQUE test

---

## 📝 CHECKLIST AVANT D'ÉCRIRE UN TEST

- [ ] Vérifier l'URL correcte (`/user/xxx` et non `/xxx`)
- [ ] Inspecter l'UI manuellement dans le navigateur
- [ ] Noter les textes exacts des boutons/liens
- [ ] Noter les IDs/names des inputs
- [ ] Créer l'utilisateur dans `before()` (pas `beforeEach()`)
- [ ] Se connecter dans `beforeEach()`
- [ ] Utiliser l'UI pour créer des données (pas `cy.request()`)
- [ ] Vérifier que les cookies sont utilisés (pas `localStorage`)

---

## 🎯 RECOMMANDATION

**Pour les tests 02-database-management.cy.ts et suivants :**

1. **NE PAS les exécuter pour l'instant**
2. **Se concentrer sur le test d'authentification** (01-authentication.cy.ts) qui FONCTIONNE
3. **Adapter UN test à la fois** :
   - Lire le code de l'UI
   - Tester manuellement
   - Écrire le test Cypress
   - Vérifier qu'il passe
   - Passer au suivant

**OU**

**Attendre que je crée des versions adaptées** une par une en inspectant votre UI réelle.

---

## 🚀 PROCHAINES ÉTAPES

### Option 1 : Vous adaptez les tests

Suivez les consignes ci-dessus pour adapter chaque test.

### Option 2 : Je les adapte pour vous

Dites-moi quel test vous voulez que j'adapte en priorité :
- Database management ?
- Backup management ?
- Schedule management ?
- History ?
- Profile ?

Je l'adapterai complètement à votre structure réelle.

---

## 📚 Documentation

- **Tests adaptés** : `01-authentication.cy.ts` ✅
- **Tests à adapter** : Tous les autres ❌

**Pour l'instant, concentrez-vous sur le test d'authentification qui FONCTIONNE !**

---

**Date** : Janvier 2026  
**Version** : 1.0.5  
**Statut** : ⚠️ Seul le test d'auth est adapté

