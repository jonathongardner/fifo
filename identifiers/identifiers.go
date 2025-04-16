package stats

import (
	"encoding/hex"

	"github.com/jonathongardner/fifo/filetype"
)

type Identifiers struct {
	Md5      string
	Sha1     string
	Sha256   string
	Sha512   string
	Entropy  float64
	Filetype filetype.Filetype
	Size     int64
}

// Identifiers returns the info of the writer that "identifies" the data
// the md5, sha1, sha256, sha512, entropy, file type, and size
// it returns an empty string if the hash is not calculated
// it returns 0 if the entropy is not calculated
// it returns nil if the file type is not calculated
func (mw *Writer) Identifiers() Identifiers {
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
	toReturn.Size = mw.size
	return toReturn
}
