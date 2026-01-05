#!/bin/bash

# ===================================
# Script de déploiement SafeBase
# ===================================

set -e  # Arrêter en cas d'erreur

# Couleurs pour les messages
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Fonction pour afficher les messages
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Vérifier que Docker est installé
if ! command -v docker &> /dev/null; then
    log_error "Docker n'est pas installé. Veuillez l'installer d'abord."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    log_error "Docker Compose n'est pas installé. Veuillez l'installer d'abord."
    exit 1
fi

# Vérifier que le fichier .env existe
if [ ! -f .env ]; then
    log_error "Le fichier .env n'existe pas. Copiez env.production.example vers .env et configurez-le."
    exit 1
fi

# Menu de sélection
echo "====================================="
echo "   Déploiement SafeBase"
echo "====================================="
echo ""
echo "Choisissez une action :"
echo "1) Build et déployer (développement)"
echo "2) Déployer depuis Docker Hub (production)"
echo "3) Arrêter tous les services"
echo "4) Voir les logs"
echo "5) Nettoyer les anciennes images"
echo "6) Backup des bases de données"
echo "7) Restaurer les bases de données"
echo "8) Quitter"
echo ""
read -p "Votre choix [1-8]: " choice

case $choice in
    1)
        log_info "Build et déploiement en mode développement..."
        
        # Sauvegarder .env
        if [ -f .env ]; then
            cp .env .env.backup
            log_success "Backup de .env créé"
        fi
        
        # Build des images
        log_info "Build des images Docker..."
        docker-compose build --no-cache
        
        # Arrêter les anciens conteneurs
        log_info "Arrêt des anciens conteneurs..."
        docker-compose down
        
        # Démarrer les nouveaux conteneurs
        log_info "Démarrage des nouveaux conteneurs..."
        docker-compose up -d
        
        # Attendre que les services soient prêts
        log_info "Attente du démarrage des services..."
        sleep 10
        
        # Vérifier l'état
        docker-compose ps
        
        log_success "Déploiement terminé!"
        log_info "Backend: http://localhost:8080"
        log_info "Frontend: http://localhost:3000"
        log_info "Grafana: http://localhost:3001"
        ;;
        
    2)
        log_info "Déploiement depuis Docker Hub (production)..."
        
        # Vérifier que DOCKER_USERNAME est défini
        if [ -z "$DOCKER_USERNAME" ]; then
            log_error "DOCKER_USERNAME n'est pas défini dans .env"
            exit 1
        fi
        
        # Sauvegarder .env
        if [ -f .env ]; then
            cp .env .env.backup
            log_success "Backup de .env créé"
        fi
        
        # Pull des dernières images
        log_info "Téléchargement des dernières images..."
        docker-compose -f docker-compose.prod.yml pull
        
        # Arrêter les anciens conteneurs
        log_info "Arrêt des anciens conteneurs..."
        docker-compose -f docker-compose.prod.yml down
        
        # Démarrer les nouveaux conteneurs
        log_info "Démarrage des nouveaux conteneurs..."
        docker-compose -f docker-compose.prod.yml up -d
        
        # Attendre que les services soient prêts
        log_info "Attente du démarrage des services..."
        sleep 10
        
        # Vérifier l'état
        docker-compose -f docker-compose.prod.yml ps
        
        # Nettoyer les anciennes images
        log_info "Nettoyage des anciennes images..."
        docker image prune -af
        
        log_success "Déploiement terminé!"
        log_info "Backend: http://localhost:8080"
        log_info "Frontend: http://localhost:80"
        ;;
        
    3)
        log_info "Arrêt de tous les services..."
        docker-compose down
        docker-compose -f docker-compose.prod.yml down 2>/dev/null || true
        log_success "Tous les services sont arrêtés"
        ;;
        
    4)
        echo ""
        echo "Choisissez le service :"
        echo "1) Tous les services"
        echo "2) Backend"
        echo "3) Frontend"
        echo "4) PostgreSQL"
        echo "5) MySQL"
        read -p "Votre choix [1-5]: " log_choice
        
        case $log_choice in
            1) docker-compose logs -f ;;
            2) docker-compose logs -f backend ;;
            3) docker-compose logs -f frontend ;;
            4) docker-compose logs -f postgres ;;
            5) docker-compose logs -f mysql ;;
            *) log_error "Choix invalide" ;;
        esac
        ;;
        
    5)
        log_info "Nettoyage des anciennes images..."
        docker system prune -af --volumes
        log_success "Nettoyage terminé"
        ;;
        
    6)
        log_info "Backup des bases de données..."
        
        # Créer le répertoire de backup
        BACKUP_DIR="./backups/manual/$(date +%Y%m%d_%H%M%S)"
        mkdir -p "$BACKUP_DIR"
        
        # Backup PostgreSQL
        log_info "Backup PostgreSQL..."
        docker-compose exec -T postgres pg_dump -U ${POSTGRES_USER:-safebase_user} ${POSTGRES_DB:-safebase} > "$BACKUP_DIR/postgres_backup.sql"
        
        # Backup MySQL
        log_info "Backup MySQL..."
        docker-compose exec -T mysql mysqldump -u ${MYSQL_USER:-safebase_user} -p${MYSQL_PASSWORD} ${MYSQL_DATABASE:-safebase} > "$BACKUP_DIR/mysql_backup.sql"
        
        # Compresser les backups
        log_info "Compression des backups..."
        tar -czf "$BACKUP_DIR.tar.gz" -C "$BACKUP_DIR" .
        rm -rf "$BACKUP_DIR"
        
        log_success "Backup créé: $BACKUP_DIR.tar.gz"
        ;;
        
    7)
        log_warning "ATTENTION: La restauration va écraser les données actuelles!"
        read -p "Êtes-vous sûr de vouloir continuer? (yes/no): " confirm
        
        if [ "$confirm" != "yes" ]; then
            log_info "Restauration annulée"
            exit 0
        fi
        
        read -p "Chemin du fichier de backup (.tar.gz): " backup_file
        
        if [ ! -f "$backup_file" ]; then
            log_error "Fichier de backup introuvable: $backup_file"
            exit 1
        fi
        
        # Extraire le backup
        RESTORE_DIR="./backups/restore_temp"
        mkdir -p "$RESTORE_DIR"
        tar -xzf "$backup_file" -C "$RESTORE_DIR"
        
        # Restaurer PostgreSQL
        if [ -f "$RESTORE_DIR/postgres_backup.sql" ]; then
            log_info "Restauration PostgreSQL..."
            docker-compose exec -T postgres psql -U ${POSTGRES_USER:-safebase_user} ${POSTGRES_DB:-safebase} < "$RESTORE_DIR/postgres_backup.sql"
            log_success "PostgreSQL restauré"
        fi
        
        # Restaurer MySQL
        if [ -f "$RESTORE_DIR/mysql_backup.sql" ]; then
            log_info "Restauration MySQL..."
            docker-compose exec -T mysql mysql -u ${MYSQL_USER:-safebase_user} -p${MYSQL_PASSWORD} ${MYSQL_DATABASE:-safebase} < "$RESTORE_DIR/mysql_backup.sql"
            log_success "MySQL restauré"
        fi
        
        # Nettoyer
        rm -rf "$RESTORE_DIR"
        
        log_success "Restauration terminée!"
        ;;
        
    8)
        log_info "Au revoir!"
        exit 0
        ;;
        
    *)
        log_error "Choix invalide"
        exit 1
        ;;
esac

echo ""
log_success "Opération terminée avec succès!"
