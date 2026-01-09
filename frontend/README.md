# SafeBase Frontend

Interface utilisateur pour la plateforme SafeBase de sauvegarde et restauration de bases de donn√©es.

## Technologies utilis√©es

- **Vue.js 3** avec Composition API
- **TypeScript** pour le typage statique
- **Vite** comme bundler de d√©veloppement
- **Tailwind CSS** pour le styling
- **Vue Router** pour la navigation
- **Pinia** pour la gestion d'√©tat
- **Prettier** pour le formatage du code

## Structure du projet

```
frontend/
 src/
Ç    assets/          # Images, styles, ressources
Ç    components/      # Composants r√©utilisables
Ç   Ç    BackupList.vue
Ç   Ç    RestoreForm.vue
Ç   Ç    AlertBox.vue
Ç    pages/          # Pages principales
Ç   Ç    Dashboard.vue
Ç   Ç    History.vue
Ç    router/         # Configuration des routes
Ç    stores/         # Gestion d'√©tat avec Pinia
Ç    views/          # Vues principales
 public/             # Ressources statiques
 dist/               # Build de production
```

## Installation

```bash
# Installation des d√©pendances
npm install

# Lancement du serveur de d√©veloppement
npm run dev

# Build pour la production
npm run build

# Pr√©visualisation du build de production
npm run preview
```

## Scripts disponibles

- `npm run dev` - Serveur de d√©veloppement
- `npm run build` - Build de production
- `npm run preview` - Pr√©visualisation du build
- `npm run format` - Formatage avec Prettier

## Configuration Tailwind

Tailwind CSS est configur√© pour scanner tous les fichiers Vue, JS, TS dans le dossier `src/`. Les classes utilitaires sont disponibles dans tous les composants.

## API Backend

Le frontend communique avec l'API Go via des requ√tes HTTP REST :

- `GET /api/databases` - Liste des bases de donn√©es
- `POST /api/backup` - Lancer une sauvegarde
- `POST /api/restore` - Restaurer une sauvegarde
- `GET /api/backups` - Historique des sauvegardes

## D√©veloppement

1. Le serveur de d√©veloppement se lance sur `http://localhost:5173`
2. Hot reload activ√© pour un d√©veloppement rapide
3. TypeScript activ√© pour la v√©rification de types
4. Prettier configur√© pour la qualit√© du code


## Build et d√©ploiement

Le build g√©n√re des fichiers statiques optimis√©s dans le dossier `dist/` qui peuvent √tre servis par n'importe quel serveur web ou int√©gr√©s dans un conteneur Docker.
