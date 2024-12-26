package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
