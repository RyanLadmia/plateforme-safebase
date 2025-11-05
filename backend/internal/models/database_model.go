package models

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// DatabaseCreateRequest is used for JSON binding during creation
type DatabaseCreateRequest struct {
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required,oneof=mysql postgresql"`
	Host     string `json:"host,omitempty"`
	Port     string `json:"port,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	DbName   string `json:"db_name,omitempty"`
	URL      string `json:"url,omitempty"` // Alternative: full connection URL
}

// DatabaseUpdateRequest is used for JSON binding during updates
type DatabaseUpdateRequest struct {
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required,oneof=mysql postgresql"`
	Host     string `json:"host,omitempty"`
	Port     string `json:"port,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	DbName   string `json:"db_name,omitempty"`
	URL      string `json:"url,omitempty"` // Alternative: full connection URL
}

type Database struct {
	Id        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Type      string         `gorm:"size:50;not null" json:"type"` // mysql, postgres
	Host      string         `gorm:"size:255;not null" json:"host"`
	Port      string         `gorm:"size:10;not null" json:"port"`
	Username  string         `gorm:"size:100;not null" json:"username"`
	Password  string         `gorm:"size:255;not null" json:"-"` // Ne pas exposer le mot de passe dans le JSON
	DbName    string         `gorm:"size:100;not null" json:"db_name"`
	URL       string         `gorm:"size:500" json:"url,omitempty"` // Full connection URL (optional)
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	UserId    uint           `gorm:"index;not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserId" json:"-"`
	Backups   []Backup       `gorm:"foreignKey:DatabaseId;constraint:OnDelete:CASCADE;" json:"backups,omitempty"`
	Restores  []Restore      `gorm:"foreignKey:DatabaseId;constraint:OnDelete:CASCADE;" json:"restores,omitempty"`
	Schedules []Schedule     `gorm:"foreignKey:DatabaseId;constraint:OnDelete:CASCADE;" json:"schedules,omitempty"`
}

// ParseDatabaseURL parses a database connection URL and extracts components
func ParseDatabaseURL(dbType, connectionURL string) (host, port, username, password, dbName string, err error) {
	if connectionURL == "" {
		return "", "", "", "", "", fmt.Errorf("URL vide")
	}

	// PostgreSQL URL format: postgresql://username:password@host:port/database
	// MySQL URL format: mysql://username:password@host:port/database
	if strings.HasPrefix(connectionURL, "postgresql://") || strings.HasPrefix(connectionURL, "mysql://") {
		parsedURL, err := url.Parse(connectionURL)
		if err != nil {
			return "", "", "", "", "", fmt.Errorf("URL invalide: %v", err)
		}

		host, port, err = net.SplitHostPort(parsedURL.Host)
		if err != nil {
			// If no port specified, use defaults
			host = parsedURL.Host
			if dbType == "postgresql" {
				port = "5432"
			} else {
				port = "3306"
			}
		}

		username = parsedURL.User.Username()
		password, _ = parsedURL.User.Password()
		dbName = strings.TrimPrefix(parsedURL.Path, "/")

		return host, port, username, password, dbName, nil
	}

	// Note: Alternative formats like mysql://user:pass@tcp(host:port)/database
	// are not supported. Use standard format: mysql://user:pass@host:port/database

	return "", "", "", "", "", fmt.Errorf("format d'URL non reconnu. Utilisez: mysql://user:pass@host:port/db ou postgresql://user:pass@host:port/db")
}

// ValidateAndNormalizeDatabaseData validates and normalizes database creation data
func ValidateAndNormalizeDatabaseData(req *DatabaseCreateRequest) error {
	// If URL is provided, parse it and override individual fields
	if req.URL != "" {
		host, port, username, password, dbName, err := ParseDatabaseURL(req.Type, req.URL)
		if err != nil {
			return fmt.Errorf("erreur dans l'URL: %v", err)
		}

		// Override individual fields with parsed values
		req.Host = host
		req.Port = port
		req.Username = username
		req.Password = password
		req.DbName = dbName
	} else {
		// Validate that all required fields are provided when no URL
		if req.Host == "" || req.Username == "" || req.Password == "" || req.DbName == "" {
			return fmt.Errorf("tous les champs (host, username, password, db_name) sont requis ou fournissez une URL complète")
		}

		// Parse host:port format (like localhost:3306) if needed
		req.Host, req.Port = ParsePHPEnvFormat(req.Host, req.Port, req.Username, req.Password, req.DbName)
	}

	// Validate port is numeric
	if portNum, err := strconv.Atoi(req.Port); err != nil || portNum < 1 || portNum > 65535 {
		return fmt.Errorf("port invalide: doit être un nombre entre 1 et 65535")
	}

	return nil
}

// ValidateAndNormalizeDatabaseUpdateData validates and normalizes database update data
func ValidateAndNormalizeDatabaseUpdateData(req *DatabaseUpdateRequest) error {
	// If URL is provided, parse it and override individual fields
	if req.URL != "" {
		host, port, username, password, dbName, err := ParseDatabaseURL(req.Type, req.URL)
		if err != nil {
			return fmt.Errorf("erreur dans l'URL: %v", err)
		}

		// Override individual fields with parsed values
		req.Host = host
		req.Port = port
		req.Username = username
		if password != "" { // Only override password if provided in URL
			req.Password = password
		}
		req.DbName = dbName
	} else {
		// For updates, individual fields are optional but if provided, validate port
		if req.Port != "" {
			if portNum, err := strconv.Atoi(req.Port); err != nil || portNum < 1 || portNum > 65535 {
				return fmt.Errorf("port invalide: doit être un nombre entre 1 et 65535")
			}
		}
	}

	return nil
}

// ParsePHPEnvFormat parses database connection info from PHP .env format
// Handles formats like DB_HOST=localhost:3306 or DB_HOST=localhost with DB_PORT=3306
func ParsePHPEnvFormat(host, port, user, pass, dbName string) (parsedHost, parsedPort string) {
	// Handle cases where host contains port (like localhost:3306)
	if strings.Contains(host, ":") {
		parts := strings.Split(host, ":")
		if len(parts) == 2 {
			parsedHost = parts[0]
			// If port is also specified separately, use the separate port
			if port != "" && port != "3306" && port != "5432" {
				parsedPort = port
			} else {
				parsedPort = parts[1]
			}
		} else {
			parsedHost = host
			parsedPort = port
		}
	} else {
		parsedHost = host
		parsedPort = port
	}

	// Set default ports if not specified
	if parsedPort == "" {
		if strings.Contains(user, "postgres") || strings.Contains(dbName, "postgres") {
			parsedPort = "5432"
		} else {
			parsedPort = "3306"
		}
	}

	return parsedHost, parsedPort
}
