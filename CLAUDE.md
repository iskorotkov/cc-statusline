# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

> **For Users:** See [README.md](README.md) for installation instructions, usage examples, and troubleshooting guidance.

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

# Run tests with race detection and coverage
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

# Format code
go fmt ./...

# Lint code (golangci-lint is available)
golangci-lint run

# Type checking/compilation verification
go build ./...

# Check for compilation errors without building
go vet ./...

# Run static analysis
staticcheck ./...

# Validate YAML configuration files
yamllint .github/dependabot.yml .github/workflows/*.yml

# Download and verify dependencies
go mod download
go mod verify

# Update dependencies (check for available updates)
go list -u -m all

# Security check for dependencies
go list -json -deps ./... | nancy sleuth
```

## Testing Hook Integration

To test the statusline with sample hook data:
```bash
# Basic test with formatted JSON
cat <<EOF | go run .
{
  "session_id": "test-session-123",
  "version": "1.0.0",
  "model": {
    "id": "claude-3-5-sonnet-20241022",
    "display_name": "Claude 3.5 Sonnet"
  },
  "output_style": {
    "name": "default"
  },
  "workspace": {
    "project_dir": "/Users/user/projects/my-project",
    "current_dir": "/Users/user/projects/my-project/src"
  },
  "cost": {
    "total_lines_added": 150,
    "total_lines_removed": 75,
    "total_api_duration_ms": 5000,
    "total_cost_usd": 1.25
  },
  "exceeds_200k_tokens": true
}
EOF

# Alternative one-liner for quick testing
echo '{"session_id":"test","version":"1.0.0","model":{"display_name":"Claude"},"output_style":{"name":"default"},"workspace":{"project_dir":"/path/to/project","current_dir":"/path/to/project/src"},"cost":{"total_lines_added":10,"total_lines_removed":5,"total_api_duration_ms":1000,"total_cost_usd":0.5}}' | go run .
```

## Environment Variables

- `CC_TASK_SERVER`: Base URL for task tracking system (e.g., "https://jira.example.com/browse"). When set, extracts task IDs from branch names and generates clickable links.

## CI/CD and Automation

The project uses GitHub Actions for automated testing, security scanning, and dependency management:

### Continuous Integration (ci.yml)
- **Format Check**: Ensures all Go code is properly formatted with `gofmt`
- **Lint**: Runs `golangci-lint` with comprehensive checks
- **Build and Test**: Cross-platform testing on Ubuntu, macOS, and Windows with race detection and coverage
- **Static Analysis**: Runs `go vet` and `staticcheck` for additional code quality checks
- **Coverage**: Uploads test coverage to Codecov for tracking

### Security Workflow (security.yml)
- Automated security scanning for vulnerabilities in dependencies
- Regular security audits of the codebase

### Release Workflow (release.yml)
- Automated release creation and binary distribution
- Version tagging and changelog generation

### Dependency Management (dependabot.yml)
- **Monthly GitHub Actions Updates**: Automatically updates all GitHub Actions to latest versions in a single grouped PR
- **Monthly Go Dependencies Updates**: Updates Go modules and Go version in a single grouped PR
- **Security Updates**: Immediate updates for security vulnerabilities
- **Auto-merge**: Enabled for patch-level updates to reduce maintenance overhead

All workflows enforce quality gates - code must pass formatting, linting, building, and testing before merging.

## Adding New Parts

### Step-by-Step Guide

1. **Create a new file** in `parts/` directory (e.g., `parts/timestamp.go`)
2. **Implement the Part function** following the established patterns
3. **Add your part** to the composition in `main.go`
4. **Test your part** with sample data

### Complete Example: Timestamp Part

Create `parts/timestamp.go`:
```go
package parts

import (
    "context"
    "time"

    "github.com/iskorotkov/cc-statusline/style"
)

// TimestampPart returns a part that displays the current timestamp
func TimestampPart() Part {
    return func(ctx context.Context, h CCHook) (string, error) {
        // Get current time
        now := time.Now()
        
        // Format timestamp - gracefully handle any formatting errors
        timeStr := now.Format("15:04:05")
        
        // Apply styling using available style functions:
        // style.Bold(), style.Italic(), style.Underline()
        // style.Color(), style.BgColor() with color constants
        formatted := style.Color("TIME", style.ColorCyan) + " " + 
                    style.Bold(timeStr)
        
        return formatted, nil
    }
}

// SessionDurationPart shows how to use CCHook data with error handling
func SessionDurationPart() Part {
    return func(ctx context.Context, h CCHook) (string, error) {
        // Access CCHook data safely
        if h.Cost.TotalAPIDurationMs == 0 {
            return "", nil // Return empty string if no data
        }
        
        // Convert milliseconds to readable duration
        duration := time.Duration(h.Cost.TotalAPIDurationMs) * time.Millisecond
        durationStr := duration.Round(time.Second).String()
        
        // Apply conditional styling based on duration
        var styledDuration string
        if duration > 5*time.Minute {
            styledDuration = style.Color(durationStr, style.ColorRed)
        } else if duration > 1*time.Minute {
            styledDuration = style.Color(durationStr, style.ColorYellow)
        } else {
            styledDuration = style.Color(durationStr, style.ColorGreen)
        }
        
        return style.Color("DURATION", style.ColorBlue) + " " + styledDuration, nil
    }
}
```

### Available Style Functions

```go
// Text styling
style.Bold("text")
style.Italic("text")
style.Underline("text")

// Colors (use with style.Color() and style.BgColor())
style.ColorRed, style.ColorGreen, style.ColorBlue, style.ColorYellow
style.ColorCyan, style.ColorMagenta, style.ColorWhite, style.ColorBlack

// Apply colors
style.Color("text", style.ColorRed)
style.BgColor("text", style.ColorBlue)
```

### Integration in main.go

Add your part to the statusline composition:
```go
// In main.go, add to the Rows() call
statusline := parts.Rows(
    parts.Row(parts.CCVersion(), parts.CCModel(), parts.CCOutputStyle(), parts.CCWorkingDir(), parts.CCSessionStats(), parts.CCContextSize()),
    parts.Row(parts.GitRemoteOrigin(), parts.GitBranch(), parts.GitStatus()),
    parts.Row(parts.GitHubPRInfo()),
    parts.Row(parts.GitHubPRURL()),
    parts.Row(parts.TaskURL()),
    parts.Row(parts.TimestampPart()), // Add your new part here
)
```

### Error Handling Patterns

- **Return empty string**: For missing optional data
- **Log and continue**: For non-critical errors
- **Graceful degradation**: Show partial info if some data is unavailable
- **Context awareness**: Check context cancellation for long operations

### Testing Your Part

```bash
# Test with your part added
cat <<EOF | go run .
{
  "session_id": "test",
  "version": "1.0.0",
  "model": {"display_name": "Claude"},
  "output_style": {"name": "default"},
  "workspace": {"current_dir": "/test"},
  "cost": {"total_api_duration_ms": 5000}
}
EOF
```

## Testing and Quality Guidelines

### Testing Patterns

#### Unit Testing Parts
```go
func TestMyPart(t *testing.T) {
    tests := []struct {
        name     string
        hook     CCHook
        expected string
    }{
        {
            name: "basic functionality",
            hook: CCHook{
                SessionID: "test",
                Version:   "1.0.0",
            },
            expected: "expected output",
        },
        {
            name:     "empty data handling",
            hook:     CCHook{},
            expected: "", // Should return empty string gracefully
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            part := MyPart()
            result, err := part(context.Background(), tt.hook)
            
            if err != nil {
                t.Errorf("unexpected error: %v", err)
            }
            if result != tt.expected {
                t.Errorf("expected %q, got %q", tt.expected, result)
            }
        })
    }
}
```

#### Integration Testing
```go
func TestStatuslineIntegration(t *testing.T) {
    // Test complete statusline composition
    hook := CCHook{
        Version: "1.0.0",
        Model: Model{DisplayName: "Claude"},
        // ... other fields
    }
    
    statusline := Rows(
        Row(CCVersion(), CCModel()),
        // ... other parts
    )
    
    result, err := statusline(context.Background(), hook)
    if err != nil {
        t.Errorf("statusline integration failed: %v", err)
    }
    
    // Verify result format and content
    lines := strings.Split(result, "\n")
    if len(lines) < 2 {
        t.Error("statusline should have multiple lines")
    }
}
```

### Quality Standards

#### Code Coverage
- **Target**: Maintain >80% code coverage for parts
- **Critical paths**: 100% coverage for error handling paths
- **Integration**: Test all part compositions in main.go

#### Performance Expectations
- **Startup time**: <50ms for complete statusline generation
- **Memory usage**: <10MB total allocation
- **External commands**: <100ms timeout for git/gh commands

#### Error Handling Standards
- **Graceful degradation**: Never crash, return empty string for missing data
- **Context awareness**: Respect context cancellation in long operations
- **Logging**: Use structured logging for debugging (avoid in production)

### Debugging Techniques

#### Local Development
```bash
# Debug with verbose output
CLAUDE_DEBUG=1 echo '{}' | go run .

# Test specific scenarios
cat test_data.json | go run .

# Check individual parts
go test -v ./parts -run TestSpecificPart
```

#### Troubleshooting Common Issues

1. **Empty Output**: Check if Part functions return non-empty strings
2. **Git Errors**: Ensure you're in a git repository with proper remotes
3. **GitHub CLI Issues**: Verify `gh auth status` and repository access
4. **Styling Issues**: Test terminal ANSI color support

#### Performance Profiling
```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=.

# Memory profiling
go test -memprofile=mem.prof -bench=.

# Analyze with pprof
go tool pprof cpu.prof
```

### Development Workflow

1. **Write tests first** for new parts
2. **Run quality checks** before committing:
   ```bash
   go test ./...
   go vet ./...
   staticcheck ./...
   golangci-lint run
   ```
3. **Verify integration** with sample hook data
4. **Check performance** impact of new parts
5. **Update documentation** for new functionality

## Shell Command Execution

Use the `shell` package for executing external commands:
- `shell.String()`: Execute command and return string output
- `shell.JSON()`: Execute command and parse JSON output into struct

Both functions properly handle errors and trim whitespace from output.