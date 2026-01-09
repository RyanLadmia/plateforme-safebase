# Stage 1: Récupérer pg_dump depuis Alpine 3.21 (pour PostgreSQL 17+ et SNI support)
FROM alpine:3.21 AS postgres-tools
RUN apk add --no-cache postgresql-client

# Stage 2: Récupérer mysqldump depuis Alpine 3.12 (pour MariaDB 10.4, compatible avec MariaDB 5.5+)
FROM alpine:3.12 AS mysql-tools
RUN apk add --no-cache mariadb-client

# Stage 3: Image principale avec Alpine 3.20
FROM alpine:3.20

WORKDIR /app

# Copier pg_dump et ses dépendances depuis Alpine 3.21
COPY --from=postgres-tools /usr/bin/pg_dump /usr/bin/pg_dump
COPY --from=postgres-tools /usr/bin/pg_restore /usr/bin/pg_restore
COPY --from=postgres-tools /usr/bin/psql /usr/bin/psql
COPY --from=postgres-tools /usr/lib/libpq.so.5* /usr/lib/
COPY --from=postgres-tools /usr/lib/libssl.so* /usr/lib/
COPY --from=postgres-tools /usr/lib/libcrypto.so* /usr/lib/
COPY --from=postgres-tools /usr/lib/libzstd.so* /usr/lib/
COPY --from=postgres-tools /usr/lib/liblz4.so* /usr/lib/

# Copier mysqldump et ses dépendances depuis Alpine 3.12
COPY --from=mysql-tools /usr/bin/mysqldump /usr/bin/mysqldump
COPY --from=mysql-tools /usr/bin/mariadb-dump /usr/bin/mariadb-dump
COPY --from=mysql-tools /usr/bin/mysql /usr/bin/mysql
COPY --from=mysql-tools /usr/lib/libmariadb.so* /usr/lib/
COPY --from=mysql-tools /usr/lib/libssl.so.1.1 /usr/lib/
COPY --from=mysql-tools /usr/lib/libcrypto.so.1.1 /usr/lib/

# Installer les dépendances de base
RUN apk add --no-cache \
    ca-certificates \
    git \
    wget \
    bash \
    zip \
    unzip \
    libstdc++ \
    libgcc

# Installer Go 1.25 manuellement
RUN wget https://go.dev/dl/go1.25.0.linux-arm64.tar.gz && \
    tar -C /usr/local -xzf go1.25.0.linux-arm64.tar.gz && \
    rm go1.25.0.linux-arm64.tar.gz

# Ajouter Go au PATH
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"
ENV PATH="${GOPATH}/bin:${PATH}"

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
