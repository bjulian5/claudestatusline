package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type TranscriptParser struct {
	GetTranscriptFile func(path string) (*os.File, error)
}

func NewTranscriptParser() *TranscriptParser {
	return &TranscriptParser{
		GetTranscriptFile: os.Open,
	}
}

func (t *TranscriptParser) ParseContextFromTranscript(transcriptPath string) (*ContextInfo, error) {
	context := &ContextInfo{
		MaxTokenCount: 200000,
	}
	transcriptFile, err := t.GetTranscriptFile(transcriptPath)
	if err != nil {
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
			continue
		}

		if entry.Type == "assistant" && entry.Message.Role == "assistant" {
			mostRecentAssistant = &entry
		}
	}

	if err := scanner.Err(); err != nil {
		context.Notes = fmt.Sprintf("Error reading transcript: %v", err)
		return context, fmt.Errorf("error reading transcript file: %w", err)
	}

	if mostRecentAssistant != nil {
		usage := mostRecentAssistant.Message.Usage
		context.InputTokenCount = usage.InputTokens +
			usage.CacheCreationInputTokens +
			usage.CacheReadInputTokens
		context.OutputTokenCount = usage.OutputTokens
	}

	return context, nil
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
