package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// GetFileExtension returns the extension of a file
func GetFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

// GetMD5Hash calculates the MD5 hash of a file
func GetMD5Hash(file *os.File) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	// Reset file position
	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// GetFileSize returns the size of a file
func GetFileSize(file *os.File) (int64, error) {
	info, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// CopyFileContent copies content from a multipart file to a local file
func CopyFileContent(src *multipart.FileHeader, dstPath string) error {
	// Open source file
	srcFile, err := src.Open()
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create destination file
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy content
	_, err = io.Copy(dstFile, srcFile)
	return err
}

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// GetMimeType tries to determine the MIME type of a file
func GetMimeType(file *os.File) (string, error) {
	// Read the first 512 bytes to determine the content type
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	// Reset file position
	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	// Use http.DetectContentType to detect the content type
	// This is a simplified version, a real implementation would use http.DetectContentType
	contentType := detectContentType(buffer)

	return contentType, nil
}

// detectContentType is a simplified version for demo purposes
// In a real app, you would use http.DetectContentType
func detectContentType(buffer []byte) string {
	// Very simplified detection for demo purposes
	if len(buffer) == 0 {
		return "application/octet-stream"
	}

	// Check for common image formats
	if buffer[0] == 0xFF && buffer[1] == 0xD8 {
		return "image/jpeg"
	}
	if buffer[0] == 0x89 && buffer[1] == 0x50 && buffer[2] == 0x4E && buffer[3] == 0x47 {
		return "image/png"
	}
	if buffer[0] == 0x47 && buffer[1] == 0x49 && buffer[2] == 0x46 {
		return "image/gif"
	}

	// Default
	return "application/octet-stream"
}
