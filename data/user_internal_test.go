package data

import (
	"testing"
)

func TestUserExistsIn(t *testing.T) {
	testSlice := []string{"1", "2", "3"}
	el := "1"

	if !elementExists(testSlice, el) {
		t.Fatalf("element %v not found", el)
	}

	el = "4"

	if elementExists(testSlice, el) {
		t.Fatalf("element %v should not be found", el)
	}
}
