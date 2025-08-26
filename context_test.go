package main

import (
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
)

func TestContextInfoToSection(t *testing.T) {
	tests := []struct {
		name           string
		context        ContextInfo
		expectedBlocks int
		expectedColor  *color.Color
		containsText   []string
	}{
		{
			name: "low usage green",
			context: ContextInfo{
				InputTokenCount:  10000,
				OutputTokenCount: 5000,
				MaxTokenCount:    200000,
			},
			expectedBlocks: 0,
			expectedColor:  color.New(color.FgGreen),
			containsText:   []string{"15k/200k", "(8%)"},
		},
		{
			name: "medium usage yellow",
			context: ContextInfo{
				InputTokenCount:  100000,
				OutputTokenCount: 40000,
				MaxTokenCount:    200000,
			},
			expectedBlocks: 7,
			expectedColor:  color.New(color.FgYellow),
			containsText:   []string{"140k/200k", "(70%)"},
		},
		{
			name: "high usage red",
			context: ContextInfo{
				InputTokenCount:  150000,
				OutputTokenCount: 30000,
				MaxTokenCount:    200000,
			},
			expectedBlocks: 9,
			expectedColor:  color.New(color.FgRed),
			containsText:   []string{"180k/200k", "(90%)"},
		},
		{
			name: "with notes",
			context: ContextInfo{
				InputTokenCount:  5000,
				OutputTokenCount: 2000,
				MaxTokenCount:    200000,
				Notes:            "cached",
			},
			expectedBlocks: 0,
			expectedColor:  color.New(color.FgGreen),
			containsText:   []string{"7k/200k", "cached"},
		},
		{
			name: "zero tokens",
			context: ContextInfo{
				MaxTokenCount: 200000,
			},
			expectedBlocks: 0,
			expectedColor:  color.New(color.FgGreen),
			containsText:   []string{"0/200k", "(0%)"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			section := tt.context.ToSection()

			assert.Empty(t, section.Icon, "Context section should have no icon")

			for _, text := range tt.containsText {
				assert.Contains(t, section.Content, text)
			}

			if tt.expectedColor != nil {
				assert.Equal(t, tt.expectedColor, section.Color)
			}
		})
	}
}

func TestContextInfoGetPercentage(t *testing.T) {
	tests := []struct {
		name     string
		context  ContextInfo
		expected float64
	}{
		{
			name: "50 percent",
			context: ContextInfo{
				InputTokenCount:  50000,
				OutputTokenCount: 50000,
				MaxTokenCount:    200000,
			},
			expected: 50.0,
		},
		{
			name: "zero percent",
			context: ContextInfo{
				InputTokenCount:  0,
				OutputTokenCount: 0,
				MaxTokenCount:    200000,
			},
			expected: 0.0,
		},
		{
			name: "100 percent",
			context: ContextInfo{
				InputTokenCount:  150000,
				OutputTokenCount: 50000,
				MaxTokenCount:    200000,
			},
			expected: 100.0,
		},
		{
			name: "zero max tokens",
			context: ContextInfo{
				InputTokenCount:  1000,
				OutputTokenCount: 1000,
				MaxTokenCount:    0,
			},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.context.getPercentage()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestContextInfoGetContextColor(t *testing.T) {
	tests := []struct {
		name     string
		context  ContextInfo
		expected *color.Color
	}{
		{
			name: "green under 60%",
			context: ContextInfo{
				InputTokenCount:  50000,
				OutputTokenCount: 10000,
				MaxTokenCount:    200000,
			},
			expected: color.New(color.FgGreen),
		},
		{
			name: "yellow at 65%",
			context: ContextInfo{
				InputTokenCount:  100000,
				OutputTokenCount: 30000,
				MaxTokenCount:    200000,
			},
			expected: color.New(color.FgYellow),
		},
		{
			name: "red at 85%",
			context: ContextInfo{
				InputTokenCount:  150000,
				OutputTokenCount: 20000,
				MaxTokenCount:    200000,
			},
			expected: color.New(color.FgRed),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.context.getContextColor()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatTokenCount(t *testing.T) {
	tests := []struct {
		tokens   int
		expected string
	}{
		{0, "0"},
		{999, "999"},
		{1000, "1k"},
		{1500, "2k"},
		{150000, "150k"},
		{200000, "200k"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatTokenCount(tt.tokens)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestContextVisualization(t *testing.T) {
	context := ContextInfo{
		InputTokenCount:  50000,
		OutputTokenCount: 50000,
		MaxTokenCount:    200000,
	}

	section := context.ToSection()

	assert.Contains(t, section.Content, "⛁⛁⛁⛁⛁", "Should have 5 filled blocks for 50%")
	assert.Contains(t, section.Content, "⛶⛶⛶⛶⛶", "Should have 5 empty blocks for 50%")
}
