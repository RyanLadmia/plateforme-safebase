# 🐳 Guide Tests E2E avec Docker - SafeBase

## ⚠️ IMPORTANT : Configuration Docker

Tous les tests ont été adaptés pour fonctionner avec votre environnement Docker !

---

## 🎯 Différences Docker vs Local

### Ports utilisés

| Service | Local (sans Docker) | Docker |
|---------|---------------------|--------|
| **Frontend** | Port 5173 (Vite) | **Port 3000** |
| **Backend** | Port 8080 | Port 8080 |
| **PostgreSQL** | Port 5432 | Port 5432 |
| **MySQL** | Port 3306 | Port 3306 |

### URLs d'accès

**Avec Docker :**
- Frontend : `http://localhost:3000`
- Backend API : `http://localhost:8080/api`

**Sans Docker (local) :**
- Frontend : `http://localhost:5173`
- Backend API : `http://localhost:8080/api`

---

## 🚀 Installation et configuration

### Option 1️⃣ : Script automatique (RECOMMANDÉ)

```bash
# À la racine du projet
./test-docker.sh
```

Ce script va :
1. ✅ Vérifier Docker et Docker Compose
2. ✅ Arrêter les conteneurs existants
3. ✅ Démarrer tous les services Docker
4. ✅ Attendre que les services soient prêts
5. ✅ Configurer Cypress pour Docker
6. ✅ Installer les dépendances si nécessaire

### Option 2️⃣ : Configuration manuelle

#### Étape 1 : Démarrer les services Docker

```bash
# Démarrer tous les services
docker-compose up -d

# Vérifier que tout est démarré
docker-compose ps

# Attendre que les services soient prêts (1-2 minutes)
# Backend : http://localhost:8080
# Frontend : http://localhost:3000
```

#### Étape 2 : Configurer Cypress pour Docker

```bash
cd tests

# Créer le fichier .env pour Docker
cat > .env << EOF
CYPRESS_BASE_URL=http://localhost:3000
CYPRESS_API_URL=http://localhost:8080/api
CYPRESS_IS_DOCKER=true
EOF

# Installer les dépendances Cypress
npm install
```

---

## 🧪 Exécution des tests avec Docker

### Prérequis

**Les services Docker DOIVENT être démarrés :**

```bash
# Vérifier l'état
docker-compose ps

# Démarrer si nécessaire
docker-compose up -d

# Attendre ~30-60 secondes que tout soit prêt
```

### Lancer les tests

```bash
cd tests

# Mode interactif (GUI) - RECOMMANDÉ pour développement
npm run cy:open

# Mode headless (CI/CD)
npm run test

# Test spécifique
npx cypress run --spec "e2E/01-authentication.cy.ts"

# Avec un navigateur spécifique
npm run cy:run:chrome
npm run cy:run:firefox
```

---

## 🔧 Configuration Cypress

### Fichier `cypress.config.ts` (déjà adapté)

La configuration détecte automatiquement l'environnement :

```typescript
baseUrl: process.env.CYPRESS_BASE_URL || 'http://localhost:3000',
env: {
  apiUrl: process.env.CYPRESS_API_URL || 'http://localhost:8080/api',
  isDocker: process.env.CYPRESS_IS_DOCKER || 'false'
}
```

### Fichiers de configuration

**Pour Docker (par défaut) :**
```bash
# tests/.env
CYPRESS_BASE_URL=http://localhost:3000
CYPRESS_API_URL=http://localhost:8080/api
CYPRESS_IS_DOCKER=true
```

**Pour développement local (sans Docker) :**
```bash
# tests/.env
CYPRESS_BASE_URL=http://localhost:5173
CYPRESS_API_URL=http://localhost:8080/api
CYPRESS_IS_DOCKER=false
```

---

## 🐛 Dépannage Docker

### Problème : Frontend non accessible sur port 3000

```bash
# Vérifier si le conteneur tourne
docker-compose ps frontend

# Voir les logs
docker-compose logs frontend

# Redémarrer le frontend
docker-compose restart frontend

# Si ça ne marche pas, reconstruire
docker-compose up -d --build frontend
```

### Problème : Backend non accessible sur port 8080

```bash
# Vérifier le conteneur
docker-compose ps backend

# Logs backend
docker-compose logs backend

# Vérifier la connexion à la BDD
docker-compose logs postgres

# Redémarrer
docker-compose restart backend
```

### Problème : PostgreSQL ne démarre pas

```bash
# Logs PostgreSQL
docker-compose logs postgres

# Nettoyer les volumes et redémarrer
docker-compose down -v
docker-compose up -d
```

### Problème : Tests timeout

```bash
# Les conteneurs peuvent être lents au démarrage
# Attendre 1-2 minutes supplémentaires

# Vérifier que tout est "healthy"
docker-compose ps

# Si un service est "unhealthy", voir ses logs
docker-compose logs [service-name]
```

### Problème : Port déjà utilisé

```bash
# Trouver quel processus utilise le port 3000
lsof -i :3000

# Ou le port 8080
lsof -i :8080

# Arrêter le processus ou changer le port dans docker-compose.yml
```

---

## 📊 Tests Backend Go avec Docker

Les tests backend Go n'ont **PAS besoin** de Docker pour s'exécuter car ils utilisent :
- SQLite en mémoire pour les tests unitaires
- PostgreSQL via `.env` pour les tests d'intégration

### Configuration Backend

Assurez-vous que `backend/.env` contient :

```bash
# Pour tests d'intégration et fonctionnels
DB_HOST=localhost          # ou "postgres" si vous testez depuis un conteneur
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=safebase
DB_SSL_MODE=disable
```

### Exécuter les tests Backend

```bash
cd backend

# Tests unitaires (pas besoin de Docker)
go test ./tests/units/... -v

# Tests d'intégration (besoin de PostgreSQL)
# Si PostgreSQL est dans Docker:
docker-compose up -d postgres
go test ./tests/integrations/... -v

# Tests fonctionnels
go test ./tests/functionals/... -v
```

---

## 🔄 Workflow complet avec Docker

### 1. Démarrage

```bash
# Démarrer tous les services
docker-compose up -d

# Attendre que tout soit prêt
./test-docker.sh
# OU attendre manuellement 1-2 minutes
```

### 2. Tests Backend

```bash
cd backend

# Unitaires (rapide, sans Docker)
go test ./tests/units/... -v

# Intégration (besoin PostgreSQL Docker)
go test ./tests/integrations/... -v

# Fonctionnels
go test ./tests/functionals/... -v
```

### 3. Tests E2E Cypress

```bash
cd tests

# Configuration (si pas déjà fait)
cat > .env << EOF
CYPRESS_BASE_URL=http://localhost:3000
CYPRESS_API_URL=http://localhost:8080/api
CYPRESS_IS_DOCKER=true
EOF

# Installation (si pas déjà fait)
npm install

# Lancer les tests
npm run cy:open    # Mode GUI
# OU
npm run test       # Mode headless
```

### 4. Arrêt

```bash
# Arrêter les services
docker-compose down

# Nettoyer complètement (volumes inclus)
docker-compose down -v
```

---

## 📝 Commandes Docker utiles

### Gestion des services

```bash
# Démarrer tous les services
docker-compose up -d

# Démarrer un service spécifique
docker-compose up -d backend

# Arrêter tous les services
docker-compose down

# Redémarrer un service
docker-compose restart backend

# Reconstruire et redémarrer
docker-compose up -d --build
```

### Logs et debug

```bash
# Voir tous les logs
docker-compose logs -f

# Logs d'un service spécifique
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f postgres

# Dernières 100 lignes
docker-compose logs --tail=100 backend
```

### État et inspection

```bash
# État des conteneurs
docker-compose ps

# Détails d'un conteneur
docker inspect safebase-backend

# Entrer dans un conteneur
docker-compose exec backend sh
docker-compose exec postgres psql -U user -d safebase
```

### Nettoyage

```bash
# Arrêter et supprimer les conteneurs
docker-compose down

# Supprimer aussi les volumes
docker-compose down -v

# Supprimer les images
docker-compose down --rmi all

# Nettoyage complet Docker
docker system prune -a --volumes
```

---

## 🎯 Configurations recommandées

### Pour le développement (avec Docker)

```bash
# tests/.env
CYPRESS_BASE_URL=http://localhost:3000
CYPRESS_API_URL=http://localhost:8080/api
CYPRESS_IS_DOCKER=true
```

### Pour CI/CD

Utiliser des variables d'environnement :

```yaml
# .github/workflows/ci.yml
env:
  CYPRESS_BASE_URL: http://localhost:3000
  CYPRESS_API_URL: http://localhost:8080/api
  CYPRESS_IS_DOCKER: true
```

---

## ⚡ Optimisations Docker pour les tests

### 1. Utiliser les caches Docker

```bash
# Construire avec cache
docker-compose build --parallel

# Pull les images avant de builder
docker-compose pull
```

### 2. Health checks

Les conteneurs ont des health checks configurés :

```yaml
healthcheck:
  test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
  interval: 30s
  timeout: 10s
  retries: 3
```

### 3. Attendre que les services soient prêts

```bash
# Script wait-for-it.sh (déjà inclus dans test-docker.sh)
while ! curl -s http://localhost:8080/api > /dev/null; do
  echo "Waiting for backend..."
  sleep 2
done
```

---

## 📚 Documentation connexe

- **Tests généraux** : `TEST_SYNTHESIS.md`
- **Guide rapide** : `QUICK_START_TESTS.md`
- **Setup Cypress** : `CYPRESS_E2E_SETUP.md`
- **README E2E** : `tests/README.md`
- **Ce guide** : `DOCKER_TESTS_GUIDE.md`

---

## ✅ Checklist avant les tests

- [ ] Docker et Docker Compose installés
- [ ] Services démarrés : `docker-compose up -d`
- [ ] Backend accessible : `curl http://localhost:8080/api`
- [ ] Frontend accessible : `curl http://localhost:3000`
- [ ] PostgreSQL running : `docker-compose ps postgres`
- [ ] Cypress configuré : `tests/.env` existe
- [ ] Dépendances installées : `tests/node_modules/` existe

---

## 🎉 Résumé

### Commande rapide pour tout tester avec Docker

```bash
# 1. Démarrer et configurer
./test-docker.sh

# 2. Lancer les tests E2E
cd tests
npm run cy:open    # ou npm run test

# 3. Tests Backend (si nécessaire)
cd backend
go test ./tests/... -v

# 4. Arrêter
docker-compose down
```

---

**Tous les tests sont maintenant configurés pour votre environnement Docker ! 🐳**

**Date** : Janvier 2026  
**Version** : 1.0.0 (Docker)  
**Statut** : ✅ Adapté pour Docker

