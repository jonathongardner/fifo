package gzip

import (
	"compress/gzip"
	"io"
)

type Reader struct {
	gr *gzip.Reader
	r  io.ReadSeeker
}

// NewReader Creates a new gzip reader
// r must be a ReadSeeker, needed to seek to reset the reader
func NewReader(r io.ReadSeeker) (*Reader, error) {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &Reader{gr: gr, r: r}, nil
}

// Read reads data from the gzip reader
func (r *Reader) Read(p []byte) (int, error) {
	return r.gr.Read(p)
}

// Close closes the gzip reader and the underlying reader if it implements io.Closer
func (r *Reader) Close() error {
	if err := r.gr.Close(); err != nil {
		return err
	}
	if c, ok := r.r.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

// Reset resets the gzip reader to the beginning of the stream
func (r *Reader) Reset() error {
	r.r.Seek(0, io.SeekStart)
	if err := r.gr.Reset(r.r); err != nil {
		return err
	}
	return nil
}
