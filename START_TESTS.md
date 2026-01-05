# 🚀 DÉMARRAGE RAPIDE - Tests avec Docker

## Pour lancer TOUS les tests avec Docker en 3 étapes :

### 1️⃣ Script automatique (LE PLUS SIMPLE)

```bash
./test-docker.sh
```

Ce script fait TOUT pour vous ! 🎉

### 2️⃣ Lancer les tests E2E

```bash
cd tests
npm run cy:open    # Interface graphique (développement)
# OU
npm run test       # Headless (automatique)
```

### 3️⃣ Tests Backend (optionnel)

```bash
cd backend
go test ./tests/... -v
```

---

## ⚡ Commandes ultra-rapides

```bash
# TOUT en une commande
./test-docker.sh && cd tests && npm run test

# Juste E2E
docker-compose up -d && cd tests && npm run cy:open

# Juste Backend
cd backend && go test ./tests/... -v
```

---

## ✅ Vérification rapide

```bash
# Docker tourne ?
docker-compose ps

# Services OK ?
curl http://localhost:8080/api   # Backend ✓
curl http://localhost:3000        # Frontend ✓

# Tout est vert ? Lancez les tests !
```

---

## 🐛 Problème ?

```bash
# Redémarrer Docker
docker-compose restart

# Voir les logs
docker-compose logs -f

# Nettoyer et recommencer
docker-compose down && docker-compose up -d
```

---

## 📚 Plus d'infos ?

- **Guide Docker complet** : `DOCKER_TESTS_GUIDE.md`
- **Mise à jour Docker** : `DOCKER_UPDATE.md`
- **Synthèse complète** : `TEST_SYNTHESIS.md`

---

## 🎯 L'essentiel

### Votre configuration :
- ✅ Docker Compose
- ✅ Frontend port **3000**
- ✅ Backend port **8080**
- ✅ ~253 tests (~200 E2E + ~53 Go)
- ✅ Couverture >90%

### La commande magique :
```bash
./test-docker.sh && cd tests && npm run cy:open
```

**C'est tout ! 🎉**

