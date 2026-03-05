package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	"cloud.google.com/go/storage"
	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/config"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func UploadToGCS(filePath string, bucket, project string) error {
	ctx := context.Background()

	var client *storage.Client
	var err error

	client, err = storage.NewClient(ctx)
	if err != nil {
		if _, statErr := os.Stat("credentials.json"); statErr == nil {
			client, err = storage.NewClient(ctx, option.WithCredentialsFile("credentials.json"))
		}
		if err != nil {
			return fmt.Errorf("failed to create GCS client: %v", err)
		}
	}
	defer client.Close()

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	objectName := fmt.Sprintf("%s_backup_%s.sql", filepath.Base(filePath), time.Now().Format("20060102_150405"))

	writer := client.Bucket(bucket).Object(objectName).NewWriter(ctx)
	writer.ObjectAttrs.ContentType = "application/sql"

	if _, err := io.Copy(writer, file); err != nil {
		return fmt.Errorf("failed to upload to GCS: %v", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close GCS writer: %v", err)
	}

	zap.L().Sugar().Infof("Uploaded %s to GCS bucket %s as %s", filePath, bucket, objectName)
	return nil
}

func EnforceRetentionGCS(cfg config.StorageConfig, dryRun bool) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	bucket := client.Bucket(cfg.Bucket)
	it := bucket.Objects(ctx, &storage.Query{Prefix: ""})

	type blobInfo struct {
		name    string
		updated time.Time
	}
	var backupBlobs []blobInfo
	for {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		if len(obj.Name) >= 4 && obj.Name[len(obj.Name)-4:] == ".sql" {
			backupBlobs = append(backupBlobs, blobInfo{
				name:    obj.Name,
				updated: obj.Updated,
			})
		}
	}

	sort.Slice(backupBlobs, func(i, j int) bool {
		return backupBlobs[i].updated.After(backupBlobs[j].updated)
	})

	if dryRun {
		zap.L().Sugar().Infof("Would keep %d newest backups in GCS bucket %s", cfg.Retain, cfg.Bucket)
		return nil
	}

	if len(backupBlobs) > cfg.Retain {
		for i := len(backupBlobs) - 1; i >= cfg.Retain; i-- {
			blobName := backupBlobs[i].name
			if err := bucket.Object(blobName).Delete(ctx); err != nil {
				return err
			}
			zap.L().Sugar().Infof("Deleted old GCS backup: gs://%s/%s", cfg.Bucket, blobName)
		}
	}

	return nil
}
