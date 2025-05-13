package filetype

import (
	"io"

	"github.com/gabriel-vasile/mimetype"
	"github.com/jonathongardner/fifo/cache"
	// log "github.com/sirupsen/logrus"
)

// https://github.com/gabriel-vasile/mimetype/blob/master/mimetype.go#L17
var maxBytesFileDetect uint32 = 3072

// SetMaxBytesFileDetect sets the maximum number of bytes to read for file type detection.
func SetMaxBytesFileDetect(maxBytes uint32) {
	if maxBytes > 0 {
		maxBytesFileDetect = maxBytes
		mimetype.SetLimit(uint32(maxBytesFileDetect))
	}
}

func MaxBytesFileDetect() uint32 {
	return maxBytesFileDetect
}

func init() {
	mimetype.SetLimit(uint32(maxBytesFileDetect))
}

// Filetype represents a file type with its extension and MIME type.
type Filetype struct {
	Extension string `json:"extension"`
	Mimetype  string `json:"mimetype"`
}

// Dir is a predefined Filetype for directories.
var Dir = Filetype{Extension: "dir", Mimetype: "directory/directory"}

// Symlink is a predefined Filetype for symbolic links.
var Symlink = Filetype{Extension: "symlink", Mimetype: "symlink/symlink"}

// newFiletype creates a new Filetype instance from a mimetype.MIME object.
func newFiletype(mtype *mimetype.MIME) Filetype {
	return Filetype{Extension: mtype.Extension(), Mimetype: mtype.String()}
}

// NewFiletypeFromCached creates a new Filetype instance from a cached writer
func NewFiletypeFromCached(wr *cache.Writer) Filetype {
	return newFiletype(mimetype.Detect(wr.Bytes()))
}

// NewFiletypeFromPath creates a new Filetype instance from a reader
// it reads maxBytesFileDetect of the reader
func NewFiletypeFromReader(reader io.Reader) (Filetype, error) {
	data := make([]byte, maxBytesFileDetect)
	w, err := reader.Read(data)
	if err != nil {
		return Filetype{}, err
	}
	// If reader doesnt have that much data truncate
	if uint32(w) < maxBytesFileDetect {
		data = data[:w]
	}
	return newFiletype(mimetype.Detect(data)), nil
}

// FiletypeFromJson creates a Filetype instance from a JSON representation.
func FiletypeFromJson(v map[string]any) Filetype {
	return Filetype{Extension: v["extension"].(string), Mimetype: v["mimetype"].(string)}
}
