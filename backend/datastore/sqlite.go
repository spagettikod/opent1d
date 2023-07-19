package datastore

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	// KeySettings is the used in the kv-table to store settings in JSON
	KeySettings = "settings"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(uri string) (SQLiteStore, error) {
	store := SQLiteStore{}
	db, err := sql.Open("sqlite3", uri)
	if err != nil {
		return store, err
	}
	db.SetConnMaxLifetime(-1)
	store.db = db

	return store, store.db.QueryRow("SELECT 1").Err()
}

func (sls SQLiteStore) Migrate(from int) error {
	for i := from; i < len(migrations); i++ {
		for _, stmt := range migrations[i] {
			if _, err := sls.db.Exec(stmt); err != nil {
				return err
			}
		}
	}
	return nil
}

func (sls SQLiteStore) Close() error {
	return sls.db.Close()
}

func (sls SQLiteStore) GetSettings() (Settings, error) {
	row := sls.db.QueryRow("SELECT value FROM kv WHERE key = ?", KeySettings)
	var jsn string
	if err := row.Scan(&jsn); err != nil {
		return Settings{}, fmt.Errorf("error while loading settings from SQLite: %w", err)
	}
	return SettingsFromJson(jsn)
}

func (sls SQLiteStore) SaveSettings(settings Settings) error {
	json, err := settings.ToJson()
	if err != nil {
		return err
	}
	tx, err := sls.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO kv (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = ?", KeySettings, json, json)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}

func (sls SQLiteStore) SaveCGM(cgms ...CGMEntry) error {
	tx, err := sls.db.Begin()
	if err != nil {
		return err
	}

	for _, cgm := range cgms {
		_, err = tx.Exec("INSERT INTO cgm VALUES (?, ?) ON CONFLICT DO NOTHING", cgm.Timestamp.Unix(), cgm.Mmoll)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}
	return tx.Commit()
}

func (sls SQLiteStore) LoadCGMInterval(from, to time.Time) ([]CGMEntry, error) {
	cgms := []CGMEntry{}
	// row := sls.db.Query("SELECT ts,mmoll FROM cgm WHERE ts  = ?", setting.key)
	return cgms, nil
}

var (
	migrations = [][]string{
		{
			`CREATE TABLE IF NOT EXISTS cgm (
	ts INTEGER PRIMARY KEY,
	mmoll REAL NOT NULL
)`,
			`CREATE TABLE IF NOT EXISTS kv (
	key TEXT PRIMARY KEY,
	value TEXT
)`,
		},
	}
)
