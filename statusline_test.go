package main

import (
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusLineString(t *testing.T) {
	tests := []struct {
		name     string
		line     StatusLine
		contains []string
	}{
		{
			name: "basic status line",
			line: StatusLine{
				Separator: " | ",
				Sections: []Section{
					{Icon: "👤", Content: "user@host"},
					{Icon: "📁", Content: "project"},
				},
			},
			contains: []string{"👤 user@host", "📁 project", " | "},
		},
		{
			name: "custom separator",
			line: StatusLine{
				Separator: " :: ",
				Sections: []Section{
					{Icon: "🤖", Content: "claude"},
					{Icon: "💰", Content: "$0.05"},
				},
			},
			contains: []string{"🤖 claude", "💰 $0.05", " :: "},
		},
		{
			name: "section without icon",
			line: StatusLine{
				Separator: " | ",
				Sections: []Section{
					{Content: "no-icon-section"},
					{Icon: "✓", Content: "with-icon"},
				},
			},
			contains: []string{"no-icon-section", "✓ with-icon"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color.NoColor = true
			defer func() { color.NoColor = false }()

			result := tt.line.String()

			for _, expected := range tt.contains {
				assert.Contains(t, result, expected, "StatusLine output should contain %q", expected)
			}
		})
	}
}

func TestSectionString(t *testing.T) {
	tests := []struct {
		name     string
		section  Section
		expected string
	}{
		{
			name:     "section with icon",
			section:  Section{Icon: "📁", Content: "folder"},
			expected: "📁 folder",
		},
		{
			name:     "section without icon",
			section:  Section{Content: "plain"},
			expected: "plain",
		},
		{
			name:     "empty section",
			section:  Section{},
			expected: "",
		},
		{
			name:     "icon without content",
			section:  Section{Icon: "🤖"},
			expected: "🤖 ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color.NoColor = true
			defer func() { color.NoColor = false }()

			result := tt.section.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSectionWithColor(t *testing.T) {
	color.NoColor = false
	defer func() { color.NoColor = true }()

	section := Section{
		Icon:    "🤖",
		Content: "colored",
		Color:   color.New(color.FgGreen),
	}

	result := section.String()
	assert.Contains(t, result, "🤖 colored", "Colored section should contain icon and content")
	assert.Contains(t, result, "\x1b[", "Colored section should contain ANSI escape codes")
}

func TestNewStatusLineFromEvent(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		event := &StatusHookEvent{
			TranscriptPath: "/tmp/nonexistent.jsonl",
			Model: Model{
				DisplayName: "Claude 3",
			},
			Workspace: Workspace{
				CurrentDir: "/home/user/project",
			},
			Cost: Cost{
				TotalCostUSD: 0.0542,
			},
		}

		statusLine, err := NewStatusLineFromEvent(event)
		require.NoError(t, err)
		require.NotNil(t, statusLine)

		assert.Equal(t, " | ", statusLine.Separator)
		assert.Len(t, statusLine.Sections, 5, "Should have 5 sections: user, dir, model, cost, context")

		assert.Contains(t, statusLine.Sections[0].Content, "@")
		assert.Equal(t, "project", statusLine.Sections[1].Content)
		assert.Equal(t, "Claude 3", statusLine.Sections[2].Content)
		assert.Equal(t, "$0.0542", statusLine.Sections[3].Content)
	})
}

func TestStatusLineDefaultSeparator(t *testing.T) {
	sl := StatusLine{
		Sections: []Section{
			{Content: "first"},
			{Content: "second"},
		},
	}

	result := sl.String()
	assert.Contains(t, result, " | ", "Should use default separator when not specified")
}
