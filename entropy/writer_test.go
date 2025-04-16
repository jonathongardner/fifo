package entropy

import (
	"testing"
)

type entTest struct {
	value   string
	entropy float64
}

func (at entTest) equal(a float64) bool {
	return int(at.entropy*100) == int(a*100)
}

func TestEntropyWriter(t *testing.T) {
	x := []byte{}
	for i := 0; i < 256; i++ {
		x = append(x, byte(i))
	}

	tests := []entTest{
		{"", 0},
		{"a", 0},
		{"coolo beans this is awesome... really awesome", 3.76},
		{"foo who boo this is nice", 3.25},
		{string(x), 8}, // binary
	}
	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			ent := NewWriter()
			ent.Write([]byte(test.value))
			if !test.equal(ent.Entropy()) {
				t.Errorf("expected %v, got %v", test.entropy, ent.Entropy())
			}
		})
	}
}
