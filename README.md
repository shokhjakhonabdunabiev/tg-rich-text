# tg-rich-text

Convert Plate editor nodes into Telegram rich Markdown.

## Use as a library

```go
import "github.com/shokhjakhonabdunabiev/tg-rich-text/telegram"

message, err := telegram.ConvertJSON(data)
if err != nil {
	return err
}
```

For already decoded input, use `telegram.Convert` with `[]plate.Node` from the
public `plate` package. Both conversion functions validate before rendering.

## Project layout

- `plate/` contains the shared Plate node model and input validation.
- `telegram/` contains the Telegram rich Markdown renderer.
- `cmd/tg-rich-text/` contains the runnable example program.
- `examples/` contains sample Plate input and expected rich text.
- `docs/` contains project notes and the format roadmap.

Additional output formats can be added as sibling renderers that consume
`plate.Node`, such as an `html/` package.
