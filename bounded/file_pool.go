package bounded

import (
	"fmt"
	"sync"
)

var ErrNoAvailableFile = fmt.Errorf("no available bounded file")

type FilePool struct {
	bfs map[*File]bool
	mu  *sync.Mutex
}

// NewFilePool creates a new shared bounded file pool
func NewFilePool(path string, count int) (*FilePool, error) {
	if count > 50 {
		return nil, fmt.Errorf("maximum number of files exceeded")
	}
	if count < 1 {
		return nil, fmt.Errorf("minimum number of files is 1")
	}
	bfp := &FilePool{bfs: map[*File]bool{}, mu: &sync.Mutex{}}
	for i := 0; i < count; i++ {
		bf := NewFile(path)
		bfp.bfs[bf] = true
		bf.pool = bfp
	}
	return bfp, nil
}

// Release the bounded file back in the pool
func (bfp *FilePool) Release(b *File) {
	bfp.mu.Lock()
	defer bfp.mu.Unlock()
	// if _, ok := bfp.bfs[b]; ok {
	// 	return fmt.Errorf("File already released")
	// }

	bfp.bfs[b] = true
}

// Get gets a bounded file from the pool
func (bfp *FilePool) Get(start, length int64) (*File, error) {
	bfp.mu.Lock()
	defer bfp.mu.Unlock()

	for bf := range bfp.bfs {
		delete(bfp.bfs, bf)
		return bf.Bound(start, length)
	}
	return nil, ErrNoAvailableFile
}

// Cleanup cleanups all the bounded files
func (bfp *FilePool) Cleanup() error {
	bfp.mu.Lock()
	defer bfp.mu.Unlock()

	errs := []error{}
	for bf := range bfp.bfs {
		err := bf.Cleanup()
		if err != nil {
			errs = append(errs, err)
		}
		delete(bfp.bfs, bf)
	}
	if len(errs) > 0 {
		return fmt.Errorf("error cleaning up files: %v", errs)
	}
	return nil
}
