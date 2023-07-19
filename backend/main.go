package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/spagettikod/opent1d/datastore"
	"github.com/spagettikod/opent1d/envctx"
	"github.com/spagettikod/opent1d/event"
	"github.com/spagettikod/opent1d/graph"
	"github.com/spagettikod/opent1d/handle"
	"github.com/spagettikod/opent1d/librelinkup"
)

const (
	DB_PATH_ENV = "OPENT1D_DBPATH"
	// DB_PATH_DIR path where the database is stored, this is added to the DB_PATH
	DB_PATH_DIR = "OpenT1D"
	// DB_FILENAME name of the database file
	DB_FILENAME = "opent1d.sqlite"

	PORT = "8080"

	LOG_KEY_DB = "database"
)

func GetDBPath() string {
	path := os.Getenv(DB_PATH_ENV)
	if path != "" {
		return path
	}
	path, err := os.UserConfigDir()
	if err != nil {
		log.Fatal().Err(err).Str(LOG_KEY_DB, path).Msg("can not determine database location path, exiting")
	}
	path = filepath.Clean(filepath.Join(path, DB_PATH_DIR))
	if err := os.MkdirAll(path, 0750); err != nil {
		log.Fatal().Err(err).Str(LOG_KEY_DB, path).Msg("can not create database path, exiting")
	}
	return fmt.Sprintf("file:%s", filepath.Join(path, DB_FILENAME))
}

func EnvOrDie(env string) string {
	if val, ok := os.LookupEnv(env); ok {
		if val != "" {
			return val
		}
	}
	log.Fatal().Msgf("could not find environment variable %s, exiting", env)
	return ""
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).Level(envctx.EnvToLogLevel())

	dbPath := GetDBPath()
	store, err := datastore.NewSQLiteStore(dbPath)
	if err != nil {
		log.Fatal().Err(err).Str(LOG_KEY_DB, dbPath).Msg("could not open OpenT1D database, exiting")
	}
	if err := store.Migrate(0); err != nil {
		log.Fatal().Err(err).Str(LOG_KEY_DB, dbPath).Msg("could not migrate database, exiting")
	}

	ctx := envctx.NewContext(store, log.Logger)

	// this event can be async
	go event.OnStartup(ctx)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Context: ctx}}))
	srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		err := graphql.DefaultErrorPresenter(ctx, e)

		if errors.Is(e, librelinkup.ErrLoginFailed) {
			err.Message = "Login failed, please verify username and password"
		}

		return err
	})

	http.Handle("/query", handle.Middleware(srv))
	http.Handle("/", http.FileServer(http.Dir("/www")))

	log.Info().Msgf("http server is listening on port %s", PORT)
	log.Fatal().Err(http.ListenAndServe(fmt.Sprintf(":%s", PORT), nil))
}
