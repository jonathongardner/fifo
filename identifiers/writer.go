package stats

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"io"

	"github.com/jonathongardner/fifo/cache"
	"github.com/jonathongardner/fifo/entropy"
	"github.com/jonathongardner/fifo/filetype"
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
	size    int64
	mw      io.Writer
}

// NewMultiWriterWithOptions creates a new Writer that writes to the given io.Writer
// and calculates info based on options
func newWriterWithOptions(o Options, w ...io.Writer) *Writer {
	toReturn := &Writer{
		size: 0,
	}
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
	if o.Filetype || o.CacheSize != 0 {
		size := o.CacheSize
		if size == 0 {
			size = uint64(filetype.MaxBytesFileDetect())
		}
		toReturn.cache = cache.NewWriter(size)
		w = append(w, toReturn.cache)
	}

	toReturn.mw = io.MultiWriter(w...)
	return toReturn
}

// NewMultiWriter creates a new Writer that writes to the given io.Writer and defaults everything
func NewWriter(w ...io.Writer) *Writer {
	return newWriterWithOptions(NewDefultOptions(), w...)
}
func (mw *Writer) Write(p []byte) (int, error) {
	n, err := mw.mw.Write(p)
	mw.size += int64(n)
	return n, err
}

func (mw *Writer) Cache() *cache.Writer {
	return mw.cache
}
