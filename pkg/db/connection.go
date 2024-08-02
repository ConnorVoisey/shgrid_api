package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"

	"github.com/joho/godotenv"
)

const (
	host   = "localhost"
	port   = 5432
	dbName = "todo"
)

func GetDb(log *zerolog.Logger) (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Err(err)
		panic("Failed to parse env file")
	}

	// Connect to database
	var connectString = fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		dbName,
	)

	db, err := sql.Open("postgres", connectString)
	if err != nil {
		log.Err(err)
		return nil, err
	}
	return db, nil
}
