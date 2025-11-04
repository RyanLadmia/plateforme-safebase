package services

import (
	"fmt"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
	"github.com/robfig/cron/v3"
)

type ScheduleService struct {
	scheduleRepo  *repositories.ScheduleRepository
	databaseRepo  *repositories.DatabaseRepository
	backupService *BackupService
	cronScheduler *cron.Cron
	jobs          map[uint]cron.EntryID // key: schedule ID
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

// Start the cron scheduler
func (s *ScheduleService) StartScheduler() {
	s.cronScheduler.Start()
}

// CreateSchedule creates a new schedule in DB and adds it to the cron scheduler
func (s *ScheduleService) CreateSchedule(databaseID uint, userID uint, cronExpression string) (*models.Schedule, error) {
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
		_, err := s.backupService.CreateBackup(db.Id, db.UserId)
		if err != nil {
			fmt.Printf("erreur lors de la création de la sauvegarde planifiée pour la base de données %s: %v\n", db.Name, err)
		} else {
			fmt.Printf("sauvegarde planifiée créée avec succès pour la base de données %s\n", db.Name)
		}
	})
	if err != nil {
		// If adding to cron fails, delete the schedule from DB
		s.scheduleRepo.Delete(schedule.Id)
		return nil, fmt.Errorf("erreur lors de l'ajout de la tâche cron: %v", err)
	}
	s.jobs[schedule.Id] = jobID
	return schedule, nil
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

// UpdateSchedule updates a schedule and reschedules the cron job
func (s *ScheduleService) UpdateSchedule(id uint, userID uint, cronExpression string, active bool) (*models.Schedule, error) {
	schedule, err := s.scheduleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if schedule.UserId != userID {
		return nil, fmt.Errorf("accès non autorisé")
	}

	// Validate cron expression if provided
	if cronExpression != "" {
		if _, err := cron.ParseStandard(cronExpression); err != nil {
			return nil, fmt.Errorf("expression cron invalide: %v", err)
		}
		schedule.CronExpression = cronExpression
	}
	schedule.Active = active

	// Remove old job if exists
	if entryID, exists := s.jobs[schedule.Id]; exists {
		s.cronScheduler.Remove(entryID)
		delete(s.jobs, schedule.Id)
	}

	// Add new job if active
	if active {
		db, err := s.databaseRepo.GetByID(schedule.DatabaseId)
		if err != nil {
			return nil, fmt.Errorf("base de données introuvable: %v", err)
		}
		jobID, err := s.cronScheduler.AddFunc(schedule.CronExpression, func() {
			_, err := s.backupService.CreateBackup(db.Id, db.UserId)
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
	return schedule, nil
}

// DeleteSchedule deletes a schedule and removes it from cron
func (s *ScheduleService) DeleteSchedule(id uint, userID uint) error {
	schedule, err := s.scheduleRepo.GetByID(id)
	if err != nil {
		return err
	}
	if schedule.UserId != userID {
		return fmt.Errorf("accès non autorisé")
	}

	// Remove from cron if exists
	if entryID, exists := s.jobs[schedule.Id]; exists {
		s.cronScheduler.Remove(entryID)
		delete(s.jobs, schedule.Id)
	}

	return s.scheduleRepo.Delete(id)
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
			_, err := s.backupService.CreateBackup(db.Id, db.UserId)
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
