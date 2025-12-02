package services

import (
	"fmt"
	"os"
	"strings"

	"github.com/t3rm1n4l/go-mega"
)

// MegaService implements CloudStorageService for Mega.nz
type MegaService struct {
	email    string
	password string
	mega     *mega.Mega
}

// MegaConfig holds Mega configuration
type MegaConfig struct {
	Email    string
	Password string
}

// NewMegaService creates a new Mega service
func NewMegaService(config MegaConfig) (*MegaService, error) {
	service := &MegaService{
		email:    config.Email,
		password: config.Password,
		mega:     mega.New(),
	}

	// Login to Mega
	err := service.mega.Login(config.Email, config.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to login to Mega: %v", err)
	}

	fmt.Printf("[MEGA] Successfully logged in as %s\n", config.Email)
	return service, nil
}

// UploadFile uploads a file to Mega
func (m *MegaService) UploadFile(localPath, remotePath string) error {
	// Parse path to create folder structure
	parts := strings.Split(strings.Trim(remotePath, "/"), "/")
	filename := parts[len(parts)-1]
	folderPath := strings.Join(parts[:len(parts)-1], "/")

	// Ensure folder structure exists
	parentNode, err := m.ensureFolderPath(folderPath)
	if err != nil {
		return fmt.Errorf("failed to create folder structure: %v", err)
	}

	// Upload file using the library
	_, err = m.mega.UploadFile(localPath, parentNode, filename, nil)
	if err != nil {
		return fmt.Errorf("failed to upload file to Mega: %v", err)
	}

	fmt.Printf("[MEGA] Successfully uploaded file: %s\n", remotePath)
	return nil
}

// DownloadFile downloads a file from Mega
func (m *MegaService) DownloadFile(remotePath string) ([]byte, error) {
	// Find the file node
	node, err := m.findNodeByPath(remotePath)
	if err != nil {
		return nil, fmt.Errorf("file not found: %v", err)
	}

	// Create a temporary file for download
	tmpFile, err := os.CreateTemp("", "mega_download_*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Download file to temp file
	err = m.mega.DownloadFile(node, tmpFile.Name(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to download file from Mega: %v", err)
	}

	// Read the downloaded file
	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to read downloaded file: %v", err)
	}

	fmt.Printf("[MEGA] Successfully downloaded file: %s\n", remotePath)
	return data, nil
}

// DeleteFile deletes a file from Mega
func (m *MegaService) DeleteFile(remotePath string) error {
	node, err := m.findNodeByPath(remotePath)
	if err != nil {
		return fmt.Errorf("file not found: %v", err)
	}

	err = m.mega.Delete(node, false)
	if err != nil {
		return fmt.Errorf("failed to delete file from Mega: %v", err)
	}

	fmt.Printf("[MEGA] Successfully deleted file: %s\n", remotePath)
	return nil
}

// FileExists checks if a file exists in Mega
func (m *MegaService) FileExists(remotePath string) (bool, error) {
	_, err := m.findNodeByPath(remotePath)
	return err == nil, nil
}

// GenerateRemotePath generates a remote path for Mega storage
func (m *MegaService) GenerateRemotePath(username, dbType, filename string) string {
	// Format: "Ryan-Ladmia/mysql/filename.zip" or "Ryan-Ladmia/postgresql/filename.zip"
	sanitizedUsername := strings.ReplaceAll(username, " ", "-")
	return fmt.Sprintf("%s/%s/%s", sanitizedUsername, dbType, filename)
}

// Helper methods

func (m *MegaService) ensureFolderPath(folderPath string) (*mega.Node, error) {
	if folderPath == "" {
		return m.mega.FS.GetRoot(), nil
	}

	folders := strings.Split(folderPath, "/")
	currentNode := m.mega.FS.GetRoot()

	for _, folderName := range folders {
		if folderName == "" {
			continue
		}

		// Check if folder exists
		folderNode, err := m.findFolderByName(folderName, currentNode)
		if err != nil {
			// Create folder
			folderNode, err = m.mega.CreateDir(folderName, currentNode)
			if err != nil {
				return nil, fmt.Errorf("failed to create folder %s: %v", folderName, err)
			}
		}

		currentNode = folderNode
	}

	return currentNode, nil
}

func (m *MegaService) findFolderByName(name string, parent *mega.Node) (*mega.Node, error) {
	children, err := m.mega.FS.GetChildren(parent)
	if err != nil {
		return nil, err
	}

	for _, child := range children {
		if child.GetName() == name && child.GetType() == mega.FOLDER {
			return child, nil
		}
	}

	return nil, fmt.Errorf("folder not found")
}

func (m *MegaService) findNodeByPath(path string) (*mega.Node, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	currentNode := m.mega.FS.GetRoot()

	for i, part := range parts {
		isLast := i == len(parts)-1

		children, err := m.mega.FS.GetChildren(currentNode)
		if err != nil {
			return nil, err
		}

		found := false
		for _, child := range children {
			expectedType := mega.FILE
			if !isLast {
				expectedType = mega.FOLDER
			}

			if child.GetName() == part && child.GetType() == expectedType {
				currentNode = child
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("path not found: %s", part)
		}
	}

	return currentNode, nil
}
