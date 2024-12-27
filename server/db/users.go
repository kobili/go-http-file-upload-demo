package db

import (
	"context"
	"database/sql"
	"fmt"
)

type UserEntity struct {
	UserId   string
	Username string
}

func GetUserById(db *sql.DB, ctx context.Context, userId string) (*UserEntity, error) {
	var user UserEntity
	err := db.QueryRowContext(
		ctx,
		`SELECT
			user_id,
			username
		FROM users
		WHERE user_id = $1`,
		userId,
	).Scan(
		&user.UserId,
		&user.Username,
	)
	if err != nil {
		return nil, fmt.Errorf("GetUserById(%s): %w", userId, err)
	}

	return &user, nil
}

type CreateUserPayload struct {
	Username string `json:"username"`
}

func CreateUser(db *sql.DB, ctx context.Context, data CreateUserPayload) (*UserEntity, error) {
	var user UserEntity
	err := db.QueryRowContext(
		ctx,
		`INSERT INTO users (username)
		VALUES ($1)
		RETURNING
			user_id,
			username`,
		data.Username,
	).Scan(
		&user.UserId,
		&user.Username,
	)
	if err != nil {
		return nil, fmt.Errorf("CreateUser - Could not create user: %w", err)
	}

	return &user, nil
}

func DeleteUser(db *sql.DB, ctx context.Context, userId string) error {
	_, err := db.ExecContext(
		ctx,
		`DELETE FROM users WHERE user_id = $1`,
		userId,
	)
	if err != nil {
		return fmt.Errorf("DeleteUser - Could not delete user: %w", err)
	}

	return nil
}
