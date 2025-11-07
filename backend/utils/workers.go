package utils

import (
	"log"
	"os"
	"time"

	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
)

// Worker pool for background tasks
type WorkerPool struct {
	workers   int
	taskQueue chan func()
}

func NewWorkerPool(workers int) *WorkerPool {
	return &WorkerPool{
		workers:   workers,
		taskQueue: make(chan func(), 100), // Buffer of 100 tasks
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		go wp.worker(i)
	}
}

func (wp *WorkerPool) worker(id int) {
	log.Printf("Worker %d started", id)
	for task := range wp.taskQueue {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Worker %d panic: %v", id, r)
				}
			}()
			task()
		}()
	}
}

func (wp *WorkerPool) Submit(task func()) {
	select {
	case wp.taskQueue <- task:
		// Task submitted successfully
	default:
		log.Println("Warning: task queue is full, executing synchronously")
		task()
	}
}

// StartSessionCleanupWorker starts the session cleanup worker
func StartSessionCleanupWorker(sessionRepo *repositories.SessionRepository) {
	log.Println("Starting session cleanup worker...")
	ticker := time.NewTicker(1 * time.Hour) // Clean every hour
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Running scheduled session cleanup...")
		if err := sessionRepo.DeleteExpiredSessions(); err != nil {
			log.Printf("Error during session cleanup: %v", err)
		} else {
			log.Println("Session cleanup completed successfully")
		}
	}
}

// StartBackupCleanupWorker starts the backup cleanup worker
func StartBackupCleanupWorker(backupRepo *repositories.BackupRepository, workerPool *WorkerPool) {
	log.Println("Starting backup cleanup worker...")
	ticker := time.NewTicker(24 * time.Hour) // Clean every 24 hours
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Running scheduled backup cleanup...")
		// Clean old backup files (older than 30 days)
		workerPool.Submit(func() {
			if err := CleanupOldBackupFiles(backupRepo); err != nil {
				log.Printf("Error during backup cleanup: %v", err)
			} else {
				log.Println("Backup cleanup completed successfully")
			}
		})
	}
}

// CleanupOldBackupFiles removes old backup files and records
func CleanupOldBackupFiles(backupRepo *repositories.BackupRepository) error {
	// Get old backups (older than 1 year)
	oldBackups, err := backupRepo.GetOldBackups(365) // 365 days = 1 year
	if err != nil {
		return err
	}

	for _, backup := range oldBackups {
		// Delete physical file
		if backup.Filepath != "" {
			if err := os.Remove(backup.Filepath); err != nil && !os.IsNotExist(err) {
				log.Printf("Warning: failed to delete backup file %s: %v", backup.Filepath, err)
			}
		}

		// Delete database record
		if err := backupRepo.Delete(backup.Id); err != nil {
			log.Printf("Warning: failed to delete backup record %d: %v", backup.Id, err)
		} else {
			log.Printf("Cleaned up old backup: %s", backup.Filename)
		}
	}

	return nil
}
