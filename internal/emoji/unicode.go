package emoji

import (
	_ "embed"
	"errors"
	"io"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const LATEST_EMOJI_LIST_URL = "https://unicode.org/Public/emoji/latest/emoji-test.txt"

func GetUnicodeEmojis() ([]Emoji, error) {
	data, err := getUnicodeData("file")
	if err != nil {
		return nil, err
	}

	emojis, err := parseUnicodeEmojis(string(data), "E15.1", "fully-qualified", "basic")
	if err != nil {
		return nil, err
	}

	return emojis, nil
}

//go:embed "data/emoji-test.txt"
var fallbackFileData []byte

func getUnicodeData(source string) ([]byte, error) {
	if source == "file" {
		return fallbackFileData, nil
	}

	resp, err := http.Get(LATEST_EMOJI_LIST_URL)
	if err != nil {
		return fallbackFileData, err
	}
	defer resp.Body.Close()

	emojiData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return emojiData, nil
}

func parseUnicodeEmojis(data string, maxVersion string, includeStatus ...EmojiStatus) ([]Emoji, error) {
	emojis := []Emoji{}

	r := regexp.MustCompile(`(?mU)^([0-9A-F ]+)\s*; (.+) +# (.+) (.+) (.+)\s*$`)
	matches := r.FindAllStringSubmatch(data, -1)
	for _, m := range matches {
		hexs := m[1]
		rawStatus := m[2]
		emoji := m[3]
		version := m[4]
		name := m[5]

		unicode, err := parseHexsToRunes(hexs)
		if err != nil {
			return nil, err
		}

		status, err := toEmojiStatus(rawStatus, unicode)
		if err != nil {
			return nil, err
		}

		if (compareEmojiVersions(version, maxVersion) > 0) ||
			(len(includeStatus) > 0 && !slices.Contains(includeStatus, status)) {
			continue
		}

		if string(emoji) != string(unicode) {
			return nil, errors.New("emoji does not match unicode")
		}

		e := Emoji{
			Name:       name,
			Unicode:    unicode,
			Status:     status,
			ShortNames: map[string][]string{},
		}

		emojis = append(emojis, e)
	}

	return emojis, nil
}

// compares two version strings in the format "Ex.y" where x and y are integers
func compareEmojiVersions(v1, v2 string) int {
	// Split the versions by dot
	parts1 := strings.Split(strings.TrimPrefix(v1, "E"), ".")
	parts2 := strings.Split(strings.TrimPrefix(v2, "E"), ".")

	// Compare major versions
	major1, _ := strconv.Atoi(parts1[0])
	major2, _ := strconv.Atoi(parts2[0])

	if major1 != major2 {
		return major1 - major2
	}

	// If major versions are equal, compare minor versions
	minor1, _ := strconv.Atoi(parts1[1])
	minor2, _ := strconv.Atoi(parts2[1])

	return minor1 - minor2
}

func parseHexsToRunes(hexsRaw string) ([]rune, error) {
	hexs := strings.Fields(hexsRaw)
	runes := make([]rune, 0, len(hexs))

	for _, hex := range hexs {
		newRune, err := parseHexToRune(hex)
		if err != nil {
			return nil, err
		}
		runes = append(runes, newRune)
	}

	return runes, nil
}

// parses a string hex value into a rune
func parseHexToRune(hexRaw string) (rune, error) {
	hexInt, err := strconv.ParseUint(hexRaw, 16, 32)
	if err != nil {
		return 0, errors.New("failed to parse hex")
	}

	return rune(hexInt), nil
}
