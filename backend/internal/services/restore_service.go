package services

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
)

// WorkerPool interface for background tasks
type RestoreWorkerPoolInterface interface {
	Submit(task func())
}

type RestoreService struct {
	restoreRepo          *repositories.RestoreRepository
	backupService        *BackupService
	databaseService      *DatabaseService
	userService          *UserService
	workerPool           RestoreWorkerPoolInterface
	actionHistoryService *ActionHistoryService
}

// Constructor for RestoreService
func NewRestoreService(restoreRepo *repositories.RestoreRepository, backupService *BackupService, databaseService *DatabaseService, userService *UserService) *RestoreService {
	return &RestoreService{
		restoreRepo:     restoreRepo,
		backupService:   backupService,
		databaseService: databaseService,
		userService:     userService,
	}
}

// SetActionHistoryService sets the action history service reference for logging
func (s *RestoreService) SetActionHistoryService(actionHistoryService *ActionHistoryService) {
	s.actionHistoryService = actionHistoryService
}

// SetWorkerPool sets the worker pool for background tasks
func (s *RestoreService) SetWorkerPool(workerPool RestoreWorkerPoolInterface) {
	s.workerPool = workerPool
}

// unzipBackupData extracts the SQL file content from the zipped backup data
func (s *RestoreService) unzipBackupData(zipData []byte) ([]byte, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'ouverture du fichier ZIP: %v", err)
	}

	// Find the SQL file in the ZIP (should be the only file)
	for _, file := range zipReader.File {
		if strings.HasSuffix(file.Name, ".sql") {
			// Open the SQL file
			rc, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("erreur lors de l'ouverture du fichier SQL dans le ZIP: %v", err)
			}
			defer rc.Close()

			// Read the content
			sqlData, err := io.ReadAll(rc)
			if err != nil {
				return nil, fmt.Errorf("erreur lors de la lecture du fichier SQL: %v", err)
			}

			return sqlData, nil
		}
	}

	return nil, fmt.Errorf("fichier SQL non trouvé dans l'archive ZIP")
}

// executeRestoreAsync executes the restore process asynchronously
func (s *RestoreService) executeRestoreAsync(restore *models.Restore, backup *models.Backup, database *models.Database) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[RESTORE] Panic in restore process for ID %d: %v\n", restore.Id, r)
			s.updateRestoreError(restore.Id)
		}
	}()

	fmt.Printf("[RESTORE] Starting asynchronous restore process for ID: %d\n", restore.Id)

	// Update status to running
	if err := s.restoreRepo.UpdateStatus(restore.Id, "running"); err != nil {
		fmt.Printf("[RESTORE] Failed to update status to running: %v\n", err)
		return
	}

	// Download and decrypt the backup file
	backupData, err := s.backupService.DownloadBackup(backup.Id, backup.UserId)
	if err != nil {
		fmt.Printf("[RESTORE] Failed to download backup file: %v\n", err)
		s.updateRestoreError(restore.Id)
		return
	}

	// Unzip the backup data to extract the SQL file
	sqlData, err := s.unzipBackupData(backupData)
	if err != nil {
		fmt.Printf("[RESTORE] Failed to unzip backup data: %v\n", err)
		s.updateRestoreError(restore.Id)
		return
	}

	// Execute restore based on database type
	switch dbType := database.Type; dbType {
	case "mysql":
		err = s.restoreMySQL(database, sqlData)
	case "postgresql":
		err = s.restorePostgreSQL(database, sqlData)
	default:
		err = fmt.Errorf("type de base de données non supporté: %s", dbType)
	}

	if err != nil {
		fmt.Printf("[RESTORE] Restore failed for ID %d: %v\n", restore.Id, err)
		s.updateRestoreError(restore.Id)
		return
	}

	// Update restore record with success
	if err := s.restoreRepo.UpdateStatus(restore.Id, "success"); err != nil {
		fmt.Printf("[RESTORE] Failed to update restore status: %v\n", err)
		return
	}

	fmt.Printf("[RESTORE] Restore process completed successfully for ID: %d\n", restore.Id)
}

// restoreMySQL restores a MySQL database from backup data
func (s *RestoreService) restoreMySQL(database *models.Database, backupData []byte) error {
	fmt.Printf("[RESTORE] Starting MySQL restore for database %s\n", database.DbName)

	// Find mysql executable
	mysqlPaths := []string{
		"/Applications/MAMP/Library/bin/mysql80/bin/mysql",
		"/Applications/MAMP/Library/bin/mysql",
		"/usr/local/mysql/bin/mysql",
		"/usr/local/bin/mysql",
		"/opt/homebrew/bin/mysql",
		"/usr/bin/mysql",
		"mysql",
	}

	var mysqlPath string
	for _, path := range mysqlPaths {
		if _, err := exec.LookPath(path); err == nil {
			mysqlPath = path
			break
		}
	}

	if mysqlPath == "" {
		return fmt.Errorf("mysql client non trouvé")
	}

	fmt.Printf("[RESTORE] Using mysql at: %s\n", mysqlPath)

	// Build mysql command
	args := []string{
		"-u", database.Username,
		"-h", database.Host,
		"-P", database.Port,
		database.DbName,
	}

	// Set environment variable for password
	env := []string{fmt.Sprintf("MYSQL_PWD=%s", database.Password)}

	// Execute mysql command with backup data as input
	cmd := exec.Command(mysqlPath, args...)
	cmd.Env = append(cmd.Env, env...)
	cmd.Stdin = strings.NewReader(string(backupData))

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("[RESTORE] MySQL restore failed: %v\n", err)
		fmt.Printf("[RESTORE] Output: %s\n", string(output))
		return fmt.Errorf("erreur lors de la restauration MySQL: %v, output: %s", err, string(output))
	}

	fmt.Printf("[RESTORE] MySQL restore completed successfully\n")
	return nil
}

// restorePostgreSQL restores a PostgreSQL database from backup data
func (s *RestoreService) restorePostgreSQL(database *models.Database, backupData []byte) error {
	fmt.Printf("[RESTORE] Starting PostgreSQL restore for database %s\n", database.DbName)

	// Find psql executable
	psqlPaths := []string{
		"/Applications/Postgres.app/Contents/Versions/latest/bin/psql",
		"/usr/local/pgsql/bin/psql",
		"/usr/local/bin/psql",
		"/opt/homebrew/bin/psql",
		"/usr/bin/psql",
		"psql",
	}

	var psqlPath string
	for _, path := range psqlPaths {
		if _, err := exec.LookPath(path); err == nil {
			psqlPath = path
			break
		}
	}

	if psqlPath == "" {
		return fmt.Errorf("psql client non trouvé")
	}

	fmt.Printf("[RESTORE] Using psql at: %s\n", psqlPath)

	// Build psql command
	args := []string{
		"-h", database.Host,
		"-p", database.Port,
		"-U", database.Username,
		"-d", database.DbName,
	}

	// Set environment variable for password
	env := []string{fmt.Sprintf("PGPASSWORD=%s", database.Password)}

	// Execute psql command with backup data as input
	cmd := exec.Command(psqlPath, args...)
	cmd.Env = append(cmd.Env, env...)
	cmd.Stdin = strings.NewReader(string(backupData))

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("[RESTORE] PostgreSQL restore failed: %v\n", err)
		fmt.Printf("[RESTORE] Output: %s\n", string(output))
		return fmt.Errorf("erreur lors de la restauration PostgreSQL: %v, output: %s", err, string(output))
	}

	fmt.Printf("[RESTORE] PostgreSQL restore completed successfully\n")
	return nil
}

// updateRestoreError updates restore status to failed
func (s *RestoreService) updateRestoreError(restoreID uint) {
	// For now, we'll just update the status to failed
	// In the future, we might want to add an error_msg field to the Restore model
	s.restoreRepo.UpdateStatus(restoreID, "failed")
}

// GetRestoresByUser returns all restores for a user
func (s *RestoreService) GetRestoresByUser(userID uint) ([]models.Restore, error) {
	return s.restoreRepo.GetByUserID(userID)
}

// GetRestoresByDatabase returns all restores for a database
func (s *RestoreService) GetRestoresByDatabase(databaseID uint) ([]models.Restore, error) {
	return s.restoreRepo.GetByDatabaseID(databaseID)
}

// GetRestoresByBackup returns all restores for a backup
func (s *RestoreService) GetRestoresByBackup(backupID uint) ([]models.Restore, error) {
	return s.restoreRepo.GetByBackupID(backupID)
}

// GetRestoreByID returns a restore by ID
func (s *RestoreService) GetRestoreByID(id uint) (*models.Restore, error) {
	return s.restoreRepo.GetByID(id)
}

// Logging methods for action history

// CreateRestore creates a restore operation and logs the action
func (s *RestoreService) CreateRestore(backupID uint, databaseID uint, userID uint, ipAddress string, userAgent string) (*models.Restore, error) {
	// Get backup info
	backup, err := s.backupService.GetBackupByID(backupID)
	if err != nil {
		return nil, fmt.Errorf("sauvegarde introuvable: %v", err)
	}

	// Verify that the backup belongs to the user
	if backup.UserId != userID {
		return nil, fmt.Errorf("accès non autorisé à cette sauvegarde")
	}

	// Get database info
	database, err := s.databaseService.GetDatabaseByID(databaseID)
	if err != nil {
		return nil, fmt.Errorf("base de données introuvable: %v", err)
	}

	// Verify that the database belongs to the user
	if database.UserId != userID {
		return nil, fmt.Errorf("accès non autorisé à cette base de données")
	}

	// Verify that backup is completed
	if backup.Status != "completed" {
		return nil, fmt.Errorf("la sauvegarde n'est pas encore terminée")
	}

	// Create restore record with pending status
	restore := &models.Restore{
		UserId:     userID,
		BackupId:   backupID,
		DatabaseId: databaseID,
		Status:     "pending",
	}

	if err := s.restoreRepo.Create(restore); err != nil {
		return nil, fmt.Errorf("échec de la création de l'enregistrement de restauration: %v", err)
	}

	// Execute restore asynchronously using worker pool
	if s.workerPool != nil {
		s.workerPool.Submit(func() {
			s.executeRestoreAsync(restore, backup, database)
		})
	} else {
		go s.executeRestoreAsync(restore, backup, database)
	}

	// Log the action
	if s.actionHistoryService != nil {
		metadata := map[string]interface{}{
			"restore_id":  restore.Id,
			"backup_id":   restore.BackupId,
			"database_id": restore.DatabaseId,
			"status":      restore.Status,
		}
		s.actionHistoryService.LogAction(userID, "create", "restore", restore.Id, "Restauration effectuée", metadata, ipAddress, userAgent)
	}

	return restore, nil
}
