# 🐳 Corrections Docker - Résumé Complet

## ✅ Tous les problèmes ont été résolus !

### 🔧 Problème 1 : Backend - "exec format error"
**Erreur initiale** :
```
exec ./main: exec format error
```

**Cause** : Le binaire était compilé pour une architecture différente (problème de cross-compilation)

**Solution appliquée** :
- ✅ Utilisation d'**Air** (hot-reload) au lieu de compiler un binaire statique
- ✅ Configuration du polling dans `.air.toml` pour Docker sur Mac/Windows
- ✅ Mise à jour vers **Go 1.25** pour supporter Air v1.63+

---

### 🔧 Problème 2 : Backend - Air nécessite Go 1.25
**Erreur** :
```
go: github.com/air-verse/air@latest: github.com/air-verse/air@v1.63.5 requires go >= 1.25 (running go 1.24.11)
```

**Solution appliquée** :
- ✅ Mise à jour du `Dockerfile` : `FROM golang:1.25-alpine`
- ✅ Mise à jour du `go.mod` : `go 1.25.0`

---

### 🔧 Problème 3 : Frontend - Version Node.js
**Erreur** :
```
You are using Node.js 18.20.8. Vite requires Node.js version 20.19+ or 22.12+
TypeError: crypto.hash is not a function
```

**Solution appliquée** :
- ✅ Mise à jour du `Dockerfile` : `FROM node:20-alpine`
- ✅ Configuration de Vite avec `host: '0.0.0.0'` et `usePolling: true`

---

## 📝 Fichiers modifiés

### Backend
1. **`backend/Dockerfile`**
   - Image : `golang:1.25-alpine` (au lieu de 1.24)
   - Installation d'Air pour le hot-reload
   - Installation des outils de sauvegarde : `postgresql-client`, `mysql-client`, `zip`, `unzip`
   - Commande : `air -c .air.toml`

2. **`backend/go.mod`**
   - Version Go : `1.25.0` (au lieu de 1.24.0)

3. **`backend/.air.toml`**
   - `tmp_dir = "tmp"` (au lieu de "cmd/tmp")
   - `bin = "./tmp/main"`
   - `cmd = "go build -o ./tmp/main ./cmd"`
   - `poll = true` (nécessaire pour Docker)
   - `poll_interval = 500`

4. **`backend/.dockerignore`**
   - Exclusion des fichiers temporaires
   - **Important** : `.air.toml` n'est PAS ignoré (nécessaire pour Air)

### Frontend
1. **`frontend/Dockerfile`**
   - Image : `node:20-alpine` (au lieu de node:18)
   - Commande : `npm run dev`

2. **`frontend/vite.config.ts`**
   - Configuration serveur :
     ```typescript
     server: {
       host: '0.0.0.0',
       port: 3000,
       watch: {
         usePolling: true, // Pour Docker sur Mac/Windows
       },
     }
     ```

3. **`frontend/.dockerignore`**
   - Exclusion de `node_modules/` et `dist/`

### Docker Compose
1. **`docker-compose.yml`**
   - Variables d'environnement complètes pour le backend :
     ```yaml
     - PORT=8080
     - JWT_SECRET=your_jwt_secret_key_change_in_production
     - DB_HOST=postgres
     - DB_PORT=5432
     - DB_USER=user
     - DB_PASSWORD=password
     - DB_NAME=safebase
     ```
   - Volume `/app/tmp` pour les binaires temporaires d'Air

---

## 🚀 Comment lancer le projet

### Option 1 : Script automatique (recommandé)
```bash
./docker-start.sh
```

### Option 2 : Commandes manuelles
```bash
# Arrêter les conteneurs existants
docker-compose down

# Reconstruire les images (IMPORTANT après les changements)
docker-compose build

# Démarrer tous les services
docker-compose up
```

### Option 3 : Mode détaché
```bash
docker-compose up -d
```

---

## 🌐 Services disponibles

| Service | URL | Description |
|---------|-----|-------------|
| Frontend | http://localhost:3000 | Interface Vue.js |
| Backend | http://localhost:8080 | API Go (Gin) |
| PostgreSQL | localhost:5432 | Base de données principale |
| MySQL | localhost:3306 | Base de données secondaire |
| Grafana | http://localhost:3001 | Monitoring |
| Prometheus | http://localhost:9090 | Métriques |

---

## 🔥 Hot-reload activé

Les deux services supportent le hot-reload en temps réel :

### Backend (Air)
- Détecte automatiquement les changements dans les fichiers `.go`
- Recompile et redémarre le serveur automatiquement
- Logs visibles avec `docker-compose logs -f backend`

### Frontend (Vite)
- Détecte automatiquement les changements dans `.vue`, `.ts`, `.css`, etc.
- Recharge instantanément le navigateur (HMR)
- Logs visibles avec `docker-compose logs -f frontend`

---

## 🛠️ Commandes utiles

### Voir les logs
```bash
# Tous les services
docker-compose logs -f

# Backend uniquement
docker-compose logs -f backend

# Frontend uniquement
docker-compose logs -f frontend

# PostgreSQL
docker-compose logs -f postgres
```

### Redémarrer un service
```bash
docker-compose restart backend
docker-compose restart frontend
```

### Reconstruire un service spécifique
```bash
docker-compose build backend
docker-compose build frontend
```

### Accéder au shell d'un conteneur
```bash
# Backend
docker-compose exec backend sh

# Frontend
docker-compose exec frontend sh

# PostgreSQL
docker-compose exec postgres psql -U user -d safebase
```

### Nettoyer complètement
```bash
# Arrêter et supprimer les conteneurs
docker-compose down

# Supprimer aussi les volumes (⚠️ supprime les données)
docker-compose down -v

# Supprimer les images
docker-compose down --rmi all
```

---

## 🐛 Résolution de problèmes

### Le backend ne démarre pas
1. Vérifiez que PostgreSQL est démarré :
   ```bash
   docker-compose ps
   ```

2. Vérifiez les logs :
   ```bash
   docker-compose logs backend
   ```

3. Vérifiez que le port 8080 n'est pas utilisé :
   ```bash
   lsof -i :8080
   ```

4. Reconstruisez l'image :
   ```bash
   docker-compose build backend
   docker-compose up backend
   ```

### Le frontend ne démarre pas
1. Vérifiez les logs :
   ```bash
   docker-compose logs frontend
   ```

2. Vérifiez que le port 3000 n'est pas utilisé :
   ```bash
   lsof -i :3000
   ```

3. Supprimez `node_modules` et reconstruisez :
   ```bash
   rm -rf frontend/node_modules
   docker-compose build frontend
   docker-compose up frontend
   ```

### Les changements ne sont pas détectés
1. Le polling est déjà activé dans la configuration
2. Redémarrez le service :
   ```bash
   docker-compose restart backend
   # ou
   docker-compose restart frontend
   ```

### Erreur "Cannot connect to database"
1. Attendez que PostgreSQL soit complètement démarré (peut prendre 10-15 secondes)
2. Vérifiez que PostgreSQL est en cours d'exécution :
   ```bash
   docker-compose ps postgres
   ```

3. Testez la connexion :
   ```bash
   docker-compose exec postgres psql -U user -d safebase -c "SELECT 1;"
   ```

---

## 📊 Vérification de l'installation

### Test du backend
```bash
curl http://localhost:8080/test
```

Réponse attendue :
```json
{"message":"Safebase API is running!"}
```

### Test du frontend
Ouvrez http://localhost:3000 dans votre navigateur.

---

## 🔐 Sécurité

⚠️ **IMPORTANT** : Avant de déployer en production :

1. Changez `JWT_SECRET` dans `docker-compose.yml`
2. Changez les mots de passe PostgreSQL et MySQL
3. Utilisez des variables d'environnement ou un fichier `.env`
4. N'exposez pas les ports de base de données publiquement

---

## 📚 Documentation

- [Docker Setup Guide](./DOCKER_SETUP.md)
- [Backend README](./backend/README.md)
- [Frontend README](./frontend/README.md)

---

## ✨ Résumé des versions

| Technologie | Version |
|-------------|---------|
| Go | 1.25 |
| Node.js | 20 |
| PostgreSQL | 15 |
| MySQL | 8 |
| Air | v1.63.5+ |
| Vite | 7.0.6 |

---

**Tout est maintenant configuré et prêt à fonctionner ! 🎉**

