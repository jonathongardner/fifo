package identifiers

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"io"
	"os"

	"github.com/jonathongardner/fifo/cache"
	"github.com/jonathongardner/fifo/entropy"
)

// Writer is a writer that calculates the md5, sha1, sha256, sha512 hashes
// and the entropy of the data written to it. It also detects the file type
type Writer struct {
	md5     hash.Hash
	sha1    hash.Hash
	sha256  hash.Hash
	sha512  hash.Hash
	entropy *entropy.Writer
	cache   *cache.Writer
	ftype   bool
	mw      io.Writer
	closed  bool
}

// NewMultiWriterWithOptions creates a new Writer that writes to the given io.Writer
// and calculates info based on options
func newWriterWithOptions(o Options, w ...io.Writer) *Writer {
	toReturn := &Writer{}
	if o.Md5 {
		toReturn.md5 = md5.New()
		w = append(w, toReturn.md5)
	}
	if o.Sha1 {
		toReturn.sha1 = sha1.New()
		w = append(w, toReturn.sha1)
	}
	if o.Sha256 {
		toReturn.sha256 = sha256.New()
		w = append(w, toReturn.sha256)
	}
	if o.Sha512 {
		toReturn.sha512 = sha512.New()
		w = append(w, toReturn.sha512)
	}
	if o.Entropy {
		toReturn.entropy = entropy.NewWriter()
		w = append(w, toReturn.entropy)
	}
	toReturn.ftype = o.Filetype
	// Always set cached cause its used to calculate the size
	toReturn.cache = cache.NewWriter(o.minCachSize())
	w = append(w, toReturn.cache)

	toReturn.mw = io.MultiWriter(w...)
	return toReturn
}

// NewMultiWriter creates a new Writer that writes to the given io.Writer and defaults everything
func NewWriter(w ...io.Writer) *Writer {
	return newWriterWithOptions(NewDefultOptions(), w...)
}
func (mw *Writer) Write(p []byte) (int, error) {
	if mw.closed {
		return 0, os.ErrClosed
	}
	return mw.mw.Write(p)
}

func (iw *Writer) AddWriter(w io.Writer) {
	// multiwriter will flatten out the writers
	iw.mw = io.MultiWriter(w, iw.mw)
}

func (mw *Writer) Cache() *cache.Writer {
	return mw.cache
}

func (mw *Writer) Close() error {
	if mw.closed {
		return os.ErrClosed
	}
	mw.closed = true

	return nil
}

func (mw *Writer) Reset(w ...io.Writer) {
	mw.closed = false
	if mw.md5 != nil {
		mw.md5.Reset()
		w = append(w, mw.md5)
	}
	if mw.sha1 != nil {
		mw.sha1.Reset()
		w = append(w, mw.sha1)
	}
	if mw.sha256 != nil {
		mw.sha256.Reset()
		w = append(w, mw.sha256)
	}
	if mw.sha512 != nil {
		mw.sha512.Reset()
		w = append(w, mw.sha512)
	}
	if mw.entropy != nil {
		mw.entropy.Reset()
		w = append(w, mw.entropy)
	}

	mw.cache.Reset()
	w = append(w, mw.cache)
	mw.mw = io.MultiWriter(w...)
}
