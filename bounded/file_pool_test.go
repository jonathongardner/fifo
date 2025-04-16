package bounded

import (
	"testing"
)

func TestFilePool(t *testing.T) {
	sbf, err := NewFilePool("../files/file.data", 2)
	if err != nil {
		t.Errorf("error creating shared bounded file %v", err)
	}
	defer sbf.Cleanup()

	bf1, err := sbf.Get(0, fileSize)
	if err != nil {
		t.Errorf("error getting bounded file 1 %v", err)
	}

	_, err = sbf.Get(1*fileSize, fileSize)
	if err != nil {
		t.Errorf("error getting bounded file 2 %v", err)
	}

	bf3, err := sbf.Get(2*fileSize, fileSize)
	if err != ErrNoAvailableFile {
		t.Errorf("expected no avialable bounded files but got %v", err)
	}
	if bf3 != nil {
		t.Errorf("expected no bounded file but got %v", bf3)
	}

	bf1.Close()

	_, err = sbf.Get(2*fileSize, fileSize)
	if err != nil {
		t.Errorf("error getting bounded file 2 %v", err)
	}

}
