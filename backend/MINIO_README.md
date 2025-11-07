# Service MinIO pour Stockage des Sauvegardes

Ce document explique comment utiliser MinIO pour stocker les dumps de base de données dans le cloud.

## Configuration

Le service MinIO se configure via des variables d'environnement :

```bash
# Configuration MinIO
MINIO_ENDPOINT=localhost:9000          # Adresse du serveur MinIO
MINIO_ACCESS_KEY=minioadmin            # Clé d'accès MinIO
MINIO_SECRET_KEY=minioadmin            # Clé secrète MinIO
MINIO_BUCKET=SafeBase                  # Nom du bucket (défaut: SafeBase)
MINIO_USE_SSL=false                    # Utiliser SSL (true/false)
```

## Installation et Configuration de MinIO

### Installation avec Docker (recommandé)

```bash
# Démarrer MinIO
docker run -d \
  --name minio-server \
  -p 9000:9000 \
  -p 9001:9001 \
  -e MINIO_ROOT_USER=minioadmin \
  -e MINIO_ROOT_PASSWORD=minioadmin \
  -v /data/minio:/data \
  minio/minio server /data --console-address ":9001"
```

### Accès à la console MinIO

- URL: http://localhost:9001
- Utilisateur: minioadmin
- Mot de passe: minioadmin

### Création du bucket SafeBase

1. Connectez-vous à la console MinIO
2. Cliquez sur "Create Bucket"
3. Nommez le bucket "SafeBase"
4. Configurez les permissions selon vos besoins

## Fonctionnement

### Stockage Hybride

Le système implémente un stockage hybride :

1. **Stockage primaire**: Les dumps sont d'abord créés localement dans `./db/backups/`
2. **Stockage cloud**: Si MinIO est configuré, les fichiers sont automatiquement uploadés vers MinIO
3. **Fallback**: En cas de problème avec MinIO, les fichiers restent disponibles localement

### Téléchargement

Lors du téléchargement d'une sauvegarde :
1. Le système essaie d'abord de récupérer depuis MinIO
2. En cas d'échec, il récupère depuis le stockage local

### Suppression

Lors de la suppression d'une sauvegarde :
1. Le fichier est supprimé de MinIO (si disponible)
2. Le fichier local est supprimé
3. L'enregistrement en base de données est supprimé

## Structure des fichiers dans MinIO

```
SafeBase/
├── backups/
│   ├── nom_base_donnees/
│   │   ├── nom_base_donnees_20241107_143052.sql
│   │   └── nom_base_donnees_20241108_091234.sql
│   └── autre_base/
│       └── autre_base_20241107_154321.sql
```

## Migration vers le Cloud

Pour migrer vos sauvegardes existantes vers MinIO :

1. Assurez-vous que MinIO est configuré et accessible
2. Redémarrez l'application
3. Les nouvelles sauvegardes seront automatiquement stockées dans MinIO
4. Les anciennes sauvegardes restent disponibles localement

## Dépannage

### Erreur de connexion MinIO

Si MinIO n'est pas disponible :
- Les sauvegardes continueront de fonctionner en stockage local uniquement
- Un avertissement sera affiché au démarrage
- Vérifiez la configuration et la connectivité réseau

### Erreur d'upload vers MinIO

- Le fichier reste disponible localement
- Vérifiez les permissions du bucket
- Vérifiez l'espace disponible dans MinIO

## Sécurité

- Utilisez des clés d'accès fortes en production
- Activez SSL (`MINIO_USE_SSL=true`) en production
- Configurez les politiques d'accès appropriées sur le bucket
- Surveillez les logs pour détecter les problèmes de stockage