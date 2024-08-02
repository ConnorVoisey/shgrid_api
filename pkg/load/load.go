package load

import (
	"database/sql"
	"github.com/connorvoisey/shgrid_api/pkg/db"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"os"
)

type State struct {
	Log *zerolog.Logger
	Db  *sql.DB
}

func InitLogger() (*zerolog.Logger, error) {
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

func Init() (*State, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	log, err := InitLogger()
	if err != nil {
		return nil, err
	}

	db, err := db.GetDb(log)
	if err != nil {
		return nil, err
	}

	state := State{
		log,
		db,
	}
	return &state, nil
}
