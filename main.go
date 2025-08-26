package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

type Pricing struct {
	Input  int
	Output int
}

type StatusHookEvent struct {
	HookEventName  string    `json:"hook_event_name"`
	SessionID      string    `json:"session_id"`
	TranscriptPath string    `json:"transcript_path"`
	CWD            string    `json:"cwd"`
	Model          Model     `json:"model"`
	Workspace      Workspace `json:"workspace"`
	Version        string    `json:"version"`
	OutputStyle    Style     `json:"output_style"`
	Cost           Cost      `json:"cost"`
}

type Model struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

type Workspace struct {
	CurrentDir string `json:"current_dir"`
	ProjectDir string `json:"project_dir"`
}

type Style struct {
	Name string `json:"name"`
}

type Cost struct {
	TotalCostUSD       float64 `json:"total_cost_usd"`
	TotalDurationMS    int64   `json:"total_duration_ms"`
	TotalAPIDurationMS int64   `json:"total_api_duration_ms"`
	TotalLinesAdded    int     `json:"total_lines_added"`
	TotalLinesRemoved  int     `json:"total_lines_removed"`
}

func main() {
	statusLine, err := buildStatusLine()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(statusLine.String())
}

func buildStatusLine() (*StatusLine, error) {
	var statusData StatusHookEvent
	if err := json.NewDecoder(os.Stdin).Decode(&statusData); err != nil {
		return nil, fmt.Errorf("failed to decode status hook event: %w", err)
	}

	user := cmp.Or(os.Getenv("USER"), "unknown")
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}

	return &StatusLine{
		Seperator: " | ",
		Sections: []Section{
			{
				Icon:    "👤",
				Content: fmt.Sprintf("%s@%s", user, hostname),
				Color:   color.New(color.FgCyan),
			},
		},
	}, nil
}

type StatusLine struct {
	Seperator string
	Sections  []Section
}

func (s *StatusLine) String() string {
	separator := cmp.Or(s.Seperator, " | ")
	parts := make([]string, len(s.Sections))
	for i, section := range s.Sections {
		if section.Color != nil {
			parts[i] = section.Color.Sprintf("%s %s", section.Icon, section.Content)
		} else {
			parts[i] = fmt.Sprintf("%s %s", section.Icon, section.Content)
		}
	}

	return strings.Join(parts, separator)
}

type Section struct {
	Icon    string
	Content string
	Color   *color.Color
}
