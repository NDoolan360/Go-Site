package emoji

import "testing"

func TestToEmojiStatus(t *testing.T) {
	// Test with valid statuses
	tests := []struct {
		input    string
		unicode  []rune
		expected EmojiStatus
	}{
		{"fully-qualified", []rune{0x263A, 0xFE0F}, EmojiStatusFullyQualified},
		{"component", []rune{0x1F9B0}, EmojiStatusComponent},
		{"minimally-qualified", []rune{0x1F636, 0x200D, 0x1F32B}, EmojiStatusMinimallyQualified},
		{"unqualified", []rune{0x1F636, 0x1F32B}, EmojiStatusUnqualified},
		{"unqualified", []rune{0x263A}, EmojiStatusPlainText},
	}

	for _, test := range tests {
		result, err := toEmojiStatus(test.input, test.unicode)
		if err != nil {
			t.Errorf("Error converting %s, %s: %v", string(test.unicode), test.input, err)
			continue
		}
		if result != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, result)
		}
	}

	// Test with an invalid status
	invalidStatus := "invalid-status"
	_, err := toEmojiStatus(invalidStatus, []rune{0x1F600})
	if err == nil {
		t.Errorf("Expected error for invalid status %s, got nil", invalidStatus)
	}
}
