package storage_backends

import (
	"bufio"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type FileSystemStorageBackend struct {
}

func NewFileSystemStorageBackend() *FileSystemStorageBackend {
	return &FileSystemStorageBackend{}
}

/*
Save the given multipart.File to the local file system under the directory given by path and with the name given by fileName.
Returns the path to the new file.
*/
func (backend *FileSystemStorageBackend) SaveFile(file multipart.File, path string, fileName string) (string, error) {
	// create the directory if it doesn't exist
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", err
	}

	// Create a new file
	newFileName := filepath.Join(path, fileName)
	newFile, err := os.Create(newFileName)
	if err != nil {
		return "", err
	}
	defer newFile.Close()

	// Copy the data from the multipart.File to the new File
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

/*
Retrieve the binary data of the file located at filePath.
*/
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

/*
Remove the file given by filePath.
*/
func (backend *FileSystemStorageBackend) DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}
