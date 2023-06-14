package main

import (
	"context"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/redis/go-redis/v9"
)

type PlistStore struct {
	db *redis.Client
}

func NewPlistStore(redisAddress, redisPassword string) *PlistStore {
	db := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: redisPassword,
	})
	return &PlistStore{
		db: db,
	}
}

func (s *PlistStore) Save(plist LaunchdPlist) (string, error) {
	plist.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	id, _ := gonanoid.New(10)
	key := "launchd_plist:" + id
	cmd := s.db.HSet(context.Background(), key, plist)
	return id, cmd.Err()
}

func (s *PlistStore) Load(id string) (LaunchdPlist, bool, error) {
	plist := LaunchdPlist{}
	result := s.db.HGetAll(context.Background(), "launchd_plist:"+id)
	if len(result.Val()) == 0 || result.Err() != nil {
		return plist, false, result.Err()
	}
	err := result.Scan(&plist)
	plist.ID = id
	return plist, true, err
}
