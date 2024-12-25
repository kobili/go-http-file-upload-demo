package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"server/utils"
)

func ConnectToDB() *sql.DB {
	host := os.Getenv("DB_HOST")
	port := utils.GetNumericEnv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connectionString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Fatalf("sql.Open: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("sql.Ping: %v", err)
	}

	return db
}
