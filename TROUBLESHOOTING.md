# 🐛 Erreurs fréquentes et solutions - Tests Cypress

## ❌ Erreur : "exports is not defined in ES module scope"

### Problème
```
ReferenceError: exports is not defined in ES module scope
Your configFile is invalid: cypress.config.ts
```

### Solution ✅
Le problème vient de `"type": "module"` dans `package.json` qui entre en conflit avec Cypress.

**Retirer la ligne `"type": "module"` du package.json** (déjà fait ✅)

```bash
cd tests
# Le fichier a été corrigé automatiquement
npm run cy:open
```

---

## ❌ Erreur : "You are attempting to run a TypeScript file, but do not have TypeScript installed"

### Problème
```
Error: You are attempting to run a TypeScript file, but do not have TypeScript installed.
```

### Solution ✅
```bash
cd tests
npm install typescript @types/node --save-dev
# ou simplement
npm install
```

TypeScript est maintenant dans les dépendances du projet !

---

## ❌ Erreur : "EPERM: operation not permitted" ou "npm install" échoue

### Problème
```
npm error code EPERM
npm error syscall open
npm error errno -1
```

### Solutions ✅

**Option 1 : Nettoyer le cache npm**
```bash
cd tests
rm -rf node_modules package-lock.json
npm cache clean --force
npm install
```

**Option 2 : Avec sudo (macOS/Linux)**
```bash
cd tests
sudo npm install
```

**Option 3 : Corriger les permissions**
```bash
# Trouver le propriétaire npm
ls -la ~/.npm

# Corriger les permissions
sudo chown -R $(whoami) ~/.npm
sudo chown -R $(whoami) /usr/local/lib/node_modules

# Réessayer
cd tests
npm install
```

---

## ❌ Erreur : "Connection refused" sur http://localhost:3000

### Problème
Les tests ne peuvent pas se connecter au frontend.

### Solutions ✅

**Vérifier que Docker tourne :**
```bash
docker-compose ps
```

**Si services arrêtés :**
```bash
docker-compose up -d
# Attendre 1-2 minutes
```

**Vérifier manuellement :**
```bash
curl http://localhost:3000      # Frontend
curl http://localhost:8080/api  # Backend
```

**Voir les logs :**
```bash
docker-compose logs -f frontend
docker-compose logs -f backend
```

---

## ❌ Erreur : "Timed out waiting for the browser to connect"

### Problème
Cypress ne peut pas démarrer le navigateur ou se connecter.

### Solutions ✅

**Option 1 : Augmenter les timeouts**
```typescript
// Dans cypress.config.ts (déjà fait)
defaultCommandTimeout: 15000,
pageLoadTimeout: 90000,
```

**Option 2 : Effacer le cache Cypress**
```bash
cd tests
npx cypress cache clear
npx cypress install
```

**Option 3 : Vérifier les navigateurs disponibles**
```bash
npx cypress info
```

---

## ❌ Erreur : Module not found ou Cannot find module

### Problème
```
Error: Cannot find module 'cypress'
Error: Cannot find module '@types/node'
```

### Solution ✅
```bash
cd tests
npm install
```

---

## ❌ Erreur : "baseUrl" is not responding

### Problème
```
Cypress cannot verify that this server is running:
> http://localhost:3000
```

### Solutions ✅

**1. Vérifier la configuration**
```bash
# Pour Docker
cat tests/.env
# Devrait contenir:
CYPRESS_BASE_URL=http://localhost:3000
```

**2. Vérifier que le frontend est accessible**
```bash
curl http://localhost:3000
# Devrait retourner du HTML
```

**3. Attendre plus longtemps**
```bash
# Docker peut prendre 1-2 minutes au démarrage
docker-compose logs -f frontend
# Attendre "ready" ou "listening"
```

**4. Redémarrer Docker**
```bash
docker-compose restart frontend
```

---

## ❌ Erreur : Tests échouent avec "database not found"

### Problème
Les tests d'intégration Backend ne trouvent pas la base de données.

### Solutions ✅

**1. Vérifier PostgreSQL Docker**
```bash
docker-compose ps postgres
# Devrait être "Up"
```

**2. Vérifier backend/.env**
```bash
cat backend/.env
# Devrait contenir:
DB_HOST=localhost  # ou "postgres" si dans Docker
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=safebase
```

**3. Tester la connexion**
```bash
docker-compose exec postgres psql -U user -d safebase
# Devrait se connecter
```

---

## ❌ Erreur : Port already in use

### Problème
```
Error: listen EADDRINUSE: address already in use :::3000
Error: listen EADDRINUSE: address already in use :::8080
```

### Solutions ✅

**Trouver le processus qui utilise le port :**
```bash
# Port 3000 (frontend)
lsof -i :3000

# Port 8080 (backend)
lsof -i :8080

# Port 5432 (PostgreSQL)
lsof -i :5432
```

**Tuer le processus :**
```bash
# Remplacer PID par le numéro affiché
kill -9 PID
```

**Ou arrêter Docker complètement :**
```bash
docker-compose down
docker-compose up -d
```

---

## ❌ Erreur : "cy.login is not a function"

### Problème
Les commandes personnalisées ne sont pas reconnues.

### Solution ✅
```bash
# Vérifier que le fichier support existe
ls tests/e2E/support/commands.ts

# Si manquant, réinstaller
cd tests
npm install
```

---

## ❌ Erreur : Vidéos/Screenshots non créés

### Problème
Pas de vidéos ou screenshots après les tests.

### Solution ✅

**Vérifier la configuration :**
```typescript
// Dans cypress.config.ts
video: true,
screenshotOnRunFailure: true,
```

**Vérifier les dossiers :**
```bash
ls tests/e2E/videos/
ls tests/e2E/screenshots/
```

**Les créer si manquants :**
```bash
mkdir -p tests/e2E/videos
mkdir -p tests/e2E/screenshots
```

---

## ❌ Erreur : Tests Go "no such table"

### Problème
```
Error: no such table: users
```

### Solution ✅

**C'est normal pour les tests unitaires** (utilisent SQLite en mémoire vide).

Les migrations sont dans les helpers :
```go
db.AutoMigrate(&models.User{}, &models.Role{}, &models.Session{})
```

Si l'erreur persiste :
```bash
cd backend
go clean -testcache
go test ./tests/units/... -v
```

---

## ❌ Erreur : "context deadline exceeded"

### Problème
Les requêtes prennent trop de temps.

### Solutions ✅

**1. Vérifier que Docker n'est pas surchargé**
```bash
docker stats
```

**2. Redémarrer Docker**
```bash
docker-compose restart
```

**3. Augmenter les timeouts Cypress**
```bash
# Dans tests/.env
CYPRESS_defaultCommandTimeout=20000
CYPRESS_pageLoadTimeout=120000
```

---

## ❌ Erreur : npm audit vulnerabilities

### Problème
```
4 high severity vulnerabilities
```

### Solution ✅

**Voir les détails :**
```bash
cd tests
npm audit
```

**Corriger automatiquement :**
```bash
npm audit fix
# Si ça ne suffit pas
npm audit fix --force
```

**Note :** Certaines vulnérabilités peuvent nécessiter des mises à jour de Cypress lui-même.

---

## 🔧 Commandes de dépannage générales

### Tout réinitialiser (tests)
```bash
cd tests
rm -rf node_modules package-lock.json
npm cache clean --force
npm install
```

### Tout réinitialiser (Docker)
```bash
docker-compose down -v
docker-compose up -d --build
```

### Vérifier l'état complet
```bash
# Docker
docker-compose ps
docker-compose logs --tail=50

# Frontend
curl http://localhost:3000

# Backend
curl http://localhost:8080/api

# PostgreSQL
docker-compose exec postgres pg_isready -U user
```

### Logs en temps réel
```bash
# Tous les services
docker-compose logs -f

# Service spécifique
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f postgres
```

---

## 📞 Besoin d'aide ?

### Checklist avant de demander de l'aide :

- [ ] J'ai vérifié que Docker tourne : `docker-compose ps`
- [ ] J'ai vérifié les logs : `docker-compose logs`
- [ ] J'ai essayé de redémarrer : `docker-compose restart`
- [ ] J'ai nettoyé npm : `rm -rf node_modules && npm install`
- [ ] J'ai attendu 2 minutes après `docker-compose up`
- [ ] J'ai vérifié les URLs manuellement avec `curl`

### Collecter les informations :

```bash
# Versions
node --version
npm --version
docker --version
docker-compose --version

# État
docker-compose ps
ls -la tests/node_modules/typescript/

# Logs récents
docker-compose logs --tail=100 > docker-logs.txt
```

---

## 🎯 Résumé des solutions rapides

| Erreur | Solution rapide |
|--------|----------------|
| TypeScript missing | `cd tests && npm install` |
| EPERM | `npm cache clean --force && npm install` |
| Connection refused | `docker-compose up -d` + attendre 2 min |
| Timeout | Augmenter timeouts dans config |
| Module not found | `cd tests && npm install` |
| Port in use | `lsof -i :PORT` puis `kill PID` |
| Database error | Vérifier `backend/.env` |

---

**Dernière mise à jour** : Janvier 2026  
**Version** : 1.0.0  
**Fichier** : TROUBLESHOOTING.md

