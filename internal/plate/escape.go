package plate

import "strings"

// var specialChars = []string{
// 	"_",
// 	"*",
// 	"[",
// 	"]",
// 	"(",
// 	")",
// 	"~",
// 	"`",
// 	">",
// 	"#",
// 	"+",
// 	"-",
// 	"=",
// 	"|",
// 	"{",
// 	"}",
// 	".",
// 	"!",
// }

func escapeText(text string) string {
	var builder strings.Builder

	for _, char := range text {
		switch char {
		case '*', '~', '`', '=', '|':
			builder.WriteByte('\\')
		case '<':
			builder.WriteString("&lt;")
			continue
		case '>':
			builder.WriteString("&gt;")
			continue
		}

		builder.WriteRune(char)
	}

	return builder.String()
}
