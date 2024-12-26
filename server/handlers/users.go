package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

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

func CreateProfilePicHandler(w http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(500 << 20)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed parsing request body: %v", err), 500)
		return
	}

	username := req.FormValue("username")
	saveUser(username)

	file, fileHeader, err := req.FormFile("profilePic")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve key `profilePic` from form data: %v", err), 500)
		return
	}
	buf := make([]byte, fileHeader.Size)

	_, err = file.Read(buf)
	if err != nil && err != io.EOF {
		http.Error(w, fmt.Sprintf("Failed to read file data: %v", err), 500)
		return
	}

	saveFile(&file)

	w.Write(buf)
}

func saveUser(username string) {
	// Placeholder function
}

func saveFile(file *multipart.File) {
	// TODO: save to disk
}
