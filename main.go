package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
)

func main() {
	color.NoColor = false

	var event StatusHookEvent
	if err := json.NewDecoder(os.Stdin).Decode(&event); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to decode status hook event: %v\n", err)
		os.Exit(1)
	}

	statusLine, err := NewStatusLineFromEvent(&event)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(statusLine.String())
}
