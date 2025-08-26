# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go build -o claudestatusline    # Build
go test -v ./...                 # Test
go fmt ./...                     # Format
```

## Architecture

**Pipeline**: JSON stdin → Parse event → Read transcript file → Calculate token usage → Output colored status line

**Key Files**:
- `statusline.go`: Composes status from `Section` structs (icon, content, color)
- `context.go` + `transcript.go`: Parses Claude Code JSONL transcript to extract token counts from latest assistant message
- `event.go`: StatusHookEvent struct matching Claude Code's JSON schema

**Token Usage**:
- Combines input + cache_creation + cache_read tokens
- Shows progress bar with 10 Unicode blocks (⛁, ⛀, ⛶)
- Colors: Green <60%, Yellow 60-80%, Red >80%
- Max context: 200k tokens (hardcoded)

**Testing**: Uses dependency injection for file operations (`TranscriptParser.GetTranscriptFile`)