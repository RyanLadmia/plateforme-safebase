# 🔐 Guide d'installation et de test du système d'authentification

## 📋 Prérequis

1. **PostgreSQL** installé et en cours d'exécution
2. **Go 1.25+** installé
3. **Node.js et npm** installés
4. **pgAdmin** (optionnel, pour visualiser la base de données)

## 🗄️ Configuration de la base de données PostgreSQL

### 1. Créer la base de données
```sql
-- Connectez-vous à PostgreSQL en tant que superutilisateur
CREATE DATABASE safebase_db;
CREATE USER postgres WITH PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE safebase_db TO postgres;
```

### 2. Créer le fichier .env
Créez un fichier `.env` dans le dossier `backend/` avec le contenu suivant :

```env
# Configuration du serveur
PORT=8080

# Configuration JWT - CHANGEZ CETTE CLÉ EN PRODUCTION !
JWT_SECRET=ma_cle_secrete_jwt_tres_longue_et_complexe_pour_la_production_2024

# Configuration PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=safebase_db

# URL complète de la base de données (alternative)
DB_URL=postgres://postgres:password@localhost:5432/safebase_db?sslmode=disable
```

## 🚀 Démarrage du backend (Go)

```bash
# Naviguez vers le dossier backend
cd backend/

# Installez les dépendances Go (si nécessaire)
go mod tidy

# Démarrez le serveur
go run cmd/main.go
```

Le serveur devrait afficher :
```
🚀 Server running on port 8080
📋 Available endpoints:
   GET  /test            - Test endpoint
   POST /auth/register   - User registration
   POST /auth/login      - User login
   POST /auth/logout     - User logout
```

## 🎨 Démarrage du frontend (Vue.js)

```bash
# Naviguez vers le dossier frontend
cd frontend/

# Installez les dépendances npm
npm install

# Démarrez le serveur de développement
npm run dev
```

## 🧪 Tests du système d'authentification

### 1. Test avec l'interface web
1. Ouvrez votre navigateur sur `http://localhost:5173` (ou le port affiché par Vite)
2. Vous verrez l'interface d'authentification avec deux onglets : "Connexion" et "Inscription"

### 2. Test d'inscription
1. Cliquez sur l'onglet "Inscription"
2. Remplissez le formulaire :
   - **Prénom** : Jean
   - **Nom** : Dupont
   - **Email** : jean.dupont@example.com
   - **Mot de passe** : MonMotDePasse123!
3. Cliquez sur "S'inscrire"
4. Vous devriez voir un message de succès

### 3. Test de connexion
1. Cliquez sur l'onglet "Connexion"
2. Utilisez les identifiants créés :
   - **Email** : jean.dupont@example.com
   - **Mot de passe** : MonMotDePasse123!
3. Cliquez sur "Se connecter"
4. Vous devriez voir l'interface utilisateur connecté avec vos informations

### 4. Test de déconnexion
1. Une fois connecté, cliquez sur "Se déconnecter"
2. Vous devriez revenir aux formulaires de connexion/inscription

## 🔧 Tests avec curl (API directe)

### Test d'inscription
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "firstname": "Marie",
    "lastname": "Martin",
    "email": "marie.martin@example.com",
    "password": "MotDePasseSecure456!"
  }'
```

### Test de connexion
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "marie.martin@example.com",
    "password": "MotDePasseSecure456!"
  }'
```

### Test de déconnexion (remplacez YOUR_JWT_TOKEN par le token reçu)
```bash
curl -X POST http://localhost:8080/auth/logout \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 🗄️ Vérification en base de données

### Avec pgAdmin
1. Connectez-vous à pgAdmin
2. Naviguez vers votre serveur PostgreSQL
3. Ouvrez la base `safebase_db`
4. Vous devriez voir les tables :
   - `roles` (avec admin et user)
   - `users` (avec les utilisateurs créés)
   - `sessions` (avec les sessions actives)

### Avec psql
```sql
-- Connectez-vous à la base
psql -U postgres -d safebase_db

-- Vérifiez les rôles
SELECT * FROM roles;

-- Vérifiez les utilisateurs
SELECT id, firstname, lastname, email, role_id, created_at FROM users;

-- Vérifiez les sessions actives
SELECT id, user_id, expires_at, created_at FROM sessions WHERE expires_at > NOW();
```

## 🔍 Structure du code expliquée

### Backend (Go)
- **`cmd/main.go`** : Point d'entrée, configuration du serveur
- **`internal/models/`** : Structures de données (User, Session, Role)
- **`internal/repositories/`** : Accès aux données (CRUD)
- **`internal/services/`** : Logique métier (validation, hachage)
- **`internal/handlers/`** : Contrôleurs HTTP (endpoints API)
- **`internal/routes/`** : Configuration des routes
- **`pkg/security/`** : Utilitaires JWT et hachage

### Frontend (Vue.js)
- **`src/services/authService.ts`** : Service d'API pour l'authentification
- **`src/components/AuthComponent.vue`** : Interface utilisateur complète
- **`src/views/HomeView.vue`** : Page d'accueil avec le composant auth

## 🛡️ Sécurité implémentée

1. **Hachage des mots de passe** avec bcrypt
2. **Validation des mots de passe** (10+ caractères, complexité)
3. **Tokens JWT** signés avec clé secrète
4. **Sessions en base** pour révocation des tokens
5. **Middleware d'authentification** pour les routes protégées
6. **CORS configuré** pour le frontend
7. **Validation des données** côté serveur

## 🚨 Problèmes courants

### Erreur de connexion à la base
- Vérifiez que PostgreSQL est démarré
- Vérifiez les paramètres dans le fichier `.env`
- Vérifiez que la base `safebase_db` existe

### Erreur CORS
- Vérifiez que le backend est démarré sur le port 8080
- Le middleware CORS est configuré dans `main.go`

### Token invalide
- Les tokens expirent après 24h
- Vérifiez que la session existe en base de données

Votre système d'authentification est maintenant prêt ! 🎉
