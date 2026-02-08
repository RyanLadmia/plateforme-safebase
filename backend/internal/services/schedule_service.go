package services

import (
	"fmt"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
	"github.com/robfig/cron/v3"
)

type ScheduleService struct {
	scheduleRepo         *repositories.ScheduleRepository
	databaseRepo         *repositories.DatabaseRepository
	backupService        *BackupService
	cronScheduler        *cron.Cron
	jobs                 map[uint]cron.EntryID // key: schedule ID
	actionHistoryService *ActionHistoryService
}

// Constructor of the ScheduleService
func NewScheduleService(scheduleRepo *repositories.ScheduleRepository, databaseRepo *repositories.DatabaseRepository, backupService *BackupService) *ScheduleService {
	s := &ScheduleService{
		scheduleRepo:  scheduleRepo,
		databaseRepo:  databaseRepo,
		backupService: backupService,
		cronScheduler: cron.New(),
		jobs:          make(map[uint]cron.EntryID),
	}
	return s
}

// SetActionHistoryService sets the action history service reference for logging
func (s *ScheduleService) SetActionHistoryService(actionHistoryService *ActionHistoryService) {
	s.actionHistoryService = actionHistoryService
}

// Start the cron scheduler
func (s *ScheduleService) StartScheduler() {
	s.cronScheduler.Start()
}

// GetSchedules returns all schedules for a user
func (s *ScheduleService) GetSchedules(userID uint) ([]models.Schedule, error) {
	return s.scheduleRepo.GetByUserID(userID)
}

// GetSchedule returns a specific schedule by ID
func (s *ScheduleService) GetSchedule(id uint, userID uint) (*models.Schedule, error) {
	schedule, err := s.scheduleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if schedule.UserId != userID {
		return nil, fmt.Errorf("accès non autorisé")
	}
	return schedule, nil
}

// LoadActiveSchedules loads all active schedules from DB and adds them to the cron scheduler
func (s *ScheduleService) LoadActiveSchedules() error {
	schedules, err := s.scheduleRepo.GetActive()
	if err != nil {
		return err
	}

	for _, schedule := range schedules {
		db, err := s.databaseRepo.GetByID(schedule.DatabaseId)
		if err != nil {
			fmt.Printf("Warning: database %d not found for schedule %d\n", schedule.DatabaseId, schedule.Id)
			continue
		}
		jobID, err := s.cronScheduler.AddFunc(schedule.CronExpression, func() {
			_, err := s.backupService.CreateBackup(db.Id, db.UserId, "127.0.0.1", "Scheduled Task")
			if err != nil {
				fmt.Printf("erreur lors de la création de la sauvegarde planifiée pour la base de données %s: %v\n", db.Name, err)
			} else {
				fmt.Printf("sauvegarde planifiée créée avec succès pour la base de données %s\n", db.Name)
			}
		})
		if err != nil {
			fmt.Printf("Warning: failed to schedule job for schedule %d: %v\n", schedule.Id, err)
			continue
		}
		s.jobs[schedule.Id] = jobID
	}
	return nil
}

// Logging methods for action history

// CreateSchedule creates a new schedule and logs the action
func (s *ScheduleService) CreateSchedule(databaseID uint, userID uint, name string, cronExpression string, ipAddress string, userAgent string) (*models.Schedule, error) {
	// Verify that the database exists and belongs to the user
	db, err := s.databaseRepo.GetByID(databaseID)
	if err != nil {
		return nil, fmt.Errorf("base de données introuvable: %v", err)
	}
	if db.UserId != userID {
		return nil, fmt.Errorf("accès non autorisé à cette base de données")
	}

	// Vérifier que l'expression cron est valide
	if _, err := cron.ParseStandard(cronExpression); err != nil {
		return nil, fmt.Errorf("expression cron invalide: %v", err)
	}

	// Create the schedule record
	schedule := &models.Schedule{
		Name:           name,
		CronExpression: cronExpression,
		Active:         true,
		UserId:         userID,
		DatabaseId:     databaseID,
	}
	if err := s.scheduleRepo.Create(schedule); err != nil {
		return nil, fmt.Errorf("erreur lors de la création du schedule: %v", err)
	}

	// Reload the schedule with preloaded relations
	schedule, err = s.scheduleRepo.GetByID(schedule.Id)
	if err != nil {
		return nil, fmt.Errorf("erreur lors du rechargement du schedule: %v", err)
	}

	// Add to cron scheduler
	jobID, err := s.cronScheduler.AddFunc(cronExpression, func() {
		_, err := s.backupService.CreateBackup(db.Id, db.UserId, "127.0.0.1", "Scheduled Task")
		if err != nil {
			fmt.Printf("erreur lors de la création de la sauvegarde planifiée pour la base de données %s: %v\n", db.Name, err)
		} else {
			fmt.Printf("sauvegarde planifiée créée avec succès pour la base de données %s\n", db.Name)
		}
	})
	if err != nil {
		// If adding to cron fails, delete the schedule from DB
		if deleteErr := s.scheduleRepo.Delete(schedule.Id); deleteErr != nil {
			// Log the error but don't fail the operation
			fmt.Printf("Warning: failed to delete schedule after cron error: %v\n", deleteErr)
		}
		return nil, fmt.Errorf("erreur lors de l'ajout de la tâche cron: %v", err)
	}
	s.jobs[schedule.Id] = jobID

	// Log the action
	if s.actionHistoryService != nil {
		metadata := map[string]interface{}{
			"schedule_id":     schedule.Id,
			"database_id":     schedule.DatabaseId,
			"database_name":   db.Name,
			"schedule_name":   schedule.Name,
			"cron_expression": schedule.CronExpression,
			"active":          schedule.Active,
		}
		description := fmt.Sprintf("Planification créée pour la base de données '%s' (%s)", db.Name, schedule.CronExpression)
		_ = s.actionHistoryService.LogAction(userID, "create", "schedule", schedule.Id, description, metadata, ipAddress, userAgent)
	}

	return schedule, nil
}

// UpdateSchedule updates a schedule and logs the action
func (s *ScheduleService) UpdateSchedule(id uint, userID uint, name string, cronExpression string, active *bool, ipAddress string, userAgent string) (*models.Schedule, error) {
	schedule, err := s.scheduleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if schedule.UserId != userID {
		return nil, fmt.Errorf("accès non autorisé")
	}

	// Capture l'état avant modification pour les changements
	oldName := schedule.Name
	oldCronExpression := schedule.CronExpression
	oldActive := schedule.Active

	// Validate and update name if provided
	if name != "" {
		schedule.Name = name
	}

	// Validate cron expression if provided
	if cronExpression != "" {
		if _, err := cron.ParseStandard(cronExpression); err != nil {
			return nil, fmt.Errorf("expression cron invalide: %v", err)
		}
		schedule.CronExpression = cronExpression
	}

	// Update active status only if explicitly provided
	if active != nil {
		schedule.Active = *active
	}

	// Remove old job if exists
	if entryID, exists := s.jobs[schedule.Id]; exists {
		s.cronScheduler.Remove(entryID)
		delete(s.jobs, schedule.Id)
	}

	// Add new job if active (either explicitly set to true, or not changed and was already active)
	jobShouldBeActive := schedule.Active
	if active != nil {
		jobShouldBeActive = *active
	}

	if jobShouldBeActive {
		db, err := s.databaseRepo.GetByID(schedule.DatabaseId)
		if err != nil {
			return nil, fmt.Errorf("base de données introuvable: %v", err)
		}
		jobID, err := s.cronScheduler.AddFunc(schedule.CronExpression, func() {
			_, err := s.backupService.CreateBackup(db.Id, db.UserId, "127.0.0.1", "Scheduled Task")
			if err != nil {
				fmt.Printf("erreur lors de la création de la sauvegarde planifiée pour la base de données %s: %v\n", db.Name, err)
			} else {
				fmt.Printf("sauvegarde planifiée créée avec succès pour la base de données %s\n", db.Name)
			}
		})
		if err != nil {
			return nil, fmt.Errorf("erreur lors de la rescheduling: %v", err)
		}
		s.jobs[schedule.Id] = jobID
	}

	if err := s.scheduleRepo.Update(schedule); err != nil {
		return nil, fmt.Errorf("erreur lors de la mise à jour: %v", err)
	}

	// Log the action with change details
	if s.actionHistoryService != nil {
		db, err := s.databaseRepo.GetByID(schedule.DatabaseId)
		if err != nil {
			// If database not found, log with available info
			metadata := map[string]interface{}{
				"schedule_id":     schedule.Id,
				"database_id":     schedule.DatabaseId,
				"database_name":   db.Name,
				"schedule_name":   schedule.Name,
				"cron_expression": schedule.CronExpression,
				"active":          schedule.Active,
				"changes":         s.buildScheduleChanges(oldName, oldCronExpression, oldActive, schedule.Name, schedule.CronExpression, schedule.Active),
			}
			s.actionHistoryService.LogAction(userID, "update", "schedule", schedule.Id, "Planification modifiée", metadata, ipAddress, userAgent)
		} else {
			metadata := map[string]interface{}{
				"schedule_id":     schedule.Id,
				"database_id":     schedule.DatabaseId,
				"database_name":   db.Name,
				"schedule_name":   schedule.Name,
				"cron_expression": schedule.CronExpression,
				"active":          schedule.Active,
				"changes":         s.buildScheduleChanges(oldName, oldCronExpression, oldActive, schedule.Name, schedule.CronExpression, schedule.Active),
			}
			description := fmt.Sprintf("Planification modifiée pour la base de données '%s'", db.Name)
			s.actionHistoryService.LogAction(userID, "update", "schedule", schedule.Id, description, metadata, ipAddress, userAgent)
		}
	}

	return schedule, nil
}

// buildScheduleChanges builds the changes metadata for schedule updates
func (s *ScheduleService) buildScheduleChanges(oldName string, oldCronExpression string, oldActive bool, newName string, newCronExpression string, newActive bool) map[string]interface{} {
	changes := make(map[string]interface{})

	// Check if name changed
	if oldName != newName {
		changes["name"] = map[string]interface{}{
			"from": oldName,
			"to":   newName,
		}
	}

	// Check if cron expression changed
	if oldCronExpression != newCronExpression {
		changes["cron_expression"] = map[string]interface{}{
			"from": oldCronExpression,
			"to":   newCronExpression,
		}
	}

	// Check if active status changed
	if oldActive != newActive {
		changes["enabled"] = map[string]interface{}{
			"from": oldActive,
			"to":   newActive,
		}
	}

	// Return nil if no changes
	if len(changes) == 0 {
		return nil
	}

	return changes
}

// DeleteSchedule deletes a schedule and logs the action
func (s *ScheduleService) DeleteSchedule(id uint, userID uint, ipAddress string, userAgent string) error {
	schedule, err := s.scheduleRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Verify user ownership
	if schedule.UserId != userID {
		return fmt.Errorf("accès non autorisé")
	}

	// Get database info for logging
	db, err := s.databaseRepo.GetByID(schedule.DatabaseId)
	if err != nil {
		return fmt.Errorf("base de données introuvable: %v", err)
	}

	// Remove from cron if exists
	if entryID, exists := s.jobs[schedule.Id]; exists {
		s.cronScheduler.Remove(entryID)
		delete(s.jobs, schedule.Id)
	}

	err = s.scheduleRepo.Delete(id)
	if err != nil {
		return err
	}

	// Log the action
	if s.actionHistoryService != nil {
		metadata := map[string]interface{}{
			"schedule_id":     schedule.Id,
			"database_id":     schedule.DatabaseId,
			"database_name":   db.Name,
			"schedule_name":   schedule.Name,
			"cron_expression": schedule.CronExpression,
			"active":          schedule.Active,
		}
		description := fmt.Sprintf("Planification supprimée pour la base de données '%s' (%s)", db.Name, schedule.CronExpression)
		s.actionHistoryService.LogAction(userID, "delete", "schedule", schedule.Id, description, metadata, ipAddress, userAgent)
	}

	return nil
}
