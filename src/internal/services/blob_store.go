package services

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type blobEntry struct {
	data        []byte
	contentType string
	expiresAt   time.Time
}

// BlobStore keeps uploaded image bytes for a limited time.
type BlobStore struct {
	mu    sync.RWMutex
	items map[string]blobEntry
	ttl   time.Duration
}

// NewBlobStore creates a new BlobStore with background cleanup.
func NewBlobStore(ttl time.Duration) *BlobStore {
	store := &BlobStore{
		items: make(map[string]blobEntry),
		ttl:   ttl,
	}
	store.startCleanup(10 * time.Minute)
	return store
}

// Put stores image bytes and returns an ID.
func (s *BlobStore) Put(data []byte, contentType string) string {
	id := newBlobID()
	s.mu.Lock()
	s.items[id] = blobEntry{
		data:        data,
		contentType: contentType,
		expiresAt:   time.Now().Add(s.ttl),
	}
	s.mu.Unlock()
	return id
}

// Get returns image bytes if not expired.
func (s *BlobStore) Get(id string) ([]byte, string, bool) {
	s.mu.RLock()
	entry, ok := s.items[id]
	s.mu.RUnlock()
	if !ok {
		return nil, "", false
	}
	if time.Now().After(entry.expiresAt) {
		s.mu.Lock()
		delete(s.items, id)
		s.mu.Unlock()
		return nil, "", false
	}
	return entry.data, entry.contentType, true
}

func (s *BlobStore) startCleanup(interval time.Duration) {
	if interval <= 0 {
		return
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			s.mu.Lock()
			for key, entry := range s.items {
				if now.After(entry.expiresAt) {
					delete(s.items, key)
				}
			}
			s.mu.Unlock()
		}
	}()
}

func newBlobID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return hex.EncodeToString([]byte(time.Now().Format("20060102150405.000000000")))
	}
	return hex.EncodeToString(buf)
}
