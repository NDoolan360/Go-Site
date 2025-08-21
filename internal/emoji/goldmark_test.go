package emoji

import "testing"

func TestEmojisAsGoldmark(t *testing.T) {
	goldmarkResults := emojisAsGoldmark([]Emoji{
		{
			Name:    "smile",
			Unicode: []rune{0x1F604},
			ShortNames: map[string][]string{
				"Custom": {"smile_custom", "smiley"},
				"Github": {"smile_gh"},
			},
		},
	})

	if _, ok := goldmarkResults.Get("smile_custom"); !ok {
		t.Errorf("Expected 'smile_custom' to be present")
	}
	if _, ok := goldmarkResults.Get("smile_gh"); !ok {
		t.Errorf("Expected 'smile_gh' to be present")
	}
	if goldmarkResult, ok := goldmarkResults.Get("smiley"); !ok {
		t.Errorf("Expected 'smiley' to be present")
	} else {
		if goldmarkResult.Name != "smile" {
			t.Errorf("Expected name to be 'smile', got '%s'", goldmarkResult.Name)
		}
		if len(goldmarkResult.Unicode) != 1 && goldmarkResult.Unicode[0] != 0x1F604 {
			t.Errorf("Expected unicode to be '😄', got '%s'", string(goldmarkResult.Unicode))
		}
		if len(goldmarkResult.ShortNames) != 3 {
			t.Errorf("Expected 3 short names, got %d", len(goldmarkResult.ShortNames))
		}
	}
}
