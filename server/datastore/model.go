package datastore

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var (
	ErrNotFound  = errors.New("item was not found")
	ErrNonUnique = errors.New("field must be unique")
)

type Store interface {
	Migrate(from int) error
	Close() error
	GetSettings() (Settings, error)
	SaveSettings(settings Settings) error
	SaveCGM(cgms ...CGMEntry) error
	LoadCGMInterval(from, to time.Time) ([]CGMEntry, error)
}

type Settings struct {
	LibreLinkUpUsername string `json:"libreLinkupUsername"`
	LibreLinkUpPassword string `json:"libreLinkupPassword"`
	LibreLinkUpEndpoint string `json:"libreLinkupEndpoint"`
}

func SettingsFromJson(jsn string) (Settings, error) {
	settings := Settings{}
	err := json.Unmarshal([]byte(jsn), &settings)
	return settings, err
}

func (s Settings) ToJson() (string, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return "", fmt.Errorf("error while marshaling Settings to JSON: %w", err)
	}
	return string(b), nil
}

type Mmoll float32

type CGMEntry struct {
	Timestamp time.Time
	Mmoll     Mmoll
}

func NewCGMEntry(timestamp time.Time, mmoll Mmoll) CGMEntry {
	return CGMEntry{Timestamp: timestamp, Mmoll: mmoll}
}

func (cgme CGMEntry) String() string {
	ts := cgme.Timestamp.Local().Format(time.RFC3339)
	return fmt.Sprintf("%v@%s", cgme.Mmoll, ts)
}
