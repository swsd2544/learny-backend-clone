package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/swsd2544/learny-backend-clone/internal/repository"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"go.uber.org/zap"
	"os"
	"syscall"
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
	logger       *zap.Logger
	validate     *validator.Validate
	repositories repository.Repositories
}

func main() {
	var config config

	flag.IntVar(&config.port, "port", 4000, "API server port")
	flag.StringVar(&config.environment, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&config.db.dsn, "db-dsn", os.Getenv("LEARNY_DB_DSN"), "Postgres DSN")

	flag.Parse()

	logger, err := newLogger(config.environment)
	if err != nil {
		panic(err)
	}

	db, err := openDB(config.db.dsn)
	if err != nil {
		logger.Fatal("error open db connections", zap.Error(err))
	}

	repositories := repository.New(db)
	validate := validator.New()

	err = validate.RegisterValidation(
		"password",
		func(fl validator.FieldLevel) bool {
			entropy := passwordvalidator.GetEntropy("pa55word")
			err := passwordvalidator.Validate("some password", entropy)
			if err != nil {
				return false
			}
			return true
		},
	)
	if err != nil {
		logger.Fatal("error register password validation", zap.Error(err))
	}

	app := application{
		config:       config,
		logger:       logger,
		validate:     validate,
		repositories: repositories,
	}

	err = app.serve()
	if err != nil {
		logger.Fatal("error closing server", zap.Error(err))
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

func newLogger(env string) (logger *zap.Logger, err error) {
	var config zap.Config
	switch env {
	case "development", "staging":
		config = zap.NewDevelopmentConfig()
	case "production":
		config = zap.NewProductionConfig()
	default:
		return nil, fmt.Errorf("misformat environment (development|staging|production): %v", env)
	}

	logger, err = config.Build()
	if err != nil {
		return nil, err
	}

	err = logger.Sync()
	if err != nil && !errors.Is(err, syscall.ENOTTY) {
		return nil, err
	}

	return logger, nil
}
