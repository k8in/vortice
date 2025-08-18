package util

import "testing"

// --- AI GENERATED CODE BEGIN ---

func TestInSlice_String(t *testing.T) {
	slice := []string{"a", "b", "c"}
	if !InSlice("a", slice) {
		t.Error("InSlice should return true for existing string")
	}
	if InSlice("d", slice) {
		t.Error("InSlice should return false for non-existing string")
	}
}

func TestInSlice_Int(t *testing.T) {
	slice := []int{1, 2, 3}
	if !InSlice(2, slice) {
		t.Error("InSlice should return true for existing int")
	}
	if InSlice(4, slice) {
		t.Error("InSlice should return false for non-existing int")
	}
}

// --- AI GENERATED CODE END ---
