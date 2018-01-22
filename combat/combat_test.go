package combat

import "testing"

// This is literally just to have a test file for the data structures
func TestFalseIsFalse(t *testing.T) {
	if false {
		t.Errorf("wat\nSomehow, false is true?")
	}
}
