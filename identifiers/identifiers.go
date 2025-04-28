package identifiers

import (
	"encoding/hex"
	"fmt"

	"github.com/jonathongardner/fifo/filetype"
)

var ErrWriterNotClosed = fmt.Errorf("writer is not closed")

type Identifiers struct {
	Md5      string            `json:"md5,omitempty"`
	Sha1     string            `json:"sha1,omitempty"`
	Sha256   string            `json:"sha256,omitempty"`
	Sha512   string            `json:"sha512,omitempty"`
	Entropy  float64           `json:"entropy,omitempty"`
	Filetype filetype.Filetype `json:"filetype,omitempty"`
	Size     int64             `json:"size,omitempty"`
}

// Identifiers returns the info of the writer that "identifies" the data
// the md5, sha1, sha256, sha512, entropy, file type, and size
// it returns an empty string if the hash is not calculated
// it returns 0 if the entropy is not calculated
// it returns nil if the file type is not calculated
func (mw *Writer) Identifiers() (Identifiers, error) {
	if !mw.closed {
		return Identifiers{}, ErrWriterNotClosed
	}

	toReturn := Identifiers{}
	if mw.md5 != nil {
		toReturn.Md5 = hex.EncodeToString(mw.md5.Sum(nil))
	}
	if mw.sha1 != nil {
		toReturn.Sha1 = hex.EncodeToString(mw.sha1.Sum(nil))
	}
	if mw.sha256 != nil {
		toReturn.Sha256 = hex.EncodeToString(mw.sha256.Sum(nil))
	}
	if mw.sha512 != nil {
		toReturn.Sha512 = hex.EncodeToString(mw.sha512.Sum(nil))
	}
	if mw.entropy != nil {
		toReturn.Entropy = mw.entropy.Entropy()
	}
	if mw.ftype {
		toReturn.Filetype = filetype.NewFiletypeFromCached(mw.cache)
	}
	toReturn.Size = mw.cache.Size()
	return toReturn, nil
}
