package cache

import (
	"fmt"
	"io"
	"os"
	// log "github.com/sirupsen/logrus"
)

var ErrAlreadyOpen = fmt.Errorf("file already open")

// File is a wrapper around os.File that caches the data read from it
type File struct {
	cache *Writer
	file  *os.File
}

// Open creates a new file object and opens the file at the given path
func Open(path string, cache int64) (*File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &File{file: file, cache: NewWriter(cache)}, nil
}

// NewFile creates a new file object but does NOT open a file
func NewFile(cache int64) *File {
	return &File{file: nil, cache: NewWriter(cache)}
}

// Open opens the file at the given path
// It returns an error if the file is already open or if there is an error opening the file
func (f *File) Open(path string) error {
	if f.file != nil {
		return ErrAlreadyOpen
	}
	var err error
	f.file, err = os.Open(path)
	return err
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

func (f *File) Reset() error {
	if err := f.cache.Reset(); err != nil {
		return err
	}
	if f.file == nil {
		return nil
	}
	defer func() { f.file = nil }()
	return f.file.Close()
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
