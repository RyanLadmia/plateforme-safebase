package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOService handles file storage operations with MinIO
type MinIOService struct {
	client     *minio.Client
	bucketName string
	endpoint   string
}

// MinIOConfig holds MinIO configuration
type MinIOConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	UseSSL          bool
}

// NewMinIOService creates a new MinIO service instance
func NewMinIOService(config MinIOConfig) (*MinIOService, error) {
	// Initialize MinIO client
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
		// Region should be specified in Options.Region if needed by your MinIO setup.
		// Do NOT pass the region as the session token (3rd param of NewStaticV4).
		Region: "",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %v", err)
	}

	service := &MinIOService{
		client:     client,
		bucketName: config.BucketName,
		endpoint:   config.Endpoint,
	}

	// Ensure bucket exists
	if err := service.ensureBucketExists(); err != nil {
		log.Printf("Warning: failed to ensure bucket exists: %v", err)
		// Don't return error, continue with the assumption that bucket exists
	}

	return service, nil
}

// ensureBucketExists creates the bucket if it doesn't exist
func (s *MinIOService) ensureBucketExists() error {
	ctx := context.Background()

	exists, err := s.client.BucketExists(ctx, s.bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %v", err)
	}

	if !exists {
		err = s.client.MakeBucket(ctx, s.bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %v", err)
		}
		log.Printf("Created MinIO bucket: %s", s.bucketName)
	}

	return nil
}

// UploadFile uploads a file to MinIO
func (s *MinIOService) UploadFile(filePath, objectName string, fileReader io.Reader, fileSize int64) error {
	ctx := context.Background()

	// Determine content type based on file extension
	contentType := s.getContentType(objectName)

	// Upload the file
	_, err := s.client.PutObject(ctx, s.bucketName, objectName, fileReader, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to MinIO: %v", err)
	}

	log.Printf("Successfully uploaded file to MinIO: %s", objectName)
	return nil
}

// UploadFileFromPath uploads a file from local path to MinIO
func (s *MinIOService) UploadFileFromPath(localPath, objectName string) error {
	// Open local file
	file, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("failed to read local file: %v", err)
	}

	reader := bytes.NewReader(file)
	return s.UploadFile(localPath, objectName, reader, int64(len(file)))
}

// DownloadFile downloads a file from MinIO
func (s *MinIOService) DownloadFile(objectName string) ([]byte, error) {
	ctx := context.Background()

	object, err := s.client.GetObject(ctx, s.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from MinIO: %v", err)
	}
	defer object.Close()

	data, err := io.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("failed to read object data: %v", err)
	}

	return data, nil
}

// DownloadFileToWriter downloads a file from MinIO and writes it to a writer
func (s *MinIOService) DownloadFileToWriter(objectName string, writer io.Writer) error {
	ctx := context.Background()

	object, err := s.client.GetObject(ctx, s.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to get object from MinIO: %v", err)
	}
	defer object.Close()

	_, err = io.Copy(writer, object)
	if err != nil {
		return fmt.Errorf("failed to copy object data: %v", err)
	}

	return nil
}

// DeleteFile deletes a file from MinIO
func (s *MinIOService) DeleteFile(objectName string) error {
	ctx := context.Background()

	err := s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object from MinIO: %v", err)
	}

	log.Printf("Successfully deleted file from MinIO: %s", objectName)
	return nil
}

// FileExists checks if a file exists in MinIO
func (s *MinIOService) FileExists(objectName string) (bool, error) {
	ctx := context.Background()

	_, err := s.client.StatObject(ctx, s.bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("failed to stat object: %v", err)
	}

	return true, nil
}

// GetFileInfo gets information about a file in MinIO
func (s *MinIOService) GetFileInfo(objectName string) (*minio.ObjectInfo, error) {
	ctx := context.Background()

	info, err := s.client.StatObject(ctx, s.bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object info: %v", err)
	}

	return &info, nil
}

// GenerateObjectName generates a unique object name for MinIO storage
func (s *MinIOService) GenerateObjectName(username, dbType, filename string) string {
	return fmt.Sprintf("%s/%s/%s", username, dbType, filename)
}

// getContentType determines content type based on file extension
func (s *MinIOService) getContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".sql":
		return "application/sql"
	case ".zip":
		return "application/zip"
	case ".gz":
		return "application/gzip"
	case ".tar":
		return "application/x-tar"
	default:
		return "application/octet-stream"
	}
}

// UploadMultipartFile uploads a multipart file to MinIO
func (s *MinIOService) UploadMultipartFile(file multipart.File, header *multipart.FileHeader, objectName string) error {
	defer file.Close()

	// Read file content
	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read multipart file: %v", err)
	}

	reader := bytes.NewReader(data)
	return s.UploadFile(header.Filename, objectName, reader, int64(len(data)))
}

// ListFiles lists files in a specific path
func (s *MinIOService) ListFiles(prefix string) ([]string, error) {
	ctx := context.Background()

	var files []string
	objectCh := s.client.ListObjects(ctx, s.bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("failed to list objects: %v", object.Err)
		}
		files = append(files, object.Key)
	}

	return files, nil
}
