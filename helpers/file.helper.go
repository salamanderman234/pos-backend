package helpers

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/salamanderman234/pos-backend/config"
)

func verifyFileMimes(file multipart.File, alloweds []string) ([]byte, error) {
	var buffer bytes.Buffer
	_, err := io.Copy(&buffer, file)
	if err != nil && err != io.EOF {
		return []byte{}, err
	}
	content := buffer.Bytes()
	mime := http.DetectContentType(content)

	ok := false

	for _, allowed := range alloweds {
		if mime == allowed {
			ok = true
		}
	}

	if !ok {
		return content, config.ErrInvalidMimeType
	}
	return content, nil
}

func verifyFileSize(file multipart.FileHeader, maxSize int64) bool {
	return file.Size <= maxSize
}

func CheckAndSaveFile(file multipart.FileHeader, path string, allowedMimes []string, maxSize int64) (string, error) {
	if !verifyFileSize(file, maxSize) {
		return "", config.ErrTooLarge
	}
	name := GenerateRandomString(50)
	storageDir := "./storage"
	dirPath := fmt.Sprintf("%s/%s", storageDir, path)

	fullfile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer fullfile.Close()

	content, err := verifyFileMimes(fullfile, allowedMimes)
	if err != nil {
		return "", err
	}

	ext := filepath.Ext(file.Filename)
	if err := os.MkdirAll(fmt.Sprintf("%s/%s", storageDir, path), os.ModePerm); err != nil {
		return "", err
	}

	fullname := fmt.Sprintf("%s%s", name, ext)
	fullpath := filepath.Join(dirPath, fullname)
	dst, err := os.Create(fullpath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = dst.Write(content); err != nil {
		return "", err
	}

	return fullname, nil
}

func RemoveFile(path string) error {
	return os.Remove(path)
}

func GetFile(path string) ([]byte, string, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, "", err
	}
	mime := http.DetectContentType(file)
	return file, mime, nil
}
