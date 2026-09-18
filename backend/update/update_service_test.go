package update

import "testing"

func TestIsNewer(t *testing.T) {
	tests := []struct {
		latest   string
		current  string
		expected bool
	}{
		{"1.4.3", "1.4.2", true},
		{"1.4.2", "1.4.2", false},
		{"1.4.1", "1.4.2", false},
		{"1.5.0", "1.4.9", true},
		{"2.0.0", "1.9.9", true},
		{"1.4.2.1", "1.4.2", true},
		{"1.4.2", "1.4.2.1", false},
	}

	for _, tt := range tests {
		got := isNewer(tt.latest, tt.current)
		if got != tt.expected {
			t.Errorf("isNewer(%q, %q) = %v, expected %v", tt.latest, tt.current, got, tt.expected)
		}
	}
}
