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

func GetUsers(db *sql.DB, ctx context.Context) ([]UserEntity, error) {
	var users []UserEntity

	rows, err := db.QueryContext(
		ctx,
		`SELECT
			user_id,
			username
		FROM users`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user UserEntity
		if err := rows.Scan(
			&user.UserId, &user.Username,
		); err != nil {
			return nil, fmt.Errorf("GetUsers: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetUsers: %v", err)
	}
	return users, nil
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

type ProfilePictureEntity struct {
	Id       int64
	FilePath string
	UserId   string
}

func CreateProfilePic(db *sql.DB, ctx context.Context, userId string, filePath string) (*ProfilePictureEntity, error) {
	var entity ProfilePictureEntity

	err := db.QueryRowContext(
		ctx,
		`INSERT INTO profile_pictures
		(file_path, user_id)
		VALUES ($1, $2)
		RETURNING
			id, file_path, user_id`,
		filePath,
		userId,
	).Scan(
		&entity.Id, &entity.FilePath, &entity.UserId,
	)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func GetProfilePic(db *sql.DB, ctx context.Context, id string) (*ProfilePictureEntity, error) {
	var entity ProfilePictureEntity

	err := db.QueryRowContext(
		ctx,
		`SELECT id, file_path, user_id FROM profile_pictures WHERE id = $1`,
		id,
	).Scan(&entity.Id, &entity.FilePath, &entity.UserId)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}
