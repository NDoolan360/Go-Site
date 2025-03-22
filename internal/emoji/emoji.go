package emoji

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Emoji struct {
	Name       string
	Status     EmojiStatus
	Unicode    []rune
	ShortNames map[string][]string
}

func (e Emoji) String() string {
	return string(e.Unicode)
}

func GetEmojis() []Emoji {
	emojis, err := GetUnicodeEmojis()
	if err != nil {
		panic(err)
	}

	for i, e := range emojis {
		emojis[i].ShortNames = addCustomShortnamesForUnicode(e)
	}

	emojis, err = AddGithubShortnames(emojis)
	if err != nil {
		panic(err)
	}

	return emojis
}

// parses a string into snake_case removing all non-alphanumeric characters
func addCustomShortnamesForUnicode(e Emoji) map[string][]string {
	shortnames := e.ShortNames

	name, _, _ := transform.String(
		transform.Chain(
			norm.NFD,
			runes.Remove(runes.In(unicode.Mn)),
			norm.NFC,
		),
		e.Name,
	)

	// preserve terms
	name = strings.ReplaceAll(name, "o’clock", "o_clock")

	// remove type prefix's
	name = strings.TrimPrefix(name, "flag:")
	name = strings.TrimPrefix(name, "keycap:")

	// remove terms
	name = strings.ReplaceAll(name, "skin tone", "")
	name = strings.TrimSuffix(name, "SAR China")

	// remove special characters
	name = strings.ReplaceAll(name, ".", "")
	name = strings.ReplaceAll(name, "’", "")

	// replace special characters
	name = strings.ReplaceAll(name, "*", "asterisk")
	name = strings.ReplaceAll(name, "#", "hash")

	// convert to lowercase
	name = strings.ToLower(name)

	// replace non-alphanumeric characters with underscores
	re := regexp.MustCompile(`[^a-z0-9]+`)
	name = re.ReplaceAllString(name, "_")

	// remove leading and trailing underscores
	name = strings.Trim(name, "_")

	if e.Status == EmojiStatusPlainText {
		name = "plaintext_" + name
	}

	shortnames["Custom"] = append(shortnames["Custom"], name)

	return shortnames
}
