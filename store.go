package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	_ "modernc.org/sqlite"
)

type PlistStore struct {
	db *sql.DB
}

func NewPlistStore(dbPath string) (*PlistStore, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set SQLite pragmas for performance and safety
	pragmas := []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA synchronous=NORMAL",
		"PRAGMA busy_timeout=5000",
		"PRAGMA foreign_keys=ON",
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			return nil, fmt.Errorf("failed to set pragma: %w", err)
		}
	}

	// Initialize schema (idempotent)
	if err := initSchema(db); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return &PlistStore{db: db}, nil
}

func initSchema(db *sql.DB) error {
	schema := `
		CREATE TABLE IF NOT EXISTS launchd_plists (
			id TEXT PRIMARY KEY,
			data TEXT NOT NULL,
			created_at TEXT NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_created_at
		ON launchd_plists(created_at);
	`

	_, err := db.Exec(schema)
	return err
}

func (s *PlistStore) Save(plist LaunchdPlist) (string, error) {
	plist.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	id, _ := gonanoid.New(10)

	data, err := json.Marshal(plist)
	if err != nil {
		return "", fmt.Errorf("failed to marshal plist: %w", err)
	}

	_, err = s.db.Exec(
		"INSERT INTO launchd_plists (id, data, created_at) VALUES (?, ?, ?)",
		id, string(data), plist.CreatedAt)

	return id, err
}

func (s *PlistStore) Load(id string) (LaunchdPlist, bool, error) {
	var data string
	var plist LaunchdPlist

	err := s.db.QueryRow(
		"SELECT data FROM launchd_plists WHERE id = ?", id).Scan(&data)

	if err == sql.ErrNoRows {
		return plist, false, nil
	}
	if err != nil {
		return plist, false, err
	}

	if err := json.Unmarshal([]byte(data), &plist); err != nil {
		return plist, false, fmt.Errorf("failed to unmarshal plist: %w", err)
	}

	plist.ID = id
	return plist, true, nil
}

func (s *PlistStore) Close() error {
	return s.db.Close()
}
