package identifiers

import (
	"io"

	"github.com/jonathongardner/fifo/filetype"
)

// Options is a struct that contains options for the Writer
type Options struct {
	Md5       bool
	Sha1      bool
	Sha256    bool
	Sha512    bool
	Entropy   bool
	Filetype  bool
	CacheSize int64 // use 0 for no cache
}

// NewDefultOptions creates a new Options struct with default values
// md5, sha1, sha256, sha512, entropy, filetype are all true
func NewDefultOptions() Options {
	return NewOptions(
		true, // md5
		true, // sha1
		true, // sha256
		true, // sha512
		true, // entropy
		true, // filetype
		0,    // cache size
	)
}

func NewChecksumOptions() Options {
	return NewOptions(
		true,  // md5
		true,  // sha1
		true,  // sha256
		true,  // sha512
		false, // entropy
		false, // filetype
		0,     // cache size
	)
}

// NewOptions creates a new Options struct with the given options
func NewOptions(md5, sha1, sha256, sha512, entropy, filetype bool, cacheSize int64) Options {
	return Options{
		Md5:       md5,
		Sha1:      sha1,
		Sha256:    sha256,
		Sha512:    sha512,
		Entropy:   entropy,
		Filetype:  filetype,
		CacheSize: cacheSize,
	}
}

// UpdateMd5 updates the md5 option of the Options struct
func (o Options) UpdateMd5(md5 bool) Options {
	o.Md5 = md5
	return o
}

// UpdateSha1 updates the sha1 option of the Options struct
func (o Options) UpdateSha1(sha1 bool) Options {
	o.Sha1 = sha1
	return o
}

// UpdateSha256 updates the sha256 option of the Options struct
func (o Options) UpdateSha256(sha256 bool) Options {
	o.Sha256 = sha256
	return o
}

// UpdateSha512 updates the sha512 option of the Options struct
func (o Options) UpdateSha512(sha512 bool) Options {
	o.Sha512 = sha512
	return o
}

// UpdateEntropy updates the entropy option of the Options struct
func (o Options) UpdateEntropy(entropy bool) Options {
	o.Entropy = entropy
	return o
}

// UpdateFiletype updates the filetype option of the Options struct
func (o Options) UpdateFiletype(filetype bool) Options {
	o.Filetype = filetype
	return o
}

// UpdateCacheSize updates the cache size of the Options struct
func (o Options) UpdateCacheSize(size int64) Options {
	o.CacheSize = size
	return o
}

// NewWriter creates a new Writer with the given options
func (o Options) NewWriter(w ...io.Writer) *Writer {
	return newWriterWithOptions(o, w...)
}

func (o Options) minCachSize() int64 {
	if !o.Filetype {
		return o.CacheSize
	}
	ftCachSize := int64(filetype.MaxBytesFileDetect())
	if ftCachSize > o.CacheSize {
		return ftCachSize
	}

	return o.CacheSize
}
