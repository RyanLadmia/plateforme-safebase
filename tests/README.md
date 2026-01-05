# Tests End-to-End (E2E) avec Cypress

## 📋 Vue d'ensemble

Cette suite de tests E2E couvre **au moins 90%** de la plateforme SafeBase, testant tous les flux utilisateur critiques de bout en bout.

**⚠️ IMPORTANT : Cette configuration supporte Docker et le développement local !**

## 🐳 Modes d'exécution

### Mode Docker (RECOMMANDÉ - par défaut)
- Frontend sur port **3000**
- Backend sur port **8080**
- Utilise `docker-compose.yml`

### Mode Local (sans Docker)
- Frontend sur port **5173** (Vite dev)
- Backend sur port **8080**
- Serveurs lancés manuellement

## 🚀 Installation

### Avec Docker (RECOMMANDÉ)

```bash
# Option 1: Script automatique (le plus simple)
./test-docker.sh

# Option 2: Manuel
docker-compose up -d
cd tests
npm install

# Créer .env pour Docker
cat > .env << EOF
CYPRESS_BASE_URL=http://localhost:3000
CYPRESS_API_URL=http://localhost:8080/api
CYPRESS_IS_DOCKER=true
EOF
```

### Sans Docker (développement local)

```bash
cd tests
npm install

# Créer .env pour local
cat > .env << EOF
CYPRESS_BASE_URL=http://localhost:5173
CYPRESS_API_URL=http://localhost:8080/api
CYPRESS_IS_DOCKER=false
EOF
```

## 🏃 Exécution des tests

### Avec Docker

```bash
# Démarrer les services d'abord
docker-compose up -d

# Mode interactif
cd tests
npm run cy:open:docker
# ou simplement (si .env est configuré)
npm run cy:open

# Mode headless
npm run test:docker
# ou
npm run test
```

### Sans Docker

```bash
# Démarrer les services manuellement
# Terminal 1: Backend
cd backend && go run cmd/main.go

# Terminal 2: Frontend
cd frontend && npm run dev

# Terminal 3: Tests
cd tests
npm run cy:open:local
# ou
npm run test:local
```

### Scripts disponibles

```bash
# Mode interactif (GUI)
npm run cy:open              # Utilise .env
npm run cy:open:docker       # Force Docker (port 3000)
npm run cy:open:local        # Force local (port 5173)

# Mode headless
npm run test                 # Utilise .env
npm run test:docker          # Force Docker
npm run test:local           # Force local
npm run test:ci              # Pour CI/CD

# Par navigateur
npm run cy:run:chrome
npm run cy:run:firefox
npm run cy:run:edge
```

## 📁 Structure des tests

```
tests/
├── e2E/
│   ├── 01-authentication.cy.ts         # Tests d'authentification (15%)
│   ├── 02-database-management.cy.ts    # Tests de gestion des BDD (25%)
│   ├── 03-backup-management.cy.ts      # Tests de sauvegardes (20%)
│   ├── 04-schedule-management.cy.ts    # Tests de planification (15%)
│   ├── 05-history.cy.ts                # Tests d'historique (10%)
│   ├── 06-profile.cy.ts                # Tests de profil utilisateur (10%)
│   ├── 07-dashboard.cy.ts              # Tests du tableau de bord (5%)
│   └── 08-complete-workflows.cy.ts     # Flux complets (10%)
│   ├── fixtures/
│   │   ├── users.json                  # Données utilisateurs de test
│   │   ├── databases.json              # Données de BDD de test
│   │   └── schedules.json              # Données de planifications
│   └── support/
│       ├── e2e.ts                      # Configuration globale
│       └── commands.ts                 # Commandes personnalisées
├── cypress.config.ts                   # Configuration Cypress
├── package.json                        # Dépendances npm
└── README.md                           # Ce fichier
```

## 🎯 Couverture des tests

### 1. Authentification (15%)
- ✅ Inscription utilisateur
- ✅ Connexion/Déconnexion
- ✅ Validation des mots de passe
- ✅ Gestion des sessions
- ✅ Gestion des erreurs

### 2. Gestion des bases de données (25%)
- ✅ CRUD complet (Create, Read, Update, Delete)
- ✅ Support MySQL et PostgreSQL
- ✅ Validation des champs
- ✅ Chiffrement des mots de passe
- ✅ Filtres et recherche
- ✅ Actions multiples

### 3. Gestion des sauvegardes (20%)
- ✅ Création manuelle de sauvegardes
- ✅ Téléchargement de sauvegardes
- ✅ Restauration de sauvegardes
- ✅ Suppression de sauvegardes
- ✅ Filtres par date, statut, BDD
- ✅ Opérations en masse
- ✅ Pagination

### 4. Planification (15%)
- ✅ Création de planifications
- ✅ Validation des expressions CRON
- ✅ Activation/Désactivation
- ✅ Modification et suppression
- ✅ Historique d'exécution
- ✅ Planifications multiples

### 5. Historique & Audit (10%)
- ✅ Traçabilité complète des actions
- ✅ Filtres par type, ressource, date
- ✅ Recherche par mots-clés
- ✅ Export CSV
- ✅ Isolation multi-utilisateurs
- ✅ Pagination

### 6. Profil utilisateur (10%)
- ✅ Affichage du profil
- ✅ Modification des informations
- ✅ Changement de mot de passe
- ✅ Statistiques utilisateur
- ✅ Préférences et paramètres
- ✅ Sécurité du compte

### 7. Tableau de bord (5%)
- ✅ Affichage des statistiques
- ✅ Activité récente
- ✅ Actions rapides
- ✅ Planifications à venir
- ✅ Notifications
- ✅ Navigation

### 8. Flux complets (10%)
- ✅ Parcours complet nouvel utilisateur
- ✅ Workflow de sauvegarde
- ✅ Gestion multi-bases
- ✅ Collaboration multi-utilisateurs
- ✅ Récupération d'erreurs
- ✅ Export de données

## 🛠️ Commandes personnalisées

### `cy.login(email, password)`
Authentifie un utilisateur.

```typescript
cy.login('user@example.com', 'password123')
```

### `cy.registerUser(userData)`
Crée un nouveau compte utilisateur.

```typescript
cy.registerUser({
  firstname: 'John',
  lastname: 'Doe',
  email: 'john@example.com',
  password: 'StrongP@ssw0rd123'
})
```

### `cy.createDatabase(dbData)`
Crée une configuration de base de données via l'API.

```typescript
cy.createDatabase({
  name: 'Test DB',
  type: 'mysql',
  host: 'localhost',
  port: '3306',
  username: 'user',
  password: 'pass',
  db_name: 'db'
})
```

### `cy.createSchedule(scheduleData)`
Crée une planification via l'API.

```typescript
cy.createSchedule({
  database_id: 1,
  name: 'Daily Backup',
  cron_expression: '0 2 * * *'
})
```

### `cy.deleteAllTestData()`
Nettoie toutes les données de test.

```typescript
cy.deleteAllTestData()
```

## 📊 Rapports et résultats

Les résultats des tests sont sauvegardés dans :
- **Vidéos** : `e2E/videos/`
- **Screenshots** : `e2E/screenshots/`
- **Screenshots d'échecs** : `e2E/screenshots/failed/`

## ⚙️ Configuration

### Variables d'environnement

Dans `cypress.config.ts` :
```typescript
env: {
  apiUrl: 'http://localhost:8080/api',  // URL de l'API backend
  coverage: true                        // Activer la couverture
}
```

### Configuration personnalisée

Pour modifier la configuration, éditez `cypress.config.ts`.

## 🔧 Prérequis

Avant d'exécuter les tests :

1. **Backend** : Le serveur Go doit être démarré
   ```bash
   cd backend
   go run cmd/main.go
   ```

2. **Frontend** : Le serveur de développement doit être actif
   ```bash
   cd frontend
   npm run dev
   ```

3. **Base de données** : PostgreSQL doit être accessible

## 🐛 Débogage

### Activer les logs détaillés
```bash
DEBUG=cypress:* npm run cy:run
```

### Exécuter un seul fichier de test
```bash
npx cypress run --spec "e2E/01-authentication.cy.ts"
```

### Exécuter des tests spécifiques
```bash
npx cypress run --spec "e2E/**/*authentication*.cy.ts"
```

## 📝 Bonnes pratiques

1. **Isolation des tests** : Chaque test doit être indépendant
2. **Nettoyage** : Utilisez `afterEach()` pour nettoyer les données
3. **Sélecteurs** : Préférez `data-cy` aux sélecteurs CSS
4. **Attentes** : Utilisez des timeouts appropriés
5. **Fixtures** : Utilisez des fixtures pour les données de test

## 🚨 Gestion des erreurs

Les tests gèrent automatiquement :
- ❌ Erreurs réseau
- ❌ Timeouts
- ❌ Erreurs de validation
- ❌ Erreurs serveur (500)
- ❌ Données invalides

## 📈 Métriques de qualité

- **Couverture** : >90% du code frontend
- **Fiabilité** : Tests réexécutés 2 fois en cas d'échec (CI)
- **Performance** : Chargement des pages <3s
- **Accessibilité** : Tests A11Y inclus

## 🔄 Intégration CI/CD

Ces tests sont intégrés dans le pipeline CI/CD et s'exécutent automatiquement sur chaque push.

Voir `.github/workflows/ci-cd.yml` pour la configuration.

## 📚 Documentation

- [Cypress Documentation](https://docs.cypress.io)
- [Best Practices](https://docs.cypress.io/guides/references/best-practices)
- [TypeScript Support](https://docs.cypress.io/guides/tooling/typescript-support)

## 🤝 Contribution

Pour ajouter de nouveaux tests :

1. Créez un nouveau fichier `.cy.ts` dans `e2E/`
2. Suivez la structure existante
3. Ajoutez des fixtures si nécessaire
4. Documentez les nouveaux tests
5. Assurez-vous du nettoyage avec `afterEach()`

## ⚡ Performance

Pour optimiser les performances :
- Utilisez `cy.session()` pour la réutilisation des sessions
- Limitez les `cy.wait()` explicites
- Utilisez des interceptions pour mocker les appels API lents
- Parallélisez les tests en CI

## 📞 Support

Pour toute question ou problème, consultez la documentation ou créez une issue.

---

**Dernière mise à jour** : Janvier 2026
**Version** : 1.0.0
**Couverture** : 90%+

