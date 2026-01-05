# SafeBase

Plateforme de gestion sécurisée de sauvegarde et restauration de bases de données MySQL et PostgreSQL.

## 📋 Table des matières

- [Description](#description)
- [Fonctionnalités](#fonctionnalités)
- [Technologies](#technologies)
- [Prérequis](#prérequis)
- [Installation](#installation)
- [Démarrage avec Docker](#démarrage-avec-docker)
- [Configuration](#configuration)
- [Sécurité](#sécurité)
- [Tests](#tests)
- [Structure du projet](#structure-du-projet)
- [API](#api)
- [Dépannage](#dépannage)

## 📝 Description

SafeBase est une plateforme complète permettant de :
- Gérer des configurations de bases de données MySQL et PostgreSQL
- Créer des sauvegardes manuelles ou automatiques
- Planifier des sauvegardes récurrentes via CRON
- Restaurer des bases de données depuis des sauvegardes
- Consulter un historique complet de toutes les opérations

## ✨ Fonctionnalités

### Gestion des bases de données
- ✅ Création, modification, suppression de configurations BDD
- ✅ Support MySQL et PostgreSQL
- ✅ Chiffrement des mots de passe de connexion
- ✅ Validation des configurations

### Sauvegardes
- ✅ Création manuelle de sauvegardes
- ✅ Compression automatique en ZIP
- ✅ Téléchargement des sauvegardes
- ✅ Statuts en temps réel (pending, completed, failed)
- ✅ Organisation par type de base de données

### Planification
- ✅ Planification automatique via expressions CRON
- ✅ Activation/désactivation de planifications
- ✅ Support de multiples planifications par base

### Historique & Audit
- ✅ Traçabilité complète de toutes les actions
- ✅ Filtres par type, ressource, date
- ✅ Export CSV
- ✅ Isolation multi-utilisateurs

### Authentification
- ✅ Inscription et connexion sécurisées
- ✅ Gestion de sessions JWT
- ✅ Cookies HTTP-only pour la sécurité
- ✅ Rôles utilisateur (user, admin)

## 🛠️ Technologies

### Backend
- **Go 1.25+** avec Gin Framework
- **PostgreSQL** pour la base de données
- **JWT** pour l'authentification
- **GORM** pour l'ORM
- **Air** pour le hot-reload

### Frontend
- **Vue.js 3** avec Composition API
- **TypeScript** pour le typage
- **Vite** comme bundler
- **Tailwind CSS** pour le styling
- **Pinia** pour la gestion d'état
- **Vue Router** pour la navigation

### Tests
- **Cypress** pour les tests E2E
- **Go testing** pour les tests unitaires et d'intégration
- **Coverage > 90%**

## 📦 Prérequis

- **Docker** et **Docker Compose** (recommandé)
- Ou :
  - Go 1.25+
  - Node.js 20+
  - PostgreSQL 14+
  - MySQL 8+ (optionnel, pour tester les sauvegardes MySQL)

## 🚀 Installation

### Option 1 : Avec Docker (Recommandé)

```bash
# Cloner le projet
git clone <repository-url>
cd plateforme-safebase

# Démarrer tous les services
docker-compose up -d

# Vérifier que tout fonctionne
docker-compose ps
```

### Option 2 : Installation locale

#### Backend

```bash
cd backend
go mod download
cp .env.example .env  # Configurer les variables d'environnement
go run cmd/main.go
```

#### Frontend

```bash
cd frontend
npm install
cp .env.example .env  # Configurer les variables d'environnement
npm run dev
```

## 🐳 Démarrage avec Docker

### Services disponibles

- **Frontend** : http://localhost:3000
- **Backend API** : http://localhost:8080
- **PostgreSQL** : localhost:5432
- **MySQL** : localhost:3306 (pour tests)
- **Grafana** : http://localhost:3001 (optionnel)
- **Prometheus** : http://localhost:9090 (optionnel)

### Commandes utiles

```bash
# Démarrer tous les services
docker-compose up -d

# Voir les logs
docker-compose logs -f

# Arrêter tous les services
docker-compose down

# Reconstruire les images
docker-compose build

# Accéder au shell d'un conteneur
docker-compose exec backend sh
docker-compose exec frontend sh

# Nettoyer complètement (supprime les volumes)
docker-compose down -v
```

### Hot-reload

Les deux services supportent le hot-reload :
- **Backend** : Air détecte les changements et recompile automatiquement
- **Frontend** : Vite détecte les changements et recharge automatiquement

## ⚙️ Configuration

### Backend (.env)

```env
# JWT (OBLIGATOIRE : changez en production !)
JWT_SECRET=votre_cle_secrete_jwt_tres_longue_et_complexe
GO_ENV=development  # ou production

# Base de données
DB_HOST=postgres     # ou localhost si sans Docker
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=safebase

# Serveur
PORT=8080

# Sauvegardes MEGA (optionnel)
MEGA_EMAIL=votre_email@example.com
MEGA_PASSWORD=votre_mot_de_passe_mega
```

### Frontend (.env)

```env
# URL de l'API
VITE_API_BASE_URL=http://localhost:8080

# Configuration
VITE_APP_NAME=SafeBase
VITE_APP_VERSION=1.0.0
```

## 🔒 Sécurité

### Mesures implémentées

- ✅ **Cookies HTTP-only** : Inaccessibles via JavaScript (protection XSS)
- ✅ **Cookies Secure** : Transmission uniquement via HTTPS en production
- ✅ **Hachage bcrypt** : Mots de passe avec salt automatique
- ✅ **JWT sécurisé** : Signature avec clé secrète, expiration 24h
- ✅ **CORS configuré** : Origines spécifiques, pas de wildcard
- ✅ **Validation des entrées** : Protection contre l'injection SQL
- ✅ **Isolation utilisateurs** : Chaque utilisateur ne voit que ses ressources

### Checklist de déploiement en production

- [ ] HTTPS obligatoire (certificat SSL/TLS valide)
- [ ] `JWT_SECRET` : Clé de 256+ bits générée aléatoirement
- [ ] `GO_ENV=production` : Active les cookies sécurisés
- [ ] CORS : URLs de production configurées
- [ ] Base de données : Credentials forts
- [ ] Firewall : Limiter l'accès aux ports nécessaires
- [ ] Rate limiting : Limiter les tentatives de connexion
- [ ] Monitoring : Logs de sécurité et alertes
- [ ] Backups : Sauvegardes chiffrées régulières

**⚠️ Ne JAMAIS stocker :**
- ❌ Tokens dans localStorage
- ❌ Mots de passe en clair
- ❌ Secrets dans le code

## 🧪 Tests

### Tests E2E avec Cypress

Les tests E2E couvrent **>90%** de l'application avec **~200 tests**.

#### Installation

```bash
cd tests
npm install
```

#### Configuration

Créer un fichier `.env` dans `tests/` :

```env
# Pour Docker
CYPRESS_BASE_URL=http://localhost:3000
CYPRESS_API_URL=http://localhost:8080/api

# Pour développement local
# CYPRESS_BASE_URL=http://localhost:5173
# CYPRESS_API_URL=http://localhost:8080/api
```

#### Exécution

```bash
# Mode interactif (GUI)
npm run cy:open

# Mode headless
npm run test

# Par navigateur
npm run cy:run:chrome
npm run cy:run:firefox

# Un fichier spécifique
npx cypress run --spec "e2E/01-authentication.cy.ts"
```

#### Structure des tests

- `01-authentication.cy.ts` - Authentification (15%)
- `02-database-management.cy.ts` - Gestion BDD (25%)
- `03-backup-management.cy.ts` - Sauvegardes (20%)
- `04-schedule-management.cy.ts` - Planifications (15%)
- `05-history.cy.ts` - Historique (10%)
- `06-profile.cy.ts` - Profil utilisateur (10%)
- `07-dashboard.cy.ts` - Tableau de bord (5%)
- `08-complete-workflows.cy.ts` - Flux complets (10%)

### Tests Backend

```bash
# Tests unitaires
cd backend
go test ./tests/units/... -v

# Tests d'intégration
go test ./tests/integrations/... -v

# Tests fonctionnels
go test ./tests/functionals/... -v
```

## 📁 Structure du projet

```
plateforme-safebase/
├── backend/                 # Application Go
│   ├── cmd/
│   │   └── main.go         # Point d'entrée
│   ├── internal/
│   │   ├── handlers/       # Gestionnaires HTTP
│   │   ├── services/       # Logique métier
│   │   ├── repositories/   # Accès aux données
│   │   ├── models/         # Modèles de données
│   │   ├── routes/         # Définition des routes
│   │   └── middlewares/    # Middlewares (auth, CORS)
│   └── tests/              # Tests backend
│       ├── units/          # Tests unitaires
│       ├── integrations/   # Tests d'intégration
│       └── functionals/    # Tests fonctionnels
│
├── frontend/               # Application Vue.js
│   ├── src/
│   │   ├── views/         # Pages principales
│   │   ├── components/    # Composants réutilisables
│   │   ├── stores/        # État Pinia
│   │   ├── api/           # Clients API
│   │   └── router/        # Routes Vue Router
│   └── dist/              # Build de production
│
├── tests/                  # Tests E2E Cypress
│   ├── e2E/               # Tests end-to-end
│   │   ├── fixtures/      # Données de test
│   │   └── support/       # Configuration Cypress
│   └── cypress.config.ts  # Config Cypress
│
├── docker-compose.yml      # Configuration Docker
└── README.md              # Ce fichier
```

## 🔌 API

### Authentification

- `POST /auth/register` - Inscription
- `POST /auth/login` - Connexion
- `POST /auth/logout` - Déconnexion
- `GET /auth/me` - Informations utilisateur

### Bases de données

- `GET /api/databases` - Liste des BDD de l'utilisateur
- `POST /api/databases` - Créer une configuration BDD
- `GET /api/databases/:id` - Détails d'une BDD
- `PUT /api/databases/:id` - Mettre à jour une BDD
- `DELETE /api/databases/:id` - Supprimer une BDD

### Sauvegardes

- `GET /api/backups` - Liste des sauvegardes
- `POST /api/backups/database/:database_id` - Créer une sauvegarde
- `GET /api/backups/:id` - Détails d'une sauvegarde
- `GET /api/backups/:id/download` - Télécharger une sauvegarde
- `DELETE /api/backups/:id` - Supprimer une sauvegarde

### Planifications

- `GET /api/schedules` - Liste des planifications
- `POST /api/schedules` - Créer une planification
- `PUT /api/schedules/:id` - Mettre à jour une planification
- `DELETE /api/schedules/:id` - Supprimer une planification

### Historique

- `GET /api/history` - Historique des actions
- `GET /api/history/:type/:id` - Historique d'une ressource

### Profil

- `PUT /api/profile` - Mettre à jour le profil
- `PUT /api/profile/password` - Changer le mot de passe

## 🔧 Dépannage

### Le backend ne démarre pas

1. Vérifier que PostgreSQL est démarré : `docker-compose ps`
2. Vérifier les logs : `docker-compose logs backend`
3. Vérifier que le port 8080 n'est pas utilisé

### Le frontend ne démarre pas

1. Vérifier les logs : `docker-compose logs frontend`
2. Vérifier que le port 3000 n'est pas utilisé
3. Réinstaller les dépendances : `docker-compose exec frontend npm install`

### Les tests Cypress échouent

1. Vérifier que les services Docker tournent : `docker-compose ps`
2. Vérifier les URLs dans `tests/.env`
3. Vérifier la connexion : `curl http://localhost:3000`

### Port déjà utilisé

```bash
# Trouver le processus
lsof -i :3000  # Frontend
lsof -i :8080  # Backend

# Arrêter le processus
kill -9 <PID>
```

### Problèmes de permissions npm

```bash
cd tests
rm -rf node_modules package-lock.json
npm cache clean --force
npm install
```

## 📚 Documentation supplémentaire

- `backend/README.md` - Documentation backend détaillée
- `frontend/README.md` - Documentation frontend
- `tests/README.md` - Guide complet des tests E2E

## 🤝 Contribution

1. Fork le projet
2. Créer une branche (`git checkout -b feature/nouvelle-fonctionnalite`)
3. Commit les changements (`git commit -m 'Ajout d'une nouvelle fonctionnalité'`)
4. Push vers la branche (`git push origin feature/nouvelle-fonctionnalite`)
5. Ouvrir une Pull Request

## 📄 Licence

[Spécifier la licence]

---

**Dernière mise à jour** : Janvier 2026  
**Version** : 1.0.0
