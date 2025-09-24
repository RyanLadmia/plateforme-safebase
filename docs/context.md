# Job 2 : Safebase

API REST pour sauvegarder des bases de données

## Préparation

### I - Présentation du projet

- **Présentation globale :**
    
    SafeBase vise à développer une solution complète de gestion de sauvegarde et de restauration de base de données sous forme d'une API REST.
    
- **Pourquoi ce projet :**
    
    Pour permettre à des utilisateurs, des applications ou des entreprises de sauvegarder et restaurer leur base de données de manière efficace rapide et sécurisée en utilisant un outils de source externe.
    
- **Objectifs :**
    
    Offrir la possibilité à travers une interface intuitive et un système bien penser de gérer sa base de données et façon indépendante et optimale sans craindre de perdre ses données et de les réemployer facilement.
    
- **MVP (Minimum Viable Product) :**
    - sauvegarde des données
    - restauration des données
    - interface simple
    
- **Valeur ajoutée / bénéfices pour l'utilisateur :**
    - gain de temps
    - sécurité
    - fiabilité
    - efficacité
    - source externe de sauvegarde

### II - Cahier des charges

La solution **SafeBase** s'adresse à plusieurs catégories de profils :

- **Entreprises** : PME, startups ou grands groupes qui gèrent des bases de données critiques (CRM, ERP, e-commerce).
- **Organisations** : associations, ONG ou structures publiques qui doivent sécuriser leurs données sensibles.
- **Applications** : éditeurs SaaS ou solutions internes nécessitant des sauvegardes automatisées et restaurables rapidement.
- **Utilisateurs particuliers** : développeurs indépendants, freelances ou étudiants qui veulent protéger leurs bases de test ou projets personnels.

#### Persona(s) : utilisateur type, contexte, frustrations, attentes

- **Entreprise :**
    - **Contexte :**
        Une PME de vente en ligne qui gère une base clients et une base de commandes.
        
    - **Problèmes rencontrés :**
        - Sauvegardes trop longues à faire à la main.
        - Risque de perte de données en cas de panne.
        - Difficile de retrouver facilement les anciennes sauvegardes.
    - **Attentes :**
        - Sauvegardes régulières sans intervention humaine.
        - Un tableau de bord simple pour tout contrôler.
        - Être averti si une sauvegarde échoue.

- **Organisation (association) :**
    - **Contexte :**
        Une association qui gère ses adhérents et ses dons dans un outil informatique.
        
    - **Problèmes rencontrés :**
        - Peu de compétences techniques dans l'équipe.
        - Dépendance à un prestataire externe pour les sauvegardes.
        - Budget limité
    - **Attentes :**
        - Interface facile à utiliser, même sans connaissances techniques.
        - Sauvegardes sécurisées et accessibles.
        - Solution économique et durable.

- **Application :**
    - **Contexte :**
        Une startup qui propose une application web à ses clients.
        
    - **Problèmes rencontrés :**
        - Gérer différents types de bases de données.
        - Besoin de restaurer les données rapidement pour ne pas interrompre le service.
        - Processus de sauvegarde compliqués à unifier.
    - **Attentes :**
        - Pouvoir sauvegarder plusieurs types de bases.
        - Un outil adaptable et flexible.
        - Un suivi clair de toutes les sauvegardes.

- **Utilisateur :**
    - **Contexte :**
        Un développeur freelance qui travaille sur plusieurs projets web.
        
    - **Problèmes rencontrés :**
        - Sauvegardes dispersées, pas toujours régulières.
        - Restauration compliquée quand un projet plante.
        - Pas de solution simple et centralisée.
    - **Attentes :**
        - Pouvoir lancer une sauvegarde ou une restauration facilement.
        - Automatiser les sauvegardes sans effort.
        - Une solution simple à installer et légère à utiliser.

#### Fonctionnalités principales

Liste de toutes les fonctionnalités prévues (MoSCoW) :

| Must | Should | Could | Won't (dans l'immédiat) |
| --- | --- | --- | --- |
| Sauvegarde automatique (CRON) | Authentification | Support d'autres bases (NoSQL) | Suppression définitive des base de données |
| Sauvegarde manuelle (MySQL/PostgreSQL) | Gestion des utilisateurs et des rôles | Sauvegardes chiffrées avec clé utilisateur (sécurité) | Duplication des bases de données |
| Restauration manuelle | Modification de la planification (changer la fréquence des sauvegardes) | Modification | Application mobile |
| Notifications de réussite/échec | Interface graphique élégante | Déploiement en production |  |
| Historique des sauvegardes | Tests unitaires, intégrations, fonctionnels, end to end | Export des journaux d'activités (sauvegarde et restaurations) |  |
| Interface graphique simple | Intégration CI/CD | Mode multi-server |  |
| Gestions des connexions DB |  |  |  |

#### Contraintes techniques :
- GO
- Conteneurisation

#### Contraintes fonctionnelles :
- Sauvegarde automatique des bases de données
- Restauration manuelle des bases de données
- Historique des versions
- Gestion des bases de données
- Authentification
- Système de Notification logs et alertes
- Tests unitaires, d'intégration fonctionnels et end to end

#### Benchmark : comparaison avec solutions existantes, forces/faiblesses

- **Node.JS**
    - **Positif :**
        - Full JS (front et back)
        - Rapide à développer
        - Très gros écosystème
    - **Négatif :**
        - Moins adapté aux tâches système
        - Mono thread
        - Moins performant que Go
        - Typage plus souple

- **Go**
    - **Positif :**
        - Bonne gestion des tâches système
        - Concurrence native (goroutines)
        - Compilation rapide, binaire léger
        - Très utilisé en production (Docker, Kubernetes…)
    - **Négatif :**
        - Syntaxe verbeuse
        - Moins d'abstractions que Node.js

- **Rust**
    - **Positif :**
        - Performance et sécurité maximale
        - Zéro bug mémoire possible
        - Fort typage et contrôle strict
    - **Négatif :**
        - Moins d'outils clé-en-main
        - Écosystème back-end encore jeune
        - Courbe d'apprentissage raide

Go est le compromis entre performances et simplicité

#### Méthodologie de développement :
- Méthode Agile Scrum (backlog, sprints, réunions, rétrospectives)

### III - Conception et design

#### Diagramme de Gantt : planning détaillé semaine par semaine

| **Semaine** | **Tâche** | **Détails** |
| --- | --- | --- |
| **S1** | Analyse et préparation | - Étude du cahier des charges<br>- Identification des personas<br>- Benchmark<br>- Préparation de l'environnement de dev (Docker, repos GitHub) |
| **S1** | Architecture technique | - Schéma global de l'architecture<br>- Diagramme Docker<br>- Diagramme de séquence pour MVP<br>- MCD et MLD des bases MySQL/PostgreSQL |
| **S2** | Backend Go – Partie MVP | - Gestion des bases (ajout, suppression)<br>- Sauvegarde manuelle / restauration<br>- API REST pour le frontend<br>- Versioning des sauvegardes |
| **S2** | Tests backend | - Tests unitaires des fonctions clés<br>- Tests de sauvegarde/restauration avec bases fictives |
| **S3** | Frontend Vue.js | - Interface de contrôle basique (MVP)<br>- Formulaire ajout base, boutons sauvegarde/restauration<br>- Affichage de l'historique des sauvegardes<br>- Notifications simples |
| **S3** | Intégration frontend ↔ backend | - Connexion API REST<br>- Tests manuels des fonctionnalités principales |
| **S3** | Automatisation | - CRON pour sauvegardes automatiques<br>- Logs et notifications |
| **S4** | Tests finaux et debug | - Tests fonctionnels complets<br>- Vérification multi-SGBD (MySQL/Postgres)<br>- Ajustements et corrections |
| **S4** | Documentation et conception finale | - Wireframes et maquette finalisés<br>- Diagrammes UML et Docker finalisés<br>- Cahier des charges complété<br>- Rédaction des variables d'environnement et sécurité |
| **S4** | Conteneurisation finale | - Docker Compose final (frontend + backend + MySQL + Postgres)<br>- Tests dans les conteneurs<br>- Préparation du dépôt GitHub public |
| **S4** | Bilan et présentation | - Bilan du projet et perspectives d'évolution<br>- Bilan personnel<br>- Préparation soutenance / démo |

#### Tableau Kanban : colonnes À faire, En cours, Terminé pour gestion des tâches
GitHub Project

#### Arborescence du projet : (dossiers backend, frontend, docker, tests)

```
safebase/
│── backend/                     # API REST en Go
│   ├── cmd/                     # Point d'entrée (main.go)
│   │   └── main.go
│   ├── internal/                # Logique interne (clean architecture)
│   │   ├── db/                  # Connexions à la DB (MySQL, Postgres, interne ⇒ contient du code)
│   │   │   ├── mysql.go
│   │   │   ├── postgres.go
│   │   │   └── safebase.go      # Base interne (historique, logs)
│   │   ├── services/            # Logique métier (sauvegarde, restauration)
│   │   │   ├── backup.go
│   │   │   ├── restore.go
│   │   │   └── scheduler.go     # Gestion CRON
│   │   ├── api/                 # Handlers HTTP (endpoints API REST)
│   │   │   ├── backup_handler.go
│   │   │   ├── restore_handler.go
│   │   │   └── db_handler.go
│   │   └── utils/               # Fonctions utilitaires (logs, sécurité…)
│   │       ├── logger.go
│   │       └── security.go
│   ├── tests/                   # Tests unitaires backend
│   │   ├── backup_test.go
│   │   └── restore_test.go
│   └── go.mod / go.sum
│
│── frontend/                    # Interface utilisateur Vue.js
│   ├── public/                  # Ressources statiques
│   ├── src/
│   │   ├── assets/              # Images, icônes
│   │   ├── components/          # Composants Vue.js
│   │   │   ├── BackupList.vue
│   │   │   ├── RestoreForm.vue
│   │   │   └── AlertBox.vue
│   │   ├── pages/               # Pages principales
│   │   │   ├── Dashboard.vue
│   │   │   └── History.vue
│   │   ├── router/              # Routes Vue.js
│   │   │   └── index.js
│   │   ├── store/               # State management (Pinia/Vuex)
│   │   │   └── backup.js
│   │   └── App.vue
│   └── package.json
│
│── docker/                      # Configurations Docker
│   ├── Dockerfile.backend
│   ├── Dockerfile.frontend
│   ├── mysql_init.sql           # Script init MySQL (tests)
│   ├── postgres_init.sql        # Script init PostgreSQL (tests)
│   └── docker-compose.yml
│
│── db/                          # Sauvegardes et historique (contient les données)
│   ├── backups/                 # Dossiers persistants de backup
│   │   ├── mysql/
│   │   └── postgresql/
│   └── safebase.db              # SQLite interne (ou Postgres interne)
│
│── tests/                       # Tests d'intégration & e2e
│   ├── integration/             # Tests API (Go + DB factices)
│   └── e2e/                     # Tests bout en bout (frontend + backend)
│
│── docs/                        # Documentation du projet
│   ├── cahier_des_charges.md
│   ├── conception.md
│   ├── architecture.png
│   └── gantt.png
│
│── .env                         # Variables d'environnement
│── .gitignore
│── README.md                    # Présentation GitHub
```

#### Wireframes : croquis des pages principales
#### Maquettes : version graphique plus aboutie
#### Charte graphique : couleurs, typographies, boutons, icônes
#### Responsive : adaptation aux écrans mobile, tablette, desktop
#### Accessibilité et ergonomie : navigation claire, contraste, labels, boutons accessibles

#### Glossaire : définitions techniques importantes (API REST, Docker, ORM, CRON, backup, etc.)

- **API REST**
    Application Programming Interface **(Representational State Transfer)** est un moyen pour une application (frontend) de communiquer avec un serveur (backend) via protocol HTTP. Les données sont généralement échangées au format JSON.
    
- **DOCKER**
    Technologie de conteneurisation qui permet de créer des environnements isolés pour les applications.
    
- **ORM :**
    Object-Relational-Mapping, outil qui permet de **manipuler une base de données relationnelle** via des objets dans le code, plutôt que d'écrire directement des requêtes SQL.
    
- **CRON job :**
    Tâche automatique ****planifiée qui s'exécute à intervalles réguliers sur un système Unix/Linux.
    
- **Dump :**
    Copie complète d'une base de données dans un fichier
    
- **Pointer :**
- **Channels :**
- **Goroutines :**

### IV - Architecture technique

#### Diagramme d'architecture logique : flux backend ↔ frontend ↔ bases (internes et externes)

```
[Utilisateur]
 │
▼
[Frontend Vue.js] ───> [Backend Go (API REST)]
 │                                             │
 │                                            ├──> [Base MySQL] (sauvegarde)
 │                                            ├──> [Base PostgreSQL] (sauvegarde)
 │                                            └──> [SafeBase DB interne] (historique, logs)
 │
└──> [Docker Compose] orchestre tous les services
```

#### Diagramme de conteneur Docker : conteneurs backend, frontend, bases internes et clientes, volumes, réseau

```
+---------------- Docker Network safebase_net ----------------+
|                                                             |
|  [frontend]  Vue.js  <-->  [backend] Go API                |
|                               │                             |
|                              ▼                             |
|                     +-------------------------+             |
|                     |   MySQL Container       |             |
|                     +-------------------------+             |
|                     | PostgreSQL Container    |             |
|                     +-------------------------+             |
|                     | SafeBase Internal DB    |             |
|                     +-------------------------+             |
|                                                             |
+-------------------------------------------------------------+
```

#### Diagramme de séquence : interactions entre utilisateur, frontend, backend et bases pour sauvegarde/restauration

```
Utilisateur       Frontend(Vue)      Backend(Go API)      DB cible (MySQL/Postgres)
│                 │                  │                    │
│   Clique "Backup" │                │                    │
├————————->│                  │                    │
│                 │   POST /backup   │                    │
│                ├—————————->│                    │
│                 │                  │  dump DB (mysqldump…)   │
│                 │                 ├—————————————->│
│                 │                  │                    │
│                 │                  │   Retour OK + fichier   │
│                 │                  │←—————————————┤
│                 │   Réponse JSON   │                    │
│                 │←—————————-┤                    │
│   Notification OK │                  │                    │
│←—————-———┤                  │                    │
```

#### Description des technologies

- **Back-end :**
    - GO
- **Front-end :**
    - Vue.js
    - Tailwind
- **Conteneurisation :**
    - Docker
- **Sécurité :**
    - JWT (authentification)
    - Bcrypt
- **Base de données :**
    - **Base de données du projet :**
        - PostgreSQL
        - Interface : pgAdmin
        - ORM GORM
    - **Base de données pouvant être sauvegardées :**
        - MySQL
        - PostgreSQL
        - SQLite (pour plus tard)
        - NoSQL (pour plus tard)
- **Tests :**
    - Go (intégré)
    - PostMan
    - Docker (intégré)
    - Cypress
- **Documentation et suivi de projet :**
    - Notion et README.md pour la documentation
    - GitHub Project pour les tâches à réaliser
    - Diagrammes UML

#### Variables d'environnement et configuration

### V - Base de données (si utilisée)

- MCD (Modèle Conceptuel de Données)
- MLD (Modèle Logique de Données)
- Choix de l'ORM et justification
    - GORM

### VI - Processus de sauvegarde / restauration

#### Diagramme de processus / flowchart :
- Sélection de la base
- Sauvegarde manuelle ou automatique
- Stockage des fichiers
- Historique / versioning
- Restauration
- Logs et alertes

#### Description textuelle détaillée pour chaque étape

### VII - Sécurité

- Gestion des mots de passe et informations sensibles
- Permissions sur fichiers de backup
- Protection des endpoints de l'API
- Sauvegardes sécurisées et fiables
- Politique de gestion des erreurs

### VIII - Tests

#### Plan de tests détaillé :
- Tests unitaires (Go backend)
- Tests fonctionnels (sauvegarde/restauration sur bases fictives)
- Tests frontend (Vue.js interaction API)
- Tests d'intégration Docker (stack complet)

#### Tableau des scénarios :
| Test | Description | Résultat attendu |

#### Stratégie de tests alignée sur **cycle en V** : conception → implémentation → tests

### IX - Déploiement et production

- Docker Compose pour l'ensemble du projet
- Scripts de démarrage / mise à jour
- Déploiement sur serveur / cloud
- Gestion des volumes persistants
- Instructions pour production
- Intégration CI/CD possible (GitHub Actions)

### X - Méthodologie et organisation

- **Scrum / Agile** : backlog, sprints, daily meetings, rétrospectives
- Gestion de projet professionnelle : GitHub Projects et Notion
- Suivi des tâches et priorisation avec MoSCoW

### XI - Problèmes rencontrés

- Difficultés techniques et solutions
- Limitations rencontrées
- Gestion des erreurs et cas exceptionnels

### XII - Bilan du projet et perspectives d'évolution

- Fonctionnalités réalisées / non réalisées
- Évolutions possibles : support NoSQL, notifications avancées, interface plus complète, sauvegarde cryptée, amélioration UX

### XIII - Bilan personnel

- Compétences acquises (techniques, DevOps, gestion de projet)
- Points d'amélioration et apprentissages
- Expérience professionnelle simulée (cycle V / Scrum, planning, tests)

### XIV - Conclusion

- Résumé global du projet, objectifs atteints, valeur pour l'utilisateur et l'entreprise

### Annexes

- Fichiers docker-compose.yml, scripts backend et frontend
- Screenshots des wireframes et maquettes
- Diagrammes complets (architecture, séquence, conteneur, processus)
- Exemple de plan de tests complet
