package cache_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/julianshen/gonp-datareader/internal/cache"
)

func TestFileCache_SetAndGet(t *testing.T) {
	// Create temporary directory for cache
	tmpDir, err := os.MkdirTemp("", "cache-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	c := cache.NewFileCache(tmpDir)

	key := "test-key"
	value := []byte("test data")

	// Set value
	err = c.Set(key, value, 1*time.Hour)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Get value
	retrieved, found := c.Get(key)
	if !found {
		t.Fatal("Expected to find cached value")
	}

	if string(retrieved) != string(value) {
		t.Errorf("Expected %q, got %q", string(value), string(retrieved))
	}
}

func TestFileCache_GetMissingKey(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cache-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	c := cache.NewFileCache(tmpDir)

	_, found := c.Get("nonexistent-key")
	if found {
		t.Error("Expected not to find nonexistent key")
	}
}

func TestFileCache_TTLExpired(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cache-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	c := cache.NewFileCache(tmpDir)

	key := "test-key"
	value := []byte("test data")

	// Set with very short TTL
	err = c.Set(key, value, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Should be available immediately
	_, found := c.Get(key)
	if !found {
		t.Fatal("Expected to find value immediately")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should be expired now
	_, found = c.Get(key)
	if found {
		t.Error("Expected value to be expired")
	}
}

func TestFileCache_Delete(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cache-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	c := cache.NewFileCache(tmpDir)

	key := "test-key"
	value := []byte("test data")

	err = c.Set(key, value, 1*time.Hour)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Delete the key
	err = c.Delete(key)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Should not find it anymore
	_, found := c.Get(key)
	if found {
		t.Error("Expected not to find deleted key")
	}
}

func TestFileCache_ZeroTTL(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cache-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	c := cache.NewFileCache(tmpDir)

	key := "test-key"
	value := []byte("test data")

	// Zero TTL means no expiration
	err = c.Set(key, value, 0)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Should still be available
	_, found := c.Get(key)
	if !found {
		t.Fatal("Expected to find value with zero TTL")
	}
}

func TestFileCache_SafeFilenames(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cache-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	c := cache.NewFileCache(tmpDir)

	// Test with special characters that need hashing
	keys := []string{
		"http://example.com/path?query=value",
		"key/with/slashes",
		"key:with:colons",
		"key*with*asterisks",
	}

	for _, key := range keys {
		value := []byte("test data for " + key)

		err = c.Set(key, value, 1*time.Hour)
		if err != nil {
			t.Errorf("Set failed for key %q: %v", key, err)
			continue
		}

		retrieved, found := c.Get(key)
		if !found {
			t.Errorf("Expected to find key %q", key)
			continue
		}

		if string(retrieved) != string(value) {
			t.Errorf("Value mismatch for key %q", key)
		}
	}
}

func TestFileCache_CreatesCacheDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cache-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Use a subdirectory that doesn't exist yet
	cacheDir := filepath.Join(tmpDir, "subdir", "cache")

	c := cache.NewFileCache(cacheDir)

	// Should create directory when setting first value
	err = c.Set("test", []byte("data"), 1*time.Hour)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Directory should exist now
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		t.Error("Expected cache directory to be created")
	}
}

func TestNilCache(t *testing.T) {
	var c *cache.FileCache

	// Nil cache should not panic
	_, found := c.Get("key")
	if found {
		t.Error("Nil cache should return not found")
	}

	err := c.Set("key", []byte("value"), 1*time.Hour)
	if err == nil {
		t.Error("Nil cache Set should return error")
	}

	err = c.Delete("key")
	if err == nil {
		t.Error("Nil cache Delete should return error")
	}
}

// Benchmark tests

func BenchmarkFileCache_Set(b *testing.B) {
	tempDir := b.TempDir()
	cache := cache.NewFileCache(tempDir)

	testData := []byte("test data for benchmarking cache set operation")
	ttl := 1 * time.Hour

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("benchmark-key-%d", i)
		err := cache.Set(key, testData, ttl)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFileCache_Get(b *testing.B) {
	tempDir := b.TempDir()
	cache := cache.NewFileCache(tempDir)

	// Pre-populate cache
	testData := []byte("test data for benchmarking cache get operation")
	ttl := 1 * time.Hour
	key := "benchmark-key"

	err := cache.Set(key, testData, ttl)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, found := cache.Get(key)
		if !found {
			b.Fatal("expected to find cached value")
		}
	}
}

func BenchmarkFileCache_SetAndGet(b *testing.B) {
	tempDir := b.TempDir()
	cache := cache.NewFileCache(tempDir)

	testData := []byte("test data for benchmarking cache set and get operations")
	ttl := 1 * time.Hour

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("benchmark-key-%d", i)

		err := cache.Set(key, testData, ttl)
		if err != nil {
			b.Fatal(err)
		}

		_, found := cache.Get(key)
		if !found {
			b.Fatal("expected to find cached value")
		}
	}
}
