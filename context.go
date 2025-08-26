package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
)

type ContextInfo struct {
	InputTokenCount  int
	OutputTokenCount int
	MaxTokenCount    int
	Notes            string
}

type TranscriptParser struct {
	// GetTranscriptFile allows us to mock file access in tests
	GetTranscriptFile func(path string) (*os.File, error)
}

func NewTranscriptParser() *TranscriptParser {
	return &TranscriptParser{
		GetTranscriptFile: os.Open,
	}
}

func (t *TranscriptParser) ParseContextFromTranscript(transcriptPath string) (*ContextInfo, error) {
	context := &ContextInfo{
		MaxTokenCount: 200000, // Default max context for Claude
	}
	transcriptFile, err := t.GetTranscriptFile(transcriptPath)
	if err != nil {
		// If the file doesn't exist, return empty context
		if os.IsNotExist(err) {
			return context, nil
		}
		return context, fmt.Errorf("failed to open transcript file: %w", err)
	}

	defer transcriptFile.Close()
	scanner := bufio.NewScanner(transcriptFile)
	var mostRecentAssistant *TranscriptEntry

	for scanner.Scan() {
		var entry TranscriptEntry
		line := scanner.Bytes()
		if err := json.Unmarshal(line, &entry); err != nil {
			// Ignore malformed lines
			continue
		}

		// Find the most recent assistant message (which has token usage info)
		if entry.Type == "assistant" && entry.Message.Role == "assistant" {
			mostRecentAssistant = &entry
		}
	}

	if err := scanner.Err(); err != nil {
		context.Notes = fmt.Sprintf("Error reading transcript: %v", err)
		return context, fmt.Errorf("error reading transcript file: %w", err)
	}

	// Use the most recent assistant message's token counts as current context
	// This represents what's currently in Claude's memory
	if mostRecentAssistant != nil {
		usage := mostRecentAssistant.Message.Usage
		// Context length = all input token types from most recent message
		// These already include system prompt and tools
		context.InputTokenCount = usage.InputTokens +
			usage.CacheCreationInputTokens +
			usage.CacheReadInputTokens
		context.OutputTokenCount = usage.OutputTokens
	}

	return context, nil
}

// ToSection creates a Section struct for the context usage
func (c *ContextInfo) ToSection() Section {
	currentTokens := c.InputTokenCount + c.OutputTokenCount
	percentage := c.getPercentage()

	// Create blocks visualization
	filledBlocks := int(percentage / 10)
	partialBlock := int(percentage)%10 >= 5

	var blocks []rune
	for i := range 10 {
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
		Icon:    "", // No icon for context section
		Content: content,
		Color:   c.getContextColor(),
	}
}

// getPercentage calculates the usage percentage
func (c *ContextInfo) getPercentage() float64 {
	currentTokens := c.InputTokenCount + c.OutputTokenCount
	if c.MaxTokenCount == 0 {
		return 0
	}
	return float64(currentTokens) / float64(c.MaxTokenCount) * 100
}

// getContextColor returns the appropriate color based on usage percentage
func (c *ContextInfo) getContextColor() *color.Color {
	percentage := c.getPercentage()
	if percentage < 60 {
		return color.New(color.FgGreen)
	} else if percentage < 80 {
		return color.New(color.FgYellow)
	} else {
		return color.New(color.FgRed)
	}
}

// formatTokenCount formats token count with k notation for thousands
func formatTokenCount(tokens int) string {
	if tokens >= 1000 {
		return fmt.Sprintf("%.0fk", float64(tokens)/1000)
	}
	return fmt.Sprintf("%d", tokens)
}

type TranscriptEntry struct {
	ParentUUID string `json:"parentUuid"`
	UUID       string `json:"uuid"`
	Type       string `json:"type"`
	Message    struct {
		Role  string `json:"role"`
		Usage struct {
			InputTokens              int `json:"input_tokens"`
			OutputTokens             int `json:"output_tokens"`
			CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
			CacheReadInputTokens     int `json:"cache_read_input_tokens"`
		} `json:"usage"`
	} `json:"message"`
}
