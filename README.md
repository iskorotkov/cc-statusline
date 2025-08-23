# cc-statusline

A customizable statusline tool for Claude Code that displays real-time session information, Git status, GitHub PR details, and task tracking links in your terminal.

## Features

- **Claude Code Session Info**: Display current model, version, output style, working directory, session stats (lines added/removed, duration, cost), and 200K+ context indicator
- **Git Integration**: Show remote origin URL, current branch, and file change status
- **GitHub PR Integration**: Display PR number, title, statistics, merge status, and URL
- **Task Tracking**: Automatically extract and link to task/issue numbers from branch names
- **Styled Output**: Rich terminal formatting with colors, bold, italic, and underline styles

## Installation

### Using go install
```bash
go install github.com/iskorotkov/cc-statusline@latest
```

### Build from source
```bash
git clone https://github.com/iskorotkov/cc-statusline.git
cd cc-statusline
go build -o cc-statusline .
```

## Configuration

### Claude Code Setup

To use this statusline with Claude Code, add the following to your Claude Code settings file (`~/.claude/settings.json`):

```json
{
  "statusLine": {
    "type": "command",
    "command": "cc-statusline",
    "padding": 0
  }
}
```

If you want to include task tracking integration, set the `CC_TASK_SERVER` environment variable in the command:

```json
{
  "statusLine": {
    "type": "command",
    "command": "CC_TASK_SERVER='https://github.com/iskorotkov/cc-statusline/issues' cc-statusline",
    "padding": 0
  }
}
```

For Jira integration:
```json
{
  "statusLine": {
    "type": "command",
    "command": "CC_TASK_SERVER='https://yourcompany.atlassian.net/browse' cc-statusline",
    "padding": 0
  }
}
```

### Environment Variables

- `CC_TASK_SERVER`: Base URL for your task tracking system (e.g., `https://jira.example.com/browse`). When set, the tool will extract task IDs from branch names and generate clickable links.

Example:
```bash
export CC_TASK_SERVER="https://github.com/iskorotkov/cc-statusline/issues"
```

## Usage

The tool reads Claude Code hook data from stdin and outputs a formatted statusline:

```bash
echo '{"session_id":"test","version":"1.0.0","model":{"display_name":"Claude 3.5 Sonnet"},"output_style":{"name":"detailed"},"workspace":{"project_dir":"/home/user/project","current_dir":"/home/user/project/src"},"cost":{"total_lines_added":150,"total_lines_removed":75,"total_api_duration_ms":5000,"total_cost_usd":1.25},"exceeds_200k_tokens":true}' | cc-statusline
```

### Output Format

The statusline displays information in multiple rows:

```
CC   v1.0.0 | Claude 3.5 Sonnet | detailed | src | +150L -75L 0.1m $1.25 | 200K+
GIT  https://github.com/iskorotkov/cc-statusline | main | M:5 A:2 D:1
PR   #42 | Fix authentication bug | +200L -50L ~8F M
PR   https://github.com/iskorotkov/cc-statusline/pull/42
TASK https://github.com/iskorotkov/cc-statusline/issues/42
```

Each row shows different information:
- **CC**: Claude Code version, model, output style, current directory, session statistics, context size indicator
- **GIT**: Remote origin URL (underlined), current branch (italic), and file status (M=modified, A=added, D=deleted, etc.)
- **PR**: Pull request number, title, statistics (lines added/removed, files changed), mergeable status
- **TASK**: Extracted task/issue URL based on branch name patterns

## Prerequisites

- Go 1.25 or later
- Git CLI (for Git status information)
- GitHub CLI (`gh`) (for PR information)
- Terminal with ANSI color support

## Input Schema

The tool expects JSON input matching the `CCHook` structure:

```json
{
  "session_id": "string",
  "version": "string",
  "model": {
    "id": "string",
    "display_name": "string"
  },
  "output_style": {
    "name": "string"
  },
  "workspace": {
    "project_dir": "string",
    "current_dir": "string"
  },
  "cost": {
    "total_lines_added": 0,
    "total_lines_removed": 0,
    "total_api_duration_ms": 0,
    "total_cost_usd": 0.0
  },
  "exceeds_200k_tokens": false
}
```

## Development

### Building
```bash
go build -o cc-statusline .
```

### Testing
```bash
go test ./...
```

### Linting
```bash
golangci-lint run
```

### Formatting
```bash
go fmt ./...
```

## Architecture

The project follows a modular architecture with composable parts:

- `main.go`: Entry point and statusline composition
- `parts/`: Individual statusline components (Git, GitHub, Claude Code info)
- `shell/`: Command execution utilities
- `style/`: Terminal formatting functions

Each statusline component is a `Part` - a function that takes context and hook data and returns a formatted string. Parts are composed using `Row()` and `Rows()` functions to build the complete statusline.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

To add a new statusline component:

1. Create a new function in the `parts/` package that returns a `Part`
2. Implement the logic to extract and format your information
3. Add your part to the composition in `main.go`
4. Submit a PR with your changes

## License

MIT License - see [LICENSE](LICENSE) file for details

## Author

Ivan Korotkov ([@iskorotkov](https://github.com/iskorotkov))