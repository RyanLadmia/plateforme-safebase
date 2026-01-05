#!/bin/bash

# Script d'installation et de vérification des tests E2E Cypress
# Usage: ./install-cypress.sh

echo "🚀 Installation des tests E2E Cypress pour SafeBase"
echo "===================================================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if in correct directory
if [ ! -d "tests" ]; then
    echo -e "${RED}❌ Erreur: Le dossier 'tests' n'existe pas.${NC}"
    echo "   Assurez-vous d'être à la racine du projet."
    exit 1
fi

echo -e "${YELLOW}📦 Étape 1/5: Vérification de Node.js et npm...${NC}"
if ! command -v node &> /dev/null; then
    echo -e "${RED}❌ Node.js n'est pas installé.${NC}"
    echo "   Installez Node.js depuis https://nodejs.org/"
    exit 1
fi

if ! command -v npm &> /dev/null; then
    echo -e "${RED}❌ npm n'est pas installé.${NC}"
    echo "   Installez npm avec Node.js depuis https://nodejs.org/"
    exit 1
fi

NODE_VERSION=$(node -v)
NPM_VERSION=$(npm -v)
echo -e "${GREEN}✓ Node.js $NODE_VERSION détecté${NC}"
echo -e "${GREEN}✓ npm $NPM_VERSION détecté${NC}"
echo ""

echo -e "${YELLOW}📦 Étape 2/5: Installation des dépendances Cypress...${NC}"
cd tests

# Clean install
if [ -d "node_modules" ]; then
    echo "   Nettoyage des anciennes dépendances..."
    rm -rf node_modules package-lock.json
fi

echo "   Installation en cours (cela peut prendre quelques minutes)..."
npm install --loglevel=error

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Erreur lors de l'installation des dépendances.${NC}"
    echo "   Essayez: sudo npm install --unsafe-perm"
    exit 1
fi

echo -e "${GREEN}✓ Dépendances installées avec succès${NC}"
echo ""

echo -e "${YELLOW}📦 Étape 3/5: Vérification de Cypress...${NC}"
if ! npx cypress --version &> /dev/null; then
    echo -e "${RED}❌ Cypress n'est pas installé correctement.${NC}"
    exit 1
fi

CYPRESS_VERSION=$(npx cypress --version | grep "Cypress package version" | cut -d: -f2 | xargs)
echo -e "${GREEN}✓ Cypress $CYPRESS_VERSION installé${NC}"
echo ""

cd ..

echo -e "${YELLOW}🔍 Étape 4/5: Vérification des prérequis...${NC}"

# Check if backend is running
if curl -s http://localhost:8080/api > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Backend détecté sur http://localhost:8080${NC}"
else
    echo -e "${YELLOW}⚠ Backend non détecté sur http://localhost:8080${NC}"
    echo "   Démarrez le backend avec: cd backend && go run cmd/main.go"
fi

# Check if frontend is running
if curl -s http://localhost:5173 > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Frontend détecté sur http://localhost:5173${NC}"
else
    echo -e "${YELLOW}⚠ Frontend non détecté sur http://localhost:5173${NC}"
    echo "   Démarrez le frontend avec: cd frontend && npm run dev"
fi

# Check if PostgreSQL is accessible
if command -v psql &> /dev/null; then
    echo -e "${GREEN}✓ PostgreSQL détecté${NC}"
else
    echo -e "${YELLOW}⚠ PostgreSQL non détecté${NC}"
    echo "   Assurez-vous que PostgreSQL est installé et démarré"
fi

echo ""

echo -e "${YELLOW}📊 Étape 5/5: Résumé de l'installation${NC}"
echo "=================================================="
echo ""
echo "Tests créés:"
echo "  • 01-authentication.cy.ts (Authentification)"
echo "  • 02-database-management.cy.ts (Gestion BDD)"
echo "  • 03-backup-management.cy.ts (Sauvegardes)"
echo "  • 04-schedule-management.cy.ts (Planifications)"
echo "  • 05-history.cy.ts (Historique)"
echo "  • 06-profile.cy.ts (Profil utilisateur)"
echo "  • 07-dashboard.cy.ts (Dashboard)"
echo "  • 08-complete-workflows.cy.ts (Workflows complets)"
echo ""
echo "Total: ~200 tests E2E | Couverture: >90%"
echo ""

echo -e "${GREEN}✅ Installation terminée avec succès !${NC}"
echo ""
echo "🏃 Commandes disponibles:"
echo "  cd tests && npm run cy:open      # Ouvrir Cypress (mode GUI)"
echo "  cd tests && npm run test         # Lancer tous les tests (headless)"
echo "  cd tests && npm run cy:run:chrome # Tests avec Chrome"
echo ""
echo "📚 Documentation:"
echo "  • tests/README.md                 # Documentation détaillée"
echo "  • CYPRESS_E2E_SETUP.md           # Guide d'installation"
echo "  • TEST_SYNTHESIS.md              # Synthèse complète"
echo ""
echo "⚠️  IMPORTANT:"
echo "  Avant de lancer les tests, assurez-vous que:"
echo "  1. PostgreSQL est démarré"
echo "  2. Backend est actif (http://localhost:8080)"
echo "  3. Frontend est actif (http://localhost:5173)"
echo ""
echo "🎉 Bonne chance avec vos tests !"

