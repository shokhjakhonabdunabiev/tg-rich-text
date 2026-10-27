package telegram

import (
	"testing"

	"github.com/shokhjakhonabdunabiev/tg-rich-text/plate"
)

func TestRenderTextMarks(t *testing.T) {
	tests := []struct {
		name string
		node plate.Node
		want string
	}{
		{name: "bold", node: plate.Node{Text: "text", Bold: true}, want: "**text**"},
		{name: "italic", node: plate.Node{Text: "text", Italic: true}, want: "_text_"},
		{name: "underline", node: plate.Node{Text: "text", Underline: true}, want: "<u>text</u>"},
		{name: "strikethrough", node: plate.Node{Text: "text", Strikethrough: true}, want: "~~text~~"},
		{name: "code (inline)", node: plate.Node{Text: "text", Code: true}, want: "`text`"},
		{name: "subscript", node: plate.Node{Text: "2", Subscript: true}, want: "<sub>2</sub>"},
		{name: "superscript", node: plate.Node{Text: "2", Superscript: true}, want: "<sup>2</sup>"},
		{name: "kbd", node: plate.Node{Text: "Ctrl+C", Kbd: true}, want: "`Ctrl+C`"},
		{name: "highlight", node: plate.Node{Text: "text", Highlight: true}, want: "==text=="},
		{name: "spoiler", node: plate.Node{Text: "text", Spoiler: true}, want: "||text||"},
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
		node plate.Node
		want string
	}{
		{name: "paragraph", node: plate.Node{Type: "p", Children: []plate.Node{{Text: "Hello world"}}}, want: "Hello world"},
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
		node plate.Node
		want string
	}{
		{
			name: "inline URL",
			node: plate.Node{
				Type: "a",
				URL:  "https://t.me/",
				Children: []plate.Node{
					{Text: "Telegram"},
				},
			},
			want: "[Telegram](https://t.me/)",
		},
		{
			name: "email",
			node: plate.Node{
				Type: "a",
				URL:  "mailto:user@example.com",
				Children: []plate.Node{
					{Text: "email"},
				},
			},
			want: "[email](mailto:user@example.com)",
		},
		{
			name: "phone number",
			node: plate.Node{
				Type: "a",
				URL:  "tel:+123456789",
				Children: []plate.Node{
					{Text: "call me"},
				},
			},
			want: "[call me](tel:+123456789)",
		},
		{
			name: "formatted link text",
			node: plate.Node{
				Type: "a",
				URL:  "https://google.com",
				Children: []plate.Node{
					{Text: "Google", Bold: true},
				},
			},
			want: "[**Google**](https://google.com)",
		},
		{
			name: "multiple children",
			node: plate.Node{
				Type: "a",
				URL:  "https://google.com",
				Children: []plate.Node{
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
		node plate.Node
		want string
	}{
		{name: "h1", node: plate.Node{Type: "h1", Children: []plate.Node{{Text: "Text"}}}, want: "# Text"},
		{name: "h2", node: plate.Node{Type: "h2", Children: []plate.Node{{Text: "Text"}}}, want: "## Text"},
		{name: "h3", node: plate.Node{Type: "h3", Children: []plate.Node{{Text: "Text"}}}, want: "### Text"},
		{name: "h4", node: plate.Node{Type: "h4", Children: []plate.Node{{Text: "Text"}}}, want: "#### Text"},
		{name: "h5", node: plate.Node{Type: "h5", Children: []plate.Node{{Text: "Text"}}}, want: "##### Text"},
		{name: "h6", node: plate.Node{Type: "h6", Children: []plate.Node{{Text: "Text"}}}, want: "###### Text"},
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
		node plate.Node
		want string
	}{
		{name: "date", node: plate.Node{Type: "date", Date: "2026-10-13"}, want: "![2026-10-13](tg://time?unix=1791849600&format=wD)"},
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

func codeLine(text string) plate.Node {
	return plate.Node{Type: "code_line", Children: []plate.Node{{Text: text}}}
}

func TestCodeBlock(t *testing.T) {
	tests := []struct {
		name string
		node plate.Node
		want string
	}{
		{name: "single text node", node: plate.Node{
			Type: "code_block",
			Lang: "go",
			Children: []plate.Node{
				{
					Text: `func main() {
	fmt.Println("Hello, world!")
}`,
				},
			},
		}, want: "```go\nfunc main() {\n\tfmt.Println(\"Hello, world!\")\n}\n```"},
		{name: "code lines are separate lines", node: plate.Node{
			Type: "code_block",
			Lang: "javascript",
			Children: []plate.Node{
				codeLine("function hello() {"),
				codeLine("  console.info('hi');"),
				codeLine("}"),
			},
		}, want: "```javascript\nfunction hello() {\n  console.info('hi');\n}\n```"},
		{name: "empty lines are kept", node: plate.Node{
			Type:     "code_block",
			Children: []plate.Node{codeLine("a := 1"), codeLine(""), codeLine("b := 2")},
		}, want: "```\na := 1\n\nb := 2\n```"},
		{name: "markdown characters are not escaped", node: plate.Node{
			Type:     "code_block",
			Lang:     "js",
			Children: []plate.Node{codeLine("if (a == b && c < d || *p) { x = `t` }")},
		}, want: "```js\nif (a == b && c < d || *p) { x = `t` }\n```"},
		{name: "fence longer than backticks in the code", node: plate.Node{
			Type:     "code_block",
			Lang:     "md",
			Children: []plate.Node{codeLine("```go"), codeLine("x := 1"), codeLine("```")},
		}, want: "````md\n```go\nx := 1\n```\n````"},
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
	got := renderNode(plate.Node{Text: "<u>literal</u>", Superscript: true})
	want := "<sup>&lt;u&gt;literal&lt;/u&gt;</sup>"
	if got != want {
		t.Errorf("renderNode() = %q, want %q", got, want)
	}
}

func TestCallout(t *testing.T) {
	tests := []struct {
		name string
		node plate.Node
		want string
	}{
		{
			name: "with icon",
			node: plate.Node{Type: "callout", Icon: "💡", Children: []plate.Node{{Text: "Tip"}}},
			want: ">💡 Tip",
		},
		{
			name: "without icon",
			node: plate.Node{Type: "callout", Children: []plate.Node{{Text: "Note"}}},
			want: ">Note",
		},
		{
			name: "line breaks stay inside the quote",
			node: plate.Node{Type: "callout", Icon: "⚠️", Children: []plate.Node{{Text: "First\nSecond"}}},
			want: ">⚠️ First<br>Second",
		},
		{
			name: "inline formatting",
			node: plate.Node{Type: "callout", Icon: "💡", Children: []plate.Node{{Text: "Bold", Bold: true}}},
			want: ">💡 **Bold**",
		},
		{
			name: "empty",
			node: plate.Node{Type: "callout", Icon: "💡", Children: []plate.Node{{Text: ""}}},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := renderNode(tt.node); got != tt.want {
				t.Errorf("renderNode() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestImage(t *testing.T) {
	tests := []struct {
		name string
		node plate.Node
		want string
	}{
		{
			name: "without caption",
			node: plate.Node{Type: "img", URL: "https://example.com/photo.jpg", Children: []plate.Node{{Text: ""}}},
			want: "![](https://example.com/photo.jpg)",
		},
		{
			name: "with caption",
			node: plate.Node{
				Type:    "img",
				URL:     "https://example.com/photo.jpg",
				Caption: []plate.Node{{Text: "A "}, {Text: "sunset", Bold: true}},
			},
			want: "![](https://example.com/photo.jpg \"A sunset\")",
		},
		{
			name: "caption quotes and URL spaces are escaped",
			node: plate.Node{
				Type:    "img",
				URL:     "https://example.com/my photo (1).jpg",
				Caption: []plate.Node{{Text: `The "best" \ shot`}},
			},
			want: "![](https://example.com/my%20photo%20%281%29.jpg \"The \\\"best\\\" \\\\ shot\")",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := renderNode(tt.node); got != tt.want {
				t.Errorf("renderNode() = %q, want %q", got, tt.want)
			}
		})
	}
}

func paragraph(text string, indent int) plate.Node {
	return plate.Node{Type: "p", Indent: indent, Children: []plate.Node{{Text: text}}}
}

func toggle(title string, indent int) plate.Node {
	return plate.Node{Type: "toggle", Indent: indent, Children: []plate.Node{{Text: title}}}
}

func TestToggle(t *testing.T) {
	tests := []struct {
		name  string
		nodes []plate.Node
		want  string
	}{
		{
			name:  "encloses following indented blocks",
			nodes: []plate.Node{toggle("More", 0), paragraph("Inside", 1), paragraph("After", 0)},
			want:  "<details><summary>More</summary>\nInside\n</details>\n\nAfter",
		},
		{
			name:  "empty toggle",
			nodes: []plate.Node{toggle("Nothing here", 0), paragraph("After", 0)},
			want:  "<details><summary>Nothing here</summary>\n</details>\n\nAfter",
		},
		{
			name:  "toggle at the end of the document",
			nodes: []plate.Node{paragraph("Before", 0), toggle("More", 0), paragraph("A", 1), paragraph("B", 2)},
			want:  "Before\n\n<details><summary>More</summary>\nA\n\nB\n</details>",
		},
		{
			name: "nested toggles",
			nodes: []plate.Node{
				toggle("Outer", 0),
				paragraph("Outer body", 1),
				toggle("Inner", 1),
				paragraph("Inner body", 2),
				paragraph("Back in outer", 1),
				paragraph("After", 0),
			},
			want: "<details><summary>Outer</summary>\nOuter body\n\n" +
				"<details><summary>Inner</summary>\nInner body\n</details>\n\n" +
				"Back in outer\n</details>\n\nAfter",
		},
		{
			name: "list items use their indent minus one",
			nodes: []plate.Node{
				toggle("Steps", 0),
				{ListStyleType: "decimal", Indent: 2, Children: []plate.Node{{Text: "Inside"}}},
				{ListStyleType: "disc", Indent: 1, Children: []plate.Node{{Text: "Outside"}}},
			},
			want: "<details><summary>Steps</summary>\n1. Inside\n</details>\n\n- Outside",
		},
		{
			name: "summary keeps inline formatting",
			nodes: []plate.Node{
				{Type: "toggle", Children: []plate.Node{{Text: "Bold", Bold: true}}},
				paragraph("Body", 1),
			},
			want: "<details><summary>**Bold**</summary>\nBody\n</details>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := renderRichMarkdown(tt.nodes); got != tt.want {
				t.Errorf("renderRichMarkdown() =\n%q\nwant\n%q", got, tt.want)
			}
		})
	}
}

func listItem(style, text string) plate.Node {
	return plate.Node{Type: "p", ListStyleType: style, Indent: 1, Children: []plate.Node{{Text: text}}}
}

func TestBlockSeparation(t *testing.T) {
	tests := []struct {
		name  string
		nodes []plate.Node
		want  string
	}{
		{
			name:  "paragraphs are separated by a blank line",
			nodes: []plate.Node{paragraph("One", 0), paragraph("Two", 0)},
			want:  "One\n\nTwo",
		},
		{
			name:  "items of one list stay together",
			nodes: []plate.Node{paragraph("Intro", 0), listItem("disc", "A"), listItem("disc", "B"), paragraph("Outro", 0)},
			want:  "Intro\n\n- A\n- B\n\nOutro",
		},
		{
			name:  "different lists are separate blocks",
			nodes: []plate.Node{listItem("disc", "A"), listItem("decimal", "B")},
			want:  "- A\n\n1. B",
		},
		{
			name:  "empty paragraphs are skipped",
			nodes: []plate.Node{paragraph("One", 0), {Type: "p", Children: []plate.Node{{Text: ""}}}, paragraph("Two", 0)},
			want:  "One\n\nTwo",
		},
		{
			name:  "divider",
			nodes: []plate.Node{paragraph("Above", 0), {Type: "hr", Children: []plate.Node{{Text: ""}}}, paragraph("Below", 0)},
			want:  "Above\n\n---\n\nBelow",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := renderRichMarkdown(tt.nodes); got != tt.want {
				t.Errorf("renderRichMarkdown() =\n%q\nwant\n%q", got, tt.want)
			}
		})
	}
}

func TestBlockquote(t *testing.T) {
	tests := []struct {
		name string
		node plate.Node
		want string
	}{
		{
			name: "inline content",
			node: plate.Node{Type: "blockquote", Children: []plate.Node{{Text: "Quoted "}, {Text: "text", Bold: true}}},
			want: ">Quoted **text**",
		},
		{
			name: "paragraphs, list and nested quote",
			node: plate.Node{Type: "blockquote", Children: []plate.Node{
				paragraph("First paragraph.", 0),
				paragraph("Second paragraph.", 0),
				listItem("disc", "Quoted item"),
				{Type: "blockquote", Children: []plate.Node{paragraph("Nested quote.", 0)}},
			}},
			want: ">First paragraph.\n>\n>Second paragraph.\n>\n>- Quoted item\n>\n>>Nested quote.",
		},
		{
			name: "nested quote with several paragraphs",
			node: plate.Node{Type: "blockquote", Children: []plate.Node{
				{Type: "blockquote", Children: []plate.Node{paragraph("A", 0), paragraph("B", 0)}},
			}},
			want: ">>A\n>>\n>>B",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := renderNode(tt.node); got != tt.want {
				t.Errorf("renderNode() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMention(t *testing.T) {
	tests := []struct {
		name string
		node plate.Node
		want string
	}{
		{
			name: "in a paragraph",
			node: plate.Node{Type: "p", Children: []plate.Node{
				{Text: "Thanks "},
				{Type: "mention", Value: "Alice", Children: []plate.Node{{Text: ""}}},
				{Text: "!"},
			}},
			want: "Thanks @Alice!",
		},
		{
			name: "name is escaped",
			node: plate.Node{Type: "mention", Value: "a*b", Children: []plate.Node{{Text: ""}}},
			want: "@a\\*b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := renderNode(tt.node); got != tt.want {
				t.Errorf("renderNode() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLineBreaks(t *testing.T) {
	tests := []struct {
		name string
		node plate.Node
		want string
	}{
		{
			name: "paragraph",
			node: plate.Node{Type: "p", Children: []plate.Node{{Text: "Line one\nLine two"}}},
			want: "Line one<br>Line two",
		},
		{
			name: "windows line ending",
			node: plate.Node{Type: "p", Children: []plate.Node{{Text: "One\r\nTwo"}}},
			want: "One<br>Two",
		},
		{
			name: "inside formatting",
			node: plate.Node{Type: "p", Children: []plate.Node{{Text: "Bold\nstill bold", Bold: true}}},
			want: "**Bold<br>still bold**",
		},
		{
			name: "escaped angle brackets are not confused with the tag",
			node: plate.Node{Type: "p", Children: []plate.Node{{Text: "a < b\nc > d"}}},
			want: "a &lt; b<br>c &gt; d",
		},
		{
			name: "inline quote",
			node: plate.Node{Type: "blockquote", Children: []plate.Node{{Text: "First\nSecond"}}},
			want: ">First<br>Second",
		},
		{
			name: "table cell",
			node: plate.Node{Type: "table", Children: []plate.Node{
				{Type: "tr", Children: []plate.Node{{Type: "th", Children: []plate.Node{paragraph("Head", 0)}}}},
				{Type: "tr", Children: []plate.Node{{Type: "td", Children: []plate.Node{paragraph("One\nTwo", 0)}}}},
			}},
			want: "| Head |\n|:---------|\n| One<br>Two |",
		},
		{
			name: "code blocks keep real newlines",
			node: plate.Node{Type: "code_block", Children: []plate.Node{codeLine("a\nb")}},
			want: "```\na\nb\n```",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := renderNode(tt.node); got != tt.want {
				t.Errorf("renderNode() = %q, want %q", got, tt.want)
			}
		})
	}
}
