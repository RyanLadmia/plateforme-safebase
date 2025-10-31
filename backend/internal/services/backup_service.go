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

	// Test database connectivity for all connections (including localhost)
	fmt.Printf("[BACKUP] Testing database connectivity to %s:%s\n", database.Host, database.Port)
	if err := s.testDatabaseConnectivity(database.Host, database.Port); err != nil {
		return "", fmt.Errorf("impossible de se connecter à la base de données MySQL. Vérifiez que MAMP est démarré et que MySQL fonctionne sur le port %s: %v", database.Port, err)
	}
	fmt.Printf("[BACKUP] Database connectivity test passed\n")

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

	// Connection configuration
	if isMAMP {
		fmt.Printf("[BACKUP] Detected MAMP installation, using TCP connection (recommended for MAMP)\n")
		// For MAMP, prefer TCP connection as socket might not be accessible
		args = append(args, "-h", database.Host, "-P", database.Port)
	} else {
		// Standard TCP connection for non-MAMP installations
		args = append(args, "-h", database.Host, "-P", database.Port)
	}

	// Add SSL options for remote connections (not localhost)
	if database.Host != "localhost" && database.Host != "127.0.0.1" {
		args = append(args,
			"--ssl-mode=REQUIRED",
			"--ssl-ca=/etc/ssl/certs/ca-certificates.crt", // Linux
			"--ssl-ca=/etc/ssl/cert.pem",                  // macOS
		)
	}

	args = append(args, database.DbName)

	fmt.Printf("[BACKUP] Command args: %v\n", args)

	// Set timeout for the command (30 minutes max for large databases)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
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
	if err := cmd.Run(); err != nil {
		fmt.Printf("[BACKUP] mysqldump failed: %v\n", err)
		fmt.Printf("[BACKUP] stderr: %s\n", stderr.String())
		os.Remove(dumpFile)
		return "", fmt.Errorf("erreur lors de l'exécution de mysqldump: %v, stderr: %s", err, stderr.String())
	}

	fmt.Printf("[BACKUP] mysqldump completed successfully\n")

	return dumpFile, nil
}

// dumpPostgreSQL creates a PostgreSQL dump
func (s *BackupService) dumpPostgreSQL(database *models.Database, outputDir string) (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	dumpFile := filepath.Join(outputDir, fmt.Sprintf("%s_%s.sql", database.Name, timestamp))

	// Test network connectivity first
	if err := s.testDatabaseConnectivity(database.Host, database.Port); err != nil {
		return "", fmt.Errorf("test de connectivité échoué: %v", err)
	}

	// Find pg_dump executable
	pgDumpPath, err := s.findExecutable(s.getPostgreSQLDumpPaths())
	if err != nil {
		return "", fmt.Errorf("pg_dump non trouvé: %v", err)
	}

	// Set environment variable for password
	env := os.Environ()
	env = append(env, fmt.Sprintf("PGPASSWORD=%s", database.Password))

	// Build pg_dump command with SSL options for remote connections
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

	// Add SSL options for remote connections (not localhost)
	if database.Host != "localhost" && database.Host != "127.0.0.1" {
		args = append(args,
			"--sslmode=require",
		)
	}

	// Set timeout for the command (30 minutes max for large databases)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, pgDumpPath, args...)

	cmd.Env = env

	if err := cmd.Run(); err != nil {
		os.Remove(dumpFile)
		return "", fmt.Errorf("erreur lors de l'exécution de pg_dump: %v", err)
	}

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
