# Tester le CI/CD localement avec Act

[Act](https://github.com/nektos/act) permet de tester vos workflows GitHub Actions localement avant de les pousser.

## Installation

### macOS
```bash
brew install act
```

### Linux
```bash
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash
```

### Windows
```bash
choco install act-cli
```

## Configuration

Créer un fichier `.actrc` à la racine du projet :

```bash
# .actrc
--container-architecture linux/amd64
--platform ubuntu-latest=catthehacker/ubuntu:act-latest
```

Créer un fichier `.secrets` pour les secrets locaux :

```bash
# .secrets (NE PAS COMMITTER)
DOCKER_USERNAME=votre-username
DOCKER_PASSWORD=votre-token
```

## Utilisation

### Lister les workflows
```bash
act -l
```

### Tester le job de tests
```bash
act -j test
```

### Tester le job de build (avec secrets)
```bash
act -j build --secret-file .secrets
```

### Tester tout le workflow
```bash
act push --secret-file .secrets
```

### Mode dry-run (voir ce qui serait exécuté)
```bash
act -n
```

## Notes

- Act utilise Docker pour simuler les runners GitHub
- Certaines actions peuvent ne pas fonctionner exactement comme sur GitHub
- Le déploiement SSH ne fonctionnera pas localement (normal)
- Utile pour tester les tests et les builds rapidement

## Exemple de workflow de développement

```bash
# 1. Faire vos modifications
git add .

# 2. Tester localement
act -j test

# 3. Si OK, commit et push
git commit -m "feat: nouvelle fonctionnalité"
git push origin main
```

Cela permet de détecter les erreurs avant de les pousser sur GitHub ! 

