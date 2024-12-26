package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	sdb "server/db"
)

type UserDto struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
}

func CreateUserHandler(db *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, req *http.Request) {
		var reqBody sdb.CreateUserPayload
		if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
			http.Error(w, fmt.Sprintf("Error reading request body: %v", err), 500)
			return
		}

		user, err := sdb.CreateUser(db, req.Context(), reqBody)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating user: %v", err), 500)
			return
		}

		response, err := json.Marshal(UserDto{
			UserId:   user.UserId,
			Username: user.Username,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Error serializing response: %v", err), 500)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(response)
	}

	return http.HandlerFunc(fn)
}

type ProfilePicDto struct {
	Id       int64  `json:"id"`
	FilePath string `json:"filePath"`
	UserId   string `json:"userId"`
}

func CreateProfilePicHandler(db *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, req *http.Request) {
		userId := chi.URLParam(req, "userId")

		err := req.ParseMultipartForm(200 << 20)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed parsing request body: %v", err), 500)
			return
		}

		file, fileHeader, err := req.FormFile("profilePic")
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve key `profilePic` from form data: %v", err), 500)
			return
		}

		fileName := fileHeader.Filename
		log.Println(fileName)

		userFileDir := getProfilePicFilePath(userId)
		log.Println(userFileDir)

		filePath, err := saveFile(file, fileHeader, userFileDir)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to save file: %v", err), 500)
			return
		}

		photo, err := sdb.CreateProfilePic(db, req.Context(), userId, filePath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to save photo to db: %v", err), 500)
			return
		}

		resBody, err := json.Marshal(ProfilePicDto{
			Id:       photo.Id,
			FilePath: photo.FilePath,
			UserId:   photo.UserId,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Error serializing response: %v", err), 500)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(resBody)
	}

	return http.HandlerFunc(fn)
}

func getProfilePicFilePath(userId string) string {
	uploadsDir, isSet := os.LookupEnv("UPLOADS_DIR")
	if !isSet {
		uploadsDir = "uploads/"
	}

	return filepath.Join(".", uploadsDir, userId, "profile_pics/")
}

func saveFile(file multipart.File, fileHeader *multipart.FileHeader, path string) (string, error) {

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", err
	}

	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)

	newFileName := filepath.Join(path, fmt.Sprintf("%s_%s", timestamp, fileHeader.Filename))
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
