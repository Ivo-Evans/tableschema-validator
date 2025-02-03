package util

import "testing"

func TestCommaSeparatedList(t *testing.T) {
	actual := CommaSeparatedList([]int{1, 2, 3})
	expected := "1, 2, 3"
	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}

	actual = CommaSeparatedList([]string{"a", "b", "c"})
	expected = "a, b, c"
	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}
