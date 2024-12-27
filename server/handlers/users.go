package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	sdb "server/db"

	"github.com/go-chi/chi/v5"
)

type UserDto struct {
	UserId      string   `json:"userId"`
	Username    string   `json:"username"`
	ProfilePics []string `json:"profilePics"`
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

func RetrieveUserHandler(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := chi.URLParam(r, "userId")

		user, err := sdb.GetUserById(db, r.Context(), userId)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve user: %v", err), 500)
			return
		}

		photos, err := sdb.GetProfilePicsForUser(db, r.Context(), userId)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve photos for user %s: %v", userId, err), 500)
			return
		}

		var photoUrls []string
		for i := 0; i < len(photos); i++ {
			photo := photos[i]
			photoUrls = append(photoUrls, fmt.Sprintf("/api/users/%s/profile_pic/%d", userId, photo.Id))
		}

		res, err := json.Marshal(UserDto{
			UserId:      user.UserId,
			Username:    user.Username,
			ProfilePics: photoUrls,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Error serializing response: %v", err), 500)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(res)
	})
}
