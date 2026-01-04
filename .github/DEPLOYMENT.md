# Guide de Déploiement CI/CD

Ce document explique comment configurer le pipeline CI/CD pour SafeBase.

## 📋 Prérequis

1. Un compte Docker Hub
2. Un serveur de production avec Docker et Docker Compose installés
3. Accès SSH au serveur de production
4. Un repository GitHub

## 🔐 Configuration des Secrets GitHub

Allez dans `Settings > Secrets and variables > Actions` de votre repository GitHub et ajoutez les secrets suivants :

### Secrets Docker Hub (obligatoires)

```
DOCKER_USERNAME     # Votre nom d'utilisateur Docker Hub
DOCKER_PASSWORD     # Votre mot de passe ou token Docker Hub
```

### Secrets de déploiement Production (obligatoires pour le déploiement)

```
DEPLOY_HOST         # Adresse IP ou domaine du serveur (ex: 123.45.67.89)
DEPLOY_USER         # Utilisateur SSH (ex: root ou ubuntu)
DEPLOY_SSH_KEY      # Clé privée SSH pour se connecter au serveur
DEPLOY_PORT         # Port SSH (optionnel, défaut: 22)
DEPLOY_PATH         # Chemin vers l'application sur le serveur (ex: /opt/safebase)
PRODUCTION_URL      # URL de production pour le health check (ex: https://api.votre-domaine.com)
```

### Secrets de déploiement Staging (optionnels)

```
STAGING_HOST        # Adresse du serveur staging
STAGING_USER        # Utilisateur SSH staging
STAGING_SSH_KEY     # Clé privée SSH staging
STAGING_PORT        # Port SSH staging (optionnel)
STAGING_PATH        # Chemin vers l'application staging
```

## 🔑 Génération de la clé SSH

Sur votre machine locale :

```bash
# Générer une nouvelle paire de clés SSH
ssh-keygen -t ed25519 -C "github-actions-deploy" -f ~/.ssh/github_deploy

# Copier la clé publique sur le serveur
ssh-copy-id -i ~/.ssh/github_deploy.pub user@votre-serveur.com

# Afficher la clé privée (à copier dans DEPLOY_SSH_KEY)
cat ~/.ssh/github_deploy
```

⚠️ **Important** : Copiez TOUTE la clé privée, y compris les lignes `-----BEGIN` et `-----END`.

## 🚀 Workflow de déploiement

### Branche `develop` → Staging
1. Push sur la branche `develop`
2. Tests automatiques
3. Build des images Docker
4. Déploiement automatique sur staging

### Branche `main` → Production
1. Push sur la branche `main`
2. Tests automatiques
3. Build des images Docker
4. Déploiement automatique sur production

### Pull Requests
- Exécute uniquement les tests
- Pas de build ni de déploiement

## 📦 Préparation du serveur de production

### 1. Installer Docker et Docker Compose

```bash
# Mettre à jour le système
sudo apt update && sudo apt upgrade -y

# Installer Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Ajouter l'utilisateur au groupe docker
sudo usermod -aG docker $USER

# Installer Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### 2. Créer le répertoire de l'application

```bash
# Créer le répertoire
sudo mkdir -p /opt/safebase
sudo chown $USER:$USER /opt/safebase
cd /opt/safebase

# Cloner le repository (première fois)
git clone https://github.com/votre-username/plateforme-safebase.git .
```

### 3. Configurer les variables d'environnement

```bash
# Créer le fichier .env
nano .env
```

Ajoutez vos variables d'environnement :

```env
# Base de données PostgreSQL
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=safebase_user
POSTGRES_PASSWORD=votre_mot_de_passe_securise
POSTGRES_DB=safebase

# Base de données MySQL
MYSQL_HOST=mysql
MYSQL_PORT=3306
MYSQL_ROOT_PASSWORD=votre_mot_de_passe_root
MYSQL_DATABASE=safebase
MYSQL_USER=safebase_user
MYSQL_PASSWORD=votre_mot_de_passe_securise

# Backend
JWT_SECRET=votre_secret_jwt_tres_long_et_securise
MEGA_EMAIL=votre_email@mega.nz
MEGA_PASSWORD=votre_mot_de_passe_mega
ENCRYPTION_KEY=votre_cle_de_chiffrement_32_caracteres

# Frontend
VITE_API_URL=https://api.votre-domaine.com
```

### 4. Créer le fichier docker-compose.yml pour production

```bash
nano docker-compose.yml
```

```yaml
version: '3.8'

services:
  backend:
    image: ${DOCKER_USERNAME}/safebase-backend:latest
    container_name: safebase-backend
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      - POSTGRES_HOST=${POSTGRES_HOST}
      - POSTGRES_PORT=${POSTGRES_PORT}
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
      - POSTGRES_DB=${POSTGRES_DB}
      - JWT_SECRET=${JWT_SECRET}
      - MEGA_EMAIL=${MEGA_EMAIL}
      - MEGA_PASSWORD=${MEGA_PASSWORD}
      - ENCRYPTION_KEY=${ENCRYPTION_KEY}
    depends_on:
      - postgres
      - mysql
    networks:
      - safebase-network

  frontend:
    image: ${DOCKER_USERNAME}/safebase-frontend:latest
    container_name: safebase-frontend
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    environment:
      - VITE_API_URL=${VITE_API_URL}
    depends_on:
      - backend
    networks:
      - safebase-network

  postgres:
    image: postgres:16-alpine
    container_name: safebase-postgres
    restart: unless-stopped
    environment:
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
      - POSTGRES_DB=${POSTGRES_DB}
    volumes:
      - postgres-data:/var/lib/postgresql/data
    networks:
      - safebase-network

  mysql:
    image: mysql:8.0
    container_name: safebase-mysql
    restart: unless-stopped
    environment:
      - MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}
      - MYSQL_DATABASE=${MYSQL_DATABASE}
      - MYSQL_USER=${MYSQL_USER}
      - MYSQL_PASSWORD=${MYSQL_PASSWORD}
    volumes:
      - mysql-data:/var/lib/mysql
    networks:
      - safebase-network

networks:
  safebase-network:
    driver: bridge

volumes:
  postgres-data:
  mysql-data:
```

### 5. Premier déploiement manuel

```bash
# Démarrer les services
docker-compose up -d

# Vérifier les logs
docker-compose logs -f

# Vérifier que tout fonctionne
docker-compose ps
```

## 🔍 Vérification et monitoring

### Vérifier les logs

```bash
# Logs de tous les services
docker-compose logs -f

# Logs du backend uniquement
docker-compose logs -f backend

# Logs du frontend uniquement
docker-compose logs -f frontend
```

### Vérifier l'état des conteneurs

```bash
docker-compose ps
```

### Redémarrer un service

```bash
docker-compose restart backend
docker-compose restart frontend
```

## 🔄 Rollback en cas de problème

Si un déploiement échoue, vous pouvez revenir à la version précédente :

```bash
# Voir les images disponibles
docker images

# Revenir à une version spécifique
docker-compose down
docker tag votre-username/safebase-backend:sha-abc123 votre-username/safebase-backend:latest
docker-compose up -d
```

## 🛡️ Sécurité

### Pare-feu

```bash
# Autoriser SSH
sudo ufw allow 22/tcp

# Autoriser HTTP et HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Activer le pare-feu
sudo ufw enable
```

### SSL/TLS avec Let's Encrypt (recommandé)

Utilisez Nginx Proxy Manager ou Traefik pour gérer automatiquement les certificats SSL.

## 📊 Monitoring (optionnel)

Le projet inclut déjà Prometheus et Grafana. Pour y accéder :

- Prometheus: `http://votre-serveur:9090`
- Grafana: `http://votre-serveur:3001`

## 🆘 Dépannage

### Les conteneurs ne démarrent pas

```bash
# Vérifier les logs
docker-compose logs

# Vérifier l'espace disque
df -h

# Nettoyer les anciennes images
docker system prune -a
```

### Problème de connexion à la base de données

```bash
# Vérifier que les conteneurs sont sur le même réseau
docker network ls
docker network inspect safebase-network

# Tester la connexion
docker-compose exec backend ping postgres
```

### Les images ne se téléchargent pas

```bash
# Vérifier la connexion à Docker Hub
docker login

# Pull manuel
docker pull votre-username/safebase-backend:latest
```

## 📝 Notes

- Le pipeline s'exécute automatiquement à chaque push sur `main` ou `develop`
- Les images Docker sont multi-architecture (amd64 et arm64)
- Les anciennes images sont automatiquement nettoyées après déploiement
- Un backup de `.env` est créé avant chaque déploiement

## 🔗 Ressources

- [Documentation Docker](https://docs.docker.com/)
- [Documentation GitHub Actions](https://docs.github.com/en/actions)
- [Documentation Docker Compose](https://docs.docker.com/compose/)

