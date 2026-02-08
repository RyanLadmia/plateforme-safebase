package services

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
)

// WorkerPool interface for background tasks
type WorkerPoolInterface interface {
	Submit(task func())
}

type BackupService struct {
	backupRepo           *repositories.BackupRepository
	databaseService      *DatabaseService
	userService          *UserService
	backupDir            string
	workerPool           WorkerPoolInterface
	cloudStorage         CloudStorageService // Generic cloud storage interface
	encryption           *EncryptionService  // Encryption service for files
	actionHistoryService *ActionHistoryService
}

// Constructor for BackupService
func NewBackupService(backupRepo *repositories.BackupRepository, databaseService *DatabaseService, userService *UserService, backupDir string) *BackupService {
	return &BackupService{
		backupRepo:      backupRepo,
		databaseService: databaseService,
		userService:     userService,
		backupDir:       backupDir,
	}
}

// SetActionHistoryService sets the action history service reference for logging
func (s *BackupService) SetActionHistoryService(actionHistoryService *ActionHistoryService) {
	s.actionHistoryService = actionHistoryService
}

// generateBackupFilename generates a consistent filename for backups
func (s *BackupService) generateBackupFilename(database *models.Database) string {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	return fmt.Sprintf("%s_%s_%s.zip", database.Name, database.Type, timestamp)
}

// zipFile compresses a SQL file into a ZIP archive
func (s *BackupService) zipFile(sqlFilePath, zipFilePath string) error {
	// Open the SQL file
	sqlFile, err := os.Open(sqlFilePath)
	if err != nil {
		return fmt.Errorf("failed to open SQL file: %v", err)
	}
	defer sqlFile.Close()

	// Create the ZIP file
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to create ZIP file: %v", err)
	}
	defer zipFile.Close()

	// Create a ZIP writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Get the filename from the path
	_, sqlFilename := filepath.Split(sqlFilePath)

	// Create a ZIP entry
	zipEntry, err := zipWriter.Create(sqlFilename)
	if err != nil {
		return fmt.Errorf("failed to create ZIP entry: %v", err)
	}

	// Copy the SQL file content to the ZIP entry
	_, err = io.Copy(zipEntry, sqlFile)
	if err != nil {
		return fmt.Errorf("failed to copy file to ZIP: %v", err)
	}

	return nil
}

// SetWorkerPool sets the worker pool for background tasks
func (s *BackupService) SetWorkerPool(workerPool WorkerPoolInterface) {
	s.workerPool = workerPool
}

// SetCloudStorage sets the cloud storage service
func (s *BackupService) SetCloudStorage(cloudStorage CloudStorageService) {
	s.cloudStorage = cloudStorage
}

// SetEncryptionService sets the encryption service
func (s *BackupService) SetEncryptionService(encryption *EncryptionService) {
	s.encryption = encryption
}

// getMySQLDumpPaths returns possible mysqldump paths based on OS
func (s *BackupService) getMySQLDumpPaths() []string {
	goos := runtime.GOOS
	var paths []string

	switch goos {
	case "darwin": // macOS
		paths = []string{
			"/Applications/MAMP/Library/bin/mysql80/bin/mysqldump", // MAMP MySQL 8.0
			"/Applications/MAMP/Library/bin/mysqldump",             // MAMP general
			"/usr/local/mysql/bin/mysqldump",                       // MySQL installed via DMG
			"/usr/local/bin/mysqldump",                             // Homebrew MySQL
			"/opt/homebrew/bin/mysqldump",                          // Homebrew Apple Silicon
			"mysqldump",                                            // System PATH
		}
	case "linux":
		paths = []string{
			"/usr/bin/mysqldump",             // Ubuntu/Debian default (priorité sur mariadb-dump)
			"/usr/local/bin/mysqldump",       // Custom install
			"/usr/local/mysql/bin/mysqldump", // MySQL tar.gz install
			"/opt/mysql/bin/mysqldump",       // Some Linux distributions
			"mysqldump",                      // System PATH
			// Note: mariadb-dump n'est pas inclus car il interroge generation_expression
			// ce qui cause des erreurs avec MariaDB 5.5 et MySQL 5.6
		}
	case "windows":
		paths = []string{
			"C:\\Program Files\\MySQL\\MySQL Server 8.0\\bin\\mysqldump.exe",
			"C:\\Program Files\\MySQL\\MySQL Server 5.7\\bin\\mysqldump.exe",
			"C:\\Program Files (x86)\\MySQL\\MySQL Server 8.0\\bin\\mysqldump.exe",
			"C:\\xampp\\mysql\\bin\\mysqldump.exe", // XAMPP
			"mysqldump.exe",                        // System PATH
			"mysqldump",                            // WSL or Git Bash
		}
	default:
		paths = []string{"mysqldump"}
	}

	return paths
}

// getPostgreSQLDumpPaths returns possible pg_dump paths based on OS
func (s *BackupService) getPostgreSQLDumpPaths() []string {
	goos := runtime.GOOS
	var paths []string

	switch goos {
	case "darwin": // macOS
		paths = []string{
			"/Applications/Postgres.app/Contents/Versions/latest/bin/pg_dump", // Postgres.app
			"/usr/local/pgsql/bin/pg_dump",                                    // PostgreSQL installed manually
			"/usr/local/bin/pg_dump",                                          // Homebrew PostgreSQL
			"/opt/homebrew/bin/pg_dump",                                       // Homebrew Apple Silicon
			"/Library/PostgreSQL/14/bin/pg_dump",                              // EnterpriseDB installer
			"/Library/PostgreSQL/13/bin/pg_dump",
			"/Library/PostgreSQL/12/bin/pg_dump",
			"pg_dump", // System PATH
		}
	case "linux":
		paths = []string{
			"/usr/bin/pg_dump",          // Ubuntu/Debian default
			"/usr/local/bin/pg_dump",    // Custom install
			"/usr/pgsql-14/bin/pg_dump", // RedHat/CentOS PGDG
			"/usr/pgsql-13/bin/pg_dump",
			"/usr/pgsql-12/bin/pg_dump",
			"/opt/postgresql/bin/pg_dump", // Some distributions
			"pg_dump",                     // System PATH
		}
	case "windows":
		paths = []string{
			"C:\\Program Files\\PostgreSQL\\14\\bin\\pg_dump.exe",
			"C:\\Program Files\\PostgreSQL\\13\\bin\\pg_dump.exe",
			"C:\\Program Files\\PostgreSQL\\12\\bin\\pg_dump.exe",
			"C:\\Program Files\\PostgreSQL\\11\\bin\\pg_dump.exe",
			"C:\\PostgreSQL\\pg14\\bin\\pg_dump.exe", // Some installers
			"C:\\PostgreSQL\\pg13\\bin\\pg_dump.exe",
			"C:\\PostgreSQL\\pg12\\bin\\pg_dump.exe",
			"pg_dump.exe", // System PATH
			"pg_dump",     // WSL or Git Bash
		}
	default:
		paths = []string{"pg_dump"}
	}

	return paths
}

// findExecutable tries to find a working executable from a list of paths
func (s *BackupService) findExecutable(paths []string) (string, error) {
	for _, path := range paths {
		if _, err := exec.LookPath(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("aucun exécutable trouvé dans les chemins testés")
}

// executeBackupAsync executes the backup process asynchronously
func (s *BackupService) executeBackupAsync(backup *models.Backup, database *models.Database) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[BACKUP] Panic in backup process for ID %d: %v\n", backup.Id, r)
			s.updateBackupError(backup.Id, fmt.Sprintf("Erreur critique: %v", r))
		}
	}()

	fmt.Printf("[BACKUP] Starting asynchronous backup process for ID: %d\n", backup.Id)

	// Update status to running
	if err := s.backupRepo.UpdateStatus(backup.Id, "running", ""); err != nil {
		fmt.Printf("[BACKUP] Failed to update status to running: %v\n", err)
		return
	}

	// Get user information for MinIO path
	user, err := s.userService.GetUserByID(database.UserId)
	if err != nil {
		fmt.Printf("[BACKUP] Failed to get user information: %v\n", err)
		s.updateBackupError(backup.Id, fmt.Sprintf("Erreur lors de la récupération des informations utilisateur: %v", err))
		return
	}

	// Create backup directory if it doesn't exist (temporary for dump creation)
	if err := os.MkdirAll(s.backupDir, 0755); err != nil {
		fmt.Printf("[BACKUP] Failed to create backup directory: %v\n", err)
		s.updateBackupError(backup.Id, fmt.Sprintf("Erreur lors de la création du répertoire: %v", err))
		return
	}

	var backupFilepath string

	// Execute backup based on database type
	switch dbType := database.Type; dbType {
	case "mysql":
		sqlFilename := strings.TrimSuffix(backup.Filename, ".zip") + ".sql"
		dbTypeDir := filepath.Join(s.backupDir, "mysql")
		backupFilepath, err = s.dumpMySQL(database, dbTypeDir, sqlFilename)
	case "postgresql":
		sqlFilename := strings.TrimSuffix(backup.Filename, ".zip") + ".sql"
		dbTypeDir := filepath.Join(s.backupDir, "postgresql")
		backupFilepath, err = s.dumpPostgreSQL(database, dbTypeDir, sqlFilename)
	default:
		err = fmt.Errorf("type de base de données non supporté: %s", dbType)
	}

	if err != nil {
		fmt.Printf("[BACKUP] Backup failed for ID %d: %v\n", backup.Id, err)
		s.updateBackupError(backup.Id, err.Error())
		return
	}

	// Compress the SQL file to ZIP
	zipFilePath := strings.TrimSuffix(backupFilepath, ".sql") + ".zip"
	fmt.Printf("[BACKUP] Compressing SQL file to ZIP: %s -> %s\n", backupFilepath, zipFilePath)

	if err := s.zipFile(backupFilepath, zipFilePath); err != nil {
		fmt.Printf("[BACKUP] Failed to compress backup file: %v\n", err)
		s.updateBackupError(backup.Id, fmt.Sprintf("Erreur lors de la compression: %v", err))
		return
	}

	// Remove the original SQL file after successful compression
	if err := os.Remove(backupFilepath); err != nil {
		fmt.Printf("[BACKUP] Warning: failed to remove original SQL file: %v\n", err)
		// Don't fail the backup for this
	}

	// Update filepath to point to the ZIP file
	backupFilepath = zipFilePath

	// Get file size
	fileInfo, err := os.Stat(backupFilepath)
	if err != nil {
		fmt.Printf("[BACKUP] Failed to get file info: %v\n", err)
		s.updateBackupError(backup.Id, fmt.Sprintf("Erreur lors de la récupération des informations du fichier: %v", err))
		return
	}

	// Upload to cloud storage (required for cloud storage)
	var remotePath string
	if s.cloudStorage != nil && s.encryption != nil {
		// Generate user-specific encryption key
		userKey := GenerateUserKey(database.UserId, "SafeBaseBackupSalt2025!")
		userEncryption := NewEncryptionService(userKey)

		// Encrypt the file before upload
		fmt.Printf("[BACKUP] Encrypting backup file before upload...\n")
		encryptedData, err := userEncryption.EncryptFile(backupFilepath)
		if err != nil {
			fmt.Printf("[BACKUP] Failed to encrypt backup file: %v\n", err)
			s.updateBackupError(backup.Id, fmt.Sprintf("Erreur lors du chiffrement: %v", err))
			return
		}

		// Create temporary encrypted file
		encryptedFilePath := backupFilepath + ".encrypted"
		err = os.WriteFile(encryptedFilePath, encryptedData, 0644)
		if err != nil {
			fmt.Printf("[BACKUP] Failed to write encrypted file: %v\n", err)
			s.updateBackupError(backup.Id, fmt.Sprintf("Erreur lors de l'écriture du fichier chiffré: %v", err))
			return
		}

		// Generate remote path
		remotePath = s.cloudStorage.GenerateRemotePath(fmt.Sprintf("%s %s", user.Firstname, user.Lastname), database.Type, backup.Filename)

		// Upload encrypted file
		if err := s.cloudStorage.UploadFile(encryptedFilePath, remotePath); err != nil {
			fmt.Printf("[BACKUP] Failed to upload to cloud storage: %v\n", err)
			s.updateBackupError(backup.Id, fmt.Sprintf("Erreur lors de l'upload vers le cloud: %v", err))
			os.Remove(encryptedFilePath) // Clean up
			return
		} else {
			fmt.Printf("[BACKUP] Successfully uploaded encrypted backup to cloud: %s\n", remotePath)
		}

		// Clean up encrypted temporary file
		os.Remove(encryptedFilePath)
	} else {
		fmt.Printf("[BACKUP] Cloud storage or encryption service not available - backup failed\n")
		s.updateBackupError(backup.Id, "Service de stockage cloud ou de chiffrement non disponible")
		return
	}

	// Clean up local temporary file after successful upload
	if err := os.Remove(backupFilepath); err != nil {
		fmt.Printf("[BACKUP] Warning: failed to remove temporary local file: %v\n", err)
		// Don't fail the backup for this
	}

	// Update backup record with success information
	backup.Status = "completed"
	backup.Filepath = remotePath // Store cloud path
	backup.Size = fileInfo.Size()

	// Update status and filepath
	if err := s.backupRepo.UpdateStatus(backup.Id, "completed", ""); err != nil {
		fmt.Printf("[BACKUP] Failed to update backup status: %v\n", err)
		return
	}

	// Update filepath and size
	if err := s.backupRepo.UpdateFileInfo(backup.Id, remotePath, fileInfo.Size()); err != nil {
		fmt.Printf("[BACKUP] Failed to update file info: %v\n", err)
		return
	}

	fmt.Printf("[BACKUP] Backup process completed successfully for ID: %d\n", backup.Id)
}

// dumpMySQL creates a MySQL dump
func (s *BackupService) dumpMySQL(database *models.Database, outputDir string, filename string) (string, error) {
	dumpFile := filepath.Join(outputDir, filename)

	// Log the start of backup process
	fmt.Printf("[BACKUP] Starting MySQL dump for database %s (host: %s, port: %s, user: %s, password_length: %d)\n",
		database.DbName, database.Host, database.Port, database.Username, len(database.Password))

	// Validate that password is not empty
	if len(database.Password) == 0 {
		return "", fmt.Errorf("mot de passe de base de données vide ou non configuré")
	}

	// Determine if this is a remote connection
	isRemote := database.Host != "localhost" && database.Host != "127.0.0.1" && !strings.HasPrefix(database.Host, "192.168.") && !strings.HasPrefix(database.Host, "10.")
	fmt.Printf("[BACKUP] Connection type analysis: host='%s', remote=%t\n", database.Host, isRemote)

	// Skip connectivity test for now - let mysqldump handle connection errors
	// This prevents false negatives with complex network configurations
	fmt.Printf("[BACKUP] Skipping connectivity test - will test during dump execution\n")

	// Find mysqldump executable
	mysqldumpPath, err := s.findExecutable(s.getMySQLDumpPaths())
	if err != nil {
		return "", fmt.Errorf("mysqldump non trouvé: %v", err)
	}
	fmt.Printf("[BACKUP] Using mysqldump at: %s\n", mysqldumpPath)

	// Check if we're using MAMP (detect by path)
	isMAMP := strings.Contains(mysqldumpPath, "/Applications/MAMP/")

	// Build mysqldump command with SSL options for remote connections
	env := os.Environ()
	env = append(env, fmt.Sprintf("MYSQL_PWD=%s", database.Password))

	args := []string{
		"-u", database.Username,
		"--single-transaction",
		"--routines",
		"--triggers",
		"--no-tablespaces", // Évite les erreurs de permissions sur les tablespaces
	}

	// Connection and SSL configuration
	if isRemote {
		fmt.Printf("[BACKUP] Configuring for remote MySQL connection\n")
		args = append(args, "-h", database.Host, "-P", database.Port)

		// For remote connections, try SSL but don't require it
		args = append(args,
			"--ssl-mode=PREFERRED", // Try SSL but don't require it
		)

		// Add SSL CA certificates if available
		sslCaPaths := []string{
			"/etc/ssl/certs/ca-certificates.crt", // Linux
			"/etc/ssl/cert.pem",                  // macOS
			"/etc/pki/tls/certs/ca-bundle.crt",   // RedHat/CentOS
		}

		for _, caPath := range sslCaPaths {
			if _, err := os.Stat(caPath); err == nil {
				args = append(args, "--ssl-ca="+caPath)
				fmt.Printf("[BACKUP] Using SSL CA: %s\n", caPath)
				break
			}
		}

	} else if isMAMP {
		fmt.Printf("[BACKUP] Detected MAMP installation, using socket connection (recommended for MAMP)\n")
		// For MAMP, use socket connection instead of TCP
		socketPath := "/Applications/MAMP/tmp/mysql/mysql.sock"
		if _, err := os.Stat(socketPath); err == nil {
			args = append(args, "--socket="+socketPath)
			fmt.Printf("[BACKUP] Using MAMP socket: %s\n", socketPath)
		} else {
			// Fallback to TCP if socket not found
			fmt.Printf("[BACKUP] MAMP socket not found at %s, falling back to TCP connection\n", socketPath)
			args = append(args, "-h", database.Host, "-P", database.Port)
		}
	} else {
		// Standard local TCP connection
		args = append(args, "-h", database.Host, "-P", database.Port)
	}

	args = append(args, database.DbName)

	// Set timeout for the command (30 minutes max for large databases)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Try with SSL options first
	fmt.Printf("[BACKUP] Command args: %v\n", args)
	cmd := exec.CommandContext(ctx, mysqldumpPath, args...)
	cmd.Env = env

	output, err := os.Create(dumpFile)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la création du fichier de dump: %v", err)
	}
	defer output.Close()

	// Capture stderr for better error reporting
	var stderr bytes.Buffer
	cmd.Stdout = output
	cmd.Stderr = &stderr

	fmt.Printf("[BACKUP] Starting mysqldump execution...\n")
	err = cmd.Run()

	if err != nil {
		stderrStr := stderr.String()
		fmt.Printf("[BACKUP] mysqldump failed: %v\n", err)
		fmt.Printf("[BACKUP] stderr: %s\n", stderrStr)

		// Build base args for retry attempts (without SSL options)
		buildBaseArgs := func() []string {
			baseArgs := []string{
				"-u", database.Username,
				"--single-transaction",
				"--routines",
				"--triggers",
				"--no-tablespaces", // Évite les erreurs de permissions sur les tablespaces
			}

			if isMAMP {
				socketPath := "/Applications/MAMP/tmp/mysql/mysql.sock"
				if _, err := os.Stat(socketPath); err == nil {
					baseArgs = append(baseArgs, "--socket="+socketPath)
				} else {
					baseArgs = append(baseArgs, "-h", database.Host, "-P", database.Port)
				}
			} else {
				baseArgs = append(baseArgs, "-h", database.Host, "-P", database.Port)
			}

			baseArgs = append(baseArgs, database.DbName)
			return baseArgs
		}

		// If SSL/connect options not recognized, try without SSL options
		if strings.Contains(stderrStr, "unknown variable") &&
			(strings.Contains(stderrStr, "ssl-mode") ||
				strings.Contains(stderrStr, "connect_timeout") ||
				strings.Contains(stderrStr, "ssl-ca")) {

			fmt.Printf("[BACKUP] SSL options not supported, retrying with basic options...\n")

			// Retry with basic args (no SSL options)
			baseArgs := buildBaseArgs()

			cmd = exec.CommandContext(context.Background(), mysqldumpPath, baseArgs...)
			cmd.Env = env
			output2, err := os.Create(dumpFile)
			if err != nil {
				return "", fmt.Errorf("erreur lors de la création du fichier de dump: %v", err)
			}
			defer output2.Close()

			var stderr2 bytes.Buffer
			cmd.Stdout = output2
			cmd.Stderr = &stderr2

			fmt.Printf("[BACKUP] Retry command args: %v\n", baseArgs)
			if err2 := cmd.Run(); err2 != nil {
				stderr2Str := stderr2.String()
				fmt.Printf("[BACKUP] mysqldump retry failed: %v\n", err2)
				fmt.Printf("[BACKUP] stderr: %s\n", stderr2Str)
				os.Remove(dumpFile)
				return "", fmt.Errorf("erreur lors de l'exécution de mysqldump (même sans SSL): %v, stderr: %s", err2, stderr2Str)
			}
		} else if strings.Contains(stderrStr, "Unknown column 'generation_expression'") {
			// If we get generation_expression error, it means the client is too recent for the server (e.g., MariaDB 5.5)
			// This should be resolved by using an older mysql-client version in Dockerfile (Alpine v3.12)
			os.Remove(dumpFile)
			return "", fmt.Errorf("incompatibilité de version client/serveur MySQL: le serveur est MariaDB 5.5 ou MySQL 5.6 (trop ancien). Le client mysqldump doit être downgraded à Alpine v3.12 (MariaDB 10.4). Veuillez rebuilder le conteneur Docker")
		} else {
			os.Remove(dumpFile)
			return "", fmt.Errorf("erreur lors de l'exécution de mysqldump: %v, stderr: %s", err, stderrStr)
		}
	}

	fmt.Printf("[BACKUP] mysqldump completed successfully\n")

	return dumpFile, nil
}

// dumpPostgreSQL creates a PostgreSQL dump
func (s *BackupService) dumpPostgreSQL(database *models.Database, outputDir string, filename string) (string, error) {
	dumpFile := filepath.Join(outputDir, filename)

	// Log the start of backup process
	fmt.Printf("[BACKUP] Starting PostgreSQL dump for database %s (host: %s, port: %s, user: %s, password_length: %d)\n",
		database.DbName, database.Host, database.Port, database.Username, len(database.Password))

	// Validate that password is not empty
	if len(database.Password) == 0 {
		return "", fmt.Errorf("mot de passe de base de données vide ou non configuré")
	}

	// Determine if this is a remote connection
	isRemote := database.Host != "localhost" && database.Host != "127.0.0.1" && !strings.HasPrefix(database.Host, "192.168.") && !strings.HasPrefix(database.Host, "10.")
	fmt.Printf("[BACKUP] Connection type analysis: host='%s', remote=%t\n", database.Host, isRemote)

	// Skip connectivity test for now - let pg_dump handle connection errors
	// This prevents false negatives with complex network configurations
	fmt.Printf("[BACKUP] Skipping connectivity test - will test during dump execution\n")

	// Find pg_dump executable
	pgDumpPath, err := s.findExecutable(s.getPostgreSQLDumpPaths())
	if err != nil {
		return "", fmt.Errorf("pg_dump non trouvé: %v", err)
	}
	fmt.Printf("[BACKUP] Using pg_dump at: %s\n", pgDumpPath)

	// Set environment variables for password and SSL mode
	env := os.Environ()
	env = append(env, fmt.Sprintf("PGPASSWORD=%s", database.Password))

	// Configure SSL mode via environment variable (not command line option)
	if isRemote {
		// For remote connections (Neon, Supabase, etc.), require SSL
		fmt.Printf("[BACKUP] Configuring SSL mode 'require' for remote PostgreSQL connection\n")
		env = append(env, "PGSSLMODE=require")
	} else {
		// For local connections, prefer SSL but allow fallback
		fmt.Printf("[BACKUP] Configuring SSL mode 'prefer' for local PostgreSQL connection\n")
		env = append(env, "PGSSLMODE=prefer")
	}

	// Build pg_dump arguments (no SSL options in args, use env vars instead)
	args := []string{
		"-h", database.Host,
		"-p", database.Port,
		"-U", database.Username,
		"-d", database.DbName,
		"--no-password",
		"--verbose",
		"--clean",
		"--no-acl",
		"--no-owner",
		"-f", dumpFile,
	}

	// Set timeout for the command (30 minutes max for large databases)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, pgDumpPath, args...)

	cmd.Env = env

	// Capture stderr for better error reporting
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	fmt.Printf("[BACKUP] Command args: %v\n", args)
	fmt.Printf("[BACKUP] Environment: PGSSLMODE=%s\n", func() string {
		if isRemote {
			return "require"
		}
		return "prefer"
	}())
	fmt.Printf("[BACKUP] Starting pg_dump execution...\n")

	err = cmd.Run()
	if err != nil {
		stderrStr := stderr.String()
		fmt.Printf("[BACKUP] pg_dump failed: %v\n", err)
		fmt.Printf("[BACKUP] stderr: %s\n", stderrStr)

		// If SSL connection failed for local connections, try without SSL
		if !isRemote && (strings.Contains(stderrStr, "SSL") || strings.Contains(stderrStr, "ssl")) {
			fmt.Printf("[BACKUP] SSL connection failed for local connection, retrying with PGSSLMODE=disable...\n")

			// Retry without SSL for local connections
			envNoSSL := os.Environ()
			envNoSSL = append(envNoSSL, fmt.Sprintf("PGPASSWORD=%s", database.Password))
			envNoSSL = append(envNoSSL, "PGSSLMODE=disable")

			cmd = exec.CommandContext(context.Background(), pgDumpPath, args...)
			cmd.Env = envNoSSL

			var stderr2 bytes.Buffer
			cmd.Stderr = &stderr2

			fmt.Printf("[BACKUP] Retry with PGSSLMODE=disable\n")
			err2 := cmd.Run()
			if err2 != nil {
				fmt.Printf("[BACKUP] pg_dump retry failed: %v\n", err2)
				fmt.Printf("[BACKUP] stderr: %s\n", stderr2.String())
				os.Remove(dumpFile)
				return "", fmt.Errorf("erreur lors de l'exécution de pg_dump (même sans SSL): %v, stderr: %s", err2, stderr2.String())
			}
		} else {
			os.Remove(dumpFile)
			return "", fmt.Errorf("erreur lors de l'exécution de pg_dump: %v, stderr: %s", err, stderrStr)
		}
	}

	fmt.Printf("[BACKUP] pg_dump completed successfully\n")

	return dumpFile, nil
}

// updateBackupError updates backup status to failed with error message
func (s *BackupService) updateBackupError(backupID uint, errorMsg string) {
	if err := s.backupRepo.UpdateStatus(backupID, "failed", errorMsg); err != nil {
		// Log the error but don't fail the operation
		fmt.Printf("Warning: failed to update backup status: %v\n", err)
	}
}

// GetBackupsByUser returns all backups for a user
func (s *BackupService) GetBackupsByUser(userID uint) ([]models.Backup, error) {
	return s.backupRepo.GetByUserID(userID)
}

// GetBackupsByDatabase returns all backups for a database
func (s *BackupService) GetBackupsByDatabase(databaseID uint) ([]models.Backup, error) {
	return s.backupRepo.GetByDatabaseID(databaseID)
}

// GetBackupByID returns a backup by ID
func (s *BackupService) GetBackupByID(id uint) (*models.Backup, error) {
	return s.backupRepo.GetByID(id)
}

// DownloadBackup downloads a backup file from cloud storage with decryption
func (s *BackupService) DownloadBackup(id uint, userID uint) ([]byte, error) {
	backup, err := s.backupRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("sauvegarde introuvable: %v", err)
	}

	// Verify user ownership
	if backup.UserId != userID {
		return nil, fmt.Errorf("accès non autorisé à cette sauvegarde")
	}

	// Check if cloud storage and encryption services are available
	if s.cloudStorage == nil || s.encryption == nil {
		return nil, fmt.Errorf("services de stockage cloud ou de chiffrement non disponibles")
	}

	// Generate user-specific encryption key
	userKey := GenerateUserKey(userID, "SafeBaseBackupSalt2025!")
	userEncryption := NewEncryptionService(userKey)

	// Download encrypted file from cloud
	encryptedData, err := s.cloudStorage.DownloadFile(backup.Filepath)
	if err != nil {
		return nil, fmt.Errorf("fichier non trouvé dans le stockage cloud: %v", err)
	}

	// Decrypt the file
	decryptedData, err := userEncryption.DecryptData(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("erreur lors du déchiffrement: %v", err)
	}

	return decryptedData, nil
}

// DeleteBackup deletes a backup file and record (from cloud storage)
func (s *BackupService) DeleteBackup(id uint, userID uint) error {
	backup, err := s.backupRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("sauvegarde introuvable: %v", err)
	}

	// Verify user ownership
	if backup.UserId != userID {
		return fmt.Errorf("accès non autorisé à cette sauvegarde")
	}

	// Delete from cloud storage
	if s.cloudStorage != nil {
		if err := s.cloudStorage.DeleteFile(backup.Filepath); err != nil {
			fmt.Printf("[BACKUP] Warning: failed to delete from cloud storage: %v\n", err)
			return fmt.Errorf("erreur lors de la suppression du fichier cloud: %v", err)
		}
	} else {
		return fmt.Errorf("service de stockage cloud non disponible")
	}

	// Delete record from database
	return s.backupRepo.Delete(id)
}

// Logging methods for action history

// CreateBackup creates a backup and logs the action
func (s *BackupService) CreateBackup(databaseID uint, userID uint, ipAddress string, userAgent string) (*models.Backup, error) {
	// Get database info with decrypted password
	database, err := s.databaseService.GetDatabaseByID(databaseID)
	if err != nil {
		return nil, fmt.Errorf("database not found: %v", err)
	}

	// Verify that the database belongs to the user
	if database.UserId != userID {
		return nil, fmt.Errorf("unauthorized: database does not belong to user")
	}

	// Create backup record with pending status
	backup := &models.Backup{
		UserId:     userID,
		DatabaseId: databaseID,
		Status:     "pending",
		Filename:   s.generateBackupFilename(database),
		Filepath:   "", // Will be set when backup is completed
		UserAgent:  userAgent,
	}

	if err := s.backupRepo.Create(backup); err != nil {
		return nil, fmt.Errorf("failed to create backup record: %v", err)
	}

	// Execute backup asynchronously using worker pool
	if s.workerPool != nil {
		s.workerPool.Submit(func() {
			s.executeBackupAsync(backup, database)
		})
	} else {
		go s.executeBackupAsync(backup, database)
	}

	// Log the action
	if s.actionHistoryService != nil {
		metadata := map[string]interface{}{
			"backup_id":     backup.Id,
			"database_id":   backup.DatabaseId,
			"database_name": database.Name,
			"filename":      backup.Filename,
			"size":          backup.Size,
			"filepath":      backup.Filepath,
		}
		description := fmt.Sprintf("Sauvegarde '%s' créée (Base de données: %s)", backup.Filename, database.Name)
		_ = s.actionHistoryService.LogAction(userID, "create", "backup", backup.Id, description, metadata, ipAddress, userAgent)
	}

	return backup, nil
}

// DeleteBackupWithLogging deletes a backup and logs the action
func (s *BackupService) DeleteBackupWithLogging(id uint, userID uint, ipAddress string, userAgent string) error {
	backup, err := s.backupRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("sauvegarde introuvable: %v", err)
	}

	// Verify user ownership
	if backup.UserId != userID {
		return fmt.Errorf("accès non autorisé à cette sauvegarde")
	}

	// Get database info for logging
	database, err := s.databaseService.GetDatabaseByID(backup.DatabaseId)
	if err != nil {
		return fmt.Errorf("database not found: %v", err)
	}

	// Delete from cloud storage
	if s.cloudStorage != nil {
		if err := s.cloudStorage.DeleteFile(backup.Filepath); err != nil {
			fmt.Printf("[BACKUP] Warning: failed to delete from cloud storage: %v\n", err)
			return fmt.Errorf("erreur lors de la suppression du fichier cloud: %v", err)
		}
	} else {
		return fmt.Errorf("service de stockage cloud non disponible")
	}

	// Delete record from database
	err = s.backupRepo.Delete(id)
	if err != nil {
		return err
	}

	// Log the action
	if s.actionHistoryService != nil {
		metadata := map[string]interface{}{
			"backup_id":     backup.Id,
			"database_id":   backup.DatabaseId,
			"database_name": database.Name,
			"filename":      backup.Filename,
			"size":          backup.Size,
			"filepath":      backup.Filepath,
		}
		description := fmt.Sprintf("Sauvegarde '%s' supprimée (Base de données: %s)", backup.Filename, database.Name)
		_ = s.actionHistoryService.LogAction(userID, "delete", "backup", backup.Id, description, metadata, ipAddress, userAgent)
	}

	return nil
}

// DownloadBackupWithLogging downloads a backup and logs the action
func (s *BackupService) DownloadBackupWithLogging(id uint, userID uint, ipAddress string, userAgent string) ([]byte, error) {
	data, err := s.DownloadBackup(id, userID)
	if err != nil {
		return nil, err
	}

	// Log the action
	if s.actionHistoryService != nil {
		backup, _ := s.backupRepo.GetByID(id)
		if backup != nil {
			// Get database info for logging
			database, err := s.databaseService.GetDatabaseByID(backup.DatabaseId)
			if err != nil {
				// If database not found, still log with available info
				metadata := map[string]interface{}{
					"backup_id":   backup.Id,
					"database_id": backup.DatabaseId,
					"filename":    backup.Filename,
					"size":        backup.Size,
					"filepath":    backup.Filepath,
				}
				_ = s.actionHistoryService.LogAction(userID, "download", "backup", backup.Id, "Sauvegarde téléchargée", metadata, ipAddress, userAgent)
			} else {
				metadata := map[string]interface{}{
					"backup_id":     backup.Id,
					"database_id":   backup.DatabaseId,
					"database_name": database.Name,
					"filename":      backup.Filename,
					"size":          backup.Size,
					"filepath":      backup.Filepath,
				}
				description := fmt.Sprintf("Sauvegarde '%s' téléchargée (Base de données: %s)", backup.Filename, database.Name)
				s.actionHistoryService.LogAction(userID, "download", "backup", backup.Id, description, metadata, ipAddress, userAgent)
			}
		}
	}

	return data, nil
}
