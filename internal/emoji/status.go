package emoji

import "errors"

type EmojiStatus string

const (
	EmojiStatusFullyQualified     = EmojiStatus("fully-qualified")
	EmojiStatusComponent          = EmojiStatus("component")
	EmojiStatusMinimallyQualified = EmojiStatus("minimally-qualified")
	EmojiStatusUnqualified        = EmojiStatus("unqualified")
	EmojiStatusPlainText          = EmojiStatus("plain-text") // Custom status: Unqualified with only one codepoint
)

func toEmojiStatus(s string, unicode []rune) (EmojiStatus, error) {
	switch s {
	case string(EmojiStatusFullyQualified):
		return EmojiStatusFullyQualified, nil
	case string(EmojiStatusComponent):
		return EmojiStatusComponent, nil
	case string(EmojiStatusMinimallyQualified):
		return EmojiStatusMinimallyQualified, nil
	case string(EmojiStatusUnqualified):
		if len(unicode) == 1 {
			return EmojiStatusPlainText, nil
		} else {
			return EmojiStatusUnqualified, nil
		}
	default:
		return "", errors.New("invalid emoji status")
	}
}
