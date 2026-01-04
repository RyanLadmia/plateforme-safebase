# Guide de démarrage avec Docker

## Problèmes résolus

### Backend
- **Erreur** : `exec format error`
- **Cause** : Compilation du binaire pour une mauvaise architecture
- **Solution** : Utilisation d'Air pour le hot-reload dans Docker (comme en développement local)
- **Mise à jour** : Utilisation de Go 1.25 pour supporter Air v1.63+

- **Erreur** : `pg_dump non trouvé: aucun exécutable trouvé dans les chemins testés`
- **Cause** : Image Alpine ne contient pas les outils de sauvegarde de bases de données
- **Solution** : Installation de `postgresql-client`, `mysql-client`, `zip`, `unzip`

### Frontend
- **Erreur** : `Node.js 18.20.8. Vite requires Node.js version 20.19+ or 22.12+`
- **Cause** : Version de Node.js trop ancienne
- **Solution** : Mise à jour vers Node.js 20

## Démarrage rapide

### 1. Arrêter les conteneurs existants (si nécessaire)
```bash
docker-compose down
```

### 2. Reconstruire les images
```bash
docker-compose build
```

### 3. Démarrer tous les services
```bash
docker-compose up
```

Ou en mode détaché :
```bash
docker-compose up -d
```

### 4. Vérifier les logs
```bash
# Tous les services
docker-compose logs -f

# Backend uniquement
docker-compose logs -f backend

# Frontend uniquement
docker-compose logs -f frontend
```

## Services disponibles

- **Frontend** : http://localhost:3000
- **Backend API** : http://localhost:8080
- **PostgreSQL** : localhost:5432
- **MySQL** : localhost:3306
- **Grafana** : http://localhost:3001
- **Prometheus** : http://localhost:9090

## Commandes utiles

### Arrêter tous les services
```bash
docker-compose down
```

### Reconstruire un service spécifique
```bash
docker-compose build backend
docker-compose build frontend
```

### Redémarrer un service
```bash
docker-compose restart backend
docker-compose restart frontend
```

### Voir les conteneurs en cours d'exécution
```bash
docker-compose ps
```

### Accéder au shell d'un conteneur
```bash
docker-compose exec backend sh
docker-compose exec frontend sh
```

### Nettoyer complètement (attention : supprime les volumes)
```bash
docker-compose down -v
```

## Configuration

### Backend
Les variables d'environnement sont définies dans `docker-compose.yml` :
- `PORT=8080`
- `JWT_SECRET` (à changer en production)
- `DB_HOST=postgres`
- `DB_PORT=5432`
- `DB_USER=user`
- `DB_PASSWORD=password`
- `DB_NAME=safebase`

**Pour les sauvegardes MEGA** (optionnel) :
- `MEGA_EMAIL=votre_email@example.com`
- `MEGA_PASSWORD=votre_mot_de_passe_mega`

Voir [DOCKER_BACKUP_MEGA.md](./DOCKER_BACKUP_MEGA.md) pour plus de détails.

### Frontend
Utilise Node.js 20 pour supporter Vite.

## Hot-reload

Les deux services supportent le hot-reload :
- **Backend** : Air détecte les changements et recompile automatiquement
- **Frontend** : Vite détecte les changements et recharge automatiquement

## Résolution de problèmes

### Le backend ne démarre pas
1. Vérifiez que PostgreSQL est bien démarré : `docker-compose ps`
2. Vérifiez les logs : `docker-compose logs backend`
3. Vérifiez que le port 8080 n'est pas déjà utilisé

### Le frontend ne démarre pas
1. Vérifiez les logs : `docker-compose logs frontend`
2. Vérifiez que le port 3000 n'est pas déjà utilisé
3. Supprimez `node_modules` et reconstruisez : `docker-compose build frontend`

### Les changements ne sont pas détectés
1. Sur Mac/Windows, le polling peut être nécessaire (déjà configuré dans `.air.toml`)
2. Redémarrez le service : `docker-compose restart backend` ou `docker-compose restart frontend`

