package parts

import (
	"context"
	"os"
	"regexp"

	"github.com/iskorotkov/cc-statusline/style"
)

var taskRegex = regexp.MustCompile(`(\w+\-?)?\d+`)

func TaskURL() Part {
	return func(ctx context.Context, h CCHook) (string, error) {
		branch, _ := gitBranchShowCurrent(ctx)
		pr, _ := ghPRViewJSON(ctx)
		taskServer := os.Getenv("CC_TASK_SERVER")
		if taskServer == "" {
			return "", nil
		}
		taskCode, ok := extractTaskCode(branch, pr.HeadRefName, pr.Title)
		if !ok {
			return "", nil
		}
		return style.Underline(taskServer + "/" + taskCode), nil
	}
}

func extractTaskCode(s ...string) (string, bool) {
	for _, s := range s {
		match := taskRegex.FindString(s)
		if match != "" {
			return match, true
		}
	}

	return "", false
}
