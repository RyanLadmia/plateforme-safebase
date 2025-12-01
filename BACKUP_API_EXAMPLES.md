# API de Sauvegarde - Exemples d'utilisation

Ce document présente des exemples d'utilisation de l'API de sauvegarde de bases de données.

## Prérequis

1. Avoir un utilisateur enregistré et être authentifié
2. Avoir `mysqldump` installé pour les sauvegardes MySQL
3. Avoir `pg_dump` installé pour les sauvegardes PostgreSQL

## 1. Authentification

Avant d'utiliser l'API de sauvegarde, vous devez vous authentifier :

```bash
# Inscription
curl -X POST http://localhost:3000/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "motdepasse123",
    "role": "user"
  }'

# Connexion
curl -X POST http://localhost:3000/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "motdepasse123"
  }'
```

Récupérez le token JWT de la réponse et utilisez-le dans les requêtes suivantes.

## 2. Gestion des bases de données

### Ajouter une base de données PostgreSQL

```bash
curl -X POST http://localhost:3000/api/databases \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Ma Base PostgreSQL",
    "type": "postgresql",
    "host": "localhost",
    "port": "5432",
    "username": "postgres",
    "password": "motdepasse",
    "db_name": "ma_base_de_donnees"
  }'
```

### Ajouter une base de données MySQL

```bash
curl -X POST http://localhost:3000/api/databases \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Ma Base MySQL",
    "type": "mysql",
    "host": "localhost",
    "port": "3306",
    "username": "root",
    "password": "motdepasse",
    "db_name": "ma_base_mysql"
  }'
```

### Lister les bases de données

```bash
curl -X GET http://localhost:3000/api/databases \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Obtenir une base de données spécifique

```bash
curl -X GET http://localhost:3000/api/databases/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Mettre à jour une base de données

```bash
curl -X PUT http://localhost:3000/api/databases/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Ma Base PostgreSQL Modifiée",
    "type": "postgresql",
    "host": "localhost",
    "port": "5432",
    "username": "postgres",
    "password": "nouveau_motdepasse",
    "db_name": "ma_base_de_donnees"
  }'
```

### Supprimer une base de données

```bash
curl -X DELETE http://localhost:3000/api/databases/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 3. Gestion des sauvegardes

### Créer une sauvegarde

```bash
curl -X POST http://localhost:3000/api/backups/database/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Cette commande :
1. Crée un enregistrement de sauvegarde avec le statut "pending"
2. Lance le processus de dump en arrière-plan
3. Compresse le dump en fichier ZIP
4. Stocke le fichier dans `db/backups/{type}/`
5. Met à jour le statut à "completed" ou "failed"

### Lister toutes les sauvegardes de l'utilisateur

```bash
curl -X GET http://localhost:3000/api/backups \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Lister les sauvegardes d'une base de données spécifique

```bash
curl -X GET http://localhost:3000/api/backups/database/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Obtenir les détails d'une sauvegarde

```bash
curl -X GET http://localhost:3000/api/backups/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Télécharger une sauvegarde

```bash
curl -X GET http://localhost:3000/api/backups/1/download \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -o ma_sauvegarde.zip
```

### Supprimer une sauvegarde

```bash
curl -X DELETE http://localhost:3000/api/backups/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 4. Structure des fichiers de sauvegarde

Les sauvegardes sont organisées comme suit :

```
db/backups/
├── mysql/
│   └── ma_base_mysql_20231030_143022.zip
└── postgresql/
    └── ma_base_postgresql_20231030_143022.zip
```

Chaque fichier ZIP contient le dump SQL de la base de données.

## 5. Statuts des sauvegardes

- `pending` : La sauvegarde est en cours de création
- `completed` : La sauvegarde a été créée avec succès
- `failed` : La sauvegarde a échoué (voir le champ `error_msg`)

## 6. Exemple de réponse JSON

### Création d'une sauvegarde

```json
{
  "message": "Sauvegarde créée avec succès. Le processus de sauvegarde a commencé.",
  "backup": {
    "id": 1,
    "filename": "ma_base_postgresql_20231030_143022.zip",
    "filepath": "db/backups/postgresql/ma_base_postgresql_20231030_143022.zip",
    "size": 0,
    "status": "pending",
    "created_at": "2023-10-30T14:30:22Z",
    "updated_at": "2023-10-30T14:30:22Z",
    "user_id": 1,
    "database_id": 1
  }
}
```

### Liste des sauvegardes

```json
{
  "backups": [
    {
      "id": 1,
      "filename": "ma_base_postgresql_20231030_143022.zip",
      "filepath": "db/backups/postgresql/ma_base_postgresql_20231030_143022.zip",
      "size": 2048576,
      "status": "completed",
      "created_at": "2023-10-30T14:30:22Z",
      "updated_at": "2023-10-30T14:30:45Z",
      "user_id": 1,
      "database_id": 1
    }
  ]
}
```

## 7. Gestion des erreurs

L'API retourne des codes d'erreur HTTP appropriés :

- `400` : Données invalides
- `401` : Non authentifié
- `403` : Accès non autorisé
- `404` : Ressource introuvable
- `500` : Erreur serveur

Exemple d'erreur :

```json
{
  "error": "Erreur lors de la création de la sauvegarde: type de base de données non supporté: oracle"
}
```
