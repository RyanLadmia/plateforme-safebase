# 🚀 CI/CD SafeBase - Guide Rapide

Ce guide vous explique comment mettre en place et utiliser le pipeline CI/CD pour SafeBase.

## 📋 Vue d'ensemble

Le pipeline CI/CD automatise :
- ✅ Tests automatiques du code
- ✅ Build des images Docker
- ✅ Push vers Docker Hub
- ✅ Déploiement automatique sur les serveurs

## 🎯 Workflow

```
┌─────────────┐
│  Git Push   │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Tests     │ ← Lint, Unit tests, Build
└──────┬──────┘
       │
       ▼
┌─────────────┐
│Docker Build │ ← Build images multi-arch
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Docker Hub  │ ← Push images
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Deploy    │ ← Déploiement automatique
└─────────────┘
```

## 🔧 Configuration initiale

### 1. Créer un compte Docker Hub

1. Allez sur [hub.docker.com](https://hub.docker.com)
2. Créez un compte gratuit
3. Créez un Access Token :
   - Account Settings → Security → New Access Token
   - Nom : `github-actions`
   - Permissions : Read, Write, Delete
   - Copiez le token (vous ne le reverrez plus !)

### 2. Configurer les secrets GitHub

Allez dans votre repository GitHub :
```
Settings → Secrets and variables → Actions → New repository secret
```

Ajoutez ces secrets :

#### Obligatoires pour le build :
- `DOCKER_USERNAME` : Votre nom d'utilisateur Docker Hub
- `DOCKER_PASSWORD` : Votre Access Token Docker Hub

#### Obligatoires pour le déploiement :
- `DEPLOY_HOST` : IP ou domaine de votre serveur (ex: `123.45.67.89`)
- `DEPLOY_USER` : Utilisateur SSH (ex: `ubuntu`)
- `DEPLOY_SSH_KEY` : Clé privée SSH (voir ci-dessous)
- `DEPLOY_PATH` : Chemin sur le serveur (ex: `/opt/safebase`)
- `PRODUCTION_URL` : URL de votre API (ex: `https://api.example.com`)

### 3. Générer une clé SSH pour le déploiement

Sur votre machine locale :

```bash
# Générer la clé
ssh-keygen -t ed25519 -C "github-deploy" -f ~/.ssh/github_deploy

# Copier la clé publique sur le serveur
ssh-copy-id -i ~/.ssh/github_deploy.pub user@votre-serveur.com

# Afficher la clé privée (à copier dans DEPLOY_SSH_KEY)
cat ~/.ssh/github_deploy
```

⚠️ **Important** : Copiez TOUTE la clé, y compris les lignes `-----BEGIN` et `-----END`.

### 4. Préparer le serveur de production

Connectez-vous à votre serveur et exécutez :

```bash
# Installer Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Installer Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Créer le répertoire de l'application
sudo mkdir -p /opt/safebase
sudo chown $USER:$USER /opt/safebase
cd /opt/safebase

# Créer le fichier .env (voir env.production.example)
nano .env

# Créer le docker-compose.yml
# Copiez le contenu de docker-compose.prod.yml
nano docker-compose.yml
```

## 🚦 Utilisation

### Déploiement automatique

Le déploiement se fait automatiquement selon les branches :

#### Branch `main` → Production
```bash
git add .
git commit -m "feat: nouvelle fonctionnalité"
git push origin main
```
→ Tests → Build → Push → **Déploiement Production**

#### Branch `develop` → Staging
```bash
git checkout develop
git add .
git commit -m "feat: test nouvelle fonctionnalité"
git push origin develop
```
→ Tests → Build → Push → **Déploiement Staging**

#### Pull Request
```bash
git checkout -b feature/ma-feature
git add .
git commit -m "feat: ma feature"
git push origin feature/ma-feature
# Créer une PR sur GitHub
```
→ **Tests uniquement** (pas de déploiement)

### Déploiement manuel local

Utilisez le script `deploy.sh` :

```bash
# Rendre le script exécutable (première fois)
chmod +x deploy.sh

# Lancer le script
./deploy.sh
```

Menu disponible :
1. Build et déployer (développement)
2. Déployer depuis Docker Hub (production)
3. Arrêter tous les services
4. Voir les logs
5. Nettoyer les anciennes images
6. Backup des bases de données
7. Restaurer les bases de données

## 📊 Monitoring du pipeline

### Voir l'état du pipeline

1. Allez sur votre repository GitHub
2. Cliquez sur l'onglet **Actions**
3. Vous verrez tous les workflows en cours et terminés

### Logs détaillés

Cliquez sur un workflow pour voir :
- ✅ Étapes réussies (vert)
- ❌ Étapes échouées (rouge)
- 📝 Logs détaillés de chaque étape

## 🔍 Vérification post-déploiement

### Sur le serveur

```bash
# SSH vers le serveur
ssh user@votre-serveur.com

# Vérifier les conteneurs
cd /opt/safebase
docker-compose ps

# Voir les logs
docker-compose logs -f backend
docker-compose logs -f frontend

# Vérifier la santé
docker-compose ps
# Tous les services doivent être "Up (healthy)"
```

### Depuis votre navigateur

- Backend API : `https://api.votre-domaine.com/health`
- Frontend : `https://votre-domaine.com`
- Grafana : `https://votre-domaine.com:3001`

## 🛠️ Commandes utiles

### Sur le serveur de production

```bash
# Redémarrer un service
docker-compose restart backend

# Voir les logs en temps réel
docker-compose logs -f

# Mettre à jour manuellement
docker-compose pull
docker-compose up -d

# Rollback vers une version précédente
docker-compose down
docker tag username/safebase-backend:sha-abc123 username/safebase-backend:latest
docker-compose up -d

# Nettoyer
docker system prune -af
```

### Localement

```bash
# Tester le build localement
docker-compose build

# Pousser manuellement vers Docker Hub
docker login
docker-compose push

# Voir les images
docker images | grep safebase
```

## 🐛 Dépannage

### Le pipeline échoue aux tests

```bash
# Tester localement
cd backend
go test -v ./...

cd ../frontend
npm test
```

### Le build Docker échoue

```bash
# Build local pour voir l'erreur
docker-compose build --no-cache backend
docker-compose build --no-cache frontend
```

### Le déploiement échoue

1. Vérifiez les secrets GitHub (Settings → Secrets)
2. Testez la connexion SSH :
   ```bash
   ssh -i ~/.ssh/github_deploy user@votre-serveur.com
   ```
3. Vérifiez les logs du workflow sur GitHub Actions

### Les conteneurs ne démarrent pas

```bash
# Sur le serveur
docker-compose logs

# Vérifier l'espace disque
df -h

# Vérifier la mémoire
free -h

# Redémarrer Docker
sudo systemctl restart docker
```

## 🔐 Sécurité

### Bonnes pratiques

✅ **À FAIRE** :
- Utiliser des secrets GitHub pour les credentials
- Changer tous les mots de passe par défaut
- Utiliser HTTPS en production
- Activer le pare-feu (UFW)
- Faire des backups réguliers
- Monitorer les logs

❌ **À NE PAS FAIRE** :
- Committer des secrets dans le code
- Utiliser des mots de passe faibles
- Exposer les bases de données publiquement
- Désactiver les health checks
- Ignorer les mises à jour de sécurité

### Pare-feu (UFW)

```bash
# Autoriser SSH
sudo ufw allow 22/tcp

# Autoriser HTTP/HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Activer
sudo ufw enable
```

## 📈 Optimisations

### Cache Docker

Le pipeline utilise le cache Docker pour accélérer les builds :
- Les layers Docker sont mis en cache
- Les dépendances Go et npm sont mises en cache
- Build multi-architecture optimisé

### Images multi-architecture

Les images sont buildées pour :
- `linux/amd64` (serveurs x86_64)
- `linux/arm64` (serveurs ARM, Apple Silicon)

## 📚 Ressources

- [Documentation complète](./.github/DEPLOYMENT.md)
- [Docker Hub](https://hub.docker.com)
- [GitHub Actions](https://docs.github.com/en/actions)
- [Docker Compose](https://docs.docker.com/compose/)

## 🆘 Support

En cas de problème :
1. Consultez les logs du workflow GitHub Actions
2. Vérifiez les logs sur le serveur
3. Consultez la documentation complète
4. Ouvrez une issue sur GitHub

---

**Bon déploiement ! 🚀**

