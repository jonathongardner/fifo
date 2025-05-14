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
	if _, err := w.Identifiers(); err != ErrWriterNotClosed {
		t.Fatalf("expected ErrWriterNotClosed, got %v", err)
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

	i, err := w.Identifiers()
	if err != nil {
		t.Fatalf("failed to get identifiers %v", err)
	}
	assertIdentifiers(t, exp, i)
}

func TestResetWriter(t *testing.T) {
	toWrite := []byte("Something cool")

	w := NewWriter()
	if _, err := io.Copy(w, bytes.NewReader(toWrite)); err != nil {
		t.Fatalf("failed to copy data %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("failed to close writer %v", err)
	}

	exp1 := Identifiers{
		Md5:      "db5ee56e2cab72f4e46bdd60965bef31",
		Sha1:     "77e9c3f828f8773ba2854731eddff8d5daef99c0",
		Sha256:   "2f693b6d594eccd0cc94682cc5690324d369186d26097f8eb611231dfd8fd8ab",
		Sha512:   "dc17c6f834c2ae2af8f25d2f275e2b7ba7280c22da1cf96475888b038ad96872e771f43b87e722507773c54973a0794f7b67ee07fcbc50bf28e3b369a4801819",
		Entropy:  3.46,
		Filetype: filetype.Filetype{Mimetype: "text/plain; charset=utf-8"},
		Size:     14,
	}

	t.Run("initial identifiers", func(t *testing.T) {
		i, err := w.Identifiers()
		if err != nil {
			t.Fatalf("failed to get identifiers %v", err)
		}
		assertIdentifiers(t, exp1, i)
	})

	w.Reset()
	if _, err := io.Copy(w, bytes.NewReader(toWrite)); err != nil {
		t.Fatalf("failed to copy data %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("failed to close writer %v", err)
	}

	t.Run("after reset", func(t *testing.T) {
		i, err := w.Identifiers()
		if err != nil {
			t.Fatalf("failed to get identifiers %v", err)
		}
		assertIdentifiers(t, exp1, i)
	})

	// Cant happen but for testing lets do it
	w.closed = false
	if _, err := io.Copy(w, bytes.NewReader(toWrite)); err != nil {
		t.Fatalf("failed to copy data %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("failed to close writer %v", err)
	}

	t.Run("no reset identifiers", func(t *testing.T) {
		exp2 := Identifiers{
			Md5:      "55449e6ae68e8bd9c902f1de73141997",
			Sha1:     "0f12b887956f1b36ab9bea834ad253c16b075100",
			Sha256:   "0c70fc53f9482bead542c4c8b94af93001e0192329357446c9d79d3c5d430e23",
			Sha512:   "19a47833ba8bb971f258b942b10bd521c5631a63a5b3b104d2287acc0466f07432ee9146fe7381f14f5a8587920627494f78e71e7b7d2b79121bbd16ae4d417c",
			Entropy:  3.46,
			Filetype: filetype.Filetype{Mimetype: "text/plain; charset=utf-8"},
			Size:     28,
		}
		if exp2.Sha512 == exp1.Sha512 {
			t.Errorf("sha512 should be different %v", exp2.Sha512)
		}

		i, err := w.Identifiers()
		if err != nil {
			t.Fatalf("failed to get identifiers %v", err)
		}
		assertIdentifiers(t, exp2, i)
	})
}

func TestAddWriter(t *testing.T) {
	toWrite := []byte("Something cool")

	var buf1 bytes.Buffer
	w := NewWriter(&buf1)
	var buf2 bytes.Buffer
	w.AddWriter(&buf2)

	if _, err := w.Write(toWrite); err != nil {
		t.Fatalf("failed to write data %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("failed to close writer %v", err)
	}

	if !bytes.Equal(buf1.Bytes(), toWrite) {
		t.Errorf("data mismatch, expected %v, got %v", toWrite, buf1.Bytes())
	}

	if !bytes.Equal(buf1.Bytes(), toWrite) {
		t.Errorf("data mismatch, expected %v, got %v", toWrite, buf1.Bytes())
	}

	exp := Identifiers{
		Md5:      "db5ee56e2cab72f4e46bdd60965bef31",
		Sha1:     "77e9c3f828f8773ba2854731eddff8d5daef99c0",
		Sha256:   "2f693b6d594eccd0cc94682cc5690324d369186d26097f8eb611231dfd8fd8ab",
		Sha512:   "dc17c6f834c2ae2af8f25d2f275e2b7ba7280c22da1cf96475888b038ad96872e771f43b87e722507773c54973a0794f7b67ee07fcbc50bf28e3b369a4801819",
		Entropy:  3.46,
		Filetype: filetype.Filetype{Mimetype: "text/plain; charset=utf-8"},
		Size:     14,
	}

	i, err := w.Identifiers()
	if err != nil {
		t.Fatalf("failed to get identifiers %v", err)
	}
	assertIdentifiers(t, exp, i)
}
