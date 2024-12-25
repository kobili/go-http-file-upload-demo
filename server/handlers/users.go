package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func CreateUserHandler(w http.ResponseWriter, req *http.Request) {
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
