package http

import (
	"bytes"
	"strings"
	"testing"
)

func TestBufferPool_GetPut(t *testing.T) {
	pool := NewBufferPool()

	// Get a buffer
	buf1 := pool.Get()
	if buf1 == nil {
		t.Fatal("Get() returned nil buffer")
	}

	// Write some data
	buf1.WriteString("test data")
	if buf1.String() != "test data" {
		t.Errorf("Expected 'test data', got %q", buf1.String())
	}

	// Return to pool
	pool.Put(buf1)

	// Get another buffer - should be the same one, but reset
	buf2 := pool.Get()
	if buf2.Len() != 0 {
		t.Errorf("Expected empty buffer after reset, got length %d", buf2.Len())
	}
}

func TestBufferPool_LargeBufferNotReturned(t *testing.T) {
	pool := NewBufferPool()

	buf := pool.Get()

	// Make buffer very large (over 1MB limit)
	largeData := make([]byte, 2*1024*1024) // 2MB
	buf.Write(largeData)

	// This should not panic, but buffer won't be returned to pool
	pool.Put(buf)

	// Get another buffer - should be a new one
	buf2 := pool.Get()
	if buf2.Cap() > 1024*1024 {
		t.Error("Expected new buffer, got large capacity buffer from pool")
	}
}

func TestBufferPool_CopyWithPool(t *testing.T) {
	pool := NewBufferPool()

	testData := "Hello, World!"
	reader := strings.NewReader(testData)

	buf, err := pool.CopyWithPool(reader)
	if err != nil {
		t.Fatalf("CopyWithPool() error = %v", err)
	}

	if buf.String() != testData {
		t.Errorf("Expected %q, got %q", testData, buf.String())
	}

	// Return to pool
	pool.Put(buf)
}

func TestBufferPool_CopyWithPool_Error(t *testing.T) {
	pool := NewBufferPool()

	// Create a reader that returns an error
	errorReader := &errorReader{}

	buf, err := pool.CopyWithPool(errorReader)
	if err == nil {
		t.Error("Expected error from CopyWithPool, got nil")
	}

	if buf != nil {
		t.Error("Expected nil buffer on error, got non-nil")
	}
}

func TestDefaultBufferPool(t *testing.T) {
	// Test default pool functions
	buf := GetBuffer()
	if buf == nil {
		t.Fatal("GetBuffer() returned nil")
	}

	buf.WriteString("test")
	PutBuffer(buf)

	// Get another buffer
	buf2 := GetBuffer()
	if buf2.Len() != 0 {
		t.Errorf("Expected empty buffer, got length %d", buf2.Len())
	}

	PutBuffer(buf2)
}

func TestBufferPool_Concurrent(t *testing.T) {
	pool := NewBufferPool()

	// Test concurrent access
	const goroutines = 10
	done := make(chan bool, goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				buf := pool.Get()
				buf.WriteString("test")
				pool.Put(buf)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < goroutines; i++ {
		<-done
	}
}

// errorReader is a test helper that always returns an error
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, bytes.ErrTooLarge
}

// Benchmark tests

func BenchmarkBufferPool_GetPut(b *testing.B) {
	pool := NewBufferPool()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf := pool.Get()
		buf.WriteString("benchmark data")
		pool.Put(buf)
	}
}

func BenchmarkBufferPool_WithoutPool(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(make([]byte, 0, 65536))
		buf.WriteString("benchmark data")
	}
}

func BenchmarkBufferPool_CopyWithPool(b *testing.B) {
	pool := NewBufferPool()
	testData := strings.Repeat("test data", 1000)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(testData)
		buf, err := pool.CopyWithPool(reader)
		if err != nil {
			b.Fatal(err)
		}
		pool.Put(buf)
	}
}
