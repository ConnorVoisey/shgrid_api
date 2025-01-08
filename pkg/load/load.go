package load

import (
	"database/sql"
	DB "github.com/connorvoisey/shgrid_api/pkg/db"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"os"
)

func InitLogger() (*zerolog.Logger, error) {
    zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}

	file, err := os.OpenFile(
		"logs/main.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		return nil, err
	}
	fileLogger := zerolog.SyncWriter(file)
	multi := zerolog.MultiLevelWriter(consoleWriter, fileLogger)

	log := zerolog.New(multi).With().Timestamp().Logger()
	return &log, nil
}

func Init() (*zerolog.Logger, *sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, nil, err
	}

	log, err := InitLogger()
	if err != nil {
		return nil, nil, err
	}

	db, err := DB.GetDb(log)
	if err != nil {
		log.Err(err).Msg("Failed to get db")
		return log, nil, err
	}

	err = DB.Migrate(log)
	if err != nil {
		log.Err(err).Msg("Failed to migrate")
		return log, nil, err
	}
	return log, db, nil
}
