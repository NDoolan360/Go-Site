package emoji

import "testing"

func TestAddGithubShortnames(t *testing.T) {
	// Test with a sample emoji
	inEmojis := []Emoji{
		{
			Unicode: []rune{0x1F600},
			Status:  EmojiStatusFullyQualified,
			ShortNames: map[string][]string{
				"Unicode": {"grinning face"},
			},
		},
		{
			Unicode: []rune{0x263A},
			Status:  EmojiStatusUnqualified,
			ShortNames: map[string][]string{
				"Unicode": {"smiling face"},
			},
		},
		{
			Unicode: []rune{0x1F642, 0x200D, 0x2194},
			Status:  EmojiStatusMinimallyQualified,
			ShortNames: map[string][]string{
				"Unicode": {"head shaking horizontally"},
			},
		},
	}

	outEmojis, err := AddGithubShortnames(inEmojis)
	if err != nil {
		t.Fatalf("Error adding GitHub shortnames: %v", err)
	}

	if len(outEmojis) != 3 {
		t.Fatalf("Expected 3 emoji, got %d", len(outEmojis))
	}

	for _, e := range outEmojis {
		if e.Status == EmojiStatusFullyQualified && len(e.ShortNames["Github"]) == 0 {
			t.Fatal("Expected GitHub shortnames to be added")
		} else if e.Status != EmojiStatusFullyQualified && len(e.ShortNames["Github"]) != 0 {
			t.Fatal("Expected no GitHub shortnames for unqualified emoji")
		}
	}
}
