package plate

type Node struct {
	ID            string `json:"id,omitempty"`
	Type          string `json:"type,omitempty"`
	Children      []Node `json:"children,omitempty"`
	Text          string `json:"text,omitempty"`
	URL           string `json:"url,omitempty"`
	Date          string `json:"date,omitempty"`
	Lang          string `json:"lang,omitempty"`
	TexExpression string `json:"texExpression,omitempty"`

	ListStyleType string `json:"listStyleType,omitempty"`
	Indent        int    `json:"indent,omitempty"`
	ListStart     int    `json:"listStart,omitempty"`
	Checked       bool   `json:"checked,omitempty"`

	Identifier string `json:"identifier,omitempty"`

	// callout
	Icon string `json:"icon,omitempty"`
	// mention: the mentioned name
	Value string `json:"value,omitempty"`
	// img (and other media): caption text nodes
	Caption []Node `json:"caption,omitempty"`

	Bold          bool `json:"bold,omitempty"`
	Italic        bool `json:"italic,omitempty"`
	Underline     bool `json:"underline,omitempty"`
	Strikethrough bool `json:"strikethrough,omitempty"`
	Code          bool `json:"code,omitempty"`
	Subscript     bool `json:"subscript,omitempty"`
	Superscript   bool `json:"superscript,omitempty"`
	Kbd           bool `json:"kbd,omitempty"`
	Highlight     bool `json:"highlight,omitempty"`

	// no alternative found for Spoiler in Platejs
	// ||spoiler|| (Telegram format)
	Spoiler bool `json:"spoiler,omitempty"`
}
