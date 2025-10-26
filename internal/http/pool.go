package http

import (
	"bytes"
	"io"
	"sync"
)

// BufferPool is a pool of byte buffers for reuse.
// This reduces memory allocations when reading HTTP responses.
type BufferPool struct {
	pool *sync.Pool
}

// NewBufferPool creates a new buffer pool.
func NewBufferPool() *BufferPool {
	return &BufferPool{
		pool: &sync.Pool{
			New: func() interface{} {
				// Pre-allocate with reasonable default size (64KB)
				return bytes.NewBuffer(make([]byte, 0, 65536))
			},
		},
	}
}

// Get retrieves a buffer from the pool.
func (p *BufferPool) Get() *bytes.Buffer {
	buf := p.pool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// Put returns a buffer to the pool for reuse.
func (p *BufferPool) Put(buf *bytes.Buffer) {
	// Don't return buffers that are too large to the pool
	// to avoid keeping excessive memory
	if buf.Cap() > 1024*1024 { // 1MB limit
		return
	}
	p.pool.Put(buf)
}

// CopyWithPool copies from reader to a buffer from the pool.
// The caller is responsible for returning the buffer to the pool.
func (p *BufferPool) CopyWithPool(r io.Reader) (*bytes.Buffer, error) {
	buf := p.Get()
	_, err := io.Copy(buf, r)
	if err != nil {
		p.Put(buf)
		return nil, err
	}
	return buf, nil
}

// defaultBufferPool is the global buffer pool instance.
var defaultBufferPool = NewBufferPool()

// GetBuffer retrieves a buffer from the default pool.
func GetBuffer() *bytes.Buffer {
	return defaultBufferPool.Get()
}

// PutBuffer returns a buffer to the default pool.
func PutBuffer(buf *bytes.Buffer) {
	defaultBufferPool.Put(buf)
}
