# ✅ MISE À JOUR : Configuration Docker pour les tests

## 🐳 Ce qui a été adapté pour Docker

Tous les tests E2E Cypress ont été **adaptés pour fonctionner avec Docker** !

### Fichiers modifiés/créés

1. ✅ **`tests/cypress.config.ts`** - Configuration adaptée pour Docker et local
   - Détection automatique de l'environnement
   - Support des variables d'environnement
   - Timeouts augmentés pour Docker

2. ✅ **`tests/package.json`** - Nouveaux scripts ajoutés
   - `npm run test:docker` - Tests avec Docker
   - `npm run test:local` - Tests sans Docker
   - `npm run cy:open:docker` - GUI avec Docker
   - `npm run cy:open:local` - GUI sans Docker

3. ✅ **`tests/env.docker.example`** - Configuration Docker
   - Frontend port 3000
   - Backend port 8080

4. ✅ **`tests/env.local.example`** - Configuration locale
   - Frontend port 5173 (Vite)
   - Backend port 8080

5. ✅ **`test-docker.sh`** - Script automatique
   - Démarre Docker Compose
   - Configure Cypress
   - Vérifie que tout est prêt

6. ✅ **`DOCKER_TESTS_GUIDE.md`** - Guide complet Docker

7. ✅ **`tests/README.md`** - Mis à jour avec infos Docker

---

## 🎯 Différences clés

### Ports

| Service | Local (sans Docker) | **Docker** |
|---------|---------------------|-----------|
| Frontend | 5173 | **3000** |
| Backend | 8080 | 8080 |
| PostgreSQL | 5432 | 5432 |

### Configuration Cypress

**Docker (par défaut) :**
```bash
CYPRESS_BASE_URL=http://localhost:3000
CYPRESS_API_URL=http://localhost:8080/api
```

**Local :**
```bash
CYPRESS_BASE_URL=http://localhost:5173
CYPRESS_API_URL=http://localhost:8080/api
```

---

## 🚀 Utilisation rapide avec Docker

### Méthode 1 : Script automatique (RECOMMANDÉ)

```bash
# À la racine du projet
./test-docker.sh

# Puis lancer les tests
cd tests
npm run cy:open    # Mode GUI
# ou
npm run test       # Mode headless
```

### Méthode 2 : Manuelle

```bash
# 1. Démarrer Docker
docker-compose up -d

# 2. Attendre 1-2 minutes que tout démarre

# 3. Configurer Cypress
cd tests
cat > .env << EOF
CYPRESS_BASE_URL=http://localhost:3000
CYPRESS_API_URL=http://localhost:8080/api
CYPRESS_IS_DOCKER=true
EOF

# 4. Installer dépendances (si pas déjà fait)
npm install

# 5. Lancer les tests
npm run cy:open    # Mode GUI
npm run test       # Mode headless
```

---

## 📊 Tests Backend Go

**Les tests Backend Go n'ont PAS besoin de Docker** car ils utilisent :
- SQLite en mémoire (tests unitaires)
- PostgreSQL via `.env` (tests intégration)

### Exécution

```bash
cd backend

# Tests unitaires (pas besoin de Docker)
go test ./tests/units/... -v

# Tests d'intégration (besoin PostgreSQL Docker)
docker-compose up -d postgres
go test ./tests/integrations/... -v

# Tests fonctionnels
go test ./tests/functionals/... -v
```

---

## 🔧 Vérification rapide

### Docker est-il prêt ?

```bash
# Vérifier les conteneurs
docker-compose ps

# Tester le backend
curl http://localhost:8080/api

# Tester le frontend
curl http://localhost:3000

# Tout est OK ? Lancer les tests !
cd tests && npm run test
```

---

## 📝 Scripts disponibles

### Docker (recommandé)

```bash
cd tests

# Mode GUI
npm run cy:open:docker
npm run cy:open          # utilise .env

# Mode headless
npm run test:docker
npm run test             # utilise .env
```

### Local (sans Docker)

```bash
cd tests

# Mode GUI
npm run cy:open:local

# Mode headless
npm run test:local
```

### Autres

```bash
# CI/CD
npm run test:ci

# Par navigateur
npm run cy:run:chrome
npm run cy:run:firefox
npm run cy:run:edge
```

---

## 🐛 Problèmes courants

### "Connection refused" ou "timeout"

```bash
# Vérifier que Docker tourne
docker-compose ps

# Redémarrer si nécessaire
docker-compose restart

# Voir les logs
docker-compose logs -f backend
docker-compose logs -f frontend
```

### Frontend sur mauvais port

```bash
# Vérifier .env dans tests/
cat tests/.env

# Devrait avoir pour Docker:
CYPRESS_BASE_URL=http://localhost:3000

# Ou pour local:
CYPRESS_BASE_URL=http://localhost:5173
```

### Tests échouent immédiatement

```bash
# Attendre que Docker soit complètement démarré
# Compter 1-2 minutes après docker-compose up

# Vérifier manuellement
curl http://localhost:8080/api   # Backend doit répondre
curl http://localhost:3000       # Frontend doit répondre
```

---

## 📚 Documentation

- **Guide Docker complet** : `DOCKER_TESTS_GUIDE.md` ⭐
- **Tests généraux** : `TEST_SYNTHESIS.md`
- **Guide rapide** : `QUICK_START_TESTS.md`
- **Setup Cypress** : `CYPRESS_E2E_SETUP.md`
- **README E2E** : `tests/README.md`

---

## ✅ Récapitulatif

### Pour Docker (votre cas) :

1. **Démarrer** : `./test-docker.sh` OU `docker-compose up -d`
2. **Attendre** : 1-2 minutes
3. **Tester** : `cd tests && npm run cy:open`

### Pour local (sans Docker) :

1. **Backend** : `cd backend && go run cmd/main.go`
2. **Frontend** : `cd frontend && npm run dev`
3. **Tester** : `cd tests && npm run cy:open:local`

---

## 🎉 Conclusion

**Tous vos tests fonctionnent maintenant avec Docker ! 🐳**

- ✅ Configuration Cypress adaptée
- ✅ Scripts npm pour Docker et local
- ✅ Script automatique `test-docker.sh`
- ✅ Documentation complète
- ✅ Support des deux modes (Docker + Local)

**Commande rapide pour tout tester :**

```bash
./test-docker.sh && cd tests && npm run test
```

---

**Date** : Janvier 2026  
**Version** : 1.0.1 (Docker Support)  
**Statut** : ✅ Prêt pour Docker

