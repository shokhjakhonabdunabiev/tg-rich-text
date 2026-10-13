package plate

import (
	"fmt"
	"strings"
	"time"
)

func RichMarkdown(nodes []Node) string {
	var result strings.Builder
	for i, node := range nodes {
		if i > 0 {
			result.WriteString("\n")
		}
		result.WriteString(renderNode(node))
	}
	return result.String()
}

func renderNode(node Node) string {
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
	default:
		return renderChildren(node.Children)
	}
}

func renderChildren(nodes []Node) string {
	var result strings.Builder
	for _, node := range nodes {
		result.WriteString(renderNode(node))
	}
	return result.String()
}

func renderText(node Node) string {
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

func renderHeader(node Node) string {
	level := int(node.Type[1] - '0')
	return strings.Repeat("#", level) + " " + renderChildren(node.Children)
}

func renderLink(node Node) string {
	text := renderChildren(node.Children)
	return "[" + text + "](" + node.URL + ")"
}

func renderDate(node Node) string {
	date, err := time.Parse(time.DateOnly, node.Date)
	if err != nil {
		return renderChildren(node.Children)
	}
	return fmt.Sprintf("![%s](tg://time?unix=%d&format=wD)", node.Date, date.Unix())
}

func renderInlineEquation(node Node) string {
	return "$$" + node.TexExpression + "$$"
}

func renderEquation(node Node) string {
	return "```math\n" + node.TexExpression + "\n```"
}

func renderCodeBlock(node Node) string {
	code := renderChildren(node.Children)
	return "```" + node.Lang + "\n" + code + "\n```"
}

func renderHr() string {
	return "\n---\n"
}

func renderList(node Node) string {
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

func renderBlockquote(node Node) string {
	text := renderChildren(node.Children)
	if text == "" {
		return ""
	}

	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = ">" + line
	}

	return strings.Join(lines, "\n")
}

func renderTable(node Node) string {
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

func renderFootnoteReference(node Node) string {
	return "[^" + node.Identifier + "]"
}

func renderFootnoteDefinition(node Node) string {
	text := renderChildren(node.Children)
	return fmt.Sprintf("[^%s]: %s", node.Identifier, text)
}
