package services

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
)

type BackupService struct {
	backupRepo   *repositories.BackupRepository
	databaseRepo *repositories.DatabaseRepository
	backupDir    string
}

// Constructor for BackupService
func NewBackupService(backupRepo *repositories.BackupRepository, databaseRepo *repositories.DatabaseRepository, backupDir string) *BackupService {
	return &BackupService{
		backupRepo:   backupRepo,
		databaseRepo: databaseRepo,
		backupDir:    backupDir,
	}
}

// CreateBackup creates a backup for the specified database
func (s *BackupService) CreateBackup(databaseID uint, userID uint) (*models.Backup, error) {
	// Get database info
	database, err := s.databaseRepo.GetByID(databaseID)
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
	var dumpFile string
	var err error

	// Create backup directory if it doesn't exist
	backupTypeDir := filepath.Join(s.backupDir, database.Type)
	if err := os.MkdirAll(backupTypeDir, 0755); err != nil {
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
		s.updateBackupError(backup.Id, fmt.Sprintf("Erreur lors du dump: %v", err))
		return
	}

	// Create ZIP file
	zipFile := backup.Filepath
	if err := s.createZipFile(dumpFile, zipFile); err != nil {
		s.updateBackupError(backup.Id, fmt.Sprintf("Erreur lors de la compression: %v", err))
		// Clean up dump file
		os.Remove(dumpFile)
		return
	}

	// Clean up dump file
	os.Remove(dumpFile)

	// Get file size
	fileInfo, err := os.Stat(zipFile)
	if err != nil {
		s.updateBackupError(backup.Id, fmt.Sprintf("Erreur lors de la lecture des informations du fichier: %v", err))
		return
	}

	// Update backup record
	if err := s.backupRepo.UpdateSize(backup.Id, fileInfo.Size()); err != nil {
		s.updateBackupError(backup.Id, fmt.Sprintf("Erreur lors de la mise à jour de la taille: %v", err))
		return
	}

	if err := s.backupRepo.UpdateStatus(backup.Id, "completed", ""); err != nil {
		s.updateBackupError(backup.Id, fmt.Sprintf("Erreur lors de la mise à jour du statut: %v", err))
		return
	}
}

// dumpMySQL creates a MySQL dump
func (s *BackupService) dumpMySQL(database *models.Database, outputDir string) (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	dumpFile := filepath.Join(outputDir, fmt.Sprintf("%s_%s.sql", database.Name, timestamp))

	cmd := exec.Command("mysqldump",
		"-h", database.Host,
		"-P", database.Port,
		"-u", database.Username,
		fmt.Sprintf("-p%s", database.Password),
		"--single-transaction",
		"--routines",
		"--triggers",
		database.DbName,
	)

	output, err := os.Create(dumpFile)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la création du fichier de dump: %v", err)
	}
	defer output.Close()

	cmd.Stdout = output
	if err := cmd.Run(); err != nil {
		os.Remove(dumpFile)
		return "", fmt.Errorf("erreur lors de l'exécution de mysqldump: %v", err)
	}

	return dumpFile, nil
}

// dumpPostgreSQL creates a PostgreSQL dump
func (s *BackupService) dumpPostgreSQL(database *models.Database, outputDir string) (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	dumpFile := filepath.Join(outputDir, fmt.Sprintf("%s_%s.sql", database.Name, timestamp))

	// Set environment variable for password
	env := os.Environ()
	env = append(env, fmt.Sprintf("PGPASSWORD=%s", database.Password))

	cmd := exec.Command("pg_dump",
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
	)

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
