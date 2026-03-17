package gorbac

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

type memoryCacheStore struct {
	data map[string][]byte
}

func newMemoryCacheStore() *memoryCacheStore {
	return &memoryCacheStore{data: make(map[string][]byte)}
}

func (s *memoryCacheStore) Get(ctx context.Context, key string) ([]byte, bool, error) {
	v, ok := s.data[key]
	return v, ok, nil
}

func (s *memoryCacheStore) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	s.data[key] = value
	return nil
}

func (s *memoryCacheStore) Del(ctx context.Context, keys ...string) error {
	for _, k := range keys {
		delete(s.data, k)
	}
	return nil
}

func TestDefaultManager_LoadFromCache_UsesCacheStoreSnapshot(t *testing.T) {
	store := newMemoryCacheStore()
	manager := NewDefaultManager(nil, true)
	manager.SetCacheStore(store, "test", time.Minute)

	snap := cacheSnapshot{
		Version: cacheSnapshotVersion,
		Items: []cacheSnapshotItem{
			{
				Name:        "admin",
				Type:        RoleType.Value(),
				Description: "管理员",
				CreateTime:  time.Unix(1, 0),
				UpdateTime:  time.Unix(2, 0),
			},
			{
				Name:        "view",
				Type:        PermissionType.Value(),
				Description: "查看",
				CreateTime:  time.Unix(3, 0),
				UpdateTime:  time.Unix(4, 0),
			},
		},
		Rules: []cacheSnapshotRule{
			{
				Name:        "r1",
				ExecuteName: "demo",
				CreateTime:  time.Unix(5, 0),
				UpdateTime:  time.Unix(6, 0),
			},
		},
		Parents: map[string][]string{
			"view": {"admin"},
		},
	}

	payload, err := json.Marshal(snap)
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}
	if err := store.Set(context.Background(), manager.cacheSnapshotKey(), payload, time.Minute); err != nil {
		t.Fatalf("set snapshot: %v", err)
	}

	manager.loadFromCache()

	if manager.cache.items["admin"] == nil || manager.cache.items["view"] == nil {
		t.Fatalf("expected items loaded from cache store")
	}
	if got := manager.cache.parents["view"]; len(got) != 1 || got[0] != "admin" {
		t.Fatalf("expected parents loaded from cache store, got=%v", got)
	}
	if manager.cache.rules["r1"] == nil || manager.cache.rules["r1"].ExecuteName != "demo" {
		t.Fatalf("expected rules loaded from cache store")
	}
}

