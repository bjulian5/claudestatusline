package main

import (
	"cmp"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/fatih/color"
)

type StatusLine struct {
	Separator string
	Sections  []Section
}

type Section struct {
	Icon    string
	Content string
	Color   *color.Color
}

func (s *StatusLine) String() string {
	separator := cmp.Or(s.Separator, " | ")
	parts := make([]string, len(s.Sections))
	for i, section := range s.Sections {
		parts[i] = section.String()
	}
	return strings.Join(parts, separator)
}

func (s Section) String() string {
	content := s.Content
	if s.Icon != "" {
		content = fmt.Sprintf("%s %s", s.Icon, s.Content)
	}

	if s.Color != nil {
		return s.Color.Sprint(content)
	}
	return content
}

func NewStatusLineFromEvent(event *StatusHookEvent) (*StatusLine, error) {
	user := cmp.Or(os.Getenv("USER"), "unknown")
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}

	tp := NewTranscriptParser()
	context, err := tp.ParseContextFromTranscript(event.TranscriptPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse context from transcript: %w", err)
	}

	return &StatusLine{
		Separator: " | ",
		Sections: []Section{
			{
				Icon:    "👤",
				Content: fmt.Sprintf("%s@%s", user, hostname),
			},
			{
				Icon:    "📁",
				Content: path.Base(event.Workspace.CurrentDir),
				Color:   color.New(color.FgCyan),
			},
			{
				Icon:    "🤖",
				Content: event.Model.DisplayName,
				Color:   color.New(color.FgGreen),
			},
			{
				Icon:    "💰",
				Content: fmt.Sprintf("$%.4f", event.Cost.TotalCostUSD),
				Color:   color.New(color.FgYellow),
			},
			context.ToSection(),
		},
	}, nil
}
