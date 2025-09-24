# SafeBase Frontend

Interface utilisateur pour la plateforme SafeBase de sauvegarde et restauration de bases de données.

## Technologies utilisées

- **Vue.js 3** avec Composition API
- **TypeScript** pour le typage statique
- **Vite** comme bundler de développement
- **Tailwind CSS** pour le styling
- **Vue Router** pour la navigation
- **Pinia** pour la gestion d'état
- **Prettier** pour le formatage du code

## Structure du projet

```
frontend/
├── src/
│   ├── assets/          # Images, styles, ressources
│   ├── components/      # Composants réutilisables
│   │   ├── BackupList.vue
│   │   ├── RestoreForm.vue
│   │   └── AlertBox.vue
│   ├── pages/          # Pages principales
│   │   ├── Dashboard.vue
│   │   └── History.vue
│   ├── router/         # Configuration des routes
│   ├── stores/         # Gestion d'état avec Pinia
│   └── views/          # Vues principales
├── public/             # Ressources statiques
└── dist/               # Build de production
```

## Installation

```bash
# Installation des dépendances
npm install

# Lancement du serveur de développement
npm run dev

# Build pour la production
npm run build

# Prévisualisation du build de production
npm run preview
```

## Scripts disponibles

- `npm run dev` - Serveur de développement
- `npm run build` - Build de production
- `npm run preview` - Prévisualisation du build
- `npm run format` - Formatage avec Prettier

## Configuration Tailwind

Tailwind CSS est configuré pour scanner tous les fichiers Vue, JS, TS dans le dossier `src/`. Les classes utilitaires sont disponibles dans tous les composants.

## API Backend

Le frontend communique avec l'API Go via des requêtes HTTP REST :

- `GET /api/databases` - Liste des bases de données
- `POST /api/backup` - Lancer une sauvegarde
- `POST /api/restore` - Restaurer une sauvegarde
- `GET /api/backups` - Historique des sauvegardes

## Développement

1. Le serveur de développement se lance sur `http://localhost:5173`
2. Hot reload activé pour un développement rapide
3. TypeScript activé pour la vérification de types
4. Prettier configuré pour la qualité du code


## Build et déploiement

Le build génère des fichiers statiques optimisés dans le dossier `dist/` qui peuvent être servis par n'importe quel serveur web ou intégrés dans un conteneur Docker.
