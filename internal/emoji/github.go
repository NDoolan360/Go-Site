package emoji

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const GITHUB_EMOJI_URL = "https://api.github.com/emojis"

func AddGithubShortnames(inEmojis []Emoji) ([]Emoji, error) {
	githubData, err := getGithubData("file")
	if err != nil {
		return nil, err
	}

	outEmojis := []Emoji{}
	for _, ue := range inEmojis {
		e := ue
		names, found := getGithubShortnames(ue, githubData)
		if found {
			for _, name := range names {
				if !slices.Contains(e.ShortNames["Unicode"], name) ||
					!slices.Contains(e.ShortNames["Github"], name) {
					e.ShortNames["Github"] = append(e.ShortNames["Github"], name)
				}
			}
		}
		outEmojis = append(outEmojis, e)
	}

	return outEmojis, nil
}

//go:embed "data/github-emojis.json"
var fallbackGithubEmojis []byte

func getGithubData(source string) (map[string][]string, error) {
	var emojiData []byte

	if source == "file" {
		emojiData = fallbackGithubEmojis
	} else {
		resp, err := http.Get(GITHUB_EMOJI_URL)
		if err == nil {
			defer resp.Body.Close()

			emojiData, err = io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			emojiData = fallbackGithubEmojis
		}
	}

	jsonData := map[string]string{}
	err := json.Unmarshal(emojiData, &jsonData)
	if err != nil {
		return nil, err
	}

	emojiNameMap := map[string][]string{}
	for name, url := range jsonData {
		unicode, err := getUnicodeFromURL(url)
		if err != nil {
			continue
		}
		if _, ok := emojiNameMap[string(unicode)]; !ok {
			emojiNameMap[string(unicode)] = []string{}
		}
		emojiNameMap[string(unicode)] = append(emojiNameMap[string(unicode)], name)
	}

	return emojiNameMap, nil
}

func getGithubShortnames(emoji Emoji, githubUrlMap map[string][]string) ([]string, bool) {
	if emoji.Status != EmojiStatusFullyQualified {
		return nil, false
	}

	runesWithoutZWJOrVS := []rune{}
	for _, r := range emoji.Unicode {
		if r != 0x200D && r != 0xFE0F {
			runesWithoutZWJOrVS = append(runesWithoutZWJOrVS, r)
		}
	}

	shortnames, ok := githubUrlMap[string(runesWithoutZWJOrVS)]

	if ok && len(shortnames) > 0 {
		return shortnames, true
	}

	return nil, false
}

// Need to get the unicode from the url tail and then convert it to a []rune
// Example: https://github.githubassets.com/images/icons/emoji/unicode/1f468-2764-1f468.png?v8
// Should return []rune{0x1f468, 0x2764, 0x1f468}
func getUnicodeFromURL(url string) ([]rune, error) {
	r := regexp.MustCompile(`unicode/([0-9a-f-]+)\.png`)
	unicode := r.FindStringSubmatch(url)
	if len(unicode) != 2 {
		return nil, fmt.Errorf("could not find unicode in url: %s", url)
	}

	unicodeParts := strings.Split(unicode[1], "-")
	runes := make([]rune, len(unicodeParts))
	for i, part := range unicodeParts {
		r, err := strconv.ParseInt(part, 16, 32)
		if err != nil {
			return nil, err
		}
		runes[i] = rune(r)
	}

	return runes, nil
}
