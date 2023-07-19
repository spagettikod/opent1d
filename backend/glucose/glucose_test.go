package glucose

import "testing"

func TestMmolToMg(t *testing.T) {
	type TestCase struct {
		value    float32
		expected int
	}

	tests := []TestCase{
		{10, 180},
		{3.9, 70},
	}

	for _, test := range tests {
		actual := MmolToMg(test.value)
		if actual != test.expected {
			t.Errorf("expected %v but got %v", test.expected, actual)
		}
	}
}

func TestMgToMmol(t *testing.T) {
	type TestCase struct {
		value    int
		expected float32
	}

	tests := []TestCase{
		{180, 10},
		{70, 3.9},
	}

	for _, test := range tests {
		actual := MgToMmol(test.value)
		if actual != test.expected {
			t.Errorf("expected %v but got %v", test.expected, actual)
		}
	}
}
