package buffer

import (
	"fmt"
	"os"
)

// BufferedFileWriter is a file writer that buffers data in memory
// until it reaches a certain threshold before writing to disk.
type FileWriter struct {
	filePath string
	buffer   []byte
	max      int
	file     *os.File
}

// NewBufferedFileWriter creates a new BufferedFileWriter.
func NewFileWriter(filePath string, max int) (*FileWriter, error) {
	if max <= 0 {
		return nil, fmt.Errorf("max size must be greater than 0")
	}

	return &FileWriter{
		filePath: filePath,
		buffer:   make([]byte, 0),
		max:      max,
		file:     nil,
	}, nil
}

func (w *FileWriter) create() error {
	if w.file != nil {
		return nil
	}
	file, err := os.Create(w.filePath)
	if err != nil {
		return err
	}

	w.file = file
	return nil
}
func (w *FileWriter) checkErr() error {
	if w.max == -1 {
		return os.ErrClosed
	}
	if w.max == -2 {
		return fmt.Errorf("filed deleted")
	}
	return nil
}

// Write writes the bassed bytes to buffer, writes to file if exceeds buffer size
func (w *FileWriter) Write(data []byte) (int, error) {
	if err := w.checkErr(); err != nil {
		return 0, err
	}

	size := len(data)
	w.buffer = append(w.buffer, data...)

	if len(w.buffer) >= w.max {
		if err := w.flush(); err != nil {
			return 0, err
		}
	}
	return size, nil
}

func (w *FileWriter) flush() error {
	if err := w.create(); err != nil {
		return fmt.Errorf("failed to create file to write to: %w", err)
	}
	_, err := w.file.Write(w.buffer)
	if err != nil {
		return err
	}
	w.buffer = w.buffer[:0] // Clear buffer
	return nil
}

// Flush write to file if buffered data
func (w *FileWriter) Flush() error {
	if len(w.buffer) > 0 {
		return w.flush()
	}
	return nil
}

// Delete the file
func (w *FileWriter) Delete() error {
	if err := w.checkErr(); err != nil {
		return err
	}
	w.max = -2 // Mark as deleted
	if w.file != nil {
		if err := w.file.Close(); err != nil {
			return fmt.Errorf("failed to close file: %w", err)
		}
		if err := os.Remove(w.filePath); err != nil {
			return fmt.Errorf("failed to delete file: %w", err)
		}
		w.file = nil
	}
	return nil
}

func (w *FileWriter) Close() error {
	if w.max == -1 {
		return os.ErrClosed
	}
	// if its been deteled than just move on
	if w.max == -2 {
		return nil
	}

	if err := w.Flush(); err != nil {
		return err
	}
	if err := w.file.Close(); err != nil {
		return err
	}
	w.max = -1 // Mark as closed
	return nil
}
