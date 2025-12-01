package services

// CloudStorageService defines the interface for cloud storage operations
type CloudStorageService interface {
	// UploadFile uploads a file to cloud storage
	UploadFile(localPath, remotePath string) error

	// DownloadFile downloads a file from cloud storage
	DownloadFile(remotePath string) ([]byte, error)

	// DeleteFile deletes a file from cloud storage
	DeleteFile(remotePath string) error

	// FileExists checks if a file exists in cloud storage
	FileExists(remotePath string) (bool, error)

	// GenerateRemotePath generates a remote path for a file
	GenerateRemotePath(username, dbType, filename string) string
}
