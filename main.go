package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/shokhjakhonabdunabiev/tg-rich-text/internal/plate"
)

var input = `
[{
    "children": [
        {
            "text": "This is reference testing"
        },
        {
            "children": [
                {
                    "text": ""
                }
            ],
            "identifier": "1",
            "type": "footnoteReference",
            "id": "OhUK-2E7jg"
        },
        {
            "text": ""
        }
    ],
    "type": "p",
    "id": "ARoHjBpjv4"
},
{
    "children": [
        {
            "children": [
                {
                    "text": "testing means testing!"
                }
            ],
            "type": "p",
            "id": "eCMjZ_9wQQ"
        }
    ],
    "identifier": "1",
    "type": "footnoteDefinition",
    "id": "X8XcHotI84"
}]
`

func main() {
	var nodes []plate.Node
	if err := json.Unmarshal([]byte(input), &nodes); err != nil {
		log.Fatal(err)
	}
	fmt.Println(plate.RichMarkdown(nodes))
}
