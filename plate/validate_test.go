package plate

import (
	"strings"
	"testing"
)

func TestValidateNewBlocks(t *testing.T) {
	tests := []struct {
		name    string
		node    Node
		wantErr string
	}{
		{name: "callout", node: Node{Type: "callout", Icon: "💡", Children: []Node{{Text: "Tip"}}}},
		{name: "toggle", node: Node{Type: "toggle", Children: []Node{{Text: "More"}}}},
		{name: "image", node: Node{Type: "img", URL: "https://example.com/a.png", Caption: []Node{{Text: "Caption"}}}},
		{name: "image with uppercase scheme", node: Node{Type: "img", URL: "HTTPS://example.com/a.png"}},
		{name: "image without URL", node: Node{Type: "img"}, wantErr: "[0].url: must not be empty"},
		{name: "image with relative URL", node: Node{Type: "img", URL: "/a.png"}, wantErr: "media URL must use http or https"},
		{name: "image with data URL", node: Node{Type: "img", URL: "data:image/png;base64,AAAA"}, wantErr: "unsupported URL scheme"},
		{
			name:    "invalid caption node",
			node:    Node{Type: "img", URL: "https://example.com/a.png", Caption: []Node{{Type: "p", Text: "x"}}},
			wantErr: "[0].caption[0]: text nodes cannot also have a type or children",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate([]Node{tt.node})
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("Validate() error = %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("Validate() error = %v, want it to contain %q", err, tt.wantErr)
			}
		})
	}
}
