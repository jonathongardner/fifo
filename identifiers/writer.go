package identifiers

import (
	"compress/gzip"
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
	md5       hash.Hash
	sha1      hash.Hash
	sha256    hash.Hash
	sha512    hash.Hash
	entropy   *entropy.Writer
	cache     *cache.Writer
	gzipCache *cache.Writer
	toClose   []io.Closer
	ftype     bool
	mw        io.Writer
}

// NewMultiWriterWithOptions creates a new Writer that writes to the given io.Writer
// and calculates info based on options
func newWriterWithOptions(o Options, w ...io.Writer) *Writer {
	toReturn := &Writer{toClose: make([]io.Closer, 0)}
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
	// Always set cached size cause its used to calculate the size
	toReturn.cache = cache.NewWriter(o.minCachSize())
	w = append(w, toReturn.cache)

	if o.GzipCacheSize != 0 {
		toReturn.gzipCache = cache.NewWriter(o.GzipCacheSize)
		writer := gzip.NewWriter(toReturn.gzipCache)
		toReturn.toClose = append(toReturn.toClose, writer)
		w = append(w, writer)
	}

	toReturn.mw = io.MultiWriter(w...)
	return toReturn
}

// NewMultiWriter creates a new Writer that writes to the given io.Writer and defaults everything
func NewWriter(w ...io.Writer) *Writer {
	return newWriterWithOptions(NewDefultOptions(), w...)
}
func (mw *Writer) Write(p []byte) (int, error) {
	if mw.toClose == nil {
		return 0, os.ErrClosed
	}

	return mw.mw.Write(p)
}

func (mw *Writer) Cache() *cache.Writer {
	return mw.cache
}

func (mw *Writer) GzipCache() *cache.Writer {
	return mw.gzipCache
}

func (mw *Writer) Sha256Bytes() []byte {
	if mw.sha256 == nil {
		return []byte{}
	}
	return mw.sha256.Sum(nil)
}

func (mw *Writer) Close() error {
	if mw.toClose == nil {
		return os.ErrClosed
	}

	for _, c := range mw.toClose {
		if err := c.Close(); err != nil {
			return err
		}
	}

	mw.toClose = nil
	return nil
}
