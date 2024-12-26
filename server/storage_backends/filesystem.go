package storage_backends

import (
	"bufio"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type StorageBackend interface {
	SaveFile(file multipart.File, path string, fileName string) (string, error)
	RetrieveFile(filePath string) ([]byte, error)
}

type FileSystemStorageBackend struct {
}

func NewFileSystemStorageBackend() *FileSystemStorageBackend {
	return &FileSystemStorageBackend{}
}

func (backend *FileSystemStorageBackend) SaveFile(file multipart.File, path string, fileName string) (string, error) {

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", err
	}

	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)

	newFileName := filepath.Join(path, fmt.Sprintf("%s_%s", timestamp, fileName))
	newFile, err := os.Create(newFileName)
	if err != nil {
		return "", err
	}
	defer newFile.Close()

	buf := make([]byte, 1024)

	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return "", err
		}

		if n == 0 {
			break
		}

		_, err = newFile.Write(buf[:n])
		if err != nil {
			return "", err
		}
	}

	return newFileName, nil
}

func (backend *FileSystemStorageBackend) RetrieveFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, fileInfo.Size())
	_, err = bufio.NewReader(file).Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return buf, nil
}
