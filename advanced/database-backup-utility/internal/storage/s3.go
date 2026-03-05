package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"go.uber.org/zap"

	"github.com/jaygaha/roadmap-go-projects/advanced/database-backup-utility/internal/config"
)

// UploadToS3 uploads a file to AWS S3
func UploadToS3(filePath string, bucket, region string) error {
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(), awsConfig.WithRegion(region))
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	svc := s3.NewFromConfig(cfg)

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	key := fmt.Sprintf("%s_backup_%s.sql", filepath.Base(filePath), time.Now().Format("20060102_150405"))

	_, err = svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("failed to upload to S3: %w", err)
	}

	zap.L().Sugar().Infof("Uploaded %s to S3 bucket %s", filePath, bucket)
	return nil
}

// EnforceRetentionS3 keeps only the newest Retain .sql backup files in S3
func EnforceRetentionS3(cfg config.StorageConfig, dryRun bool) error {
	awsCfg, err := awsConfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	svc := s3.NewFromConfig(awsCfg)

	paginator := s3.NewListObjectsV2Paginator(svc, &s3.ListObjectsV2Input{
		Bucket: aws.String(cfg.Bucket),
	})

	type backupFile struct {
		Key          string
		LastModified time.Time
	}

	var backups []backupFile
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return fmt.Errorf("failed to list S3 objects: %w", err)
		}

		for _, obj := range page.Contents {
			if obj.Key != nil && strings.HasSuffix(*obj.Key, ".sql") && obj.LastModified != nil {
				backups = append(backups, backupFile{
					Key:          *obj.Key,
					LastModified: *obj.LastModified,
				})
			}
		}
	}

	if dryRun {
		zap.L().Sugar().Infof("Would keep %d newest backups in S3 bucket %s", cfg.Retain, cfg.Bucket)
		return nil
	}

	sort.Slice(backups, func(i, j int) bool {
		return backups[i].LastModified.After(backups[j].LastModified)
	})

	if len(backups) > cfg.Retain {
		objectsToDelete := make([]types.ObjectIdentifier, 0, len(backups)-cfg.Retain)
		for i := cfg.Retain; i < len(backups); i++ {
			objectsToDelete = append(objectsToDelete, types.ObjectIdentifier{
				Key: aws.String(backups[i].Key),
			})
		}

		_, err = svc.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
			Bucket: aws.String(cfg.Bucket),
			Delete: &types.Delete{
				Objects: objectsToDelete,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to delete old backups: %w", err)
		}

		for i := cfg.Retain; i < len(backups); i++ {
			zap.L().Sugar().Infof("Deleted old S3 backup: s3://%s/%s", cfg.Bucket, backups[i].Key)
		}
	}

	return nil
}
