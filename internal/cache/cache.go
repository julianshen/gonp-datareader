// Package cache provides response caching functionality.
package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

var (
	// ErrNilCache is returned when operating on a nil cache
	ErrNilCache = errors.New("cache is nil")
)

// cacheEntry represents a cached item with metadata.
type cacheEntry struct {
	Data      []byte    `json:"data"`
	ExpiresAt time.Time `json:"expires_at"`
}

// FileCache implements file-based caching.
type FileCache struct {
	dir string
}

// NewFileCache creates a new file-based cache in the specified directory.
func NewFileCache(dir string) *FileCache {
	return &FileCache{
		dir: dir,
	}
}

// Get retrieves a value from the cache.
// Returns the value and true if found and not expired, or nil and false otherwise.
func (c *FileCache) Get(key string) ([]byte, bool) {
	if c == nil {
		return nil, false
	}

	filename := c.filename(key)

	// Read file
	// #nosec G304 - File path is constructed from hashed key, scoped to cache directory
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, false
	}

	// Decode entry
	var entry cacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false
	}

	// Check expiration (zero time means no expiration)
	if !entry.ExpiresAt.IsZero() && time.Now().After(entry.ExpiresAt) {
		// Expired, delete it (ignore error as cleanup is best-effort)
		_ = os.Remove(filename)
		return nil, false
	}

	return entry.Data, true
}

// Set stores a value in the cache with the specified TTL.
// A TTL of 0 means no expiration.
func (c *FileCache) Set(key string, value []byte, ttl time.Duration) error {
	if c == nil {
		return ErrNilCache
	}

	// Ensure cache directory exists
	// #nosec G301 - Cache directory needs to be readable by owner, group, and others for flexibility
	if err := os.MkdirAll(c.dir, 0755); err != nil {
		return err
	}

	// Calculate expiration time
	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}
	// Zero time for no expiration

	entry := cacheEntry{
		Data:      value,
		ExpiresAt: expiresAt,
	}

	// Encode entry
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	// Write to file
	filename := c.filename(key)
	// #nosec G306 - Cache files need to be readable by owner, group, and others for flexibility
	return os.WriteFile(filename, data, 0644)
}

// Delete removes a value from the cache.
func (c *FileCache) Delete(key string) error {
	if c == nil {
		return ErrNilCache
	}

	filename := c.filename(key)
	err := os.Remove(filename)
	if os.IsNotExist(err) {
		return nil // Already deleted
	}
	return err
}

// filename generates a safe filename for the given key by hashing it.
func (c *FileCache) filename(key string) string {
	hash := sha256.Sum256([]byte(key))
	hashStr := hex.EncodeToString(hash[:])
	return filepath.Join(c.dir, hashStr+".cache")
}
