package cache

import (
	"os"
	"testing"
)

func TestFile(t *testing.T) {

	t.Run("reads from cache", func(t *testing.T) {
		f, err := Open("testdata/foo", 10)
		if err != nil {
			t.Errorf("expected nil error for open, got %v", err)
		}

		// read 10 bytes so smaller then the cache
		data1 := make([]byte, 10)
		n, err := f.Read(data1)
		if n != 10 {
			t.Errorf("expected 10 bytes read from file, got %d", n)
		}
		if err != nil {
			t.Errorf("expected nil error for read from file, got %v", err)
		}

		// check new reader
		c, err := f.NewReader()
		if err != nil {
			t.Errorf("expected nil error for new reader, got %v", err)
		}
		if !c.IsCached() {
			t.Errorf("expected cached but its not")
		}

		// make sure we get error if try to read again
		_, err = f.Read(data1)
		if err != os.ErrClosed {
			t.Errorf("expected closed error for read from file, got %v", err)
		}

		// read from cache
		data2 := make([]byte, 10)
		n, err = c.Read(data2)
		if n != 10 {
			t.Errorf("expected 10 bytes reading from cache, got %d", n)
		}
		if err != nil {
			t.Errorf("expected nil error for read from cache, got %v", err)
		}
		if string(data1) != string(data2) {
			t.Errorf("expected %s, got %s", string(data1), string(data2))
		}
	})

	t.Run("reads from file", func(t *testing.T) {
		f, err := Open("testdata/foo", 5)
		if err != nil {
			t.Errorf("expected nil error for open, got %v", err)
		}

		// read 10 bytes so smaller then the cache
		data1 := make([]byte, 10)
		n, err := f.Read(data1)
		if n != 10 {
			t.Errorf("expected 10 bytes read from file, got %d", n)
		}
		if err != nil {
			t.Errorf("expected nil error for read from file, got %v", err)
		}

		// check new reader
		c, err := f.NewReader()
		if err != nil {
			t.Errorf("expected nil error for new reader, got %v", err)
		}
		if c.IsCached() {
			t.Errorf("expected not cached but is")
		}

		// make sure we get error if try to read again
		_, err = f.Read(data1)
		if err != os.ErrClosed {
			t.Errorf("expected closed error for read from file, got %v", err)
		}

		// read from cache
		data2 := make([]byte, 10)
		n, err = c.Read(data2)
		if n != 10 {
			t.Errorf("expected 10 bytes reading from cache, got %d", n)
		}
		if err != nil {
			t.Errorf("expected nil error for read from cache, got %v", err)
		}
		if string(data1) != string(data2) {
			t.Errorf("expected %s, got %s", string(data1), string(data2))
		}
	})

	t.Run("new file stuff", func(t *testing.T) {
		f := NewFile(5)

		// read 10 bytes so smaller then the cache
		data1 := make([]byte, 10)
		n, err := f.Read(data1)
		if n != 0 {
			t.Errorf("expected 0 bytes read from file since nothing open, got %d", n)
		}
		if err != os.ErrClosed {
			t.Errorf("expected error closed for read from file, got %v", err)
		}

		err = f.Open("testdata/foo")
		if err != nil {
			t.Errorf("expected nil error for open, got %v", err)
		}

		err = f.Open("testdata/foo")
		if err != ErrAlreadyOpen {
			t.Errorf("expected already open error for open, got %v", err)
		}

		// read from cache
		data2 := make([]byte, 10)
		n, err = f.Read(data2)
		if n != 10 {
			t.Errorf("expected 10 bytes reading from cache, got %d", n)
		}
		if err != nil {
			t.Errorf("expected nil error for read from cache, got %v", err)
		}
	})
}
