package bounded

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"sync"
	"testing"
)

const fileSize = 40_000

// func TestMain(m *testing.M) {
// 	// Setup code here
// 	setup()
// 	// Run the tests
// 	exitCode := m.Run()
// 	// Teardown code here
// 	teardown()
// 	// Exit with the test's exit code
// 	os.Exit(exitCode)
// }

func openFile(f int64) (*os.File, *File, error) {
	file, err := os.Open(fmt.Sprintf("../files/file%d.data", f))
	if err != nil {
		return nil, nil, fmt.Errorf("error opening file %v", err)
	}

	bf, err := NewFile("../files/file.data").Bound(f*fileSize, fileSize)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating bounded file %v", err)
	}

	return file, bf, nil
}

func assertSeek(t *testing.T, exp, act io.Seeker, offset int64, whence int) {
	p1, err1 := exp.Seek(offset, whence)
	p2, err2 := act.Seek(offset, whence)

	if err1 != err2 {
		t.Errorf("For %d - %d expected error %v, got %v", offset, whence, err1, err2)
	}

	if p1 != p2 {
		t.Errorf("For %d - %d expected position %d, got %d", offset, whence, p1, p2)
	}
}

func assertPosition(t *testing.T, exp, act io.Seeker, pos int64) {
	p1, err1 := exp.Seek(0, io.SeekCurrent)
	if err1 != nil {
		t.Errorf("For position act %d expected no error %v", pos, err1)
	}
	if p1 != pos {
		t.Errorf("For position exp expected position %d, got %d", pos, p1)
	}

	p2, err2 := act.Seek(0, io.SeekCurrent)
	if err2 != nil {
		t.Errorf("For position act %d expected no error %v", pos, err2)
	}
	if p2 != pos {
		t.Errorf("For position act expected position %d, got %d", pos, p1)
	}
}

func assertRead(t *testing.T, exp, act io.Reader, size int64) {
	b1 := make([]byte, size)
	b2 := make([]byte, size)
	p1, err1 := exp.Read(b1)
	p2, err2 := act.Read(b2)

	if err1 != err2 {
		t.Errorf("For %d expected error %v, got %v", size, err1, err2)
	}

	if p1 != p2 {
		t.Errorf("For %d expected bytes %d, got %d", size, p1, p2)
	}
}

func asserPathErr(t *testing.T, op string, exp error, act error) {
	if act == nil {
		t.Errorf("expected error %v %v, got nil", op, exp)
		return
	}

	actPE, ok := act.(*fs.PathError)
	if !ok {
		t.Errorf("expected path error %v, got %v", op, act)
		return
	}

	if actPE.Op != op {
		t.Errorf("expected op %s, got %s", op, actPE.Op)
	}
	if exp == nil {
		if actPE.Err == nil {
			t.Errorf("expected error got nil")
		}
	} else {
		if actPE.Err != exp {
			t.Errorf("expected error %v, got %v", exp, actPE.Err)
		}
	}
}

func assertSeekErr(t *testing.T, exp, act io.Seeker, offset int64, whence int) {
	p1, err1 := exp.Seek(offset, whence)
	// asserPathErr(t, "seek", fs.ErrInvalid, err1)
	asserPathErr(t, "seek", nil, err1)

	p2, err2 := act.Seek(offset, whence)
	asserPathErr(t, "seek", fs.ErrInvalid, err2)

	if p1 != p2 {
		t.Errorf("For %d - %d expected position %d, got %d", offset, whence, p1, p2)
	}
}

// func assertNoErrorf(t *testing.T, v string, err error) {
// 	if err != nil {
// 		t.Errorf(v, err)
// 	}
// }

func TestCompareFileAndFileBound(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("File-%d", i), func(t *testing.T) {
			file, bf, err := openFile(int64(i))
			if err != nil {
				t.Errorf("error opening file %v", err)
			}

			defer file.Close()
			defer bf.Cleanup()

			for i := 0; i < 5; i++ {
				file.Seek(0, io.SeekStart)
				bf.Seek(0, io.SeekStart)

				hasherF := sha256.New()
				if _, err := io.Copy(hasherF, file); err != nil {
					t.Errorf("error copying file %v", err)
				}

				hasherBF := sha256.New()
				if _, err := io.Copy(hasherBF, bf); err != nil {
					t.Errorf("error copying bfile %v", err)
				}

				hf := hex.EncodeToString(hasherF.Sum(nil))
				hb := hex.EncodeToString(hasherBF.Sum(nil))
				if hf != hb {
					t.Errorf("File hash mismatch %d: %s != %s", i, hf, hb)
				}
			}
		})
	}
}

func TestSeekAndRead(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("File-%d", i), func(t *testing.T) {
			file, bf, err := openFile(int64(i))
			if err != nil {
				t.Errorf("error opening file %v", err)
			}

			defer file.Close()
			defer bf.Cleanup()

			assertSeek(t, file, bf, 0, io.SeekStart)
			assertPosition(t, file, bf, 0)

			assertRead(t, file, bf, 0)
			assertPosition(t, file, bf, 0)

			assertRead(t, file, bf, 10)
			assertPosition(t, file, bf, 10)

			assertSeek(t, file, bf, 5, io.SeekStart)
			assertPosition(t, file, bf, 5)

			assertRead(t, file, bf, 320)
			assertPosition(t, file, bf, 325)

			assertSeek(t, file, bf, 12, io.SeekCurrent)
			assertPosition(t, file, bf, 337)

			assertSeek(t, file, bf, 5, io.SeekEnd)
			assertPosition(t, file, bf, fileSize+5)

			assertRead(t, file, bf, 100)
			assertPosition(t, file, bf, fileSize+5)

			assertSeek(t, file, bf, -10, io.SeekEnd)
			assertPosition(t, file, bf, fileSize-10)

			assertRead(t, file, bf, 5)
			assertPosition(t, file, bf, fileSize-5)

			assertSeekErr(t, file, bf, -10, io.SeekStart)
			assertPosition(t, file, bf, fileSize-5)

			assertRead(t, file, bf, 7)
			assertPosition(t, file, bf, fileSize)

			assertSeek(t, file, bf, 5, io.SeekStart)
			assertPosition(t, file, bf, 5)

			assertRead(t, file, bf, 2)
			assertPosition(t, file, bf, 7)

			assertSeekErr(t, file, bf, -10, io.SeekCurrent)
			assertPosition(t, file, bf, 7)

			assertSeekErr(t, file, bf, -1*(fileSize+5), io.SeekEnd)
			assertPosition(t, file, bf, 7)
		})
	}
}

// Cant use same underlying file for multiple threads
func TestFileClose(t *testing.T) {
	file, bf, err := openFile(2)
	if err != nil {
		t.Errorf("error opening file %v", err)
	}
	defer bf.Cleanup()

	t.Run("file", func(t *testing.T) {
		err = file.Close()
		if err != nil {
			t.Errorf("error closing file %v", err)
		}

		_, err = file.Seek(0, io.SeekStart)
		asserPathErr(t, "seek", fs.ErrClosed, err)

		_, err = file.Read(make([]byte, 5))
		asserPathErr(t, "read", fs.ErrClosed, err)

		err = file.Close()
		asserPathErr(t, "close", fs.ErrClosed, err)
	})

	// test close files again
	t.Run("bounded file", func(t *testing.T) {
		err = bf.Close()
		if err != nil {
			t.Errorf("error closing bounded file %v", err)
		}

		_, err = bf.Seek(0, io.SeekStart)
		asserPathErr(t, "seek", fs.ErrClosed, err)

		_, err = bf.Read(make([]byte, 5))
		asserPathErr(t, "read", fs.ErrClosed, err)

		err = bf.Close()
		asserPathErr(t, "close", fs.ErrClosed, err)
	})
}

// Cant use same underlying file for multiple threads
func TestFileThreading(t *testing.T) {
	file, err := os.Open("../files/file.data")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	_, err = file.Seek(5, io.SeekStart)
	if err != nil {
		t.Error(err)
	}

	pos, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		t.Error(err)
	}
	if pos != 5 {
		t.Errorf("Expected position 5, got %d", pos)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		pos, err := file.Seek(10, io.SeekStart)
		if err != nil {
			t.Error(err)
		}
		if pos != 10 {
			t.Errorf("Expected position 10, got %d", pos)
		}
	}()

	wg.Wait()

	pos, err = file.Seek(0, io.SeekCurrent)
	if err != nil {
		t.Fatal(err)
	}
	if pos != 10 {
		t.Errorf("Expected position 10 still, got %d", pos)
	}
}
