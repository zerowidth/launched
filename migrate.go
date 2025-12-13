// +build ignore

// migrate.go - One-time migration script from Redis to SQLite
// Usage: go run migrate.go
//
// Required environment variables:
//   REDIS_URL - Full Redis connection URL (e.g., redis://user:password@host:6379/0)
//   SQLITE_PATH - Output SQLite database path (e.g., ./migrated.db)

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"iter"
	"log"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"
	_ "modernc.org/sqlite"
)

type LaunchdPlist struct {
	Name              string `json:"name"`
	Command           string `json:"command"`
	StartInterval     string `json:"start_interval,omitempty"`
	Minute            string `json:"minute,omitempty"`
	Hour              string `json:"hour,omitempty"`
	DayOfMonth        string `json:"day_of_month,omitempty"`
	Month             string `json:"month,omitempty"`
	Weekday           string `json:"weekday,omitempty"`
	RunAtLoad         string `json:"run_at_load,omitempty"`
	RestartOnCrash    string `json:"restart_on_crash,omitempty"`
	StartOnMount      string `json:"start_on_mount,omitempty"`
	QueueDirectories  string `json:"queue_directories,omitempty"`
	Environment       string `json:"environment,omitempty"`
	User              string `json:"user,omitempty"`
	Group             string `json:"group,omitempty"`
	WorkingDirectory  string `json:"working_directory,omitempty"`
	RootDirectory     string `json:"root_directory,omitempty"`
	StandardOutPath   string `json:"standard_out_path,omitempty"`
	StandardErrorPath string `json:"standard_error_path,omitempty"`
	CreatedAt         string `json:"created_at"`
}

func main() {
	// Get configuration from environment
	redisURL := os.Getenv("REDIS_URL")
	sqlitePath := os.Getenv("SQLITE_PATH")

	if redisURL == "" || sqlitePath == "" {
		log.Fatal("REDIS_URL and SQLITE_PATH environment variables are required")
	}

	fmt.Printf("Migration configuration:\n")
	fmt.Printf("  Redis: %s\n", redisURL)
	fmt.Printf("  SQLite: %s\n", sqlitePath)
	fmt.Println()

	// Parse Redis URL
	redisOpts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}
	redisOpts.Network = "tcp4" // Force IPv4

	// Connect to Redis
	fmt.Println("Connecting to Redis...")
	redisClient := redis.NewClient(redisOpts)
	defer redisClient.Close()

	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("✓ Connected to Redis")

	// Connect to SQLite
	fmt.Println("Connecting to SQLite...")
	db, err := sql.Open("sqlite", sqlitePath)
	if err != nil {
		log.Fatalf("Failed to open SQLite database: %v", err)
	}
	defer db.Close()

	// Set pragmas
	db.Exec("PRAGMA journal_mode=WAL")
	db.Exec("PRAGMA synchronous=NORMAL")

	// Initialize schema
	if err := initSchema(db); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}
	fmt.Println("✓ SQLite database initialized")

	// Migrate plists using iterator with batched Redis reads
	fmt.Println("\nMigrating plists from Redis to SQLite...")
	fmt.Println("(Using pipelined reads for efficiency)")
	fmt.Println("(Skipping any that already exist in SQLite)")
	fmt.Println()

	migrated := 0
	skipped := 0
	failed := 0
	count := 0

	// Process keys in batches of 100
	for batch := range batches(redisKeys(ctx, redisClient, "launchd_plist:*"), 100) {
		stats := processBatch(ctx, redisClient, db, batch, count)
		count += stats.processed
		migrated += stats.migrated
		skipped += stats.skipped
		failed += stats.failed
	}

	// Summary
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Printf("Migration complete!\n")
	fmt.Printf("  Migrated: %d\n", migrated)
	fmt.Printf("  Skipped: %d\n", skipped)
	fmt.Printf("  Failed: %d\n", failed)
	fmt.Printf("  Total: %d\n", count)
	fmt.Println(strings.Repeat("=", 50))

	if migrated > 0 {
		fmt.Printf("\nDatabase saved to: %s\n", sqlitePath)
		fmt.Println("You can now upload this file to your Fly.io volume.")
	}
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

type batchStats struct {
	processed int
	migrated  int
	skipped   int
	failed    int
}

func processBatch(ctx context.Context, redisClient *redis.Client, db *sql.DB, keys []string, startCount int) batchStats {
	stats := batchStats{processed: len(keys)}

	// Extract IDs from keys
	ids := make([]string, len(keys))
	for i, key := range keys {
		ids[i] = key[len("launchd_plist:"):]
	}

	// Batch-check which IDs already exist in SQLite
	existingIDs := checkExistingIDs(db, ids)

	// Filter out already-migrated keys
	keysToFetch := make([]string, 0, len(keys))
	idsToFetch := make([]string, 0, len(keys))

	for i, key := range keys {
		id := ids[i]
		if existingIDs[id] {
			fmt.Printf("[%d] %s SKIPPED (already migrated)\n", startCount+i+1, id)
			stats.skipped++
		} else {
			keysToFetch = append(keysToFetch, key)
			idsToFetch = append(idsToFetch, id)
		}
	}

	// If all were skipped, we're done
	if len(keysToFetch) == 0 {
		return stats
	}

	// Pipeline HGETALL for remaining keys
	pipe := redisClient.Pipeline()
	cmds := make([]*redis.MapStringStringCmd, len(keysToFetch))

	for i, key := range keysToFetch {
		cmds[i] = pipe.HGetAll(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Printf("Pipeline error: %v", err)
		stats.failed = len(keysToFetch)
		return stats
	}

	// Process pipelined results
	for i, id := range idsToFetch {
		idx := 0
		for j, origID := range ids {
			if origID == id {
				idx = j
				break
			}
		}
		fmt.Printf("[%d] Processing %s... ", startCount+idx+1, id)

		data := cmds[i].Val()
		if len(data) == 0 {
			fmt.Println("SKIPPED (empty)")
			stats.skipped++
			continue
		}

		// Convert Redis hash to plist struct
		plist := LaunchdPlist{
			Name:              data["name"],
			Command:           data["command"],
			StartInterval:     data["start_interval"],
			Minute:            data["minute"],
			Hour:              data["hour"],
			DayOfMonth:        data["day_of_month"],
			Month:             data["month"],
			Weekday:           data["weekday"],
			RunAtLoad:         data["run_at_load"],
			RestartOnCrash:    data["restart_on_crash"],
			StartOnMount:      data["start_on_mount"],
			QueueDirectories:  data["queue_directories"],
			Environment:       data["environment"],
			User:              data["user"],
			Group:             data["group"],
			WorkingDirectory:  data["working_directory"],
			RootDirectory:     data["root_directory"],
			StandardOutPath:   data["standard_out_path"],
			StandardErrorPath: data["standard_error_path"],
			CreatedAt:         data["created_at"],
		}

		// Marshal to JSON
		jsonData, err := json.Marshal(plist)
		if err != nil {
			fmt.Printf("ERROR: failed to marshal JSON: %v\n", err)
			stats.failed++
			continue
		}

		// Insert into SQLite
		_, err = db.Exec(
			"INSERT INTO launchd_plists (id, data, created_at) VALUES (?, ?, ?)",
			id, string(jsonData), plist.CreatedAt)

		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			stats.failed++
		} else {
			fmt.Println("✓")
			stats.migrated++
		}
	}

	return stats
}

// checkExistingIDs returns a map of IDs that already exist in SQLite
func checkExistingIDs(db *sql.DB, ids []string) map[string]bool {
	if len(ids) == 0 {
		return map[string]bool{}
	}

	existing := make(map[string]bool)

	// Build WHERE IN clause
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf("SELECT id FROM launchd_plists WHERE id IN (%s)",
		strings.Join(placeholders, ","))

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Printf("Error checking existing IDs: %v", err)
		return existing
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err == nil {
			existing[id] = true
		}
	}

	return existing
}

// batches converts an iterator into batches of a given size
func batches[T any](seq iter.Seq[T], size int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		batch := make([]T, 0, size)
		for item := range seq {
			batch = append(batch, item)
			if len(batch) >= size {
				if !yield(batch) {
					return
				}
				batch = make([]T, 0, size)
			}
		}
		// Yield remaining items
		if len(batch) > 0 {
			yield(batch)
		}
	}
}

// redisKeys returns an iterator over Redis keys matching the pattern
func redisKeys(ctx context.Context, client *redis.Client, pattern string) iter.Seq[string] {
	return func(yield func(string) bool) {
		var cursor uint64
		for {
			keys, newCursor, err := client.Scan(ctx, cursor, pattern, 100).Result()
			if err != nil {
				log.Printf("Error scanning Redis: %v", err)
				return
			}

			for _, key := range keys {
				if !yield(key) {
					return
				}
			}

			cursor = newCursor
			if cursor == 0 {
				return
			}
		}
	}
}
