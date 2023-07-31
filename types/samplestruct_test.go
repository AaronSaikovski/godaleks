/*
A Sample test harness.
*/

package types

import (
	"testing"
)

// A testing function.
func TestSampleStructString(t *testing.T) {

	expected := "test data"
	ateststruct := Sample{SampleString: "test data", SampleInt: 1}

	if ateststruct.SampleString != expected {
		t.Errorf("struct expected '%s' but got '%s'", expected, ateststruct.SampleString)
	}
}
