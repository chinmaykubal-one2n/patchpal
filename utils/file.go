package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// SaveUploadedFile saves the uploaded file to /tmp and returns the file path
func SaveUploadedFile(file *multipart.FileHeader) (string, error) {
	const tmpDir = "/home/one2n/Desktop/patchpal_files" // trivy is failing in dir /tmp

	dst := filepath.Join(tmpDir, fmt.Sprintf("patchpal-%d-%s", time.Now().UnixNano(), filepath.Base(file.Filename)))
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return "", err
	}
	if err := saveFile(file, dst); err != nil {
		return "", err
	}
	return dst, nil
}

func saveFile(file *multipart.FileHeader, dst string) error {
	return os.WriteFile(dst, readUploadedFile(file), 0644)
}

func readUploadedFile(file *multipart.FileHeader) []byte {
	f, _ := file.Open()
	defer f.Close()

	data := make([]byte, file.Size)
	f.Read(data)
	return data
}
