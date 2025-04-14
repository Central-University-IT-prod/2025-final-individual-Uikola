package s3

import "context"

// FileDataType represents the structure of a file for upload or retrieval.
type FileDataType struct {
	FileName string
	Data     []byte
}

// Client defines the interface for interacting with an S3-compatible storage.
type Client interface {
	// CreateOne uploads a single object to the MinIO bucket.
	CreateOne(ctx context.Context, file FileDataType, contentType string) (string, string, error)

	// GetOne retrieves a single object from the MinIO bucket by its ID.
	GetOne(ctx context.Context, objectID string) (string, error)

	// DeleteOne deletes a single object from the MinIO bucket by its ID.
	DeleteOne(ctx context.Context, objectID string) error
}
