package filetype

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io"
	"testing"

	"github.com/gabriel-vasile/mimetype"
	"github.com/jonathongardner/fifo/cache"
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

func zlibCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := zlib.NewWriter(&buf)
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

func asssertFiletype(t *testing.T, w *cache.Writer, mime, ext string) {
	tp := NewFiletypeFromCached(w)
	if tp.Extension != ext {
		t.Errorf("expected %s extension, got %s", ext, tp.Extension)
	}
	if tp.Mimetype != mime {
		t.Errorf("expected %s mime type, got %s", mime, tp.Mimetype)
	}
}

func TestFiletypeWriter(t *testing.T) {
	toWrite := []byte("Something cool")

	gzV, err := gzipCompress(toWrite)
	if err != nil {
		t.Fatalf("failed to gz compress data %v", err)
	}

	zlibV, err := zlibCompress(toWrite)
	if err != nil {
		t.Fatalf("failed to xz compress data %v", err)
	}

	t.Run("reset", func(t *testing.T) {
		w := cache.NewWriter(int64(MaxBytesFileDetect()))
		if _, err := io.Copy(w, bytes.NewReader(gzV)); err != nil {
			t.Fatalf("failed to copy gz data %v", err)
		}
		asssertFiletype(t, w, "application/gzip", ".gz")

		w.Reset()

		if _, err := io.Copy(w, bytes.NewReader(zlibV)); err != nil {
			t.Fatalf("failed to copy xz data %v", err)
		}
		asssertFiletype(t, w, "application/octet-stream", "")
	})

	t.Run("no-reset", func(t *testing.T) {
		w := cache.NewWriter(int64(MaxBytesFileDetect()))
		if _, err := io.Copy(w, bytes.NewReader(gzV)); err != nil {
			t.Fatalf("failed to copy gz data %v", err)
		}
		asssertFiletype(t, w, "application/gzip", ".gz")

		if _, err := io.Copy(w, bytes.NewReader(zlibV)); err != nil {
			t.Fatalf("failed to copy xz data %v", err)
		}
		asssertFiletype(t, w, "application/gzip", ".gz")
	})
}

func TestFiletypeWriterCustom(t *testing.T) {
	sig := []byte{0x03, 0x32, 0x45, 0x67, 0x89}
	detect := func(raw []byte, limit uint32) bool {
		return bytes.Equal(raw[0:5], sig)
	}
	mimetype.Extend(detect, "application/foo", ".foo")

	t.Run("finds", func(t *testing.T) {
		w := cache.NewWriter(int64(MaxBytesFileDetect()))
		if _, err := w.Write(sig); err != nil {
			t.Fatalf("failed to write sig data %v", err)
		}
		if _, err := w.Write([]byte("1234567890")); err != nil {
			t.Fatalf("failed to write other data %v", err)
		}
		asssertFiletype(t, w, "application/foo", ".foo")
	})

	t.Run("no-find", func(t *testing.T) {
		w := cache.NewWriter(int64(MaxBytesFileDetect()))
		if _, err := w.Write([]byte("1234567890")); err != nil {
			t.Fatalf("failed to write other data %v", err)
		}
		if _, err := w.Write(sig); err != nil {
			t.Fatalf("failed to write sig data %v", err)
		}
		asssertFiletype(t, w, "application/octet-stream", "")
	})
}
