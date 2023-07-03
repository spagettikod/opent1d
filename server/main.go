package main

import (
	"fmt"
	"log"
	"opent1d/datastore"
	"opent1d/librelinkup"
	"os"
	"path/filepath"
)

const (
	DB_PATH_ENV = "OPENT1D_DBPATH"
	DEBUG_ENV   = "OPENT1D_DEBUG"
	// DB_PATH_DIR path where the database is stored, this is added to the DB_PATH
	DB_PATH_DIR = "OpenT1D"
	// DB_FILENAME name of the database file
	DB_FILENAME = "opent1d.db"
)

func GetDBPath() string {
	path := os.Getenv(DB_PATH_ENV)
	if path != "" {
		return path
	}
	path, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("can not determine database location path, exiting: %s", err)
	}
	path = filepath.Clean(filepath.Join(path, DB_PATH_DIR))
	if err := os.MkdirAll(path, 0750); err != nil {
		log.Fatalf("can not create database path, exiting: %s", err)
	}
	return fmt.Sprintf("file:%s", filepath.Join(path, DB_FILENAME))
}

func EnvOrDie(env string) string {
	if val, ok := os.LookupEnv(env); ok {
		if val != "" {
			return val
		}
	}
	log.Fatalf("could not find environment variable %s, exiting", env)
	return ""
}

func main() {
	dbPath := GetDBPath()
	store, err := datastore.NewSQLiteStore(dbPath)
	if err != nil {
		log.Fatalf("could not open OpenT1D database at %s, exiting: %s", dbPath, err)
	}
	if err := store.Migrate(0); err != nil {
		log.Fatalf("could not migrate database, exiting: %s", err)
	}

	email := EnvOrDie("GOLIBRE_EMAIL")
	password := EnvOrDie("GOLIBRE_PASSWORD")
	_, err = librelinkup.Login(email, password, librelinkup.EndpointUS)
	if err != nil {
		log.Fatalf("could not connect, exiting: %s", err)
	}

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
