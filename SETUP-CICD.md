# 🚀 Setup CI/CD - Checklist Rapide

## ✅ Étape 1 : Préparer Docker Hub (5 min)

1. Créer un compte sur [hub.docker.com](https://hub.docker.com)
2. Créer un Access Token :
   - Account Settings → Security → New Access Token
   - Nom : `github-actions`
   - Copier le token

## ✅ Étape 2 : Configurer GitHub Secrets (5 min)

Dans votre repo GitHub : `Settings → Secrets and variables → Actions`

Ajouter ces secrets :

```
DOCKER_USERNAME=votre-username
DOCKER_PASSWORD=votre-token-dockerhub
```

**Pour le déploiement automatique (optionnel)** :
```
DEPLOY_HOST=123.45.67.89
DEPLOY_USER=ubuntu
DEPLOY_SSH_KEY=<contenu de la clé privée SSH>
DEPLOY_PATH=/opt/safebase
PRODUCTION_URL=https://api.votre-domaine.com
```

## ✅ Étape 3 : Générer la clé SSH (si déploiement auto)

```bash
# Générer la clé
ssh-keygen -t ed25519 -C "github-deploy" -f ~/.ssh/github_deploy

# Copier sur le serveur
ssh-copy-id -i ~/.ssh/github_deploy.pub user@votre-serveur.com

# Afficher la clé privée (copier dans DEPLOY_SSH_KEY)
cat ~/.ssh/github_deploy
```

## ✅ Étape 4 : Préparer le serveur (10 min)

```bash
# Installer Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Installer Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Créer le répertoire
sudo mkdir -p /opt/safebase
sudo chown $USER:$USER /opt/safebase
cd /opt/safebase

# Créer .env (copier depuis env.production.example)
nano .env

# Créer docker-compose.yml (copier depuis docker-compose.prod.yml)
nano docker-compose.yml
```

## ✅ Étape 5 : Premier push (déclenche le CI/CD)

```bash
git add .
git commit -m "ci: setup CI/CD pipeline"
git push origin main
```

## 🎉 C'est tout !

Le pipeline va automatiquement :
1. ✅ Tester le code
2. ✅ Builder les images Docker
3. ✅ Les pousser vers Docker Hub
4. ✅ Déployer sur le serveur (si configuré)

## 📊 Vérifier le déploiement

1. GitHub : Onglet **Actions** → Voir le workflow
2. Serveur : `docker-compose ps`
3. Browser : `https://votre-domaine.com`

## 📚 Documentation complète

- [Guide CI/CD détaillé](./CI-CD-README.md)
- [Guide de déploiement](./.github/DEPLOYMENT.md)
- [Script de déploiement](./deploy.sh)

## 🆘 Problèmes ?

### Le pipeline échoue
→ Vérifiez les logs dans GitHub Actions

### Le déploiement échoue
→ Vérifiez les secrets GitHub
→ Testez la connexion SSH manuellement

### Les conteneurs ne démarrent pas
→ `docker-compose logs`
→ Vérifiez le fichier `.env`

---

**Temps total : ~20-30 minutes** ⏱️

