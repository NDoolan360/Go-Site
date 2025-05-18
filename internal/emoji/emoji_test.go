package emoji

import "testing"

func TestGetEmojis(t *testing.T) {
	outEmojis := GetEmojis()

	if len(outEmojis) != 3773 {
		t.Fatalf("Expected 3773 emoji, got %d", len(outEmojis))
	}

	for _, e := range outEmojis {
		if e.Status == EmojiStatusUnqualified {
			t.Errorf("Emoji %s is unqualified", e.String())
		}
		if len(e.ShortNames) == 0 {
			t.Errorf("Emoji %s has no short names", e.String())
		}
	}
}

func TestEmojiToString(t *testing.T) {
	// Test with a sample emoji
	emoji := Emoji{
		Unicode: []rune{0x1F600},
		Status:  EmojiStatusFullyQualified,
		ShortNames: map[string][]string{
			"Unicode": {"grinning face"},
		},
	}

	expected := "😀"
	if emoji.String() != expected {
		t.Errorf("Expected %s, got %s", expected, emoji.String())
	}
}

func TestPlainTextCustomShortname(t *testing.T) {
	// Test with a sample emoji
	emoji := Emoji{
		Name:    "victory hand",
		Status:  EmojiStatusPlainText,
		Unicode: []rune{0x270C},
		ShortNames: map[string][]string{
			"Unicode": {"victory hand"},
		},
	}

	expected := "plaintext_victory_hand"
	shortnames := addCustomShortnamesForUnicode(emoji)
	if shortnames["Custom"][0] != expected {
		t.Errorf("Expected %s, got %s", expected, shortnames["Custom"][0])
	}
}
