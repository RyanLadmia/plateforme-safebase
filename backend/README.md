# SafeBase baackend

## Stacks :
Go and Gin
PostgreSQl and pgAdmin

Pour utiliser le terminal, toutjours être dans le dossier backend !!

## Initialisation :
Go mod:
`bash
go mod init github/UserName/ProjectName
`
Verify go verison :
`bash
go version
`

Install Gin :
`bash
go get github.com/gin-gonic/gin
`

To use .env variables and godotenv :
`bash
go get github.com/joho/godotenv
`

Install Air :
`bash
go install github.com/air-verse/air@latest
`

Verify air install (computer terminal) :
`
ls $(go env GOPATH)/bin | grep air
`

then (computer terminal)
`
which air
`

if it return nothing (computer terminal):
`bash
export PATH=$PATH:$(go env GOPATH)/bin
`

Verify air version (computer terminal) :
`bash
air -v
`

Initialize air :
`bash
air init
`
Crée un fichier de config .air.toml à la racine de backend

Installer robfig/cron/v3 :
`bash
go get github.com/robfig/cron/v3
`

## Lancer le serveur :
`bash
air
`
On n'utilise plus go run main.go