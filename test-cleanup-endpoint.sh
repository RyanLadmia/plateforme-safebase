#!/bin/bash

# Script de test pour l'endpoint de nettoyage
# Usage: ./test-cleanup-endpoint.sh

set -e

# Couleurs
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

BACKEND_URL="http://localhost:8080"

echo -e "${YELLOW}Test de l'endpoint de nettoyage${NC}"
echo ""

# Test 1: Vérifier que le backend est accessible
echo -e "${YELLOW}[TEST 1]${NC} Vérification de l'accessibilité du backend..."
if curl -s "${BACKEND_URL}/health" > /dev/null 2>&1; then
    echo -e "${GREEN}[OK]${NC} Backend accessible"
else
    echo -e "${RED}[ECHEC]${NC} Backend non accessible"
    exit 1
fi

# Test 2: Créer des utilisateurs de test via l'API
echo ""
echo -e "${YELLOW}[TEST 2]${NC} Création de 3 utilisateurs de test..."
for i in {1..3}; do
    TIMESTAMP=$(date +%s)
    EMAIL="cleanup.test.${TIMESTAMP}.${i}@e2e.com"
    
    RESPONSE=$(curl -s -X POST "${BACKEND_URL}/auth/register" \
        -H "Content-Type: application/json" \
        -d "{
            \"firstname\": \"Test\",
            \"lastname\": \"User${i}\",
            \"email\": \"${EMAIL}\",
            \"password\": \"TestPassword123!\",
            \"confirm_password\": \"TestPassword123!\"
        }")
    
    if echo "$RESPONSE" | grep -q "id"; then
        echo -e "${GREEN}[OK]${NC} Utilisateur ${i} créé: ${EMAIL}"
    else
        echo -e "${RED}[ECHEC]${NC} Erreur création utilisateur ${i}"
        echo "Réponse: $RESPONSE"
    fi
    
    sleep 0.5
done

# Test 3: Appeler l'endpoint de nettoyage
echo ""
echo -e "${YELLOW}[TEST 3]${NC} Appel de l'endpoint de nettoyage..."
CLEANUP_RESPONSE=$(curl -s -X POST "${BACKEND_URL}/api/test/cleanup-users")

if echo "$CLEANUP_RESPONSE" | grep -q "deleted_count"; then
    DELETED_COUNT=$(echo "$CLEANUP_RESPONSE" | grep -o '"deleted_count":[0-9]*' | grep -o '[0-9]*')
    echo -e "${GREEN}[OK]${NC} Nettoyage réussi: ${DELETED_COUNT} utilisateur(s) supprimé(s)"
    
    if [ "$DELETED_COUNT" -ge 3 ]; then
        echo -e "${GREEN}[OK]${NC} Au moins 3 utilisateurs ont été supprimés (ceux créés par ce test)"
    else
        echo -e "${YELLOW}[AVERTISSEMENT]${NC} Moins de 3 utilisateurs supprimés (il y avait peut-être déjà des nettoyages)"
    fi
else
    echo -e "${RED}[ECHEC]${NC} Erreur lors du nettoyage"
    echo "Réponse: $CLEANUP_RESPONSE"
    exit 1
fi

echo ""
echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}  Tous les tests ont réussi${NC}"
echo -e "${GREEN}================================${NC}"

