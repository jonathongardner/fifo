package buffer

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestFiletypeWriter(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "tmp")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)
	toWrite1 := []byte("Something cool")

	t.Run("writes no data", func(t *testing.T) {
		file := filepath.Join(tmpDir, "test-no-data")

		assertFileDoesNotExist(t, file)

		w, err := NewFileWriter(file, len(toWrite1)+1)
		if err != nil {
			t.Fatalf("Failed to create new writer %v", err)
		}
		assertFileDoesNotExist(t, file)

		if err := w.Close(); err != nil {
			t.Fatalf("Failed to close %v", err)
		}
		assertFileExists(t, file, []byte{})
	})

	t.Run("writes small data", func(t *testing.T) {
		file := filepath.Join(tmpDir, "test1")

		assertFileDoesNotExist(t, file)

		w, err := NewFileWriter(file, len(toWrite1)+1)
		if err != nil {
			t.Fatalf("Failed to create new writer %v", err)
		}
		assertFileDoesNotExist(t, file)

		if _, err := w.Write(toWrite1); err != nil {
			t.Fatalf("Failed to write %v", err)
		}
		assertFileDoesNotExist(t, file)

		if err := w.Close(); err != nil {
			t.Fatalf("Failed to close %v", err)
		}
		assertFileExists(t, file, toWrite1)
	})

	t.Run("writes small data delete", func(t *testing.T) {
		file := filepath.Join(tmpDir, "test2")

		assertFileDoesNotExist(t, file)

		w, err := NewFileWriter(file, len(toWrite1)+1)
		if err != nil {
			t.Fatalf("Failed to create new writer %v", err)
		}
		assertFileDoesNotExist(t, file)

		if _, err := w.Write(toWrite1); err != nil {
			t.Fatalf("Failed to write %v", err)
		}
		assertFileDoesNotExist(t, file)

		if err := w.Delete(); err != nil {
			t.Fatalf("Failed to delete %v", err)
		}
		assertFileDoesNotExist(t, file)

		// should not write on close if deleted
		if err := w.Close(); err != nil {
			t.Fatalf("Failed to close %v", err)
		}
		assertFileDoesNotExist(t, file)
	})

	t.Run("writes medium data", func(t *testing.T) {
		file := filepath.Join(tmpDir, "test3")

		assertFileDoesNotExist(t, file)

		w, err := NewFileWriter(file, len(toWrite1)+1)
		if err != nil {
			t.Fatalf("Failed to create new writer %v", err)
		}
		assertFileDoesNotExist(t, file)

		if _, err := w.Write(toWrite1); err != nil {
			t.Fatalf("Failed to write first %v", err)
		}
		assertFileDoesNotExist(t, file)

		exp := append(toWrite1, toWrite1...)
		if _, err := w.Write(toWrite1); err != nil {
			t.Fatalf("Failed to write second %v", err)
		}
		assertFileExists(t, file, exp)

		if err := w.Close(); err != nil {
			t.Fatalf("Failed to close %v", err)
		}
		assertFileExists(t, file, exp)
	})

	t.Run("writes medium data with delete", func(t *testing.T) {
		file := filepath.Join(tmpDir, "test4")

		assertFileDoesNotExist(t, file)

		w, err := NewFileWriter(file, len(toWrite1)+1)
		if err != nil {
			t.Fatalf("Failed to create new writer %v", err)
		}
		assertFileDoesNotExist(t, file)

		if _, err := w.Write(toWrite1); err != nil {
			t.Fatalf("Failed to write first %v", err)
		}
		assertFileDoesNotExist(t, file)

		if _, err := w.Write(toWrite1); err != nil {
			t.Fatalf("Failed to write second %v", err)
		}
		assertFileExists(t, file, append(toWrite1, toWrite1...))

		if err := w.Delete(); err != nil {
			t.Fatalf("Failed to delete %v", err)
		}
		assertFileDoesNotExist(t, file)

		if err := w.Close(); err != nil {
			t.Fatalf("Failed to close %v", err)
		}
		assertFileDoesNotExist(t, file)
	})
}

func assertFileExists(t *testing.T, filename string, data []byte) {
	t.Helper()
	b, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Expect %s to return no error, but returned %v", filename, err)
	}
	if !bytes.Equal(b, data) {
		t.Fatalf("Expect %s to return %v, but returned %v", filename, data, b)
	}
}

func assertFileDoesNotExist(t *testing.T, filename string) {
	t.Helper()
	_, err := os.Stat(filename)
	if !os.IsNotExist(err) {
		t.Fatalf("Expect %s to return file not found, but returned %v", filename, err)
	}
}
