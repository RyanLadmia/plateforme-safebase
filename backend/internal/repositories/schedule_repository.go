package repositories

import (
	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"gorm.io/gorm"
)

// Keep an instance of the database connection
type ScheduleRepository struct {
	db *gorm.DB
}

// Constructor of the ScheduleRepository, can be used to create a new ScheduleRepository
func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

// Create a new schedule record
func (r *ScheduleRepository) Create(schedule *models.Schedule) error {
	return r.db.Create(schedule).Error
}

// Get all schedules
func (r *ScheduleRepository) GetAll() ([]models.Schedule, error) {
	var schedules []models.Schedule
	err := r.db.Preload("Database").Preload("User").Find(&schedules).Error
	return schedules, err
}

// Get all active schedules
func (r *ScheduleRepository) GetActive() ([]models.Schedule, error) {
	var schedules []models.Schedule
	err := r.db.Preload("Database").Preload("User").Where("active = ?", true).Find(&schedules).Error
	return schedules, err
}

// Get schedule by ID
func (r *ScheduleRepository) GetByID(id uint) (*models.Schedule, error) {
	var schedule models.Schedule
	err := r.db.Preload("Database").Preload("User").First(&schedule, id).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

// Get all schedules for a user
func (r *ScheduleRepository) GetByUserID(userID uint) ([]models.Schedule, error) {
	var schedules []models.Schedule
	err := r.db.Preload("Database").Where("user_id = ?", userID).Find(&schedules).Error
	return schedules, err
}

// Get all schedules for a databse
func (r *ScheduleRepository) GetByDatabaseID(databaseID uint) ([]models.Schedule, error) {
	var schedules []models.Schedule
	err := r.db.Preload("Database").Where("database_id = ?", databaseID).Find(&schedules).Error
	return schedules, err
}

// Update a Cron Expression
func (r *ScheduleRepository) UpdateCronExpression(id uint, cronExpression string) error {
	return r.db.Model(&models.Schedule{}).Where("id = ?", id).Update("cron_expression", cronExpression).Error
}

// Update an status
func (r *ScheduleRepository) UpdateStatus(id uint, active bool) error {
	return r.db.Model(&models.Schedule{}).Where("id = ?", id).Update("active", active).Error
}

// Delete a schedule
func (r *ScheduleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Schedule{}, id).Error
}
