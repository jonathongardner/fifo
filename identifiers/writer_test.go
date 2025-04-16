package stats

import (
	"bytes"
	"compress/gzip"
	"io"
	"testing"

	"github.com/jonathongardner/fifo/filetype"
)

func gzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func TestWriter(t *testing.T) {
	toWrite := []byte("Something cool")
	gzV, err := gzipCompress(toWrite)
	if err != nil {
		t.Fatalf("failed to compress data %v", err)
	}

	var buf bytes.Buffer
	w := NewWriter(&buf)
	if _, err := io.Copy(w, bytes.NewReader(gzV)); err != nil {
		t.Fatalf("failed to copy gz data %v", err)
	}

	if !bytes.Equal(buf.Bytes(), gzV) {
		t.Fatalf("data mismatch, expected %v, got %v", gzV, buf.Bytes())
	}

	exp := Identifiers{
		Md5:      "7122f4be9338ce4c875e994db73366bb",
		Sha1:     "d8d3397929f6f34c455bf7ef3bfbfb7dd1f5bd67",
		Sha256:   "89fff6aff4d0d43239531da8b10f920c37ee6efa3b7540db5050f9ed9b5d904a",
		Sha512:   "bf9fc66a432a42dceb2ce9e29f0f2f2f9700a1375305f9cc2f2f72ce7faae85568221cb0db941001aaf79f8bc4dbca53166630d81c6a92c377e9af76ceda52b0",
		Entropy:  3.94,
		Filetype: filetype.Filetype{Mimetype: "application/gzip"},
		Size:     int64(len(gzV)),
	}

	info := w.Identifiers()
	if info.Md5 != exp.Md5 {
		t.Fatalf("md5 mismatch, expected %v, got %v", exp.Md5, info.Md5)
	}
	if info.Sha1 != exp.Sha1 {
		t.Fatalf("sha1 mismatch, expected %v, got %v", exp.Sha1, info.Sha1)
	}
	if info.Sha256 != exp.Sha256 {
		t.Fatalf("sha256 mismatch, expected %v, got %v", exp.Sha256, info.Sha256)
	}
	if info.Sha512 != exp.Sha512 {
		t.Fatalf("sha512 mismatch, expected %v, got %v", exp.Sha512, info.Sha512)
	}
	if int(info.Entropy*100) != int(exp.Entropy*100) {
		t.Fatalf("entropy mismatch, expected %v, got %v", exp.Entropy, info.Entropy)
	}
	if info.Filetype.Mimetype != exp.Filetype.Mimetype {
		t.Fatalf("mimitype mismatch, expected %v, got %v", exp.Filetype.Mimetype, info.Filetype.Mimetype)
	}
	if info.Size != exp.Size {
		t.Fatalf("size mismatch, expected %v, got %v", exp.Size, info.Size)
	}
}
