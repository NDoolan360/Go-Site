---
Title: Emojis, shortcodes and zero width joiners
Description:
  Thoughts on working with emojis and zero width joiners as runes in go.
Created: 2025-03-08
Modified: 2025-03-14
Image:
  Src: /articles/ZWJ_Emojis/hero.webp
  Alt: A collection of emojis arranged in a spiral pattern.
Tags:
  - Emojis
  - Unicode
Styles:
  - /static/styles/article.css
  - /static/styles/code.css
Scripts:
  - /static/scripts/copy-code.js
---

![{{.Image.Alt}}]({{.Image.Src}})

# {{.Title}}

This site uses the [goldmark markdown emoji extension][goldmark-emoji-extension] to render the markdown files into html. The goldmark markdown processor uses Github's markdown shortcodes to render emojis in the markdown files. Unfortunately, the shortcodes provided by Github are out of date and don't include all the emojis that unicode has to offer. Naturally, I wanted to have them all, so I set out to add the full list of unicode emojis and their shortcodes to my site.

## [Goldmark Markdown Processor](#goldmark-markdown-processor) {#goldmark-markdown-processor}

Thankfully, the goldmark markdown processor is highly extensible and allows for custom extensions to be added. Even better, the goldmark-emoji extension allows for custom emoji shortcodes to be added on initialisation.

```go
// import "github.com/yuin/goldmark-emoji"

emoji.New(emoji.WithEmojis(< custom emoji map >))
```

## [Unicode Emojis](#unicode-emojis) {#unicode-emojis}

Unicode has a wide range of emojis that have been in the standard since [Unicode v6.0][unicode-6.0].
Unicode emojis are not as simple as other unicode characters, as they can be composed of multiple unicode code points.
Some emojis first proposed in [Technical Standard (#51-R2)][emoji-zwj-earliest-draft] and introduced in [Emoji 2.0 (Unicode v8.0)][unicode-8.0] are made up of multiple emojis joined together using the Zero Width Joiner (ZWJ) character.
This allows for a wide range of emojis to be created by combining existing emojis in different ways as well as providing a way to easily add variations to existing emojis.

## [Emojis in Go](#emojis-in-go) {#emojis-in-go}

In Go, strings are made up of runes, each rune representing an individual unicode code point.
This means that for glyphs that are composed of multiple unicode code points (like joined emojis), some care needs to be taken when parsing and manipulating them.

The Go blog has a great post on [strings][go-strings-post] that goes into more detail on how strings are handled and why runes are preferred over characters.

### Example

```go
// :technologist:
runes := []rune{0x1F9D1, 0x200D, 0x1F4BB}
```

The `0x200D` rune is the Zero Width Joiner that combines the :adult: `:adult:` and :laptop: `:laptop:` emojis to create the :technologist: `:technologist:` emoji.

### Parsing Hex Strings to Runes

The [Unicode standard emoji list][unicode-public-emoji-list] represents emojis as groups of hex strings which can be parsed to runes using the `strconv` package in Go.
This is because the rune type is just an alias for the 32-bit integer type and is equivalent to it in all ways.

```go
// parses a string of hex values separated by spaces into a slice of runes
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
```

### Outputting Runes as Strings

To output runes as a string, you can use the `string` type conversion to convert a slice of runes to a string.

```go
runes := []rune{0x1F469, 0x200D, 0x1F4BB}
fmt.Println(string(runes))

// result: 👩‍💻
```

## [Results](#results) {#results}

For my own convenience, I've added a page that lists all the emojis, their Github shortcodes, and the custom shortcodes I've added for them (based on tehir unicode name). You can find it [here](/other/emojis.html) if you're curious.


[goldmark-emoji-extension]: https://github.com/yuin/goldmark-emoji

[unicode-6.0]: https://blog.unicode.org/2010/10/unicode-version-60-support-for-popular.html
[emoji-zwj-earliest-draft]: https://www.unicode.org/reports/tr51/tr51-2-archive.html
[unicode-8.0]: https://blog.unicode.org/2015/06/announcing-unicode-standard-version-80.html

[go-strings-post]: https://go.dev/blog/strings
[unicode-public-emoji-list]: https://unicode.org/Public/emoji/latest/emoji-test.txt
