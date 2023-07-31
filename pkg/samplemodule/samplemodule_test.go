package samplemodule

/*
A Sample test harness.
*/

import (
	"testing"
)

// A testing function.
func TestSampleFunction(t *testing.T) {

	msg := SampleFunction()
	expected := "OK"

	if msg != expected {
		t.Errorf("Module expected '%q' but got '%q'", expected, msg)
	}
}
