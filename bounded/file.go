package bounded

import (
	"fmt"
	"io"
	"io/fs"
	"os"
)

type File struct {
	path string
	file *os.File
	s    int64
	p    int64
	l    int64
	pool *FilePool
}

func NewFile(path string) *File {
	return &File{
		path: path,
		file: nil,
		s:    -1,
		p:    0,
		l:    0,
	}
}

// openFileIfNot checks if the file is already open, and if not, opens it.
func (b *File) openFileIfNot() error {
	// Check if the file is already open
	if b.file != nil {
		return nil
	}

	file, err := os.Open(b.path)
	if err != nil {
		return err
	}

	b.file = file
	return nil
}

func (b *File) Bound(start, length int64) (*File, error) {
	if start < 0 {
		return nil, fmt.Errorf("start position cannot be negative")
	}
	if length < 0 {
		return nil, fmt.Errorf("length cannot be negative")
	}

	if err := b.openFileIfNot(); err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}

	_, err := b.file.Seek(start, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("error seeking file: %v", err)
	}

	b.s = start
	b.p = 0
	b.l = length

	return b, nil
}

func (b *File) left() (n int64) {
	return b.l - b.p
}

func (b *File) error(op string, err error) error {
	return &fs.PathError{Op: op, Path: b.path, Err: err}
}

func (b *File) Read(p []byte) (n int, err error) {
	if b.s == -1 {
		return 0, b.error("read", fs.ErrClosed)
	}

	if b.p >= b.l {
		return 0, io.EOF
	}
	s := b.left()
	if int64(len(p)) > s {
		p = p[:s]
	}
	n, err = b.file.Read(p)
	b.p += int64(n)
	return n, err
}

func (b *File) Seek(offset int64, whence int) (int64, error) {
	if b.s == -1 {
		return 0, b.error("seek", fs.ErrClosed)
	}

	var newPos int64
	if whence == io.SeekStart {
		newPos = offset
	} else if whence == io.SeekCurrent {
		newPos = b.p + offset
	} else if whence == io.SeekEnd {
		newPos = b.l + offset
	} else {
		return 0, b.error("seek", fs.ErrInvalid)
	}

	if newPos < 0 {
		return 0, b.error("seek", fs.ErrInvalid)
	}

	relativeOffset, err := b.file.Seek(newPos+b.s, io.SeekStart)
	b.p = relativeOffset - b.s
	return b.p, err
}
func (b *File) Close() error {
	if b.s == -1 {
		return b.error("close", fs.ErrClosed)
	}
	b.s = -1
	b.l = 0
	b.p = 0
	if b.pool != nil {
		b.pool.Release(b)
	}
	return nil
}
func (b *File) Cleanup() error {
	file := b.file
	b.file = nil
	if b.s == -1 {
		return nil
	}

	return file.Close()
}
