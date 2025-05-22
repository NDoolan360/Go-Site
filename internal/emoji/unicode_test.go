package emoji

import (
	"testing"
)

func TestGetUnicodeEmojis(t *testing.T) {
	//check using the backup file that some emojis are retrieved
	emojis, err := GetUnicodeEmojis()
	if err != nil {
		t.Fatalf("Error getting emojis: %v", err)
	}

	if len(emojis) == 0 {
		t.Fatal("No emojis found")
	}

	//check for some emojis
	emojisToCheck := []string{
		"😀",  // grinning face
		"😂",  // face with tears of joy
		"🥺",  // pleading face
		"❤️", // red heart
		"👍",  // thumbs up
		"🎉",  // party popper
		"🌏",  // globe showing Asia-Australia
		"🚀",  // rocket
		"🍕",  // pizza
		"🐶",  // dog face
	}
	for _, emoji := range emojisToCheck {
		found := false
		for _, e := range emojis {
			if string(e.Unicode) == emoji {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Emoji %s not found in the list", emoji)
		}
	}
}
