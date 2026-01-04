#!/bin/bash

echo "🔍 Test des outils de sauvegarde dans le conteneur backend..."
echo ""

# Couleurs
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Vérifier si le conteneur backend est en cours d'exécution
if ! docker-compose ps backend | grep -q "Up"; then
    echo -e "${RED}✗${NC} Le conteneur backend n'est pas en cours d'exécution"
    echo ""
    echo "Démarrez-le avec : docker-compose up -d backend"
    exit 1
fi

echo -e "${GREEN}✓${NC} Le conteneur backend est en cours d'exécution"
echo ""

# Fonction pour tester une commande dans le conteneur
test_command() {
    local cmd=$1
    local name=$2
    
    echo -e "${BLUE}Testing:${NC} $name"
    
    if docker-compose exec -T backend sh -c "command -v $cmd > /dev/null 2>&1"; then
        local version=$(docker-compose exec -T backend sh -c "$cmd --version 2>&1 | head -n 1")
        echo -e "${GREEN}✓${NC} $cmd trouvé"
        echo "  Version: $version"
    else
        echo -e "${RED}✗${NC} $cmd non trouvé"
        return 1
    fi
    echo ""
}

# Tester chaque outil
echo "🗄️  Outils de sauvegarde PostgreSQL:"
test_command "pg_dump" "PostgreSQL dump"
test_command "pg_restore" "PostgreSQL restore"

echo "🗄️  Outils de sauvegarde MySQL:"
test_command "mysqldump" "MySQL dump"
test_command "mysql" "MySQL client"

echo "📦 Outils de compression:"
test_command "zip" "ZIP compression"
test_command "unzip" "ZIP decompression"

echo ""
echo "📊 Résumé:"
echo ""

# Compter les outils installés
TOOLS_COUNT=0
TOOLS_MISSING=0

for cmd in pg_dump pg_restore mysqldump mysql zip unzip; do
    if docker-compose exec -T backend sh -c "command -v $cmd > /dev/null 2>&1"; then
        ((TOOLS_COUNT++))
    else
        ((TOOLS_MISSING++))
    fi
done

if [ $TOOLS_MISSING -eq 0 ]; then
    echo -e "${GREEN}✓${NC} Tous les outils de sauvegarde sont installés ($TOOLS_COUNT/6)"
    echo ""
    echo -e "${GREEN}🎉 Votre conteneur est prêt pour les sauvegardes !${NC}"
else
    echo -e "${RED}✗${NC} $TOOLS_MISSING outil(s) manquant(s)"
    echo ""
    echo -e "${YELLOW}⚠${NC}  Vous devez reconstruire le conteneur backend:"
    echo "  docker-compose down"
    echo "  docker-compose build backend"
    echo "  docker-compose up -d"
fi

echo ""
echo "📁 Vérification du répertoire de sauvegarde:"
if docker-compose exec -T backend sh -c "test -d /app/db/backups"; then
    echo -e "${GREEN}✓${NC} /app/db/backups existe"
    
    # Lister les sauvegardes existantes
    echo ""
    echo "📦 Sauvegardes PostgreSQL:"
    docker-compose exec -T backend sh -c "ls -lh /app/db/backups/postgresql/ 2>/dev/null | tail -n +2" || echo "  (aucune sauvegarde)"
    
    echo ""
    echo "📦 Sauvegardes MySQL:"
    docker-compose exec -T backend sh -c "ls -lh /app/db/backups/mysql/ 2>/dev/null | tail -n +2" || echo "  (aucune sauvegarde)"
else
    echo -e "${YELLOW}⚠${NC}  /app/db/backups n'existe pas encore (sera créé à la première sauvegarde)"
fi

echo ""

