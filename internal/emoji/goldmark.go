package emoji

import (
	"github.com/yuin/goldmark"
	g_emoji "github.com/yuin/goldmark-emoji"
	g_emoji_def "github.com/yuin/goldmark-emoji/definition"
)

func GoldMarkCustomEmojiExtension() goldmark.Extender {
	return g_emoji.New(g_emoji.WithEmojis(emojisAsGoldmark(GetEmojis())))
}

func emojisAsGoldmark(emojis []Emoji) g_emoji_def.Emojis {
	goldmarkEmojis := []g_emoji_def.Emoji{}
	for _, e := range emojis {
		goldmarkEmojis = append(goldmarkEmojis, g_emoji_def.Emoji{
			Name:       e.Name,
			Unicode:    e.Unicode,
			ShortNames: append(e.ShortNames["Custom"], e.ShortNames["Github"]...),
		})
	}

	return g_emoji_def.NewEmojis(goldmarkEmojis...)
}
