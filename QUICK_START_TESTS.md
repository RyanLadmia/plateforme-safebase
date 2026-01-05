# 🎯 Guide Rapide - Tests Plateforme SafeBase

## 📦 Installation initiale

### Backend (Tests Go - déjà installés)
```bash
cd backend
go mod download
```

### Frontend (Tests E2E Cypress)
```bash
# Option 1: Script automatique
./install-cypress.sh

# Option 2: Manuel
cd tests
npm install
```

---

## 🧪 Exécution des tests

### Tests Backend (Go)

#### Tous les tests
```bash
cd backend
go test ./tests/... -v
```

#### Par catégorie
```bash
# Tests unitaires uniquement
go test ./tests/units/... -v

# Tests d'intégration uniquement
go test ./tests/integrations/... -v

# Tests fonctionnels uniquement
go test ./tests/functionals/... -v
```

#### Tests spécifiques
```bash
# Test d'authentification
go test ./tests/units/auth_unit_test.go -v

# Test de gestion des BDD
go test ./tests/units/database_unit_test.go -v

# Test de sauvegardes
go test ./tests/units/backup_unit_test.go -v
```

#### Avec couverture
```bash
go test ./tests/units/... -v -cover
go test ./tests/integrations/... -v -cover
go test ./tests/functionals/... -v -cover
```

#### Avec race detection
```bash
go test ./tests/... -v -race
```

---

### Tests E2E (Cypress)

#### Prérequis (services à démarrer)
```bash
# Terminal 1: PostgreSQL
docker-compose up -d postgres

# Terminal 2: Backend
cd backend
go run cmd/main.go

# Terminal 3: Frontend
cd frontend
npm run dev
```

#### Mode développement (GUI)
```bash
cd tests
npm run cy:open
```

#### Mode CI/CD (Headless)
```bash
cd tests
npm run test
```

#### Tests spécifiques
```bash
# Par fichier
npx cypress run --spec "e2E/01-authentication.cy.ts"

# Par pattern
npx cypress run --spec "e2E/**/*database*.cy.ts"
```

#### Par navigateur
```bash
npm run cy:run:chrome     # Chrome
npm run cy:run:firefox    # Firefox
npm run cy:run:edge       # Edge
```

#### Debug mode
```bash
DEBUG=cypress:* npm run cy:run
```

---

## 📊 Résultats et rapports

### Backend
```bash
# Avec verbose
go test ./tests/... -v

# Avec benchmarks
go test ./tests/... -v -bench=.

# Sauvegarder les résultats
go test ./tests/... -v > test-results.txt
```

### Cypress
- **Vidéos** : `tests/e2E/videos/`
- **Screenshots** : `tests/e2E/screenshots/`
- **Screenshots d'échecs** : `tests/e2E/screenshots/failed/`

---

## 🐛 Dépannage

### Tests Backend échouent

#### Vérifier PostgreSQL
```bash
# Status
docker-compose ps

# Logs
docker-compose logs postgres

# Redémarrer
docker-compose restart postgres
```

#### Vérifier les variables d'environnement
```bash
cat backend/.env
```

#### Nettoyer et relancer
```bash
cd backend
go clean -testcache
go test ./tests/... -v
```

### Tests Cypress échouent

#### Vérifier les services
```bash
# Backend
curl http://localhost:8080/api

# Frontend
curl http://localhost:5173
```

#### Nettoyer cache Cypress
```bash
cd tests
rm -rf node_modules package-lock.json
npm cache clean --force
npm install
npx cypress cache clear
npx cypress install
```

#### Voir les logs détaillés
```bash
DEBUG=cypress:* npm run cy:run
```

---

## 📁 Structure des fichiers de tests

```
backend/tests/
├── units/                    # Tests unitaires (isolation)
│   ├── auth_unit_test.go
│   ├── database_unit_test.go
│   ├── backup_unit_test.go
│   ├── schedule_unit_test.go
│   └── history_unit_test.go
├── integrations/             # Tests d'intégration (composants)
│   ├── db_integration_test.go
│   ├── auth_integration_test.go
│   ├── database_integration_test.go
│   ├── schedule_integration_test.go
│   └── history_integration_test.go
└── functionals/              # Tests fonctionnels (workflows)
    ├── user_workflows_test.go
    ├── schedule_workflows_test.go
    └── audit_workflows_test.go

tests/e2E/                    # Tests E2E Cypress
├── 01-authentication.cy.ts
├── 02-database-management.cy.ts
├── 03-backup-management.cy.ts
├── 04-schedule-management.cy.ts
├── 05-history.cy.ts
├── 06-profile.cy.ts
├── 07-dashboard.cy.ts
└── 08-complete-workflows.cy.ts
```

---

## 🎯 Commandes par cas d'usage

### Avant un commit
```bash
# Tests rapides backend
cd backend
go test ./tests/units/... -v

# Tests rapides E2E (smoke tests)
cd tests
npx cypress run --spec "e2E/01-authentication.cy.ts"
```

### Avant un merge/PR
```bash
# Tous les tests backend
cd backend
go test ./tests/... -v -cover

# Tous les tests E2E
cd tests
npm run test
```

### Pour un module spécifique

#### Authentification
```bash
# Backend
go test ./tests/units/auth_unit_test.go -v
go test ./tests/integrations/auth_integration_test.go -v

# E2E
npx cypress run --spec "e2E/01-authentication.cy.ts"
```

#### Bases de données
```bash
# Backend
go test ./tests/units/database_unit_test.go -v
go test ./tests/integrations/database_integration_test.go -v

# E2E
npx cypress run --spec "e2E/02-database-management.cy.ts"
```

#### Sauvegardes
```bash
# Backend
go test ./tests/units/backup_unit_test.go -v

# E2E
npx cypress run --spec "e2E/03-backup-management.cy.ts"
```

---

## ⚡ Optimisations

### Tests Backend
```bash
# Parallélisation
go test ./tests/... -v -parallel=4

# Cache
go test ./tests/... -v -count=1  # Force sans cache
```

### Tests Cypress
```bash
# Headless plus rapide
npm run cy:run:headless

# Sans vidéos (plus rapide)
npx cypress run --config video=false
```

---

## 📚 Documentation

- **README backend tests** : `backend/tests/README.md`
- **README E2E** : `tests/README.md`
- **Guide Cypress** : `CYPRESS_E2E_SETUP.md`
- **Synthèse complète** : `TEST_SYNTHESIS.md`
- **Ce guide** : `QUICK_START_TESTS.md`

---

## 🔄 Intégration CI/CD

### Exemple GitHub Actions
```yaml
- name: Run Backend Tests
  run: |
    cd backend
    go test ./tests/... -v -cover

- name: Run E2E Tests
  run: |
    cd tests
    npm install
    npm run test:ci
```

---

## ✅ Checklist avant production

- [ ] Tous les tests unitaires passent
- [ ] Tous les tests d'intégration passent
- [ ] Tous les tests fonctionnels passent
- [ ] Tous les tests E2E passent
- [ ] Couverture de code >80%
- [ ] Aucun test flakey
- [ ] Documentation à jour
- [ ] CI/CD configuré

---

## 📞 Support

### Erreurs fréquentes

**"database not found"**
→ Vérifier PostgreSQL et `.env`

**"connection refused"**
→ Démarrer backend et frontend

**"timeout"**
→ Augmenter les timeouts dans la config

**"module not found"**
→ `go mod download` ou `npm install`

---

## 🎉 Stats finales

| Type | Nombre | Temps | Couverture |
|------|--------|-------|------------|
| Unitaires | ~25 | ~3s | Core backend |
| Intégration | ~19 | ~8s | Inter-composants |
| Fonctionnels | ~9 | ~8s | Workflows backend |
| E2E | ~200 | ~20min | Frontend+Backend |
| **TOTAL** | **~253** | **~21min** | **>90%** |

---

**Dernière mise à jour** : Janvier 2026  
**Version** : 1.0.0  
**Statut** : ✅ Prêt pour la production

