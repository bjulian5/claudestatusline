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
		color.New(color.FgRed).Fprintf(os.Stdout, "Error decoding event JSON: %v\n", err)
		return
	}

	statusLine, err := NewStatusLineFromEvent(&event)
	if err != nil {
		color.New(color.FgRed).Fprintf(os.Stdout, "Error creating status line: %v\n", err)
		return
	}

	fmt.Println(statusLine.String())
}
