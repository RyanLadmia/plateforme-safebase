# 🗄️ Configuration des sauvegardes avec MEGA dans Docker

## 🔧 Problème résolu : "pg_dump non trouvé"

### Cause
L'image Docker `golang:1.25-alpine` ne contient pas les outils nécessaires pour effectuer les sauvegardes de bases de données :
- `pg_dump` pour PostgreSQL
- `mysqldump` pour MySQL
- `zip`/`unzip` pour la compression

### Solution appliquée
Installation des packages Alpine Linux nécessaires dans le `Dockerfile` :

```dockerfile
RUN apk add --no-cache \
    postgresql-client \
    mysql-client \
    zip \
    unzip
```

---

## 🚀 Configuration de MEGA pour les sauvegardes

### 1. Variables d'environnement

Ajoutez les variables MEGA dans `docker-compose.yml` :

```yaml
backend:
  environment:
    # ... autres variables ...
    - MEGA_EMAIL=votre_email@example.com
    - MEGA_PASSWORD=votre_mot_de_passe_mega
```

### 2. Redémarrer le backend

```bash
docker-compose down
docker-compose build backend
docker-compose up -d backend
```

---

## 📦 Outils installés dans le conteneur

| Outil | Package | Utilisation |
|-------|---------|-------------|
| `pg_dump` | postgresql-client | Sauvegarde PostgreSQL |
| `pg_restore` | postgresql-client | Restauration PostgreSQL |
| `mysqldump` | mysql-client | Sauvegarde MySQL |
| `mysql` | mysql-client | Restauration MySQL |
| `zip` | zip | Compression des sauvegardes |
| `unzip` | unzip | Décompression des sauvegardes |

---

## 🔍 Vérification de l'installation

### Accéder au conteneur backend

```bash
docker-compose exec backend sh
```

### Vérifier que les outils sont installés

```bash
# Vérifier pg_dump
which pg_dump
pg_dump --version

# Vérifier mysqldump
which mysqldump
mysqldump --version

# Vérifier zip
which zip
zip --version
```

Sortie attendue :
```
/usr/bin/pg_dump
pg_dump (PostgreSQL) 15.x

/usr/bin/mysqldump
mysqldump  Ver 8.x.x

/usr/bin/zip
Copyright (c) 1990-2008 Info-ZIP
```

---

## 🗄️ Test de sauvegarde

### 1. Créer une base de données de test

Via l'interface frontend ou l'API :

```bash
curl -X POST http://localhost:8080/api/databases \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "Test DB",
    "type": "postgresql",
    "host": "postgres",
    "port": "5432",
    "db_name": "safebase",
    "username": "user",
    "password": "password"
  }'
```

### 2. Lancer une sauvegarde

Via l'interface frontend ou l'API :

```bash
curl -X POST http://localhost:8080/api/backups \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "database_id": 1,
    "description": "Test backup"
  }'
```

### 3. Vérifier les logs

```bash
docker-compose logs -f backend
```

Vous devriez voir :
```
[BACKUP] Starting PostgreSQL dump for database safebase...
[BACKUP] Using pg_dump at: /usr/bin/pg_dump
[BACKUP] pg_dump completed successfully
[BACKUP] Compressing SQL file to ZIP...
[BACKUP] Uploading to MEGA...
[BACKUP] Backup completed successfully
```

---

## 🐛 Résolution de problèmes

### Erreur : "pg_dump non trouvé"

**Cause** : Le conteneur n'a pas été reconstruit avec les nouveaux outils.

**Solution** :
```bash
docker-compose down
docker-compose build backend
docker-compose up -d
```

### Erreur : "MEGA login failed"

**Cause** : Identifiants MEGA incorrects ou non configurés.

**Solution** :
1. Vérifiez les variables d'environnement dans `docker-compose.yml`
2. Vérifiez que votre compte MEGA est actif
3. Redémarrez le backend :
```bash
docker-compose restart backend
```

### Erreur : "Connection refused" lors de la sauvegarde

**Cause** : Le conteneur backend ne peut pas accéder à la base de données.

**Solution** :
1. Vérifiez que PostgreSQL/MySQL est démarré :
```bash
docker-compose ps
```

2. Utilisez le nom du service Docker comme host :
   - Pour PostgreSQL : `host: postgres` (pas `localhost`)
   - Pour MySQL : `host: mysql` (pas `localhost`)

### La sauvegarde fonctionne mais n'est pas uploadée sur MEGA

**Cause** : Service MEGA non initialisé ou erreur de connexion.

**Solution** :
1. Vérifiez les logs backend :
```bash
docker-compose logs backend | grep -i mega
```

2. Vous devriez voir :
```
Service Mega initialisé avec succès
```

3. Si vous voyez "Configuration Mega manquante", ajoutez les variables d'environnement MEGA.

---

## 📊 Flux de sauvegarde complet

```
1. Utilisateur demande une sauvegarde
   ↓
2. Backend vérifie la connexion à la base de données
   ↓
3. Exécution de pg_dump ou mysqldump
   ↓
4. Compression du fichier SQL en ZIP
   ↓
5. Chiffrement AES-256 (si configuré)
   ↓
6. Upload vers MEGA (si configuré)
   ↓
7. Sauvegarde locale dans /app/db/backups/
   ↓
8. Mise à jour du statut dans la base de données
```

---

## 🔐 Sécurité des sauvegardes

### Chiffrement
Les sauvegardes sont chiffrées avec AES-256 avant l'upload vers MEGA.

### Stockage local
Les sauvegardes sont également conservées localement dans le volume Docker :
```
/app/db/backups/
├── postgresql/
│   └── backup_2026-01-04_123456.zip
└── mysql/
    └── backup_2026-01-04_123457.zip
```

### Accéder aux sauvegardes locales

```bash
# Depuis l'hôte
docker-compose exec backend ls -lh /app/db/backups/postgresql/
docker-compose exec backend ls -lh /app/db/backups/mysql/

# Copier une sauvegarde vers l'hôte
docker cp $(docker-compose ps -q backend):/app/db/backups/postgresql/backup.zip ./
```

---

## 📝 Configuration recommandée pour la production

### docker-compose.yml

```yaml
backend:
  environment:
    # Configuration de base
    - PORT=8080
    - JWT_SECRET=${JWT_SECRET}  # Utiliser un fichier .env
    
    # Base de données
    - DB_HOST=postgres
    - DB_PORT=5432
    - DB_USER=${DB_USER}
    - DB_PASSWORD=${DB_PASSWORD}
    - DB_NAME=${DB_NAME}
    
    # MEGA (sauvegardes cloud)
    - MEGA_EMAIL=${MEGA_EMAIL}
    - MEGA_PASSWORD=${MEGA_PASSWORD}
  
  volumes:
    - ./backend:/app
    - /app/tmp
    - backup_data:/app/db/backups  # Volume persistant pour les sauvegardes

volumes:
  backup_data:  # Ajouter ce volume
```

### Fichier .env (à créer)

```bash
# Ne pas commiter ce fichier !
JWT_SECRET=votre_secret_jwt_tres_securise
DB_USER=user
DB_PASSWORD=mot_de_passe_securise
DB_NAME=safebase
MEGA_EMAIL=votre_email@example.com
MEGA_PASSWORD=votre_mot_de_passe_mega
```

---

## ✅ Checklist de vérification

- [ ] Les outils de sauvegarde sont installés (`pg_dump`, `mysqldump`, `zip`)
- [ ] Les variables d'environnement MEGA sont configurées
- [ ] Le conteneur backend a été reconstruit
- [ ] Les logs backend montrent "Service Mega initialisé avec succès"
- [ ] Une sauvegarde de test fonctionne
- [ ] Les sauvegardes apparaissent sur MEGA
- [ ] Les sauvegardes locales sont accessibles dans `/app/db/backups/`

---

## 📚 Commandes utiles

```bash
# Reconstruire le backend avec les nouveaux outils
docker-compose build backend

# Redémarrer uniquement le backend
docker-compose restart backend

# Voir les logs de sauvegarde en temps réel
docker-compose logs -f backend | grep BACKUP

# Accéder au shell du backend
docker-compose exec backend sh

# Lister les sauvegardes PostgreSQL
docker-compose exec backend ls -lh /app/db/backups/postgresql/

# Lister les sauvegardes MySQL
docker-compose exec backend ls -lh /app/db/backups/mysql/

# Vérifier l'espace disque utilisé par les sauvegardes
docker-compose exec backend du -sh /app/db/backups/
```

---

**Tout est maintenant configuré pour les sauvegardes avec MEGA ! 🎉**

