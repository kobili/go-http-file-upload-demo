package db

import (
	"context"
	"database/sql"
)

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

func GetProfilePicsForUser(db *sql.DB, ctx context.Context, userId string) ([]ProfilePictureEntity, error) {
	var pictures []ProfilePictureEntity

	rows, err := db.QueryContext(
		ctx,
		`SELECT
			id,
			file_path,
			user_id
		FROM profile_pictures
		WHERE user_id = $1`,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var picture ProfilePictureEntity

		err := rows.Scan(&picture.Id, &picture.FilePath, &picture.UserId)
		if err != nil {
			return nil, err
		}
		pictures = append(pictures, picture)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pictures, nil
}

func DeleteProfilePic(db *sql.DB, ctx context.Context, id string) error {
	_, err := db.ExecContext(
		ctx,
		`DELETE FROM profile_pictures WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}
