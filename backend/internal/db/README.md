# Database :

PostgreSQL 17
pgAdmin interface

## Installer GORM

le coeur de GORM :
`bash
go get -u gorm.io/gorm
`

le driver postgresql :
`bahs
go get -u gorm.io/driver/postgres
`

Vérifier l'installation :
`bash
go list -m all | grep gorm
`
-> réponse :
    gorm.io/driver/postgres v1.6.0
    gorm.io/driver/sqlite v1.6.0
    gorm.io/gorm v1.31.0