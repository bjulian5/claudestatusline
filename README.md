# Claude Status Line

A status line binary for [Claude Code](https://github.com/anthropics/claude-code) that displays session information, git branch, cost tracking, and context usage in your terminal.

## Features

- **Session Info**: Shows user, hostname, and current directory
- **Git Branch**: Displays the current git branch or commit hash
- **Model Display**: Shows the active Claude model
- **Cost Tracking**: Displays cumulative session cost in USD
- **Context Usage**: Visual representation of token usage with color-coded warnings
  - Green: < 60% usage
  - Yellow: 60-80% usage
  - Red: > 80% usage

## Installation

### Option 1: Using `go install`

```bash
go install github.com/bjulian5/claudestatusline@latest
```

### Option 2: Build from source

```bash
git clone https://github.com/bjulian5/claudestatusline.git
cd claudestatusline
go build -o claudestatusline
# Add the binary to your PATH
```

## Configuration

Configure Claude Code to use this status line by adding the following to your `~/.claude/settings.json` file:

```json
{
  "statusLine": {
    "type": "command",
    "command": "claudestatusline"
  }
}
```

After adding this configuration, the status line will automatically appear in your Claude Code sessions. The binary reads Claude's status hook events from stdin and outputs a formatted status line.

## Requirements

- Go 1.24.4 or later
- Claude Code CLI

## License

MIT