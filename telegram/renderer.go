package telegram

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/shokhjakhonabdunabiev/tg-rich-text/plate"
)

func renderRichMarkdown(nodes []plate.Node) string {
	return renderBlocks(nodes)
}

// renderBlocks renders a list of sibling blocks. Plate toggles are not containers:
// a toggle encloses the blocks that follow it with a greater indent, so those are
// grouped into the toggle's <details> block here.
func renderBlocks(nodes []plate.Node) string {
	var result strings.Builder
	for i := 0; i < len(nodes); {
		if i > 0 {
			result.WriteString("\n")
		}

		node := nodes[i]
		if node.Type == "toggle" {
			end := toggleEnd(nodes, i)
			result.WriteString(renderToggle(node, nodes[i+1:end]))
			i = end
			continue
		}

		result.WriteString(renderNode(node))
		i++
	}
	return result.String()
}

// toggleEnd returns the index just past the blocks enclosed by the toggle at nodes[start].
func toggleEnd(nodes []plate.Node, start int) int {
	toggleIndent := nodes[start].Indent
	end := start + 1
	for end < len(nodes) && enclosingIndent(nodes[end]) > toggleIndent {
		end++
	}
	return end
}

// enclosingIndent mirrors Plate's toggle index: list items use their indent minus one,
// because a list's first level already has indent 1.
func enclosingIndent(node plate.Node) int {
	if node.ListStyleType != "" && node.Indent > 0 {
		return node.Indent - 1
	}
	return node.Indent
}

// Convert validates Plate nodes and renders them as Telegram rich Markdown.
// The validation error identifies the path to the invalid node.
func Convert(nodes []plate.Node) (string, error) {
	if err := plate.Validate(nodes); err != nil {
		return "", err
	}
	return renderRichMarkdown(nodes), nil
}

// ConvertJSON decodes PlateJS nodes from JSON, validates them, and renders
// Telegram rich markdown message.
func ConvertJSON(data []byte) (string, error) {
	var nodes []plate.Node
	if err := json.Unmarshal(data, &nodes); err != nil {
		return "", fmt.Errorf("decode Plate nodes: %w", err)
	}
	return Convert(nodes)
}

func renderNode(node plate.Node) string {
	if node.Text != "" {
		return renderText(node)
	}

	if node.ListStyleType != "" {
		return renderList(node)
	}

	switch node.Type {
	case "p":
		return renderChildren(node.Children)
	case "a":
		return renderLink(node)
	case "date":
		return renderDate(node)
	case "h1", "h2", "h3", "h4", "h5", "h6":
		return renderHeader(node)
	case "inline_equation":
		return renderInlineEquation(node)
	case "equation":
		return renderEquation(node)
	case "code_line":
		return renderChildren(node.Children)
	case "code_block":
		return renderCodeBlock(node)
	case "hr":
		return renderHr()
	case "blockquote":
		return renderBlockquote(node)
	case "table":
		return renderTable(node)
	case "footnoteReference":
		return renderFootnoteReference(node)
	case "footnoteDefinition":
		return renderFootnoteDefinition(node)
	case "callout":
		return renderCallout(node)
	case "toggle":
		// Only reached for a toggle outside a block list; it has no enclosed content.
		return renderToggle(node, nil)
	case "img":
		return renderImage(node)
	default:
		return renderChildren(node.Children)
	}
}

func renderChildren(nodes []plate.Node) string {
	var result strings.Builder
	for _, node := range nodes {
		result.WriteString(renderNode(node))
	}
	return result.String()
}

func renderText(node plate.Node) string {
	text := escapeText(node.Text)
	if node.Bold {
		// "__bold text__" (also valid)
		text = "**" + text + "**"
	}
	if node.Italic {
		// "*italic text*" (also valid)
		text = "_" + text + "_"
	}
	if node.Underline {
		text = "<u>" + text + "</u>"
	}
	if node.Strikethrough {
		text = "~~" + text + "~~"
	}
	if node.Code || node.Kbd {
		return "`" + node.Text + "`"
	}
	if node.Superscript {
		text = "<sup>" + text + "</sup>"
	}
	if node.Subscript {
		text = "<sub>" + text + "</sub>"
	}
	if node.Highlight {
		text = "==" + text + "=="
	}
	if node.Spoiler {
		text = "||" + text + "||"
	}
	return text
}

func renderHeader(node plate.Node) string {
	level := int(node.Type[1] - '0')
	return strings.Repeat("#", level) + " " + renderChildren(node.Children)
}

func renderLink(node plate.Node) string {
	text := renderChildren(node.Children)
	return "[" + text + "](" + node.URL + ")"
}

func renderDate(node plate.Node) string {
	date, err := time.Parse(time.DateOnly, node.Date)
	if err != nil {
		return renderChildren(node.Children)
	}
	return fmt.Sprintf("![%s](tg://time?unix=%d&format=wD)", node.Date, date.Unix())
}

func renderInlineEquation(node plate.Node) string {
	return "$$" + node.TexExpression + "$$"
}

func renderEquation(node plate.Node) string {
	return "```math\n" + node.TexExpression + "\n```"
}

func renderCodeBlock(node plate.Node) string {
	code := renderChildren(node.Children)
	return "```" + node.Lang + "\n" + code + "\n```"
}

func renderHr() string {
	return "\n---\n"
}

func renderList(node plate.Node) string {
	switch node.ListStyleType {
	case "decimal":
		text := renderChildren(node.Children)
		i := node.ListStart
		if i == 0 {
			i = 1
		}
		return fmt.Sprintf("%d. %s", i, text)
	case "todo":
		if node.Checked {
			return "- [x] " + renderChildren(node.Children)
		}
		return "- [ ] " + renderChildren(node.Children)
	default:
		// disc list type
		return "- " + renderChildren(node.Children)
	}
}

func renderBlockquote(node plate.Node) string {
	return quoteLines(renderChildren(node.Children))
}

func quoteLines(text string) string {
	if text == "" {
		return ""
	}

	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = ">" + line
	}

	return strings.Join(lines, "\n")
}

// renderCallout renders a callout as a block quotation led by its icon.
// Telegram has no callout block; a quote keeps it visually set apart.
func renderCallout(node plate.Node) string {
	text := renderChildren(node.Children)
	if text == "" {
		return ""
	}
	if node.Icon != "" {
		text = node.Icon + " " + text
	}
	return quoteLines(text)
}

// renderToggle renders a toggle and its enclosed blocks as a collapsed details block.
func renderToggle(node plate.Node, content []plate.Node) string {
	var b strings.Builder
	b.WriteString("<details><summary>")
	b.WriteString(renderChildren(node.Children))
	b.WriteString("</summary>\n")
	if body := renderBlocks(content); body != "" {
		b.WriteString(body)
		b.WriteString("\n")
	}
	b.WriteString("</details>")
	return b.String()
}

var mediaURLEscaper = strings.NewReplacer(" ", "%20", "(", "%28", ")", "%29")

var mediaCaptionEscaper = strings.NewReplacer(`\`, `\\`, `"`, `\"`)

// renderImage renders an image as a separate media block; the caption becomes the title.
// Blank lines around it keep it from being parsed as part of a neighbouring paragraph.
func renderImage(node plate.Node) string {
	image := "![](" + mediaURLEscaper.Replace(node.URL)
	if caption := strings.TrimSpace(plainText(node.Caption)); caption != "" {
		image += ` "` + mediaCaptionEscaper.Replace(caption) + `"`
	}
	return "\n" + image + ")\n"
}

// plainText concatenates the text of nodes, ignoring formatting.
func plainText(nodes []plate.Node) string {
	var b strings.Builder
	for _, node := range nodes {
		b.WriteString(node.Text)
		b.WriteString(plainText(node.Children))
	}
	return b.String()
}

func renderTable(node plate.Node) string {
	var b strings.Builder

	for i, row := range node.Children {
		b.WriteString("|")

		for _, cell := range row.Children {
			content := strings.TrimSpace(renderChildren(cell.Children))
			b.WriteString(" ")
			b.WriteString(content)
			b.WriteString(" |")
		}
		b.WriteString("\n")

		// Add the Markdown separator after the header row.
		if i == 0 {
			b.WriteString("|")

			for _, cell := range row.Children {
				switch cell.Type {
				case "th":
					// Default to left alignment.
					b.WriteString(":---------|")
				default:
					b.WriteString("---------|")
				}
			}
			b.WriteString("\n")
		}
	}

	return strings.TrimSuffix(b.String(), "\n")
}

func renderFootnoteReference(node plate.Node) string {
	return "[^" + node.Identifier + "]"
}

func renderFootnoteDefinition(node plate.Node) string {
	text := renderChildren(node.Children)
	return fmt.Sprintf("[^%s]: %s", node.Identifier, text)
}
