package main

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
)

// Image represents the image uploaded by the user.
type Image struct {
	// ImageUUID is the unique ID of the image.
	ImageUUID string

	// LastModified is the timestamp when the image was last modified.
	LastModified time.Time

	// FileName is the original file name
	FileName string

	// FileSize is the size of the file in bytes
	FileSize int64

	// ContentType is the MIME type of the file
	ContentType string

	// ProcessedAt is when the image was processed by our system
	ProcessedAt time.Time

	// Status indicates processing status (uploaded, processed, error)
	Status string

	// Tags for categorization
	Tags []string
}

// NewImage creates a new image with enhanced metadata.
func NewImage(fileName string, fileSize int64, lastModified time.Time) *Image {
	// Generate a new UUID for the image.
	id := uuid.New().String()

	// Determine content type based on file extension
	contentType := "application/octet-stream"
	if strings.HasSuffix(strings.ToLower(fileName), ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(strings.ToLower(fileName), ".jpg") || strings.HasSuffix(strings.ToLower(fileName), ".jpeg") {
		contentType = "image/jpeg"
	} else if strings.HasSuffix(strings.ToLower(fileName), ".gif") {
		contentType = "image/gif"
	}

	// Generate tags based on file name and properties
	tags := []string{"uploaded"}
	if fileSize > 1024*1024 { // > 1MB
		tags = append(tags, "large")
	} else {
		tags = append(tags, "small")
	}

	// Create an image with enhanced metadata
	image := &Image{
		ImageUUID:    id,
		LastModified: lastModified,
		FileName:     fileName,
		FileSize:     fileSize,
		ContentType:  contentType,
		ProcessedAt:  time.Now(),
		Status:       "processed",
		Tags:         tags,
	}

	return image
}

// Save inserts a newly generated image with enhanced metadata into the Postgres database.
func Save(c *Image, table string, dbpool *pgxpool.Pool, m *metrics, ctx context.Context) error {
	// Create a new CHILD span to record and trace the request.
	ctx, span := tracer.Start(ctx, "SQL INSERT")
	defer span.End()

	// Get the current time to record the duration of the request.
	now := time.Now()

	// Prepare the database query to insert a record with enhanced metadata.
	query := fmt.Sprintf(`INSERT INTO %s (
		image_uuid, last_modified, file_name, file_size, 
		content_type, processed_at, status, tags
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, table)

	// Convert tags to a comma-separated string for storage
	tagsStr := ""
	if len(c.Tags) > 0 {
		tagsStr = strings.Join(c.Tags, ",")
	}

	// Execute the query to create a new image record.
	_, err := dbpool.Exec(context.Background(), query,
		c.ImageUUID, c.LastModified, c.FileName, c.FileSize,
		c.ContentType, c.ProcessedAt, c.Status, tagsStr)
	if err != nil {
		return fmt.Errorf("dbpool.Exec failed: %w", err)
	}

	// Record the duration of the insert query.
	m.duration.With(prometheus.Labels{"op": "db"}).Observe(time.Since(now).Seconds())

	return nil
}

// download downloads S3 image and returns enhanced metadata.
func download(sess *session.Session, bucket string, key string, m *metrics, ctx context.Context) (*time.Time, int64, context.Context, error) {
	// Create a new CHILD span to record and trace the request.
	ctx, span := tracer.Start(ctx, "S3 GET")
	defer span.End()

	// Get the current time to record the duration of the request.
	now := time.Now()

	// Create a new S3 session.
	svc := s3.New(sess)

	// Prepare the request for the S3 bucket.
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	// Send the request to the S3 object store to download the image.
	output, err := svc.GetObject(input)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("svc.GetObject failed: %w", err)
	}

	// Read all the image bytes returned by AWS.
	data, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("io.ReadAll failed: %w", err)
	}

	// Get file size
	fileSize := int64(len(data))

	// Record the duration of the request to S3.
	m.duration.With(prometheus.Labels{"op": "s3"}).Observe(time.Since(now).Seconds())

	return output.LastModified, fileSize, ctx, nil
}
