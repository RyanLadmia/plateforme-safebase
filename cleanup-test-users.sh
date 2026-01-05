#!/bin/bash

# Script pour nettoyer les utilisateurs de test de la base de données
# Usage: ./cleanup-test-users.sh

set -e

# Couleurs pour l'affichage
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}==================================${NC}"
echo -e "${YELLOW}  Nettoyage des utilisateurs E2E${NC}"
echo -e "${YELLOW}==================================${NC}"

# Vérifier si le backend est en cours d'exécution
BACKEND_URL="http://localhost:8080"
if ! curl -s "${BACKEND_URL}/health" > /dev/null 2>&1; then
    echo -e "${RED}[ERREUR]${NC} Le backend ne répond pas sur ${BACKEND_URL}"
    echo -e "${YELLOW}[INFO]${NC} Assurez-vous que le backend est démarré avec Docker:"
    echo -e "        docker-compose up -d backend"
    exit 1
fi

echo -e "${GREEN}[OK]${NC} Backend accessible"

# Appeler l'endpoint de nettoyage
echo -e "${YELLOW}[INFO]${NC} Suppression des utilisateurs de test (@e2e.com)..."
RESPONSE=$(curl -s -X POST "${BACKEND_URL}/api/test/cleanup-users")

# Vérifier la réponse
if echo "$RESPONSE" | grep -q "deleted_count"; then
    DELETED_COUNT=$(echo "$RESPONSE" | grep -o '"deleted_count":[0-9]*' | grep -o '[0-9]*')
    echo -e "${GREEN}[OK]${NC} ${DELETED_COUNT} utilisateur(s) de test supprimé(s)"
else
    echo -e "${RED}[ERREUR]${NC} Échec du nettoyage"
    echo "Réponse: $RESPONSE"
    exit 1
fi

echo -e "${GREEN}==================================${NC}"
echo -e "${GREEN}  Nettoyage terminé avec succès${NC}"
echo -e "${GREEN}==================================${NC}"

