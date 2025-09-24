# SafeBase baackend

## Stacks :
Go and Gin
PostgreSQl and pgAdmin

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