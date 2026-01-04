#!/bin/bash

echo "🐳 Démarrage de SafeBase avec Docker..."
echo ""

# Arrêter les conteneurs existants
echo "📦 Arrêt des conteneurs existants..."
docker-compose down

# Reconstruire les images
echo ""
echo "🔨 Reconstruction des images..."
docker-compose build

# Démarrer tous les services
echo ""
echo "🚀 Démarrage des services..."
docker-compose up -d

# Attendre que les services démarrent
echo ""
echo "⏳ Attente du démarrage des services..."
sleep 5

# Afficher le statut
echo ""
echo "📊 Statut des services :"
docker-compose ps

echo ""
echo "✅ SafeBase est prêt !"
echo ""
echo "🌐 Services disponibles :"
echo "   - Frontend : http://localhost:3000"
echo "   - Backend  : http://localhost:8080"
echo "   - Grafana  : http://localhost:3001"
echo "   - Prometheus : http://localhost:9090"
echo ""
echo "📝 Pour voir les logs : docker-compose logs -f"
echo "🛑 Pour arrêter : docker-compose down"

