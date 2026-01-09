# Job 2 : Safebase

API REST pour sauvegarder des bases de donn√©es

## Pr√©paration

### I - Pr√©sentation du projet

- **Pr√©sentation globale :**
    
    SafeBase vise √ d√©velopper une solution compl√te de gestion de sauvegarde et de restauration de base de donn√©es sous forme d'une API REST.
    
- **Pourquoi ce projet :**
    
    Pour permettre √ des utilisateurs, des applications ou des entreprises de sauvegarder et restaurer leur base de donn√©es de mani√re efficace rapide et s√©curis√©e en utilisant un outils de source externe.
    
- **Objectifs :**
    
    Offrir la possibilit√© √ travers une interface intuitive et un syst√me bien penser de g√©rer sa base de donn√©es et fa√on ind√©pendante et optimale sans craindre de perdre ses donn√©es et de les r√©employer facilement.
    
- **MVP (Minimum Viable Product) :**
    - sauvegarde des donn√©es
    - restauration des donn√©es
    - interface simple
    
- **Valeur ajout√©e / b√©n√©fices pour l'utilisateur :**
    - gain de temps
    - s√©curit√©
    - fiabilit√©
    - efficacit√©
    - source externe de sauvegarde

### II - Cahier des charges

La solution **SafeBase** s'adresse √ plusieurs cat√©gories de profils :

- **Entreprises** : PME, startups ou grands groupes qui g√rent des bases de donn√©es critiques (CRM, ERP, e-commerce).
- **Organisations** : associations, ONG ou structures publiques qui doivent s√©curiser leurs donn√©es sensibles.
- **Applications** : √©diteurs SaaS ou solutions internes n√©cessitant des sauvegardes automatis√©es et restaurables rapidement.
- **Utilisateurs particuliers** : d√©veloppeurs ind√©pendants, freelances ou √©tudiants qui veulent prot√©ger leurs bases de test ou projets personnels.

#### Persona(s) : utilisateur type, contexte, frustrations, attentes

- **Entreprise :**
    - **Contexte :**
        Une PME de vente en ligne qui g√re une base clients et une base de commandes.
        
    - **Probl√mes rencontr√©s :**
        - Sauvegardes trop longues √ faire √ la main.
        - Risque de perte de donn√©es en cas de panne.
        - Difficile de retrouver facilement les anciennes sauvegardes.
    - **Attentes :**
        - Sauvegardes r√©guli√res sans intervention humaine.
        - Un tableau de bord simple pour tout contr√¥ler.
        - √tre averti si une sauvegarde √©choue.

- **Organisation (association) :**
    - **Contexte :**
        Une association qui g√re ses adh√©rents et ses dons dans un outil informatique.
        
    - **Probl√mes rencontr√©s :**
        - Peu de comp√©tences techniques dans l'√©quipe.
        - D√©pendance √ un prestataire externe pour les sauvegardes.
        - Budget limit√©
    - **Attentes :**
        - Interface facile √ utiliser, m√me sans connaissances techniques.
        - Sauvegardes s√©curis√©es et accessibles.
        - Solution √©conomique et durable.

- **Application :**
    - **Contexte :**
        Une startup qui propose une application web √ ses clients.
        
    - **Probl√mes rencontr√©s :**
        - G√©rer diff√©rents types de bases de donn√©es.
        - Besoin de restaurer les donn√©es rapidement pour ne pas interrompre le service.
        - Processus de sauvegarde compliqu√©s √ unifier.
    - **Attentes :**
        - Pouvoir sauvegarder plusieurs types de bases.
        - Un outil adaptable et flexible.
        - Un suivi clair de toutes les sauvegardes.

- **Utilisateur :**
    - **Contexte :**
        Un d√©veloppeur freelance qui travaille sur plusieurs projets web.
        
    - **Probl√mes rencontr√©s :**
        - Sauvegardes dispers√©es, pas toujours r√©guli√res.
        - Restauration compliqu√©e quand un projet plante.
        - Pas de solution simple et centralis√©e.
    - **Attentes :**
        - Pouvoir lancer une sauvegarde ou une restauration facilement.
        - Automatiser les sauvegardes sans effort.
        - Une solution simple √ installer et l√©g√re √ utiliser.

#### Fonctionnalit√©s principales

Liste de toutes les fonctionnalit√©s pr√©vues (MoSCoW) :

| Must | Should | Could | Won't (dans l'imm√©diat) |
| --- | --- | --- | --- |
| Sauvegarde automatique (CRON) | Authentification | Support d'autres bases (NoSQL) | Suppression d√©finitive des base de donn√©es |
| Sauvegarde manuelle (MySQL/PostgreSQL) | Gestion des utilisateurs et des r√¥les | Sauvegardes chiffr√©es avec cl√© utilisateur (s√©curit√©) | Duplication des bases de donn√©es |
| Restauration manuelle | Modification de la planification (changer la fr√©quence des sauvegardes) | Modification | Application mobile |
| Notifications de r√©ussite/√©chec | Interface graphique √©l√©gante | D√©ploiement en production |  |
| Historique des sauvegardes | Tests unitaires, int√©grations, fonctionnels, end to end | Export des journaux d'activit√©s (sauvegarde et restaurations) |  |
| Interface graphique simple | Int√©gration CI/CD | Mode multi-server |  |
| Gestions des connexions DB |  |  |  |

#### Contraintes techniques :
- GO
- Conteneurisation

#### Contraintes fonctionnelles :
- Sauvegarde automatique des bases de donn√©es
- Restauration manuelle des bases de donn√©es
- Historique des versions
- Gestion des bases de donn√©es
- Authentification
- Syst√me de Notification logs et alertes
- Tests unitaires, d'int√©gration fonctionnels et end to end

#### Benchmark : comparaison avec solutions existantes, forces/faiblesses

- **Node.JS**
    - **Positif :**
        - Full JS (front et back)
        - Rapide √ d√©velopper
        - Tr√s gros √©cosyst√me
    - **N√©gatif :**
        - Moins adapt√© aux t√¢ches syst√me
        - Mono thread
        - Moins performant que Go
        - Typage plus souple

- **Go**
    - **Positif :**
        - Bonne gestion des t√¢ches syst√me
        - Concurrence native (goroutines)
        - Compilation rapide, binaire l√©ger
        - Tr√s utilis√© en production (Docker, Kubernetes)
    - **N√©gatif :**
        - Syntaxe verbeuse
        - Moins d'abstractions que Node.js

- **Rust**
    - **Positif :**
        - Performance et s√©curit√© maximale
        - Z√©ro bug m√©moire possible
        - Fort typage et contr√¥le strict
    - **N√©gatif :**
        - Moins d'outils cl√©-en-main
        - √cosyst√me back-end encore jeune
        - Courbe d'apprentissage raide

Go est le compromis entre performances et simplicit√©

#### M√©thodologie de d√©veloppement :
- M√©thode Agile Scrum (backlog, sprints, r√©unions, r√©trospectives)

### III - Conception et design

#### Diagramme de Gantt : planning d√©taill√© semaine par semaine

| **Semaine** | **T√¢che** | **D√©tails** |
| --- | --- | --- |
| **S1** | Analyse et pr√©paration | - √tude du cahier des charges<br>- Identification des personas<br>- Benchmark<br>- Pr√©paration de l'environnement de dev (Docker, repos GitHub) |
| **S1** | Architecture technique | - Sch√©ma global de l'architecture<br>- Diagramme Docker<br>- Diagramme de s√©quence pour MVP<br>- MCD et MLD des bases MySQL/PostgreSQL |
| **S2** | Backend Go  Partie MVP | - Gestion des bases (ajout, suppression)<br>- Sauvegarde manuelle / restauration<br>- API REST pour le frontend<br>- Versioning des sauvegardes |
| **S2** | Tests backend | - Tests unitaires des fonctions cl√©s<br>- Tests de sauvegarde/restauration avec bases fictives |
| **S3** | Frontend Vue.js | - Interface de contr√¥le basique (MVP)<br>- Formulaire ajout base, boutons sauvegarde/restauration<br>- Affichage de l'historique des sauvegardes<br>- Notifications simples |
| **S3** | Int√©gration frontend  backend | - Connexion API REST<br>- Tests manuels des fonctionnalit√©s principales |
| **S3** | Automatisation | - CRON pour sauvegardes automatiques<br>- Logs et notifications |
| **S4** | Tests finaux et debug | - Tests fonctionnels complets<br>- V√©rification multi-SGBD (MySQL/Postgres)<br>- Ajustements et corrections |
| **S4** | Documentation et conception finale | - Wireframes et maquette finalis√©s<br>- Diagrammes UML et Docker finalis√©s<br>- Cahier des charges compl√©t√©<br>- R√©daction des variables d'environnement et s√©curit√© |
| **S4** | Conteneurisation finale | - Docker Compose final (frontend + backend + MySQL + Postgres)<br>- Tests dans les conteneurs<br>- Pr√©paration du d√©p√¥t GitHub public |
| **S4** | Bilan et pr√©sentation | - Bilan du projet et perspectives d'√©volution<br>- Bilan personnel<br>- Pr√©paration soutenance / d√©mo |

#### Tableau Kanban : colonnes √ faire, En cours, Termin√© pour gestion des t√¢ches
GitHub Project

#### Arborescence du projet : (dossiers backend, frontend, docker, tests)

```
safebase/
Ç backend/                     # API REST en Go
Ç    cmd/                     # Point d'entr√©e (main.go)
Ç   Ç    main.go
Ç    internal/                # Logique interne (clean architecture)
Ç   Ç    db/                  # Connexions √ la DB (MySQL, Postgres, interne á contient du code)
Ç   Ç   Ç    mysql.go
Ç   Ç   Ç    postgres.go
Ç   Ç   Ç    safebase.go      # Base interne (historique, logs)
Ç   Ç    services/            # Logique m√©tier (sauvegarde, restauration)
Ç   Ç   Ç    backup.go
Ç   Ç   Ç    restore.go
Ç   Ç   Ç    scheduler.go     # Gestion CRON
Ç   Ç    api/                 # Handlers HTTP (endpoints API REST)
Ç   Ç   Ç    backup_handler.go
Ç   Ç   Ç    restore_handler.go
Ç   Ç   Ç    db_handler.go
Ç   Ç    utils/               # Fonctions utilitaires (logs, s√©curit√©)
Ç   Ç        logger.go
Ç   Ç        security.go
Ç    tests/                   # Tests unitaires backend
Ç   Ç    backup_test.go
Ç   Ç    restore_test.go
Ç    go.mod / go.sum
Ç
Ç frontend/                    # Interface utilisateur Vue.js
Ç    public/                  # Ressources statiques
Ç    src/
Ç   Ç    assets/              # Images, ic√¥nes
Ç   Ç    components/          # Composants Vue.js
Ç   Ç   Ç    BackupList.vue
Ç   Ç   Ç    RestoreForm.vue
Ç   Ç   Ç    AlertBox.vue
Ç   Ç    pages/               # Pages principales
Ç   Ç   Ç    Dashboard.vue
Ç   Ç   Ç    History.vue
Ç   Ç    router/              # Routes Vue.js
Ç   Ç   Ç    index.js
Ç   Ç    store/               # State management (Pinia/Vuex)
Ç   Ç   Ç    backup.js
Ç   Ç    App.vue
Ç    package.json
Ç
Ç docker/                      # Configurations Docker
Ç    Dockerfile.backend
Ç    Dockerfile.frontend
Ç    mysql_init.sql           # Script init MySQL (tests)
Ç    postgres_init.sql        # Script init PostgreSQL (tests)
Ç    docker-compose.yml
Ç
Ç db/                          # Sauvegardes et historique (contient les donn√©es)
Ç    backups/                 # Dossiers persistants de backup
Ç   Ç    mysql/
Ç   Ç    postgresql/
Ç    safebase.db              # SQLite interne (ou Postgres interne)
Ç
Ç tests/                       # Tests d'int√©gration & e2e
Ç    integration/             # Tests API (Go + DB factices)
Ç    e2e/                     # Tests bout en bout (frontend + backend)
Ç
Ç docs/                        # Documentation du projet
Ç    cahier_des_charges.md
Ç    conception.md
Ç    architecture.png
Ç    gantt.png
Ç
Ç .env                         # Variables d'environnement
Ç .gitignore
Ç README.md                    # Pr√©sentation GitHub
```

#### Wireframes : croquis des pages principales
#### Maquettes : version graphique plus aboutie
#### Charte graphique : couleurs, typographies, boutons, ic√¥nes
#### Responsive : adaptation aux √©crans mobile, tablette, desktop
#### Accessibilit√© et ergonomie : navigation claire, contraste, labels, boutons accessibles

#### Glossaire : d√©finitions techniques importantes (API REST, Docker, ORM, CRON, backup, etc.)

- **API REST**
    Application Programming Interface **(Representational State Transfer)** est un moyen pour une application (frontend) de communiquer avec un serveur (backend) via protocol HTTP. Les donn√©es sont g√©n√©ralement √©chang√©es au format JSON.
    
- **DOCKER**
    Technologie de conteneurisation qui permet de cr√©er des environnements isol√©s pour les applications.
    
- **ORM :**
    Object-Relational-Mapping, outil qui permet de **manipuler une base de donn√©es relationnelle** via des objets dans le code, plut√¥t que d'√©crire directement des requ√tes SQL.
    
- **CRON job :**
    T√¢che automatique ****planifi√©e qui s'ex√©cute √ intervalles r√©guliers sur un syst√me Unix/Linux.
    
- **Dump :**
    Copie compl√te d'une base de donn√©es dans un fichier
    
- **Pointer :**
- **Channels :**
- **Goroutines :**

### IV - Architecture technique

#### Diagramme d'architecture logique : flux backend  frontend  bases (internes et externes)

```
[Utilisateur]
 Ç
ñº
[Frontend Vue.js] > [Backend Go (API REST)]
 Ç                                             Ç
 Ç                                            > [Base MySQL] (sauvegarde)
 Ç                                            > [Base PostgreSQL] (sauvegarde)
 Ç                                            > [SafeBase DB interne] (historique, logs)
 Ç
> [Docker Compose] orchestre tous les services
```

#### Diagramme de conteneur Docker : conteneurs backend, frontend, bases internes et clientes, volumes, r√©seau

```
+---------------- Docker Network safebase_net ----------------+
|                                                             |
|  [frontend]  Vue.js  <-->  [backend] Go API                |
|                               Ç                             |
|                              ñº                             |
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

#### Diagramme de s√©quence : interactions entre utilisateur, frontend, backend et bases pour sauvegarde/restauration

```
Utilisateur       Frontend(Vue)      Backend(Go API)      DB cible (MySQL/Postgres)
Ç                 Ç                  Ç                    Ç
Ç   Clique "Backup" Ç                Ç                    Ç
->Ç                  Ç                    Ç
Ç                 Ç   POST /backup   Ç                    Ç
Ç                ->Ç                    Ç
Ç                 Ç                  Ç  dump DB (mysqldump)   Ç
Ç                 Ç                 ->Ç
Ç                 Ç                  Ç                    Ç
Ç                 Ç                  Ç   Retour OK + fichier   Ç
Ç                 Ç                  Ç
Ç                 Ç   R√©ponse JSON   Ç                    Ç
Ç                 Ç-                    Ç
Ç   Notification OK Ç                  Ç                    Ç
Ç-                  Ç                    Ç
```

#### Description des technologies

- **Back-end :**
    - GO
- **Front-end :**
    - Vue.js
    - Tailwind
- **Conteneurisation :**
    - Docker
- **S√©curit√© :**
    - JWT (authentification)
    - Bcrypt
- **Base de donn√©es :**
    - **Base de donn√©es du projet :**
        - PostgreSQL
        - Interface : pgAdmin
        - ORM GORM
    - **Base de donn√©es pouvant √tre sauvegard√©es :**
        - MySQL
        - PostgreSQL
        - SQLite (pour plus tard)
        - NoSQL (pour plus tard)
- **Tests :**
    - Go (int√©gr√©)
    - PostMan
    - Docker (int√©gr√©)
    - Cypress
- **Documentation et suivi de projet :**
    - Notion et README.md pour la documentation
    - GitHub Project pour les t√¢ches √ r√©aliser
    - Diagrammes UML

#### Variables d'environnement et configuration

### V - Base de donn√©es (si utilis√©e)

- MCD (Mod√le Conceptuel de Donn√©es)
- MLD (Mod√le Logique de Donn√©es)
- Choix de l'ORM et justification
    - GORM

### VI - Processus de sauvegarde / restauration

#### Diagramme de processus / flowchart :
- S√©lection de la base
- Sauvegarde manuelle ou automatique
- Stockage des fichiers
- Historique / versioning
- Restauration
- Logs et alertes

#### Description textuelle d√©taill√©e pour chaque √©tape

### VII - S√©curit√©

- Gestion des mots de passe et informations sensibles
- Permissions sur fichiers de backup
- Protection des endpoints de l'API
- Sauvegardes s√©curis√©es et fiables
- Politique de gestion des erreurs

### VIII - Tests

#### Plan de tests d√©taill√© :
- Tests unitaires (Go backend)
- Tests fonctionnels (sauvegarde/restauration sur bases fictives)
- Tests frontend (Vue.js interaction API)
- Tests d'int√©gration Docker (stack complet)

#### Tableau des sc√©narios :
| Test | Description | R√©sultat attendu |

#### Strat√©gie de tests align√©e sur **cycle en V** : conception  impl√©mentation  tests

### IX - D√©ploiement et production

- Docker Compose pour l'ensemble du projet
- Scripts de d√©marrage / mise √ jour
- D√©ploiement sur serveur / cloud
- Gestion des volumes persistants
- Instructions pour production
- Int√©gration CI/CD possible (GitHub Actions)

### X - M√©thodologie et organisation

- **Scrum / Agile** : backlog, sprints, daily meetings, r√©trospectives
- Gestion de projet professionnelle : GitHub Projects et Notion
- Suivi des t√¢ches et priorisation avec MoSCoW

### XI - Probl√mes rencontr√©s

- Difficult√©s techniques et solutions
- Limitations rencontr√©es
- Gestion des erreurs et cas exceptionnels

### XII - Bilan du projet et perspectives d'√©volution

- Fonctionnalit√©s r√©alis√©es / non r√©alis√©es
- √volutions possibles : support NoSQL, notifications avanc√©es, interface plus compl√te, sauvegarde crypt√©e, am√©lioration UX

### XIII - Bilan personnel

- Comp√©tences acquises (techniques, DevOps, gestion de projet)
- Points d'am√©lioration et apprentissages
- Exp√©rience professionnelle simul√©e (cycle V / Scrum, planning, tests)

### XIV - Conclusion

- R√©sum√© global du projet, objectifs atteints, valeur pour l'utilisateur et l'entreprise

### Annexes

- Fichiers docker-compose.yml, scripts backend et frontend
- Screenshots des wireframes et maquettes
- Diagrammes complets (architecture, s√©quence, conteneur, processus)
- Exemple de plan de tests complet
