package parts

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/iskorotkov/cc-statusline/shell"
	"github.com/iskorotkov/cc-statusline/style"
)

var gitRemoteGetURLOrigin = func() func(ctx context.Context) (string, error) {
	var remote string
	var err error
	var once sync.Once
	return func(ctx context.Context) (string, error) {
		once.Do(func() {
			remote, err = shell.String(ctx, "git", "ls-remote", "--get-url", "origin")
		})
		return remote, err
	}
}()

var gitBranchShowCurrent = func() func(ctx context.Context) (string, error) {
	var branch string
	var err error
	var once sync.Once
	return func(ctx context.Context) (string, error) {
		once.Do(func() {
			branch, err = shell.String(ctx, "git", "branch", "--show-current")
		})
		return branch, err
	}
}()

var gitStatusPorcelain = func() func(ctx context.Context) (string, error) {
	var status string
	var err error
	var once sync.Once
	return func(ctx context.Context) (string, error) {
		once.Do(func() {
			status, err = shell.String(ctx, "git", "status", "--porcelain")
		})
		return status, err
	}
}()

var gitDiffNumstat = func() func(ctx context.Context) (string, error) {
	var diff string
	var err error
	var once sync.Once
	return func(ctx context.Context) (string, error) {
		once.Do(func() {
			diff, err = shell.String(ctx, "git", "diff", "HEAD", "--numstat")
		})
		return diff, err
	}
}()

func GitRemoteOrigin() Part {
	return func(ctx context.Context, h CCHook) (string, error) {
		remote, _ := gitRemoteGetURLOrigin(ctx)
		if remote != "" {
			remote = strings.TrimSuffix(remote, ".git")
			return style.Underline(limit(remote, 60)), nil
		}
		return "", nil
	}
}

func GitBranch() Part {
	return func(ctx context.Context, h CCHook) (string, error) {
		branch, _ := gitBranchShowCurrent(ctx)
		if branch != "" {
			return style.Italic(limit(branch, 60)), nil
		}
		return "", nil
	}
}

func GitStatus() Part {
	return func(ctx context.Context, h CCHook) (string, error) {
		files, _ := gitStatusPorcelain(ctx)
		if len(files) == 0 {
			return "", nil
		}
		fileCount := count(files)
		if len(fileCount) == 0 {
			return "", nil
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%s:%d", fileCount[0].k, fileCount[0].v)
		for _, p := range fileCount[1:] {
			fmt.Fprintf(&sb, " %s:%d", p.k, p.v)
		}
		return sb.String(), nil
	}
}

func GitDiffStats() Part {
	return func(ctx context.Context, h CCHook) (string, error) {
		diff, _ := gitDiffNumstat(ctx)
		if diff == "" {
			return "", nil
		}
		var added, removed int
		for line := range strings.Lines(diff) {
			fields := strings.Fields(line)
			if len(fields) < 2 {
				continue
			}
			if a, err := strconv.Atoi(fields[0]); err == nil {
				added += a
			}
			if r, err := strconv.Atoi(fields[1]); err == nil {
				removed += r
			}
		}
		if added == 0 && removed == 0 {
			return "", nil
		}
		var sb strings.Builder
		if added > 0 {
			fmt.Fprintf(&sb, style.RGB("+%dL", 127, 255, 127), added)
		}
		if removed > 0 {
			if added > 0 {
				sb.WriteString(" ")
			}
			fmt.Fprintf(&sb, style.RGB("-%dL", 255, 127, 127), removed)
		}
		return sb.String(), nil
	}
}
