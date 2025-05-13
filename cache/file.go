package cache

import (
	"io"
	"os"
	// log "github.com/sirupsen/logrus"
)

// Writer is a writer that caches the data in it, up to a certain size
// io.Copy() or io.MultiWriter() to detect the file type of a stream
type File struct {
	cache *Writer
	file  *os.File
}

// NewWriter creates a new Cached writer
func Open(path string, cache int64) (*File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &File{file: file, cache: NewWriter(cache)}, nil
}

func (f *File) Read(p []byte) (int, error) {
	if f.file == nil {
		return 0, os.ErrClosed
	}

	n1, err1 := f.file.Read(p)
	_, err2 := f.cache.Write(p[:n1])
	if err2 != nil {
		return n1, err2
	}

	return n1, err1
}

type ReadSeekCloseCacher interface {
	io.ReadSeekCloser
	IsCached() bool
}
type fileWrapper struct {
	*os.File
}

func (f *fileWrapper) IsCached() bool {
	return false
}

// NewReader Create a new reader from the current data
func (f *File) NewReader() (ReadSeekCloseCacher, error) {
	if f.file == nil {
		return nil, os.ErrClosed
	}

	defer func() { f.file = nil }()

	if f.cache.IsCached() {
		err := f.file.Close()
		if err != nil {
			return nil, err
		}
		return f.cache.NewReader(), nil
	}

	f.file.Seek(0, io.SeekStart)
	return &fileWrapper{f.file}, nil
}
