package main

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testStore(t *testing.T) *PlistStore {
	t.Helper()
	store, err := NewPlistStore(filepath.Join(t.TempDir(), "test.db"))
	require.NoError(t, err)
	t.Cleanup(func() { store.Close() })
	return store
}

func TestPlistStore_SaveAndLoad(t *testing.T) {
	store := testStore(t)

	id, err := store.Save(LaunchdPlist{Name: "backup", Command: "rsync -a ~/src /backup"})
	require.NoError(t, err)
	assert.Len(t, id, 10)

	plist, found, err := store.Load(id)
	require.NoError(t, err)
	require.True(t, found)
	assert.Equal(t, "backup", plist.Name)
	assert.Equal(t, "rsync -a ~/src /backup", plist.Command)
}

func TestPlistStore_LoadMissingIsNotAnError(t *testing.T) {
	// handlers distinguish "not found" (404) from a real failure (500), so a
	// missing row has to come back as ok=false with a nil error
	store := testStore(t)

	plist, found, err := store.Load("nonexistent")
	assert.NoError(t, err)
	assert.False(t, found)
	assert.Empty(t, plist.Name)
}

func TestPlistStore_RoundTripsEveryField(t *testing.T) {
	// the struct is stored as one JSON blob, so a field missing a json tag would
	// silently vanish on save rather than failing loudly
	store := testStore(t)

	in := LaunchdPlist{
		Name:              "My Backup Job",
		Command:           "rsync -a ~/src /backup",
		StartInterval:     "300",
		Minute:            "0,30",
		Hour:              "*/4",
		DayOfMonth:        "1-5",
		Month:             "*",
		Weekday:           "1",
		RunAtLoad:         "on",
		RestartOnCrash:    "on",
		StartOnMount:      "on",
		QueueDirectories:  "/tmp/a,/tmp/b",
		Environment:       "PATH=/usr/bin\r\nFOO=bar",
		User:              "nathan",
		Group:             "staff",
		WorkingDirectory:  "/tmp",
		RootDirectory:     "/",
		StandardOutPath:   "/tmp/out.log",
		StandardErrorPath: "/tmp/err.log",
	}

	id, err := store.Save(in)
	require.NoError(t, err)
	out, found, err := store.Load(id)
	require.NoError(t, err)
	require.True(t, found)

	// ID and CreatedAt are assigned by the store, so compare the rest as a whole
	in.ID, in.CreatedAt = out.ID, out.CreatedAt
	assert.Equal(t, in, out)

	// environment is split on CRLF, which JSON must not normalize away
	assert.Equal(t, map[string]string{"PATH": "/usr/bin", "FOO": "bar"}, out.EnvironmentMap())
	assert.Equal(t, in.PlistXML(""), out.PlistXML(""))
}

func TestPlistStore_LoadSetsID(t *testing.T) {
	// ID isn't part of the stored JSON; Load backfills it for template URLs
	store := testStore(t)

	id, err := store.Save(LaunchdPlist{Name: "x", Command: "true"})
	require.NoError(t, err)

	plist, _, err := store.Load(id)
	require.NoError(t, err)
	assert.Equal(t, id, plist.ID)
}

func TestPlistStore_SaveSetsCreatedAt(t *testing.T) {
	store := testStore(t)

	original := LaunchdPlist{Name: "x", Command: "true"}
	id, err := store.Save(original)
	require.NoError(t, err)

	plist, _, err := store.Load(id)
	require.NoError(t, err)
	assert.NotEmpty(t, plist.CreatedAt)
	assert.Empty(t, original.CreatedAt, "Save takes a copy and must not mutate the caller's plist")
}

func TestPlistStore_GeneratesDistinctIDs(t *testing.T) {
	store := testStore(t)

	seen := map[string]bool{}
	for range 50 {
		id, err := store.Save(LaunchdPlist{Name: "x", Command: "true"})
		require.NoError(t, err)
		assert.False(t, seen[id], "duplicate id %q", id)
		seen[id] = true
	}
}

func TestPlistStore_ReopenPersistsRows(t *testing.T) {
	// initSchema runs on every startup, so it has to be idempotent and must not
	// clobber an existing database
	path := filepath.Join(t.TempDir(), "reopen.db")

	first, err := NewPlistStore(path)
	require.NoError(t, err)
	id, err := first.Save(LaunchdPlist{Name: "persisted", Command: "true"})
	require.NoError(t, err)
	require.NoError(t, first.Close())

	second, err := NewPlistStore(path)
	require.NoError(t, err)
	defer second.Close()

	plist, found, err := second.Load(id)
	require.NoError(t, err)
	require.True(t, found)
	assert.Equal(t, "persisted", plist.Name)
}

func TestPlistStore_AppliesPragmas(t *testing.T) {
	store := testStore(t)

	var journal string
	require.NoError(t, store.db.QueryRow("PRAGMA journal_mode").Scan(&journal))
	assert.Equal(t, "wal", journal, "litestream replication depends on WAL mode")

	var foreignKeys int
	require.NoError(t, store.db.QueryRow("PRAGMA foreign_keys").Scan(&foreignKeys))
	assert.Equal(t, 1, foreignKeys)
}
