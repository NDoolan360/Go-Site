package build

import (
	"bufio"
	"bytes"
	"testing"

	highlighting "github.com/yuin/goldmark-highlighting/v2"

	"internal/emoji"
)

func TestMarkdownTransformer_Transform(t *testing.T) {
	transformer := MarkdownTransformer{}

	t.Run("transforms .md to .html", func(t *testing.T) {
		asset := &Asset{
			Path:       "test.md",
			Data:       []byte("# Hello"),
			SourceRoot: "content",
		}
		err := transformer.Transform(asset)
		if err != nil {
			t.Fatalf("Transform() error = %v", err)
		}
		if asset.Path != "test.html" {
			t.Errorf("Expected path to be 'test.html', got '%s'", asset.Path)
		}

		expectedHTML := "<h1>Hello</h1>\n"
		if string(asset.Data) != expectedHTML {
			t.Errorf("Expected HTML to be %q, got %q", expectedHTML, string(asset.Data))
		}
	})

	t.Run("ignores non-.md files", func(t *testing.T) {
		asset := &Asset{
			Path: "test.txt",
			Data: []byte("Just some text."),
		}
		originalData := make([]byte, len(asset.Data))
		copy(originalData, asset.Data)
		originalPath := asset.Path

		err := transformer.Transform(asset)
		if err != nil {
			t.Fatalf("Transform() error = %v", err)
		}
		if asset.Path != originalPath {
			t.Errorf("Expected path to be '%s', got '%s'", originalPath, asset.Path)
		}
		if !bytes.Equal(asset.Data, originalData) {
			t.Errorf("Expected data to be unchanged, got '%s'", string(asset.Data))
		}
	})
}

func TestEmojisAsGoldmark(t *testing.T) {
	goldmarkResults := emojisAsGoldmark([]emoji.Emoji{
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

type MockCodeBlockContext struct {
	highlighting.CodeBlockContext
	lang []byte
}

func (m *MockCodeBlockContext) Language() ([]byte, bool) {
	if m.lang == nil {
		return nil, false
	}
	return m.lang, true
}
func (m *MockCodeBlockContext) Highlighted() bool                            { return false }
func (m *MockCodeBlockContext) Attributes() highlighting.ImmutableAttributes { return nil }

func TestCodeWrapperRenderer(t *testing.T) {
	tests := []struct {
		name     string
		context  *MockCodeBlockContext
		entering bool
		expected string
	}{
		{
			name: "Entering with language",
			context: &MockCodeBlockContext{
				lang: []byte("go"),
			},
			entering: true,
			expected: `<figure class="codeblock" data-lang="go"><figcaption><button class="copycode" disabled>Copy</button></figcaption>`,
		},
		{
			name: "Exiting with language",
			context: &MockCodeBlockContext{
				lang: []byte("go"),
			},
			entering: false,
			expected: `</figure>`,
		},
		{
			name: "Entering without language",
			context: &MockCodeBlockContext{
				lang: nil,
			},
			entering: true,
			expected: `<figure class="codeblock"><figcaption><button class="copycode" disabled>Copy</button></figcaption><pre class="chroma"><code>`,
		},
		{
			name: "Exiting without language",
			context: &MockCodeBlockContext{
				lang: nil,
			},
			entering: false,
			expected: `</code></pre></figure>`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)
			codeWrapperRenderer(writer, test.context, test.entering)
			writer.Flush()

			if buf.String() != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, buf.String())
			}
		})
	}
}
