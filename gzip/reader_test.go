package gzip

import (
	"bytes"
	"compress/gzip"
	"io"
	"testing"
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

func TestGzipReader(t *testing.T) {
	toWrite := []byte("Your String Here")
	compressed, err := gzipCompress(toWrite)
	if err != nil {
		t.Fatalf("failed to compress data %v", err)
	}
	r := bytes.NewReader(compressed)
	r.Seek(0, io.SeekStart)
	gzReader, err := NewReader(r)
	if err != nil {
		t.Fatalf("failed to create new gzip reader %v", err)
	}

	v1 := make([]byte, len(toWrite))
	n1, err := gzReader.Read(v1)
	if err != nil && err != io.EOF {
		t.Fatal("failed to read gzip", err)
	}

	if n1 != len(toWrite) {
		t.Fatalf("expected %d bytes but got %d", len(toWrite), n1)
	}
	if string(v1) != string(toWrite) {
		t.Fatalf("expected %s but got %s", string(toWrite), string(v1))
	}

	v2 := make([]byte, len(toWrite))
	n2, err := gzReader.Read(v2)
	if err != io.EOF {
		t.Fatal("expected EOF", err)
	}
	if n2 != 0 {
		t.Fatalf("expected 0 bytes but got %d", n2)
	}

	gzReader.Reset()

	v3 := make([]byte, len(toWrite))
	n3, err := gzReader.Read(v3)
	if err != nil && err != io.EOF {
		t.Fatal("failed to read gzip after reset", err)
	}

	if n3 != len(toWrite) {
		t.Fatalf("expected %d bytes but got %d after reset", len(toWrite), n3)
	}
	if string(v3) != string(toWrite) {
		t.Fatalf("expected %s but got %s after reset", string(toWrite), string(v3))
	}
}
