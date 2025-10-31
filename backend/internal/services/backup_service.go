package services

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
)

type BackupService struct {
	backupRepo      *repositories.BackupRepository
	databaseRepo    *repositories.DatabaseRepository
	databaseService *DatabaseService
	backupDir       string
}

// Constructor for BackupService
func NewBackupService(backupRepo *repositories.BackupRepository, databaseRepo *repositories.DatabaseRepository, backupDir string) *BackupService {
	return &BackupService{
		backupRepo:   backupRepo,
		databaseRepo: databaseRepo,
		backupDir:    backupDir,
	}
}

// SetDatabaseService sets the database service (to avoid circular dependency)
func (s *BackupService) SetDatabaseService(databaseService *DatabaseService) {
	s.databaseService = databaseService
}

// getMySQLDumpPaths returns possible mysqldump paths based on OS
func (s *BackupService) getMySQLDumpPaths() []string {
	goos := runtime.GOOS
	paths := []string{}

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
			"/usr/bin/mysqldump",             // Ubuntu/Debian default
			"/usr/local/bin/mysqldump",       // Custom install
			"/usr/local/mysql/bin/mysqldump", // MySQL tar.gz install
			"/opt/mysql/bin/mysqldump",       // Some Linux distributions
			"mysqldump",                      // System PATH
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
	paths := []string{}

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

// testDatabaseConnectivity tests if we can connect to the database host and port
func (s *BackupService) testDatabaseConnectivity(host, port string) error {
	timeout := 10 * time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return fmt.Errorf("impossible de se connecter à %s:%s: %v", host, port, err)
	}
	conn.Close()
	return nil
}

// CreateBackup creates a backup for the specified database
func (s *BackupService) CreateBackup(databaseID uint, userID uint) (*models.Backup, error) {
	// Get database info with decrypted password
	database, err := s.databaseService.GetDatabaseByIDForBackup(databaseID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération de la base de données: %v", err)
	}

	// Verify user ownership
	if database.UserId != userID {
		return nil, fmt.Errorf("accès non autorisé à cette base de données")
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s_%s.zip", database.Name, database.Type, timestamp)

	// Create backup record
	backup := &models.Backup{
		Filename:   filename,
		Filepath:   filepath.Join(s.backupDir, database.Type, filename),
		Status:     "pending",
		UserId:     userID,
		DatabaseId: databaseID,
	}

	// Save backup record to database
	if err := s.backupRepo.Create(backup); err != nil {
		return nil, fmt.Errorf("erreur lors de la création de l'enregistrement de sauvegarde: %v", err)
	}

	// Start backup process in goroutine
	go s.performBackup(backup, database)

	return backup, nil
}

// performBackup executes the actual backup process
func (s *BackupService) performBackup(backup *models.Backup, database *models.Database) {
	// Protect against panics in goroutine
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[BACKUP] Panic in backup process: %v\n", r)
			s.updateBackupError(backup.Id, fmt.Sprintf("Erreur critique lors de la sauvegarde: %v", r))
		}
	}()

	fmt.Printf("[BACKUP] Starting backup process for backup ID: %d\n", backup.Id)

	var dumpFile string
	var err error

	// Create backup directory if it doesn't exist
	backupTypeDir := filepath.Join(s.backupDir, database.Type)
	if err := os.MkdirAll(backupTypeDir, 0755); err != nil {
		fmt.Printf("[BACKUP] Failed to create directory: %v\n", err)
		s.updateBackupError(backup.Id, fmt.Sprintf("Erreur lors de la création du répertoire: %v", err))
		return
	}

	// Perform database dump based on type
	switch database.Type {
	case "mysql":
		dumpFile, err = s.dumpMySQL(database, backupTypeDir)
	case "postgres", "postgresql":
		dumpFile, err = s.dumpPostgreSQL(database, backupTypeDir)
	default:
		err = fmt.Errorf("type de base de données non supporté: %s", database.Type)
	}

	if err != nil {
		fmt.Printf("[BACKUP] Dump failed: %v\n", err)
		s.updateBackupError(backup.Id, fmt.Sprintf("Erreur lors du dump: %v", err))
		return
	}

	// Create ZIP file
	zipFile := backup.Filepath
	fmt.Printf("[BACKUP] Creating ZIP file: %s\n", zipFile)
	if err := s.createZipFile(dumpFile, zipFile); err != nil {
		fmt.Printf("[BACKUP] ZIP creation failed: %v\n", err)
		s.updateBackupError(backup.Id, fmt.Sprintf("Erreur lors de la compression: %v", err))
		// Clean up dump file
		os.Remove(dumpFile)
		return
	}

	// Clean up dump file
	fmt.Printf("[BACKUP] Cleaning up dump file: %s\n", dumpFile)
	os.Remove(dumpFile)

	// Get file size
	fileInfo, err := os.Stat(zipFile)
	if err != nil {
		fmt.Printf("[BACKUP] Failed to get file info: %v\n", err)
		s.updateBackupError(backup.Id, fmt.Sprintf("Erreur lors de la lecture des informations du fichier: %v", err))
		return
	}

	fmt.Printf("[BACKUP] Backup file size: %d bytes\n", fileInfo.Size())

	// Update backup record
	if err := s.backupRepo.UpdateSize(backup.Id, fileInfo.Size()); err != nil {
		fmt.Printf("[BACKUP] Failed to update size: %v\n", err)
		s.updateBackupError(backup.Id, fmt.Sprintf("Erreur lors de la mise à jour de la taille: %v", err))
		return
	}

	if err := s.backupRepo.UpdateStatus(backup.Id, "completed", ""); err != nil {
		fmt.Printf("[BACKUP] Failed to update status: %v\n", err)
		s.updateBackupError(backup.Id, fmt.Sprintf("Erreur lors de la mise à jour du statut: %v", err))
		return
	}

	fmt.Printf("[BACKUP] Backup process completed successfully for ID: %d\n", backup.Id)
}

// dumpMySQL creates a MySQL dump
func (s *BackupService) dumpMySQL(database *models.Database, outputDir string) (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	dumpFile := filepath.Join(outputDir, fmt.Sprintf("%s_%s.sql", database.Name, timestamp))

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

	if isRemote {
		fmt.Printf("[BACKUP] Detected remote MySQL connection to %s:%s\n", database.Host, database.Port)
		// For remote connections, skip the connectivity test as it might fail due to firewalls
		// or require specific SSL configurations
		fmt.Printf("[BACKUP] Skipping connectivity test for remote connection (will test during dump)\n")
	} else {
		// Test database connectivity for local connections
		fmt.Printf("[BACKUP] Testing database connectivity to %s:%s\n", database.Host, database.Port)
		if err := s.testDatabaseConnectivity(database.Host, database.Port); err != nil {
			return "", fmt.Errorf("impossible de se connecter à la base de données MySQL. Vérifiez que MySQL est accessible sur le port %s: %v", database.Port, err)
		}
		fmt.Printf("[BACKUP] Database connectivity test passed\n")
	}

	// Find mysqldump executable
	mysqldumpPath, err := s.findExecutable(s.getMySQLDumpPaths())
	if err != nil {
		return "", fmt.Errorf("mysqldump non trouvé: %v", err)
	}
	fmt.Printf("[BACKUP] Using mysqldump at: %s\n", mysqldumpPath)

	// Check if we're using MAMP (detect by path)
	isMAMP := strings.Contains(mysqldumpPath, "/Applications/MAMP/")

	// Build mysqldump command with SSL options for remote connections
	args := []string{
		"-u", database.Username,
		"-p" + database.Password, // Correctly concatenate password
		"--single-transaction",
		"--routines",
		"--triggers",
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
		fmt.Printf("[BACKUP] Detected MAMP installation, using TCP connection (recommended for MAMP)\n")
		args = append(args, "-h", database.Host, "-P", database.Port)
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

		// If SSL/connect options not recognized, try without SSL options
		if strings.Contains(stderrStr, "unknown variable") &&
			(strings.Contains(stderrStr, "ssl-mode") ||
				strings.Contains(stderrStr, "connect_timeout") ||
				strings.Contains(stderrStr, "ssl-ca")) {

			fmt.Printf("[BACKUP] SSL options not supported, retrying with basic options...\n")

			// Retry with basic args (no SSL options)
			baseArgs := []string{
				"-u", database.Username,
				"-p" + database.Password,
				"--single-transaction",
				"--routines",
				"--triggers",
			}

			if isRemote || isMAMP {
				baseArgs = append(baseArgs, "-h", database.Host, "-P", database.Port)
			} else {
				baseArgs = append(baseArgs, "-h", database.Host, "-P", database.Port)
			}

			baseArgs = append(baseArgs, database.DbName)

			cmd = exec.CommandContext(context.Background(), mysqldumpPath, baseArgs...)
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
				fmt.Printf("[BACKUP] mysqldump retry failed: %v\n", err2)
				fmt.Printf("[BACKUP] stderr: %s\n", stderr2.String())
				os.Remove(dumpFile)
				return "", fmt.Errorf("erreur lors de l'exécution de mysqldump (même sans SSL): %v, stderr: %s", err2, stderr2.String())
			}
		} else {
			os.Remove(dumpFile)
			return "", fmt.Errorf("erreur lors de l'exécution de mysqldump: %v, stderr: %s", err, stderrStr)
		}
	}

	fmt.Printf("[BACKUP] mysqldump completed successfully\n")

	return dumpFile, nil
}

// dumpPostgreSQL creates a PostgreSQL dump
func (s *BackupService) dumpPostgreSQL(database *models.Database, outputDir string) (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	dumpFile := filepath.Join(outputDir, fmt.Sprintf("%s_%s.sql", database.Name, timestamp))

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

	if isRemote {
		fmt.Printf("[BACKUP] Detected remote PostgreSQL connection to %s:%s\n", database.Host, database.Port)
		// For remote connections, skip the connectivity test as it might fail due to firewalls
		fmt.Printf("[BACKUP] Skipping connectivity test for remote connection (will test during dump)\n")
	} else {
		// Test database connectivity for local connections
		fmt.Printf("[BACKUP] Testing database connectivity to %s:%s\n", database.Host, database.Port)
		if err := s.testDatabaseConnectivity(database.Host, database.Port); err != nil {
			return "", fmt.Errorf("test de connectivité échoué: %v", err)
		}
		fmt.Printf("[BACKUP] Database connectivity test passed\n")
	}

	// Find pg_dump executable
	pgDumpPath, err := s.findExecutable(s.getPostgreSQLDumpPaths())
	if err != nil {
		return "", fmt.Errorf("pg_dump non trouvé: %v", err)
	}
	fmt.Printf("[BACKUP] Using pg_dump at: %s\n", pgDumpPath)

	// Set environment variable for password
	env := os.Environ()
	env = append(env, fmt.Sprintf("PGPASSWORD=%s", database.Password))

	// Try pg_dump with SSL options first, fallback to basic options if SSL not supported
	baseArgs := []string{
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

	var finalArgs []string
	if isRemote {
		// Try with SSL options for remote connections
		fmt.Printf("[BACKUP] Attempting SSL for remote PostgreSQL connection\n")
		finalArgs = append(baseArgs, "--sslmode=require")
	} else {
		// For local connections, try SSL first
		finalArgs = append(baseArgs, "--sslmode=prefer")
	}

	// Set timeout for the command (30 minutes max for large databases)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, pgDumpPath, finalArgs...)

	cmd.Env = env

	// Capture stderr for better error reporting
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	fmt.Printf("[BACKUP] Command args: %v\n", finalArgs)
	fmt.Printf("[BACKUP] Starting pg_dump execution...\n")

	err = cmd.Run()
	if err != nil {
		stderrStr := stderr.String()
		fmt.Printf("[BACKUP] pg_dump failed: %v\n", err)
		fmt.Printf("[BACKUP] stderr: %s\n", stderrStr)

		// If SSL option not recognized, try without SSL
		if strings.Contains(stderrStr, "unrecognized option") && strings.Contains(stderrStr, "sslmode") {
			fmt.Printf("[BACKUP] SSL options not supported, retrying without SSL...\n")

			// Retry with basic args (no SSL options)
			finalArgs = baseArgs
			cmd = exec.CommandContext(context.Background(), pgDumpPath, finalArgs...)
			cmd.Env = env

			var stderr2 bytes.Buffer
			cmd.Stderr = &stderr2

			fmt.Printf("[BACKUP] Retry command args: %v\n", finalArgs)
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

// createZipFile compresses a file into a ZIP archive
func (s *BackupService) createZipFile(sourceFile, zipFile string) error {
	// Create ZIP file
	zipFileHandle, err := os.Create(zipFile)
	if err != nil {
		return fmt.Errorf("erreur lors de la création du fichier ZIP: %v", err)
	}
	defer zipFileHandle.Close()

	// Create ZIP writer
	zipWriter := zip.NewWriter(zipFileHandle)
	defer zipWriter.Close()

	// Open source file
	sourceFileHandle, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture du fichier source: %v", err)
	}
	defer sourceFileHandle.Close()

	// Get source file info
	sourceInfo, err := sourceFileHandle.Stat()
	if err != nil {
		return fmt.Errorf("erreur lors de la lecture des informations du fichier source: %v", err)
	}

	// Create file header
	header, err := zip.FileInfoHeader(sourceInfo)
	if err != nil {
		return fmt.Errorf("erreur lors de la création de l'en-tête ZIP: %v", err)
	}

	// Set compression method
	header.Method = zip.Deflate

	// Create writer for file in ZIP
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("erreur lors de la création de l'entrée ZIP: %v", err)
	}

	// Copy file content to ZIP
	_, err = io.Copy(writer, sourceFileHandle)
	if err != nil {
		return fmt.Errorf("erreur lors de la copie du contenu vers le ZIP: %v", err)
	}

	return nil
}

// updateBackupError updates backup status to failed with error message
func (s *BackupService) updateBackupError(backupID uint, errorMsg string) {
	s.backupRepo.UpdateStatus(backupID, "failed", errorMsg)
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

// DeleteBackup deletes a backup file and record
func (s *BackupService) DeleteBackup(id uint, userID uint) error {
	backup, err := s.backupRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("sauvegarde introuvable: %v", err)
	}

	// Verify user ownership
	if backup.UserId != userID {
		return fmt.Errorf("accès non autorisé à cette sauvegarde")
	}

	// Delete file if it exists
	if _, err := os.Stat(backup.Filepath); err == nil {
		if err := os.Remove(backup.Filepath); err != nil {
			return fmt.Errorf("erreur lors de la suppression du fichier: %v", err)
		}
	}

	// Delete record from database
	return s.backupRepo.Delete(id)
}
