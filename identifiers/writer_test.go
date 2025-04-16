package identifiers

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

func assertIdentifiers(t *testing.T, exp, act Identifiers) {
	if act.Md5 != exp.Md5 {
		t.Errorf("md5 mismatch, expected %v, got %v", exp.Md5, act.Md5)
	}
	if act.Sha1 != exp.Sha1 {
		t.Errorf("sha1 mismatch, expected %v, got %v", exp.Sha1, act.Sha1)
	}
	if act.Sha256 != exp.Sha256 {
		t.Errorf("sha256 mismatch, expected %v, got %v", exp.Sha256, act.Sha256)
	}
	if act.Sha512 != exp.Sha512 {
		t.Errorf("sha512 mismatch, expected %v, got %v", exp.Sha512, act.Sha512)
	}
	if int(act.Entropy*100) != int(exp.Entropy*100) {
		t.Errorf("entropy mismatch, expected %v, got %v", exp.Entropy, act.Entropy)
	}
	if act.Filetype.Mimetype != exp.Filetype.Mimetype {
		t.Errorf("mimitype mismatch, expected %v, got %v", exp.Filetype.Mimetype, act.Filetype.Mimetype)
	}
	if act.Size != exp.Size {
		t.Errorf("size mismatch, expected %v, got %v", exp.Size, act.Size)
	}
	if act.GzipSize != exp.GzipSize {
		t.Errorf("gzip size mismatch, expected %v, got %v", exp.GzipSize, act.GzipSize)
	}
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
	if err := w.Close(); err != nil {
		t.Fatalf("failed to close writer %v", err)
	}

	if !bytes.Equal(buf.Bytes(), gzV) {
		t.Errorf("data mismatch, expected %v, got %v", gzV, buf.Bytes())
	}

	exp := Identifiers{
		Md5:      "7122f4be9338ce4c875e994db73366bb",
		Sha1:     "d8d3397929f6f34c455bf7ef3bfbfb7dd1f5bd67",
		Sha256:   "89fff6aff4d0d43239531da8b10f920c37ee6efa3b7540db5050f9ed9b5d904a",
		Sha512:   "bf9fc66a432a42dceb2ce9e29f0f2f2f9700a1375305f9cc2f2f72ce7faae85568221cb0db941001aaf79f8bc4dbca53166630d81c6a92c377e9af76ceda52b0",
		Entropy:  3.94,
		Filetype: filetype.Filetype{Mimetype: "application/gzip"},
		Size:     38,
	}

	assertIdentifiers(t, exp, w.Identifiers())
}

func TestGzipWriter(t *testing.T) {
	toWrite := []byte("Something cool")

	var buf bytes.Buffer
	w := NewOptions(
		true,
		true,
		true,
		true,
		true,
		true,
		100,
		100,
	).NewWriter(&buf)
	if _, err := io.Copy(w, bytes.NewReader(toWrite)); err != nil {
		t.Fatalf("failed to copy gz data %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("failed to close writer %v", err)
	}

	if !bytes.Equal(buf.Bytes(), toWrite) {
		t.Fatalf("data mismatch, expected %v, got %v", toWrite, buf.Bytes())
	}

	exp := Identifiers{
		Md5:      "db5ee56e2cab72f4e46bdd60965bef31",
		Sha1:     "77e9c3f828f8773ba2854731eddff8d5daef99c0",
		Sha256:   "2f693b6d594eccd0cc94682cc5690324d369186d26097f8eb611231dfd8fd8ab",
		Sha512:   "dc17c6f834c2ae2af8f25d2f275e2b7ba7280c22da1cf96475888b038ad96872e771f43b87e722507773c54973a0794f7b67ee07fcbc50bf28e3b369a4801819",
		Entropy:  3.46,
		Filetype: filetype.Filetype{Mimetype: "text/plain; charset=utf-8"},
		Size:     14,
		GzipSize: 38,
	}

	assertIdentifiers(t, exp, w.Identifiers())
}
