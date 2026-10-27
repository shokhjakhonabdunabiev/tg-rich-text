package plate

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

const maxNodeDepth = 100

// Validate checks Plate nodes for malformed structures and values that cannot
// be rendered reliably. Paths use zero-based child indexes (for example, [1].children[0]).
func Validate(nodes []Node) error {
	for i, node := range nodes {
		if err := validateNode(node, fmt.Sprintf("[%d]", i), 0); err != nil {
			return err
		}
	}
	return nil
}

func validateNode(node Node, path string, depth int) error {
	if depth > maxNodeDepth {
		return fmt.Errorf("%s: maximum nesting depth of %d exceeded", path, maxNodeDepth)
	}
	if node.Text != "" {
		if node.Type != "" || len(node.Children) != 0 {
			return fmt.Errorf("%s: text nodes cannot also have a type or children", path)
		}
	} else {
		switch node.Type {
		case "", "p", "a", "date", "h1", "h2", "h3", "h4", "h5", "h6",
			"inline_equation", "equation", "code_line", "code_block", "hr",
			"blockquote", "table", "tr", "th", "td", "li", "footnoteReference", "footnoteDefinition",
			"callout", "toggle", "img":
		default:
			return fmt.Errorf("%s.type: unsupported node type %q", path, node.Type)
		}
	}

	if node.ListStyleType != "" && node.ListStyleType != "disc" && node.ListStyleType != "decimal" && node.ListStyleType != "todo" {
		return fmt.Errorf("%s.listStyleType: unsupported value %q (expected disc, decimal, or todo)", path, node.ListStyleType)
	}
	if node.Indent < 0 {
		return fmt.Errorf("%s.indent: must be zero or greater", path)
	}
	if node.ListStart < 0 {
		return fmt.Errorf("%s.listStart: must be zero or greater", path)
	}
	if node.Type == "a" {
		if err := validateURL(node.URL); err != nil {
			return fmt.Errorf("%s.url: %w", path, err)
		}
	}
	if node.Type == "img" {
		if err := validateMediaURL(node.URL); err != nil {
			return fmt.Errorf("%s.url: %w", path, err)
		}
		for i, child := range node.Caption {
			if err := validateNode(child, fmt.Sprintf("%s.caption[%d]", path, i), depth+1); err != nil {
				return err
			}
		}
	}
	if node.Type == "date" {
		if _, err := time.Parse(time.DateOnly, node.Date); err != nil {
			return fmt.Errorf("%s.date: expected YYYY-MM-DD, got %q", path, node.Date)
		}
	}
	if (node.Type == "inline_equation" || node.Type == "equation") && strings.TrimSpace(node.TexExpression) == "" {
		return fmt.Errorf("%s.texExpression: must not be empty", path)
	}
	if node.Type == "footnoteReference" && strings.TrimSpace(node.Identifier) == "" {
		return fmt.Errorf("%s.identifier: must not be empty for a footnote reference", path)
	}
	for i, child := range node.Children {
		if err := validateNode(child, fmt.Sprintf("%s.children[%d]", path, i), depth+1); err != nil {
			return err
		}
	}
	return nil
}

func validateURL(raw string) error {
	if strings.TrimSpace(raw) == "" {
		return fmt.Errorf("must not be empty")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	// Relative references (paths, query strings, and fragments) are valid link
	// targets too, for example "/help", "../page", or "#section".
	if u.Scheme == "" {
		return nil
	}
	switch strings.ToLower(u.Scheme) {
	case "https", "http":
		if u.Host == "" {
			return fmt.Errorf("must include a host")
		}
	case "mailto", "tel", "tg":
		if u.Opaque == "" && u.Path == "" {
			return fmt.Errorf("must include a destination")
		}
	default:
		return fmt.Errorf("unsupported URL scheme %q", u.Scheme)
	}
	return nil
}

// validateMediaURL checks a media block URL. Telegram rich messages only
// accept HTTP and HTTPS URLs for media blocks.
func validateMediaURL(raw string) error {
	if err := validateURL(raw); err != nil {
		return err
	}
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	if scheme := strings.ToLower(u.Scheme); scheme != "http" && scheme != "https" {
		return fmt.Errorf("media URL must use http or https, got %q", u.Scheme)
	}
	return nil
}
