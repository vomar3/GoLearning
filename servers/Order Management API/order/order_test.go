package order

import "testing"

func TestValidateString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid Name", "Andrew", true},
		{"Name with space", "Andrew Petrov", true},
		{"Name with digit", "Andrew1", false},
		{"Name with other symbol", "Andrew&", false},
		{"Russian Name", "Андрей", true},
		{"Empty string", "", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ValidateString(test.input)
			if got != test.expected {
				t.Errorf("ValidateString(%q) = %v, want %v", test.input, got, test.expected)
			}
		})
	}
}
