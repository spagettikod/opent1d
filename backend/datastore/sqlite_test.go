package datastore

import (
	"testing"
	"time"
)

func setupStore() (Store, error) {
	store, err := NewSQLiteStore("file::memory:")
	if err != nil {
		return nil, err
	}
	if err := store.Migrate(0); err != nil {
		return nil, err
	}
	return store, nil
}

func TestSaveSetting(t *testing.T) {
	store, err := setupStore()
	if err != nil {
		t.Fatalf("failed to setup store: %v", err)
	}
	defer store.Close()
	expectedSettings := Settings{}
	expectedSettings.LibreLinkUpUsername = "foo@bar.com"
	expectedSettings.LibreLinkUpPassword = "secret"
	expectedSettings.LibreLinkUpRegion = "eu"
	if err := store.SaveSettings(expectedSettings); err != nil {
		t.Fatalf("failed to save setting: %v", err)
	}

	actualSettings, err := store.GetSettings()
	if err != nil {
		t.Fatalf("failed to get setting: %v", err)
	}

	if expectedSettings.LibreLinkUpUsername != actualSettings.LibreLinkUpUsername {
		t.Fatalf("expected LibreLinkUpUsername value %s but got %s", expectedSettings.LibreLinkUpUsername, actualSettings.LibreLinkUpUsername)
	}
	if expectedSettings.LibreLinkUpPassword != actualSettings.LibreLinkUpPassword {
		t.Fatalf("expected LibreLinkUpPassword value %s but got %s", expectedSettings.LibreLinkUpPassword, actualSettings.LibreLinkUpPassword)
	}
	if expectedSettings.LibreLinkUpRegion != actualSettings.LibreLinkUpRegion {
		t.Fatalf("expected LibreLinkUpRegion value %s but got %s", expectedSettings.LibreLinkUpRegion, actualSettings.LibreLinkUpRegion)
	}
}

func TestSaveCGM(t *testing.T) {
	store, err := setupStore()
	if err != nil {
		t.Fatalf("failed to setup store: %v", err)
	}
	defer store.Close()
	tests := []CGMEntry{
		NewCGMEntry(time.Date(2023, 06, 01, 10, 05, 35, 0, time.UTC), 3.7),
		NewCGMEntry(time.Date(2023, 06, 01, 10, 10, 35, 0, time.UTC), 4.8),
		NewCGMEntry(time.Date(2023, 06, 01, 10, 05, 35, 0, time.UTC), 10.45), // should replace the first entry
	}

	for _, test := range tests {
		if err := store.SaveCGM(test); err != nil {
			t.Fatalf("failed to save, %s: %v", test, err)
		}
	}
}
