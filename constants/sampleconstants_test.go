package constants

/*
A Sample test harness.
*/

import (
	"testing"
)

// A testing function.
func TestConstant(t *testing.T) {

	expected := 10
	if LoopConstant != expected {
		t.Errorf("const expected '%d' but got '%d'", expected, LoopConstant)
	}
}
