package utils

import (
	"os"
	"testing"
)

func TestCompressFile(t *testing.T) {
	// Create test file
	testFile, _ := os.Create("test.txt")
	testFile.Write([]byte("test data"))
	testFile.Close()

	// Run compression
	if err := CompressFile("test.txt", "test.gz"); err != nil {
		t.Errorf("Compression failed: %v", err)
	}

	// Verify output file exists
	if _, err := os.Stat("test.gz"); os.IsNotExist(err) {
		t.Error("Output file not found")
	}
}
