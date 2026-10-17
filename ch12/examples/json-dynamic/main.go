// JSON whose shape you do not know in advance: decode into any and look around
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	raw := `{"name":"xkcd","comics":[{"num":1,"title":"Barrel"},{"num":2,"title":"Petit Trees"}],"active":true}`

	var v any
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	// objects become map[string]any, arrays []any, numbers float64, and you assert as you go
	obj := v.(map[string]any)
	fmt.Println("name:", obj["name"])
	for _, c := range obj["comics"].([]any) {
		comic := c.(map[string]any)
		fmt.Printf("#%.0f %s\n", comic["num"], comic["title"])
	}
	if active, ok := obj["active"].(bool); ok && active {
		fmt.Println("still publishing")
	}
}
