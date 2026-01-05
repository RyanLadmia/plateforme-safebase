#!/bin/bash

echo "Vérification de la configuration Docker..."
echo ""

# Couleurs
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Fonction pour vérifier un fichier
check_file() {
    if [ -f "$1" ]; then
        echo -e "${GREEN}[OK]${NC} $1 existe"
        return 0
    else
        echo -e "${RED}[KO]${NC} $1 manquant"
        return 1
    fi
}

# Fonction pour vérifier une version dans un fichier
check_version() {
    local file=$1
    local pattern=$2
    local expected=$3
    
    if grep -q "$pattern" "$file"; then
        echo -e "${GREEN}[OK]${NC} $file contient $expected"
        return 0
    else
        echo -e "${RED}[KO]${NC} $file ne contient pas $expected"
        return 1
    fi
}

echo "Vérification des fichiers Docker..."
check_file "backend/Dockerfile"
check_file "frontend/Dockerfile"
check_file "docker-compose.yml"
check_file "backend/.dockerignore"
check_file "frontend/.dockerignore"
check_file "backend/.air.toml"
check_file "frontend/vite.config.ts"

echo ""
echo "Vérification des versions..."
check_version "backend/Dockerfile" "golang:1.25" "Go 1.25"
check_version "backend/go.mod" "go 1.25" "Go 1.25 dans go.mod"
check_version "frontend/Dockerfile" "node:20" "Node.js 20"
check_version "backend/.air.toml" "poll = true" "Polling activé"
check_version "frontend/vite.config.ts" "usePolling: true" "Polling Vite activé"

echo ""
echo "Vérification des outils de sauvegarde..."
check_version "backend/Dockerfile" "postgresql-client" "PostgreSQL client (pg_dump)"
check_version "backend/Dockerfile" "mysql-client" "MySQL client (mysqldump)"
check_version "backend/Dockerfile" "zip" "Compression ZIP"

echo ""
echo "Vérification de la configuration..."
check_version "docker-compose.yml" "PORT=8080" "PORT backend"
check_version "docker-compose.yml" "JWT_SECRET" "JWT_SECRET"
check_version "docker-compose.yml" "DB_HOST=postgres" "DB_HOST"

echo ""
echo "Vérification de Docker..."
if command -v docker &> /dev/null; then
    echo -e "${GREEN}[OK]${NC} Docker est installé"
    docker --version
else
    echo -e "${RED}[KO]${NC} Docker n'est pas installé"
fi

if command -v docker-compose &> /dev/null; then
    echo -e "${GREEN}[OK]${NC} Docker Compose est installé"
    docker-compose --version
else
    echo -e "${RED}[KO]${NC} Docker Compose n'est pas installé"
fi

echo ""
echo "État des conteneurs Docker..."
if docker-compose ps 2>/dev/null; then
    echo ""
    echo -e "${GREEN}[OK]${NC} Docker Compose fonctionne"
else
    echo -e "${YELLOW}[WARN]${NC} Aucun conteneur en cours d'exécution"
fi

echo ""
echo "Vérification terminée !"
echo ""
echo "Pour démarrer le projet :"
echo "  ./docker-start.sh"
echo ""
echo "Pour voir les logs :"
echo "  docker-compose logs -f"
