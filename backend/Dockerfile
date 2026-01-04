# Utiliser une image Go 1.25 pour supporter Air v1.63+
FROM golang:1.25-alpine

WORKDIR /app

# Installer les outils nécessaires pour les sauvegardes de bases de données
# - postgresql-client : contient pg_dump pour les sauvegardes PostgreSQL
# - mysql-client : contient mysqldump pour les sauvegardes MySQL
# - zip/unzip : pour la compression des sauvegardes
RUN apk add --no-cache \
    postgresql-client \
    mysql-client \
    zip \
    unzip

# Installer Air pour le hot-reload
RUN go install github.com/air-verse/air@latest

# Copier les fichiers de dépendances
COPY go.mod go.sum ./
RUN go mod download

# Copier tout le code source
COPY . .

# Exposer le port
EXPOSE 8080

# Utiliser Air pour le développement avec hot-reload
CMD ["air", "-c", ".air.toml"]