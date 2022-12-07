package main

import (
	"context"
	"flag"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/swsd2544/learny-backend-clone/internal/repository"
	"os"
	"time"
)

type config struct {
	port        int
	environment string
	db          struct {
		dsn string
	}
}

type application struct {
	config       config
	logger       zerolog.Logger
	repositories repository.Repositories
}

func main() {
	var config config

	flag.IntVar(&config.port, "port", 4000, "API server port")
	flag.StringVar(&config.environment, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&config.db.dsn, "db-dsn", os.Getenv("LEARNY_DB_DSN"), "Postgres DSN")

	flag.Parse()

	logger := zerolog.New(os.Stdout)

	logger.Info().
		Str("db-dsn", config.db.dsn).
		Msg("connecting to the db server")
	db, err := openDB(config.db.dsn)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msgf("error opening db connection")
	}

	repositories := repository.New(db)

	app := application{
		config:       config,
		logger:       logger,
		repositories: repositories,
	}

	err = app.serve()
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("error closing server")
	}
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
