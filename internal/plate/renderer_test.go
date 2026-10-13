package plate

import (
	"testing"
)

func TestRenderTextMarks(t *testing.T) {
	tests := []struct {
		name string
		node Node
		want string
	}{
		{name: "bold", node: Node{Text: "text", Bold: true}, want: "**text**"},
		{name: "italic", node: Node{Text: "text", Italic: true}, want: "_text_"},
		{name: "underline", node: Node{Text: "text", Underline: true}, want: "<u>text</u>"},
		{name: "strikethrough", node: Node{Text: "text", Strikethrough: true}, want: "~~text~~"},
		{name: "code (inline)", node: Node{Text: "text", Code: true}, want: "`text`"},
		{name: "subscript", node: Node{Text: "2", Subscript: true}, want: "<sub>2</sub>"},
		{name: "superscript", node: Node{Text: "2", Superscript: true}, want: "<sup>2</sup>"},
		{name: "kbd", node: Node{Text: "Ctrl+C", Kbd: true}, want: "`Ctrl+C`"},
		{name: "highlight", node: Node{Text: "text", Highlight: true}, want: "==text=="},
		{name: "spoiler", node: Node{Text: "text", Spoiler: true}, want: "||text||"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := renderNode(tt.node); got != tt.want {
				t.Errorf("renderNode() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParagraph(t *testing.T) {
	tests := []struct {
		name string
		node Node
		want string
	}{
		{name: "paragraph", node: Node{Type: "p", Children: []Node{{Text: "Hello world"}}}, want: "Hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderNode(tt.node)

			if got != tt.want {
				t.Errorf("renderNode() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLink(t *testing.T) {
	tests := []struct {
		name string
		node Node
		want string
	}{
		{
			name: "inline URL",
			node: Node{
				Type: "a",
				URL:  "https://t.me/",
				Children: []Node{
					{Text: "Telegram"},
				},
			},
			want: "[Telegram](https://t.me/)",
		},
		{
			name: "email",
			node: Node{
				Type: "a",
				URL:  "mailto:user@example.com",
				Children: []Node{
					{Text: "email"},
				},
			},
			want: "[email](mailto:user@example.com)",
		},
		{
			name: "phone number",
			node: Node{
				Type: "a",
				URL:  "tel:+123456789",
				Children: []Node{
					{Text: "call me"},
				},
			},
			want: "[call me](tel:+123456789)",
		},
		{
			name: "formatted link text",
			node: Node{
				Type: "a",
				URL:  "https://google.com",
				Children: []Node{
					{Text: "Google", Bold: true},
				},
			},
			want: "[**Google**](https://google.com)",
		},
		{
			name: "multiple children",
			node: Node{
				Type: "a",
				URL:  "https://google.com",
				Children: []Node{
					{Text: "Go", Italic: true},
					{Text: "ogle", Bold: true, Strikethrough: true},
				},
			},
			want: "[_Go_~~**ogle**~~](https://google.com)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderLink(tt.node)

			if got != tt.want {
				t.Errorf("renderLink() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHeading(t *testing.T) {
	tests := []struct {
		name string
		node Node
		want string
	}{
		{name: "h1", node: Node{Type: "h1", Children: []Node{{Text: "Text"}}}, want: "# Text"},
		{name: "h2", node: Node{Type: "h2", Children: []Node{{Text: "Text"}}}, want: "## Text"},
		{name: "h3", node: Node{Type: "h3", Children: []Node{{Text: "Text"}}}, want: "### Text"},
		{name: "h4", node: Node{Type: "h4", Children: []Node{{Text: "Text"}}}, want: "#### Text"},
		{name: "h5", node: Node{Type: "h5", Children: []Node{{Text: "Text"}}}, want: "##### Text"},
		{name: "h6", node: Node{Type: "h6", Children: []Node{{Text: "Text"}}}, want: "###### Text"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderHeader(tt.node)

			if got != tt.want {
				t.Errorf("renderHeader() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDate(t *testing.T) {
	tests := []struct {
		name string
		node Node
		want string
	}{
		{name: "date", node: Node{Type: "date", Date: "2026-10-13"}, want: "![2026-10-13](tg://time?unix=1791849600&format=wD)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderDate(tt.node)

			if got != tt.want {
				t.Errorf("renderHeader() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCodeBlock(t *testing.T) {
	tests := []struct {
		name string
		node Node
		want string
	}{
		{name: "date", node: Node{
			Type: "code_block",
			Lang: "go",
			Children: []Node{
				{
					Text: `func main() {
	fmt.Println("Hello, world!")
}`,
				},
			},
		}, want: "```go\nfunc main() {\n\tfmt.Println(\"Hello, world!\")\n}\n```"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderCodeBlock(tt.node)

			if got != tt.want {
				t.Errorf("renderCodeBlock() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRenderGeneratedTagsAndEscapedText(t *testing.T) {
	got := renderNode(Node{Text: "<u>literal</u>", Superscript: true})
	want := "<sup>&lt;u&gt;literal&lt;/u&gt;</sup>"
	if got != want {
		t.Errorf("renderNode() = %q, want %q", got, want)
	}
}
