package plate

import "testing"

func TestEscapeText(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "plain text",
			in:   "hello world",
			want: "hello world",
		},
		{
			name: "asterisk",
			in:   "hello * world",
			want: `hello \* world`,
		},
		{
			name: "strikethrough",
			in:   "hello ~ world",
			want: `hello \~ world`,
		},
		{
			name: "backtick",
			in:   "hello ` world",
			want: "hello \\` world",
		},
		{
			name: "highlight",
			in:   "hello = world",
			want: `hello \= world`,
		},
		{
			name: "spoiler",
			in:   "hello | world",
			want: `hello \| world`,
		},
		{
			name: "angle brackets",
			in:   "hello <u>world</u>",
			want: "hello &lt;u&gt;world&lt;/u&gt;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := escapeText(tt.in)

			if got != tt.want {
				t.Errorf("escapeText() = %q, want %q", got, tt.want)
			}
		})
	}
}
