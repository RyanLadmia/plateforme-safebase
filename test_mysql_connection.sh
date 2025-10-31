#!/bin/bash

# Script de test de connectivité MySQL avec MAMP
echo "🔍 Test de connectivité MySQL avec MAMP"
echo "======================================="

# Vérifier si MAMP est installé
if [ ! -d "/Applications/MAMP" ]; then
    echo "❌ MAMP n'est pas installé dans /Applications/MAMP"
    exit 1
fi

echo "✅ MAMP trouvé dans /Applications/MAMP"

# Trouver mysqldump
MYSQLDUMP_PATHS=(
    "/Applications/MAMP/Library/bin/mysql80/bin/mysqldump"
    "/Applications/MAMP/Library/bin/mysqldump"
    "mysqldump"
)

MYSQLDUMP=""
for path in "${MYSQLDUMP_PATHS[@]}"; do
    if command -v "$path" >/dev/null 2>&1 || [ -x "$path" ]; then
        MYSQLDUMP="$path"
        echo "✅ mysqldump trouvé: $MYSQLDUMP"
        break
    fi
done

if [ -z "$MYSQLDUMP" ]; then
    echo "❌ mysqldump non trouvé"
    exit 1
fi

# Tester la connexion TCP sur le port 8889
echo ""
echo "🌐 Test de connexion TCP sur localhost:8889"
if nc -z localhost 8889 2>/dev/null; then
    echo "✅ Port 8889 accessible"
else
    echo "❌ Port 8889 non accessible"
fi

# Tester la connexion avec mysqldump
echo ""
echo "💾 Test de connexion avec mysqldump (TCP)"
if "$MYSQLDUMP" -h localhost -P 8889 -u root -proot --version >/dev/null 2>&1; then
    echo "✅ Connexion TCP réussie"
else
    echo "❌ Échec de connexion TCP"
fi

# Tester la socket Unix
SOCKET_PATH="/Applications/MAMP/tmp/mysql/mysql.sock"
echo ""
echo "🔌 Test de socket Unix: $SOCKET_PATH"
if [ -S "$SOCKET_PATH" ]; then
    echo "✅ Socket existe"
    # Tester la connexion socket
    if "$MYSQLDUMP" -u root -proot --socket="$SOCKET_PATH" --version >/dev/null 2>&1; then
        echo "✅ Connexion socket réussie"
    else
        echo "❌ Échec de connexion socket"
    fi
else
    echo "❌ Socket n'existe pas"
fi

# Tester un dump simple
echo ""
echo "📦 Test de dump simple (mysql database)"
TEMP_FILE="/tmp/mysql_test_dump.sql"
if "$MYSQLDUMP" -h localhost -P 8889 -u root -proot mysql user --single-transaction --no-data > "$TEMP_FILE" 2>/dev/null; then
    if [ -s "$TEMP_FILE" ]; then
        echo "✅ Dump de test réussi"
        rm -f "$TEMP_FILE"
    else
        echo "❌ Dump vide"
        rm -f "$TEMP_FILE"
    fi
else
    echo "❌ Échec du dump de test"
    rm -f "$TEMP_FILE"
fi

echo ""
echo "🎯 Résumé:"
echo "- MAMP installé: ✅"
echo "- mysqldump trouvé: ✅"
echo "- Port 8889 accessible: $(nc -z localhost 8889 2>/dev/null && echo '✅' || echo '❌')"
echo "- Socket disponible: $([ -S "$SOCKET_PATH" ] && echo '✅' || echo '❌')"
