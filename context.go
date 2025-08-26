package main

import (
	"fmt"

	"github.com/fatih/color"
)

const (
	BlocksFull    = 10
	ThresholdWarn = 60
	ThresholdCrit = 80
)

type ContextInfo struct {
	InputTokenCount  int
	OutputTokenCount int
	MaxTokenCount    int
	Notes            string
}

func (c *ContextInfo) ToSection() Section {
	currentTokens := c.InputTokenCount + c.OutputTokenCount
	percentage := c.getPercentage()

	// Create blocks visualization
	filledBlocks := int(percentage / BlocksFull)
	partialBlock := int(percentage)%BlocksFull >= 5

	var blocks []rune
	for i := range BlocksFull {
		if i < filledBlocks {
			blocks = append(blocks, '⛁')
		} else if i == filledBlocks && partialBlock {
			blocks = append(blocks, '⛀')
		} else {
			blocks = append(blocks, '⛶')
		}
	}

	// Format token counts with k notation
	currentK := formatTokenCount(currentTokens)
	maxK := formatTokenCount(c.MaxTokenCount)

	content := fmt.Sprintf("%s %s/%s (%.0f%%) %s",
		string(blocks), currentK, maxK, percentage, c.Notes)

	return Section{
		Content: content,
		Color:   c.getContextColor(),
	}
}

func (c *ContextInfo) getPercentage() float64 {
	currentTokens := c.InputTokenCount + c.OutputTokenCount
	if c.MaxTokenCount == 0 {
		return 0
	}
	return float64(currentTokens) / float64(c.MaxTokenCount) * 100
}

func (c *ContextInfo) getContextColor() *color.Color {
	percentage := c.getPercentage()
	if percentage < ThresholdWarn {
		return color.New(color.FgGreen)
	} else if percentage < ThresholdCrit {
		return color.New(color.FgYellow)
	} else {
		return color.New(color.FgRed)
	}
}

func formatTokenCount(tokens int) string {
	if tokens >= 1000 {
		return fmt.Sprintf("%.0fk", float64(tokens)/1000)
	}
	return fmt.Sprintf("%d", tokens)
}
