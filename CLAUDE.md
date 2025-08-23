# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based statusline tool for Claude Code that displays contextual information from the current Claude Code session, Git repository state, GitHub pull requests, and task tracking systems. The tool receives JSON input via stdin from Claude Code hooks and outputs formatted statusline text.

## Architecture

The codebase follows a modular architecture:

- **main.go**: Entry point that reads hook data from stdin, processes it through the statusline parts, and outputs formatted text
- **parts/**: Contains all statusline components (Part functions) that render different pieces of information
  - `Part` type: Core abstraction - a function that takes context and CCHook data, returns formatted string
  - Each part handles one specific piece of information (branch, PR stats, remote origin, etc.)
  - Parts are composable via `Row()` and `Rows()` functions
  - Available parts include:
    - `GitRemoteOrigin()`: Displays the Git remote origin URL (without .git suffix), underlined and limited to 60 chars
    - `GitBranch()`: Shows current Git branch name in italics
    - `GitStatus()`: Shows count of modified/staged files by status type
    - GitHub PR parts for pull request information
    - Claude Code session parts for model, cost, and workspace info
- **shell/**: Command execution utilities for running git and gh commands
- **style/**: ANSI terminal formatting functions for text styling

## Key Patterns

1. **Part Composition**: Parts are composed using `Rows()` and `Row()` functions to build the complete statusline
2. **Lazy Evaluation**: Git and GitHub data is fetched once and cached using sync.Once pattern
3. **Error Tolerance**: Parts gracefully return empty strings on errors rather than failing
4. **Hook Data Structure**: The `CCHook` struct in parts/cc.go defines all available Claude Code session data

## Development Commands

```bash
# Build the project
go build -o cc-statusline .

# Run tests (currently no tests exist)
go test ./...

# Format code
go fmt ./...

# Lint code (golangci-lint is available)
golangci-lint run

# Type checking/compilation verification
go build ./...

# Check for compilation errors without building
go vet ./...
```

## Testing Hook Integration

To test the statusline with sample hook data:
```bash
echo '{"session_id":"test","version":"1.0.0","model":{"display_name":"Claude"},"output_style":{"name":"default"},"workspace":{"project_dir":"/path/to/project","current_dir":"/path/to/project/src"},"cost":{"total_lines_added":10,"total_lines_removed":5,"total_api_duration_ms":1000,"total_cost_usd":0.5}}' | go run .
```

## Environment Variables

- `CC_TASK_SERVER`: Base URL for task tracking system (e.g., "https://jira.example.com/browse"). When set, extracts task IDs from branch names and generates clickable links.

## Adding New Parts

1. Create a new function in `parts/` that returns a `Part`
2. The function should accept `context.Context` and `CCHook` parameters
3. Return formatted string using functions from `style/` package
4. Add the new part to the composition in `main.go`

Example:
```go
func MyNewPart() Part {
    return func(ctx context.Context, h CCHook) (string, error) {
        // Your logic here
        return style.Bold("output"), nil
    }
}
```

## Shell Command Execution

Use the `shell` package for executing external commands:
- `shell.String()`: Execute command and return string output
- `shell.JSON()`: Execute command and parse JSON output into struct

Both functions properly handle errors and trim whitespace from output.