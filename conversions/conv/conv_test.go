package conv

import "testing"

func TestDecToBase(t *testing.T) {
	tests := []struct {
		num, base uint
		expected  string
	}{
		{17, 16, "11"},
		{16, 16, "10"},
		{10, 16, "A"},
		{16, 15, "11"},
		{16, 32, "G"},
		{35, 16, "23"},
		{20, 32, "K"},
	}

	for _, test := range tests {
		result, err := DecToBase(test.num, test.base)
		if err != nil {
			t.Errorf("%v", err)
		} else if result != test.expected {
			t.Errorf("(%v, %v) Expected %v, got %v\n", test.num, test.base, test.expected, result)
		}
	}
}

func TestBaseToDec(t *testing.T) {
	tests := []struct {
		num            string
		base, expected uint
	}{
		{"11", 16, 17},
		{"10", 16, 16},
		{"A", 16, 10},
		{"11", 15, 16},
		{"G", 32, 16},
		{"23", 16, 35},
		{"K", 32, 20},
	}

	for _, test := range tests {
		if result, err := BaseToDec(test.num, test.base); err != nil {
			t.Errorf("%v", err)
		} else if result != test.expected {
			t.Errorf("(%v, %v) Expected %v, got %v", test.num, test.base, test.expected, result)
		}
	}
}
