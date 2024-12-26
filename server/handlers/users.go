package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	sdb "server/db"
	storage "server/storage_backends"
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

func CreateProfilePicHandler(db *sql.DB, storage_backend storage.StorageBackend) http.HandlerFunc {
	fn := func(w http.ResponseWriter, req *http.Request) {
		userId := chi.URLParam(req, "userId")

		// parse the multipart/form-data
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

		// persist the file
		userFileDir := getProfilePicFilePath(userId)

		timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
		newFileName := fmt.Sprintf("%s_%s", timestamp, fileHeader.Filename)

		filePath, err := storage_backend.SaveFile(file, userFileDir, newFileName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to save file: %v", err), 500)
			return
		}

		// Save a record of the photo to the database
		photo, err := sdb.CreateProfilePic(db, req.Context(), userId, filePath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to save photo to db: %v", err), 500)
			return
		}

		// generate the response
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

func RetrieveProfilePicHandler(db *sql.DB, storage_backend storage.StorageBackend) http.HandlerFunc {
	fn := func(w http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")

		// Retrieve photo information from db
		picEntity, err := sdb.GetProfilePic(db, req.Context(), id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve image info: %v", err), 500)
			return
		}

		// Retrieve the photo data from the storage backend
		buf, err := storage_backend.RetrieveFile(picEntity.FilePath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading file data: %v", err), 500)
			return
		}

		// send the data in the response
		w.Write(buf)
	}

	return http.HandlerFunc(fn)
}

func DeleteProfilePicHandler(db *sql.DB, storage_backend storage.StorageBackend) http.HandlerFunc {
	fn := func(w http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")

		// Retrieve photo info from db
		picEntity, err := sdb.GetProfilePic(db, req.Context(), id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve image info: %v", err), 500)
			return
		}

		// Delete the file from persistent storage
		err = storage_backend.DeleteFile(picEntity.FilePath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error removing file: %v", err), 500)
			return
		}

		// Remove the db record
		err = sdb.DeleteProfilePic(db, req.Context(), id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to delete db entry for photo: %v", err), 500)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}

	return http.HandlerFunc(fn)
}
