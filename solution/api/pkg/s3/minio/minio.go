package minio

import (
	"api/pkg/s3"
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	publicReadPolicy = `{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"Principal": "*",
			"Action": [
				"s3:GetObject"
			],
			"Resource": [
				"arn:aws:s3:::%s/*"
			]
		}
	]
}`
)

// Client represents a client for interacting with a MinIO S3.
type Client struct {
	mc             *minio.Client
	endpoint       string
	publicEndpoint string
	bucketName     string
}

// NewClient creates a new MinIO client and ensures the specified bucket exists.
// If the bucket does not exist, it will be created.
func NewClient(endpoint, publicEndpoint string, user, pass string, useSSL bool, bucketName string) (*Client, error) {
	ctx := context.Background()

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(user, pass, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	policy := fmt.Sprintf(publicReadPolicy, bucketName)
	err = client.SetBucketPolicy(ctx, bucketName, policy)
	if err != nil {
		return nil, err
	}

	return &Client{
		mc:             client,
		endpoint:       endpoint,
		publicEndpoint: publicEndpoint,
		bucketName:     bucketName,
	}, nil
}

// CreateOne uploads a single object to the MinIO bucket.
func (c *Client) CreateOne(ctx context.Context, file s3.FileDataType, contentType string) (string, string, error) {
	objectID := uuid.New().String()

	reader := bytes.NewReader(file.Data)

	_, err := c.mc.PutObject(ctx, c.bucketName, objectID, reader, int64(len(file.Data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", "", err
	}

	link, err := c.mc.PresignedGetObject(ctx, c.bucketName, objectID, time.Second*24*60*60, nil)
	if err != nil {
		return "", "", err
	}

	publicLink := strings.ReplaceAll(link.String(), c.endpoint, c.publicEndpoint)

	return strings.ReplaceAll(publicLink, "&", "%26"), objectID, nil
}

// GetOne retrieves a single object from the MinIO bucket by its ID.
func (c *Client) GetOne(ctx context.Context, objectID string) (string, error) {
	link, err := c.mc.PresignedGetObject(ctx, c.bucketName, objectID, time.Second*24*60*60, nil)
	if err != nil {
		return "", nil
	}

	publicLink := strings.ReplaceAll(link.String(), c.endpoint, c.publicEndpoint)

	return strings.ReplaceAll(publicLink, "&", "%26"), nil
}

// DeleteOne deletes a single object from the MinIO bucket by its ID.
func (c *Client) DeleteOne(ctx context.Context, objectID string) error {
	_ = c.mc.RemoveObject(ctx, c.bucketName, objectID, minio.RemoveObjectOptions{})
	return nil
}
