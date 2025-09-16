package util

import (
	"slices"
	"testing"
)

func Test_RemoveFromSlice(t *testing.T) {
	initial := []string{"Alice", "Bob", "Charlie", "David"}
	expected := []string{"Alice", "Bob", "David"}
	actual := Remove(initial, "Charlie")

	if !slices.Equal(expected, actual) {
		t.Errorf("expected and actual slices do not match")
	}
}

func Test_RemoveFromSliceNotFound(t *testing.T) {
	initial := []string{"Alice", "Bob", "Charlie", "David"}
	expected := []string{"Alice", "Bob", "Charlie", "David"}
	actual := Remove(initial, "Elaine")

	if !slices.Equal(expected, actual) {
		t.Errorf("expected and actual slices do not match")
	}
}
