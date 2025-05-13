package cache

import (
	"testing"
)

func TestWriter(t *testing.T) {
	toWrite1 := []byte("Something cool")
	toWrite2 := []byte("Something else cool")

	t.Run("writes and resets", func(t *testing.T) {
		w := NewWriter(100)
		w.Write(toWrite1)
		if !w.IsCached() {
			t.Errorf("expected cached but its not")
		}
		if w.size != int64(len(toWrite1)) {
			t.Errorf("expected %d, got %d", len(toWrite1), w.size)
		}
		if string(w.Bytes()) != string(toWrite1) {
			t.Errorf("expected %s, got %s", string(toWrite1), string(w.Bytes()))
		}

		w.Reset()
		if w.size != 0 {
			t.Errorf("expected 0, got %d", w.size)
		}

		w.Write(toWrite2)
		if !w.IsCached() {
			t.Errorf("expected cached but its not")
		}
		if w.size != int64(len(toWrite2)) {
			t.Errorf("expected %d, got %d", len(toWrite2), w.size)
		}
		if string(w.Bytes()) != string(toWrite2) {
			t.Errorf("expected %s, got %s", string(toWrite2), string(w.Bytes()))
		}
	})

	t.Run("writes and no-reset", func(t *testing.T) {
		w := NewWriter(100)
		w.Write(toWrite1)
		if !w.IsCached() {
			t.Errorf("expected cached but its not")
		}
		if w.size != int64(len(toWrite1)) {
			t.Errorf("expected %d, got %d", len(toWrite1), w.size)
		}
		if string(w.Bytes()) != string(toWrite1) {
			t.Errorf("expected %s, got %s", string(toWrite1), string(w.Bytes()))
		}

		w.Write(toWrite2)
		if !w.IsCached() {
			t.Errorf("expected cached but its not")
		}
		if w.size != int64(len(toWrite1)+len(toWrite2)) {
			t.Errorf("expected %d, got %d", len(toWrite1)+len(toWrite2), w.size)
		}
		if string(w.Bytes()) != string(toWrite1)+string(toWrite2) {
			t.Errorf("expected %s, got %s", string(toWrite1)+string(toWrite2), string(w.Bytes()))
		}
	})

	t.Run("limits", func(t *testing.T) {
		w := NewWriter(2)
		w.Write(toWrite1)
		if w.IsCached() {
			t.Errorf("expected not cached but it is")
		}
		// Should still have full size
		if w.size != int64(len(toWrite1)) {
			t.Errorf("expected %d, got %d", len(toWrite1), w.size)
		}
		if string(w.Bytes()) != string(toWrite1[0:2]) {
			t.Errorf("expected %s, got %s", string(toWrite1[0:2]), string(w.Bytes()))
		}
	})
}
