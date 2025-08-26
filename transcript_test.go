package main

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTranscriptParserParseContextFromTranscript(t *testing.T) {
	tests := []struct {
		name              string
		transcriptContent string
		fileError         error
		expectedInput     int
		expectedOutput    int
		expectedNotes     string
		expectError       bool
	}{
		{
			name: "successful parse with assistant message",
			transcriptContent: `{"type":"user","message":{"role":"user"}}
{"type":"assistant","message":{"role":"assistant","usage":{"input_tokens":5000,"output_tokens":2000,"cache_creation_input_tokens":1000,"cache_read_input_tokens":500}}}`,
			expectedInput:  6500,
			expectedOutput: 2000,
		},
		{
			name: "multiple assistant messages uses most recent",
			transcriptContent: `{"type":"assistant","message":{"role":"assistant","usage":{"input_tokens":1000,"output_tokens":500}}}
{"type":"assistant","message":{"role":"assistant","usage":{"input_tokens":5000,"output_tokens":2000}}}`,
			expectedInput:  5000,
			expectedOutput: 2000,
		},
		{
			name:              "empty transcript",
			transcriptContent: "",
			expectedInput:     0,
			expectedOutput:    0,
		},
		{
			name: "malformed json lines are ignored",
			transcriptContent: `not json
{"type":"assistant","message":{"role":"assistant","usage":{"input_tokens":3000,"output_tokens":1000}}}
also not json`,
			expectedInput:  3000,
			expectedOutput: 1000,
		},
		{
			name:           "file not found returns empty context",
			fileError:      os.ErrNotExist,
			expectedInput:  0,
			expectedOutput: 0,
			expectError:    false,
		},
		{
			name:        "other file error",
			fileError:   errors.New("permission denied"),
			expectError: true,
		},
		{
			name:              "user messages are ignored",
			transcriptContent: `{"type":"user","message":{"role":"user","usage":{"input_tokens":1000,"output_tokens":500}}}`,
			expectedInput:     0,
			expectedOutput:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &TranscriptParser{
				GetTranscriptFile: func(path string) (*os.File, error) {
					if tt.fileError != nil {
						return nil, tt.fileError
					}

					tempFile, err := os.CreateTemp("", "test-transcript-*.jsonl")
					if err != nil {
						t.Fatal(err)
					}
					t.Cleanup(func() { os.Remove(tempFile.Name()) })

					_, err = tempFile.WriteString(tt.transcriptContent)
					if err != nil {
						t.Fatal(err)
					}

					_, err = tempFile.Seek(0, 0)
					if err != nil {
						t.Fatal(err)
					}

					return tempFile, nil
				},
			}

			context, err := parser.ParseContextFromTranscript("/test/path")

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedInput, context.InputTokenCount)
			assert.Equal(t, tt.expectedOutput, context.OutputTokenCount)
			assert.Equal(t, 200000, context.MaxTokenCount)
			if tt.expectedNotes != "" {
				assert.Contains(t, context.Notes, tt.expectedNotes)
			}
		})
	}
}

func TestTranscriptParserWithRealFile(t *testing.T) {
	tempContent := `{"type":"assistant","message":{"role":"assistant","usage":{"input_tokens":10000,"output_tokens":5000,"cache_creation_input_tokens":2000,"cache_read_input_tokens":1000}}}`

	tempFile, err := os.CreateTemp("", "transcript-*.jsonl")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString(tempContent)
	require.NoError(t, err)
	tempFile.Close()

	parser := NewTranscriptParser()
	context, err := parser.ParseContextFromTranscript(tempFile.Name())

	require.NoError(t, err)
	assert.Equal(t, 13000, context.InputTokenCount)
	assert.Equal(t, 5000, context.OutputTokenCount)
}

func TestTranscriptEntryStructure(t *testing.T) {
	jsonData := `{
		"parentUuid": "parent-123",
		"uuid": "uuid-456",
		"type": "assistant",
		"message": {
			"role": "assistant",
			"usage": {
				"input_tokens": 100,
				"output_tokens": 50,
				"cache_creation_input_tokens": 20,
				"cache_read_input_tokens": 10
			}
		}
	}`

	var entry TranscriptEntry
	err := json.Unmarshal([]byte(jsonData), &entry)

	require.NoError(t, err)
	assert.Equal(t, "parent-123", entry.ParentUUID)
	assert.Equal(t, "uuid-456", entry.UUID)
	assert.Equal(t, "assistant", entry.Type)
	assert.Equal(t, "assistant", entry.Message.Role)
	assert.Equal(t, 100, entry.Message.Usage.InputTokens)
	assert.Equal(t, 50, entry.Message.Usage.OutputTokens)
	assert.Equal(t, 20, entry.Message.Usage.CacheCreationInputTokens)
	assert.Equal(t, 10, entry.Message.Usage.CacheReadInputTokens)
}
