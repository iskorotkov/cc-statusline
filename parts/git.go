package parts

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/iskorotkov/cc-statusline/shell"
	"github.com/iskorotkov/cc-statusline/style"
)

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
