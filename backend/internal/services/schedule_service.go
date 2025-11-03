package services

import (
	"fmt"

	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
	"github.com/robfig/cron/v3"
)

type ScheduleService struct {
	scheduleRepo  *repositories.ScheduleRepository
	databaseRepo  *repositories.DatabaseRepository
	backupService *BackupService
	cronScheduler *cron.Cron
	jobs          map[uint]cron.EntryID
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

// CreateSchedule creates a new schedule and adds it to the cron scheduler
func (s *ScheduleService) CreateCron(databaseID uint, cronExpression string) error {
	// Verify that the database exists
	db, err := s.databaseRepo.GetByID(databaseID)
	if err != nil {
		return fmt.Errorf("base de données introuvable: %v", err)
	}

	// Vérifier que l'expression cron est valide
	if _, err := cron.ParseStandard(cronExpression); err != nil {
		return fmt.Errorf("expression cron invalide: %v", err)
	}

	// Create the schedule record
	jobID, err := s.cronScheduler.AddFunc(cronExpression, func() {
		_, err := s.backupService.CreateBackup(db.Id, db.UserId)
		if err != nil {
			fmt.Printf("erreur lors de la création de la sauvegarde planifiée pour la base de données %s: %v\n", db.Name, err)
		} else {
			fmt.Printf("sauvegarde planifiée créée avec succès pour la base de données %s\n", db.Name)
		}
	})
	if err != nil {
		return fmt.Errorf("erreur lors de l'ajout de la tâche cron: %v", err)
	}
	s.jobs[databaseID] = jobID
	return nil
}
