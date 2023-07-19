package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/spagettikod/opent1d/datastore"
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
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	dbPath := GetDBPath()
	store, err := datastore.NewSQLiteStore(dbPath)
	if err != nil {
		log.Fatal().Err(err).Str(LOG_KEY_DB, dbPath).Msg("could not open OpenT1D database, exiting")
	}
	if err := store.Migrate(0); err != nil {
		log.Fatal().Err(err).Str(LOG_KEY_DB, dbPath).Msg("could not migrate database, exiting")
	}

	if err := librelinkup.StartScraper(store, log.Logger); err != nil {
		log.Err(err).Msg("failed to start LibreLinkUp scraper")
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Store: store}}))

	// http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", handle.Middleware(srv))
	http.Handle("/", http.FileServer(http.Dir("/www")))

	// http.Handle("/query", handle.Middleware(handle.GraphqlHandler(store)))
	log.Fatal().Err(http.ListenAndServe(":8080", nil))

	// email := EnvOrDie("GOLIBRE_EMAIL")
	// password := EnvOrDie("GOLIBRE_PASSWORD")
	// _, err = librelinkup.Login(email, password, librelinkup.EndpointUS)
	// if err != nil {
	// 	log.Fatalf("could not connect, exiting: %s", err)
	// }

	// connections, err := ticket.Connections()
	// if err != nil {
	// 	log.Fatalf("could not connect, exiting: %s", err)
	// }
	// if len(connections) != 1 {
	// 	log.Fatalf("expected %v connections but found %v", 1, len(connections))
	// }
	// _, _, err = ticket.Graph(connections[0].PatienID)
	// if err != nil {
	// 	log.Fatalf("erro fetching graph values, exiting: %s", err)
	// }
	// gm := connections[0].GlucoseMeasurement
	// ts, err := librelinkup.ToTime(gm.Timestamp)
	// if err != nil {
	// 	log.Fatalf("could not parse timestamp, exiting: %s", err)
	// }
	// fmt.Println(ts, gm.Value)
}
