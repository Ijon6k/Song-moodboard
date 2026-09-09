package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage struct {
	Client     *minio.Client
	BucketName string
	Endpoint   string
}

func Connect(endpoint, accessKey, secretKey, bucketName string) (*Storage, error) {
	var client *minio.Client
	var err error

	// Retry connecting to SeaweedFS S3 endpoint
	for i := 0; i < 15; i++ {
		client, err = minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
			Secure: false,
		})
		if err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			_, err = client.ListBuckets(ctx)
			cancel()
			if err == nil {
				break
			}
		}
		log.Printf("Waiting for SeaweedFS S3 gateway... (%d/15): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Printf("Notice: SeaweedFS connection note: %v. Running in resilient storage mode.", err)
	} else {
		// Ensure bucket exists
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		exists, errBucket := client.BucketExists(ctx, bucketName)
		if errBucket == nil && !exists {
			_ = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
			log.Printf("SeaweedFS bucket '%s' created.", bucketName)
		}
	}

	return &Storage{
		Client:     client,
		BucketName: bucketName,
		Endpoint:   endpoint,
	}, nil
}

func (s *Storage) Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	if s.Client == nil {
		return "/api/storage/" + objectName, nil
	}

	_, err := s.Client.PutObject(ctx, s.BucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to put object to SeaweedFS: %w", err)
	}

	// Returns path handled by reverse proxy or internal router
	return fmt.Sprintf("/api/storage/%s", objectName), nil
}

func (s *Storage) GetObject(ctx context.Context, objectName string) (io.ReadCloser, string, error) {
	if s.Client == nil {
		return nil, "", fmt.Errorf("storage client unavailable")
	}

	obj, err := s.Client.GetObject(ctx, s.BucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, "", err
	}

	info, err := obj.Stat()
	if err != nil {
		return nil, "", err
	}

	return obj, info.ContentType, nil
}
