# 🐳 SafeBase - Configuration Docker Complète

## ✅ Tous les problèmes ont été corrigés !

### 🎯 Résumé des corrections

| Problème | Solution |
|----------|----------|
| Backend : `exec format error` | Utilisation d'Air (hot-reload) au lieu de binaire compilé |
| Backend : Air nécessite Go 1.25 | Mise à jour vers Go 1.25 dans Dockerfile et go.mod |
| Backend : `pg_dump non trouvé` | Installation de postgresql-client et mysql-client |
| Frontend : Node.js trop ancien | Mise à jour vers Node.js 20 |
| Frontend : `crypto.hash is not a function` | Configuration Vite avec polling pour Docker |

---

## 🚀 Démarrage rapide (3 étapes)

### 1️⃣ Vérifier la configuration
```bash
./docker-check.sh
```

### 2️⃣ Démarrer le projet
```bash
./docker-start.sh
```

### 3️⃣ Accéder aux services
- **Frontend** : http://localhost:3000
- **Backend** : http://localhost:8080/test
- **Grafana** : http://localhost:3001 (admin/admin)
- **Prometheus** : http://localhost:9090

---

## 📋 Commandes essentielles

### Démarrage
```bash
# Démarrage automatique (recommandé)
./docker-start.sh

# Démarrage manuel
docker-compose up

# Démarrage en arrière-plan
docker-compose up -d
```

### Logs
```bash
# Tous les services
docker-compose logs -f

# Backend uniquement
docker-compose logs -f backend

# Frontend uniquement
docker-compose logs -f frontend
```

### Arrêt
```bash
# Arrêter les conteneurs
docker-compose down

# Arrêter et supprimer les volumes (⚠️ supprime les données)
docker-compose down -v
```

### Reconstruction
```bash
# Reconstruire tout
docker-compose build

# Reconstruire un service spécifique
docker-compose build backend
docker-compose build frontend
```

---

## 🔥 Hot-reload activé

Les modifications de code sont détectées automatiquement :

### Backend (Go + Air)
- Modifiez un fichier `.go`
- Air recompile et redémarre automatiquement
- Pas besoin de redémarrer le conteneur

### Frontend (Vue + Vite)
- Modifiez un fichier `.vue`, `.ts`, `.css`
- Vite recharge instantanément le navigateur
- Hot Module Replacement (HMR) activé

---

## 📊 Architecture des services

```
┌─────────────────────────────────────────────────────────┐
│                    Docker Compose                        │
├─────────────────────────────────────────────────────────┤
│                                                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │   Frontend   │  │   Backend    │  │  PostgreSQL  │  │
│  │   Node 20    │→ │   Go 1.25    │→ │     v15      │  │
│  │   Vite       │  │   Air        │  │              │  │
│  │   :3000      │  │   :8080      │  │   :5432      │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
│                                                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │    MySQL     │  │   Grafana    │  │  Prometheus  │  │
│  │     v8       │  │   Latest     │  │   Latest     │  │
│  │   :3306      │  │   :3001      │  │   :9090      │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
│                                                           │
└─────────────────────────────────────────────────────────┘
```

---

## 🛠️ Dépannage

### ❌ Erreur : "port is already allocated"
```bash
# Trouver le processus utilisant le port
lsof -i :8080  # ou :3000

# Arrêter le processus
kill -9 <PID>
```

### ❌ Backend ne démarre pas
```bash
# Vérifier les logs
docker-compose logs backend

# Reconstruire l'image
docker-compose build backend
docker-compose up backend
```

### ❌ Frontend ne démarre pas
```bash
# Vérifier les logs
docker-compose logs frontend

# Supprimer node_modules et reconstruire
rm -rf frontend/node_modules
docker-compose build frontend
docker-compose up frontend
```

### ❌ Base de données non accessible
```bash
# Attendre que PostgreSQL soit prêt (10-15 secondes)
docker-compose logs postgres

# Tester la connexion
docker-compose exec postgres psql -U user -d safebase -c "SELECT 1;"
```

### ❌ Hot-reload ne fonctionne pas
Le polling est déjà activé dans la configuration. Si le problème persiste :
```bash
# Redémarrer le service
docker-compose restart backend
# ou
docker-compose restart frontend
```

---

## 📁 Structure des fichiers Docker

```
plateforme-safebase/
├── docker-compose.yml              # Configuration des services
├── docker-start.sh                 # Script de démarrage automatique
├── docker-check.sh                 # Script de vérification
├── docker-test-backup-tools.sh    # Script de test des outils de sauvegarde
├── DOCKER_FIXES.md                 # Documentation détaillée
├── DOCKER_SETUP.md                 # Guide de configuration
├── DOCKER_BACKUP_MEGA.md           # Configuration sauvegardes MEGA
├── LIRE_MOI_DOCKER.md              # Ce fichier
│
├── backend/
│   ├── Dockerfile                  # Image Go 1.25 + Air + outils DB
│   ├── .dockerignore               # Fichiers à exclure
│   ├── .air.toml                   # Configuration Air (hot-reload)
│   ├── go.mod                      # Dépendances Go 1.25
│   └── db/backups/                 # Sauvegardes locales
│       ├── postgresql/
│       └── mysql/
│
└── frontend/
    ├── Dockerfile                  # Image Node 20
    ├── .dockerignore               # Fichiers à exclure
    └── vite.config.ts              # Configuration Vite (polling activé)
```

---

## 🔐 Variables d'environnement

Les variables sont définies dans `docker-compose.yml` :

```yaml
Backend:
  - PORT=8080
  - JWT_SECRET=your_jwt_secret_key_change_in_production
  - DB_HOST=postgres
  - DB_PORT=5432
  - DB_USER=user
  - DB_PASSWORD=password
  - DB_NAME=safebase
```

⚠️ **IMPORTANT** : Changez ces valeurs en production !

---

## 📚 Documentation complète

- **[DOCKER_FIXES.md](./DOCKER_FIXES.md)** - Détails techniques des corrections
- **[DOCKER_SETUP.md](./DOCKER_SETUP.md)** - Guide de configuration complet
- **[DOCKER_BACKUP_MEGA.md](./DOCKER_BACKUP_MEGA.md)** - Configuration des sauvegardes avec MEGA
- **[Backend README](./backend/README.md)** - Documentation backend
- **[Frontend README](./frontend/README.md)** - Documentation frontend

---

## ✨ Versions utilisées

| Technologie | Version | Raison |
|-------------|---------|--------|
| Go | 1.25 | Requis par Air v1.63+ |
| Node.js | 20 | Requis par Vite 7+ |
| PostgreSQL | 15 | Stable et performant |
| MySQL | 8 | Dernière version stable |
| Air | v1.63.5+ | Hot-reload Go |
| Vite | 7.0.6 | Build tool moderne |

---

## 🗄️ Sauvegardes avec MEGA

### Configuration rapide

1. Ajoutez vos identifiants MEGA dans `docker-compose.yml` :
```yaml
backend:
  environment:
    - MEGA_EMAIL=votre_email@example.com
    - MEGA_PASSWORD=votre_mot_de_passe_mega
```

2. Reconstruisez le backend :
```bash
docker-compose build backend
docker-compose up -d backend
```

3. Testez les outils de sauvegarde :
```bash
./docker-test-backup-tools.sh
```

**Documentation complète** : [DOCKER_BACKUP_MEGA.md](./DOCKER_BACKUP_MEGA.md)

---

## 🎉 Tout est prêt !

Votre projet est maintenant complètement dockerisé avec :
- ✅ Hot-reload backend (Air)
- ✅ Hot-reload frontend (Vite)
- ✅ Outils de sauvegarde (pg_dump, mysqldump)
- ✅ Support MEGA pour sauvegardes cloud
- ✅ Toutes les dépendances configurées
- ✅ Scripts de démarrage automatiques
- ✅ Documentation complète

**Lancez simplement** : `./docker-start.sh` 🚀

