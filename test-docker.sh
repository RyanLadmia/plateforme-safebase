#!/bin/bash

# Script de démarrage et test pour environnement Docker
# Usage: ./test-docker.sh

echo "Tests E2E SafeBase - Environnement Docker"
echo "=============================================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
DOCKER_COMPOSE_FILE="docker-compose.yml"
BACKEND_PORT=8080
FRONTEND_PORT=3000

echo -e "${YELLOW}Étape 1/6: Vérification de Docker...${NC}"
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Docker n'est pas installé${NC}"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}Docker Compose n'est pas installé${NC}"
    exit 1
fi

echo -e "${GREEN}Docker et Docker Compose détectés${NC}"
echo ""

echo -e "${YELLOW}Étape 2/6: Arrêt des conteneurs existants...${NC}"
docker-compose down
echo -e "${GREEN}Conteneurs arrêtés${NC}"
echo ""

echo -e "${YELLOW}Étape 3/6: Démarrage des services Docker...${NC}"
echo "   Démarrage de: PostgreSQL, MySQL, Backend, Frontend..."
docker-compose up -d

if [ $? -ne 0 ]; then
    echo -e "${RED}Échec du démarrage des conteneurs${NC}"
    exit 1
fi

echo -e "${GREEN}Conteneurs démarrés${NC}"
echo ""

echo -e "${YELLOW}Étape 4/6: Attente de la disponibilité des services...${NC}"

# Attendre le backend
echo "   Vérification du backend (http://localhost:$BACKEND_PORT)..."
MAX_ATTEMPTS=30
ATTEMPT=0

while [ $ATTEMPT -lt $MAX_ATTEMPTS ]; do
    if curl -s http://localhost:$BACKEND_PORT/api > /dev/null 2>&1; then
        echo -e "${GREEN}Backend prêt${NC}"
        break
    fi
    ATTEMPT=$((ATTEMPT+1))
    echo "   Tentative $ATTEMPT/$MAX_ATTEMPTS..."
    sleep 2
done

if [ $ATTEMPT -eq $MAX_ATTEMPTS ]; then
    echo -e "${RED}Timeout: Backend non accessible${NC}"
    echo "   Vérifiez les logs: docker-compose logs backend"
    exit 1
fi

# Attendre le frontend
echo "   Vérification du frontend (http://localhost:$FRONTEND_PORT)..."
ATTEMPT=0

while [ $ATTEMPT -lt $MAX_ATTEMPTS ]; do
    if curl -s http://localhost:$FRONTEND_PORT > /dev/null 2>&1; then
        echo -e "${GREEN}Frontend prêt${NC}"
        break
    fi
    ATTEMPT=$((ATTEMPT+1))
    echo "   Tentative $ATTEMPT/$MAX_ATTEMPTS..."
    sleep 2
done

if [ $ATTEMPT -eq $MAX_ATTEMPTS ]; then
    echo -e "${RED}Timeout: Frontend non accessible${NC}"
    echo "   Vérifiez les logs: docker-compose logs frontend"
    exit 1
fi

echo ""

echo -e "${YELLOW}Étape 5/6: Configuration de Cypress pour Docker...${NC}"
cd tests

# Créer le fichier .env pour Docker si nécessaire
if [ ! -f ".env" ]; then
    cat > .env << EOF
CYPRESS_BASE_URL=http://localhost:$FRONTEND_PORT
CYPRESS_API_URL=http://localhost:$BACKEND_PORT/api
CYPRESS_IS_DOCKER=true
EOF
    echo -e "${GREEN}Fichier .env créé${NC}"
else
    echo -e "${BLUE}Fichier .env existant conservé${NC}"
fi

# Vérifier que node_modules existe
if [ ! -d "node_modules" ]; then
    echo "   Installation des dépendances Cypress..."
    npm install --loglevel=error
    if [ $? -ne 0 ]; then
        echo -e "${RED}Échec de l'installation des dépendances${NC}"
        cd ..
        exit 1
    fi
fi

echo -e "${GREEN}Cypress configuré pour Docker${NC}"
echo ""

echo -e "${YELLOW}Étape 6/6: Lancement des tests E2E...${NC}"
echo ""
echo -e "${BLUE}Vous pouvez maintenant lancer les tests:${NC}"
echo ""
echo "  • Mode interactif (GUI):   npm run cy:open"
echo "  • Mode headless:           npm run test"
echo "  • Test spécifique:         npx cypress run --spec \"e2E/01-authentication.cy.ts\""
echo ""
echo -e "${GREEN}Environnement Docker prêt pour les tests E2E !${NC}"
echo ""
echo "État des services:"
docker-compose ps
echo ""
echo "Logs utiles:"
echo "  • Backend:    docker-compose logs -f backend"
echo "  • Frontend:   docker-compose logs -f frontend"
echo "  • PostgreSQL: docker-compose logs -f postgres"
echo "  • Tous:       docker-compose logs -f"
echo ""
echo "Pour arrêter les services:"
echo "  docker-compose down"
