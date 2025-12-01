package utils

import (
	"fmt"

	"github.com/RyanLadmia/plateforme-safebase/internal/config"
)

// DisplayEndpoints prints all available API endpoints
func DisplayEndpoints(port string) {
	fmt.Printf(config.Cyan + "Available endpoints:\n")
	fmt.Printf("   GET  /test                              - Test endpoint\n")
	fmt.Printf("   POST /auth/register                     - User registration\n")
	fmt.Printf("   POST /auth/login                        - User login\n")
	fmt.Printf("   POST /auth/logout                       - User logout\n")
	fmt.Printf("   GET  /auth/me                           - Get current user\n")
	fmt.Printf("   POST /api/databases                     - Create database\n")
	fmt.Printf("   GET  /api/databases                     - Get user databases\n")
	fmt.Printf("   GET  /api/databases/:id                 - Get database by ID\n")
	fmt.Printf("   PUT  /api/databases/:id                 - Update database\n")
	fmt.Printf("   DELETE /api/databases/:id               - Delete database\n")
	fmt.Printf("   POST /api/backups/database/:database_id - Create backup\n")
	fmt.Printf("   GET  /api/backups                       - Get user backups\n")
	fmt.Printf("   GET  /api/backups/:id                   - Get backup by ID\n")
	fmt.Printf("   GET  /api/backups/:id/download          - Download backup\n")
	fmt.Printf("   DELETE /api/backups/:id                 - Delete backup\n")
	fmt.Printf("   POST /api/schedules                     - Create schedule\n")
	fmt.Printf("   GET  /api/schedules                     - Get user schedules\n")
	fmt.Printf("   GET  /api/schedules/:id                 - Get schedule by ID\n")
	fmt.Printf("   PUT  /api/schedules/:id                 - Update schedule\n")
	fmt.Printf("   DELETE /api/schedules/:id               - Delete schedule\n")
	fmt.Printf("   GET  /api/admin/users                   - Get all users (admin)\n")
	fmt.Printf("   GET  /api/admin/users/active            - Get active users (admin)\n")
	fmt.Printf("   GET  /api/admin/users/:id               - Get user by ID (admin)\n")
	fmt.Printf("   PUT  /api/admin/users/:id               - Update user (admin)\n")
	fmt.Printf("   PUT  /api/admin/users/:id/role          - Change user role (admin)\n")
	fmt.Printf("   PUT  /api/admin/users/:id/deactivate    - Deactivate user (admin)\n")
	fmt.Printf("   PUT  /api/admin/users/:id/activate      - Activate user (admin)\n" + config.Reset)
}
