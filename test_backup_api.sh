#!/bin/bash

# Script de test pour l'API de sauvegarde
# Ce script teste les fonctionnalités de base de l'API

BASE_URL="http://localhost:3000"
EMAIL="test@example.com"
PASSWORD="testpassword123"

echo "🚀 Test de l'API de sauvegarde SafeBase"
echo "======================================="

# Fonction pour afficher les réponses JSON de manière lisible
pretty_json() {
    if command -v jq &> /dev/null; then
        echo "$1" | jq .
    else
        echo "$1"
    fi
}

# 1. Inscription d'un utilisateur de test
echo "📝 1. Inscription d'un utilisateur de test..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\",
    \"role\": \"user\"
  }")

echo "Réponse d'inscription:"
pretty_json "$REGISTER_RESPONSE"
echo ""

# 2. Connexion pour obtenir le token
echo "🔐 2. Connexion pour obtenir le token JWT..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\"
  }")

echo "Réponse de connexion:"
pretty_json "$LOGIN_RESPONSE"

# Extraire le token (si jq est disponible)
if command -v jq &> /dev/null; then
    TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token // empty')
else
    # Extraction basique sans jq
    TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
fi

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
    echo "❌ Erreur: Impossible d'obtenir le token JWT"
    exit 1
fi

echo "✅ Token obtenu: ${TOKEN:0:20}..."
echo ""

# 3. Ajouter une base de données MySQL de test (MAMP)
echo "🗄️ 3. Ajout d'une base de données MySQL de test (MAMP)..."
DB_RESPONSE=$(curl -s -X POST "$BASE_URL/api/databases" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Base de test MySQL MAMP",
    "type": "mysql",
    "host": "localhost",
    "port": "8889",
    "username": "root",
    "password": "root",
    "db_name": "mysql"
  }')

echo "Réponse d'ajout de base de données:"
pretty_json "$DB_RESPONSE"

# Extraire l'ID de la base de données
if command -v jq &> /dev/null; then
    DB_ID=$(echo "$DB_RESPONSE" | jq -r '.database.id // empty')
else
    DB_ID=$(echo "$DB_RESPONSE" | grep -o '"id":[0-9]*' | cut -d':' -f2)
fi

if [ -z "$DB_ID" ] || [ "$DB_ID" = "null" ]; then
    echo "❌ Erreur: Impossible d'obtenir l'ID de la base de données"
    exit 1
fi

echo "✅ Base de données créée avec l'ID: $DB_ID"
echo ""

# 4. Lister les bases de données
echo "📋 4. Liste des bases de données..."
LIST_DB_RESPONSE=$(curl -s -X GET "$BASE_URL/api/databases" \
  -H "Authorization: Bearer $TOKEN")

echo "Liste des bases de données:"
pretty_json "$LIST_DB_RESPONSE"
echo ""

# 5. Créer une sauvegarde MySQL (MAMP)
echo "💾 5. Création d'une sauvegarde..."
BACKUP_RESPONSE=$(curl -s -X POST "$BASE_URL/api/backups/database/$DB_ID" \
  -H "Authorization: Bearer $TOKEN")

echo "Réponse de création de sauvegarde:"
pretty_json "$BACKUP_RESPONSE"

# Extraire l'ID de la sauvegarde
if command -v jq &> /dev/null; then
    BACKUP_ID=$(echo "$BACKUP_RESPONSE" | jq -r '.backup.id // empty')
else
    BACKUP_ID=$(echo "$BACKUP_RESPONSE" | grep -o '"id":[0-9]*' | cut -d':' -f2)
fi

echo ""

# 6. Lister les sauvegardes
echo "📋 6. Liste des sauvegardes..."
LIST_BACKUP_RESPONSE=$(curl -s -X GET "$BASE_URL/api/backups" \
  -H "Authorization: Bearer $TOKEN")

echo "Liste des sauvegardes:"
pretty_json "$LIST_BACKUP_RESPONSE"
echo ""

# 7. Si une sauvegarde a été créée, obtenir ses détails
if [ ! -z "$BACKUP_ID" ] && [ "$BACKUP_ID" != "null" ]; then
    echo "🔍 7. Détails de la sauvegarde $BACKUP_ID..."
    BACKUP_DETAIL_RESPONSE=$(curl -s -X GET "$BASE_URL/api/backups/$BACKUP_ID" \
      -H "Authorization: Bearer $TOKEN")
    
    echo "Détails de la sauvegarde:"
    pretty_json "$BACKUP_DETAIL_RESPONSE"
    echo ""
fi

# 8. Test de l'endpoint de téléchargement (même si le fichier n'existe pas)
if [ ! -z "$BACKUP_ID" ] && [ "$BACKUP_ID" != "null" ]; then
    echo "⬇️ 8. Test de téléchargement de la sauvegarde..."
    DOWNLOAD_RESPONSE=$(curl -s -w "HTTP_CODE:%{http_code}" -X GET "$BASE_URL/api/backups/$BACKUP_ID/download" \
      -H "Authorization: Bearer $TOKEN")
    
    HTTP_CODE=$(echo "$DOWNLOAD_RESPONSE" | grep -o "HTTP_CODE:[0-9]*" | cut -d':' -f2)
    echo "Code de réponse du téléchargement: $HTTP_CODE"
    echo ""
fi

echo "✅ Test terminé!"
echo ""
echo "📝 Notes:"
echo "- Les sauvegardes MySQL utilisent maintenant les chemins MAMP (/Applications/MAMP/Library/bin/mysql80/bin/mysqldump)"
echo "- Assurez-vous que MAMP est démarré et que MySQL fonctionne sur le port 8889"
echo "- Les paramètres par défaut de MAMP sont: host=localhost, port=8889, user=root, password=root"
echo "- Les fichiers de sauvegarde sont stockés dans db/backups/mysql/"
echo "- Vérifiez les logs du serveur pour plus de détails sur les erreurs"
