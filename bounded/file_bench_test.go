package bounded

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"
)

func copyFile() {
	srcFile, err := os.Open("../files/file.data")
	if err != nil {
		fmt.Printf("error opening file %v", err)
	}
	defer srcFile.Close()
	file, err := os.Create("../files/file-to-write.data")
	if err != nil {
		fmt.Println("Error creating/truncating file:", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, srcFile)
	if err != nil {
		fmt.Printf("copying to temp file %v", err)
	}
}

func writeToFile(path string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("error opening file %v", err)
	}
	defer f.Close()

	file, err := os.Open("../files/to-copy.data")
	if err != nil {
		return fmt.Errorf("error opening to copy %v", err)
	}
	defer file.Close()

	if _, err = io.Copy(f, file); err != nil {
		return fmt.Errorf("error copying %v", err)
	}
	return nil
}

func multipleChecksums(f io.ReadSeeker, n int) error {
	for i := 0; i < n; i++ {
		hasher := sha256.New()
		if _, err := io.Copy(hasher, f); err != nil {
			return err
		}
		f.Seek(0, io.SeekStart)
	}
	return nil
}

func BenchmarkFile(b *testing.B) {
	for j := 0; j < b.N; j++ {
		for i := 0; i < 10; i++ {
			file, err := os.Open(fmt.Sprintf("../files/file%d.data", i))
			if err != nil {
				fmt.Printf("Error opening file: %v\n", err)
			}
			defer file.Close()

			err = multipleChecksums(file, b.N)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
	}
}

func BenchmarkBoundedFile(b *testing.B) {
	for j := 0; j < b.N; j++ {
		b.StopTimer()
		copyFile()
		b.StartTimer()
		nbf := NewFile("../files/file-to-write.data")
		defer nbf.Cleanup()

		for i := 0; i < 10; i++ {
			file, err := nbf.Bound(int64(i)*fileSize, fileSize)
			if err != nil {
				fmt.Printf("Error opening bounded: %v\n", err)
			}
			defer file.Close()

			err = multipleChecksums(file, b.N)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
	}
}

func BenchmarkFileGoRoutine(b *testing.B) {
	for j := 0; j < b.N; j++ {
		wg := &sync.WaitGroup{}

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				file, err := os.Open(fmt.Sprintf("../files/file%d.data", i))
				if err != nil {
					fmt.Printf("Error opening file: %v\n", err)
				}
				defer file.Close()

				err = multipleChecksums(file, b.N)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkBoundedFileGoRoutine(b *testing.B) {
	for j := 0; j < b.N; j++ {
		b.StopTimer()
		copyFile()
		b.StartTimer()
		wg := &sync.WaitGroup{}

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				file, err := NewFile("../files/file-to-write.data").Bound(int64(i)*fileSize, fileSize)
				if err != nil {
					fmt.Printf("Error opening bounded: %v\n", err)
				}
				defer file.Cleanup()

				err = multipleChecksums(file, b.N)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkFileWithWrite(b *testing.B) {
	for j := 0; j < b.N; j++ {
		wg := &sync.WaitGroup{}

		wg.Add(1)
		go func() {
			err := writeToFile("../files/new.data")
			os.Remove("../files/new.data")
			if err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
			}
			wg.Done()
		}()

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				file, err := os.Open(fmt.Sprintf("../files/file%d.data", i))
				if err != nil {
					fmt.Printf("Error opening file: %v\n", err)
				}
				defer file.Close()

				err = multipleChecksums(file, b.N)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkBoundedFileWithWrite(b *testing.B) {
	for j := 0; j < b.N; j++ {
		b.StopTimer()
		copyFile()
		b.StartTimer()
		wg := &sync.WaitGroup{}

		wg.Add(1)
		go func() {
			err := writeToFile("../files/file-to-write.data")
			if err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
			}
			wg.Done()
		}()

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				file, err := NewFile("../files/file-to-write.data").Bound(int64(i)*fileSize, fileSize)
				if err != nil {
					fmt.Printf("Error opening bounded: %v\n", err)
				}
				defer file.Cleanup()

				err = multipleChecksums(file, b.N)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}
