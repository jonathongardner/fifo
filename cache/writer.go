package cache

import (
	"bytes"
	"io"
	// log "github.com/sirupsen/logrus"
)

// Writer is a writer that caches the data in it, up to a certain size
// io.Copy() or io.MultiWriter() to detect the file type of a stream
type Writer struct {
	size int64
	max  int64
	data []byte
}

// NewWriter creates a new Cached writer
func NewWriter(max int64) *Writer {
	return &Writer{size: 0, data: make([]byte, 0), max: max}
}

// Write writes data to the writer
func (mw *Writer) Write(p []byte) (n int, err error) {
	len := len(p)
	toCopy := 0

	if mw.size < mw.max {
		toCopy = int(mw.max - mw.size)

		if toCopy > len {
			toCopy = len
		}

		mw.data = append(mw.data, p[:toCopy]...)
	}

	mw.size += int64(len)
	return len, nil
}

func (mw *Writer) IsCached() bool {
	return mw.max >= mw.size
}

func (mw *Writer) Size() int64 {
	return mw.size
}

// Data returns the data written to the writer
func (mw *Writer) Bytes() []byte {
	return mw.data
}

// NewReader Create a new reader from the current data
func (mw *Writer) NewReader() io.Reader {
	return bytes.NewReader(mw.data)
}

// Reset resets the writer
func (mw *Writer) Reset() error {
	mw.size = 0
	mw.data = make([]byte, 0)
	return nil
}
