package utils

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

func CompressFile(inputPath, outputPath string) error {
	return CompressFileWithLevel(inputPath, outputPath, gzip.DefaultCompression)
}

func CompressFileWithLevel(inputPath, outputPath string, level int) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %v", err)
	}
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	if level < gzip.DefaultCompression || level > gzip.BestCompression {
		level = gzip.DefaultCompression
	}

	writer, err := gzip.NewWriterLevel(outFile, level)
	if err != nil {
		return fmt.Errorf("failed to create gzip writer: %v", err)
	}
	defer writer.Close()

	if _, err = io.Copy(writer, inFile); err != nil {
		return fmt.Errorf("compression failed: %v", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush compressed data: %v", err)
	}

	inputInfo, _ := os.Stat(inputPath)
	outputInfo, _ := os.Stat(outputPath)

	if inputInfo != nil && outputInfo != nil {
		compressionRatio := float64(outputInfo.Size()) / float64(inputInfo.Size()) * 100
		levelName := getCompressionLevelName(level)
		zap.L().Sugar().Infof("Compressed %s → %s (%.1f%% of original size, %s)",
			filepath.Base(inputPath), filepath.Base(outputPath), compressionRatio, levelName)
	}

	return nil
}

func getCompressionLevelName(level int) string {
	switch level {
	case gzip.NoCompression:
		return "no compression"
	case gzip.BestSpeed:
		return "best speed"
	case gzip.BestCompression:
		return "best compression"
	case gzip.DefaultCompression:
		return "default"
	default:
		return fmt.Sprintf("level %d", level)
	}
}

func DecompressFile(inputPath, outputPath string) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open compressed file: %v", err)
	}
	defer inFile.Close()

	reader, err := gzip.NewReader(inFile)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer reader.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	if _, err = io.Copy(outFile, reader); err != nil {
		return fmt.Errorf("decompression failed: %v", err)
	}

	zap.L().Sugar().Infof("Decompressed %s → %s",
		filepath.Base(inputPath), filepath.Base(outputPath))

	return nil
}

func IsCompressedFile(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	magic := make([]byte, 2)
	if _, err := file.Read(magic); err != nil {
		return false, err
	}

	return magic[0] == 0x1f && magic[1] == 0x8b, nil
}

func GetCompressedExtension(originalPath string) string {
	ext := filepath.Ext(originalPath)
	base := strings.TrimSuffix(originalPath, ext)
	return base + ".gz"
}
