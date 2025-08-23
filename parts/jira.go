package parts

import (
	"context"
	"os"
	"regexp"

	"github.com/iskorotkov/cc-statusline/style"
)

var jiraCodeRegex = regexp.MustCompile(`\w+\-\d+`)

func JiraURL() Part {
	return func(ctx context.Context, h CCHook) (string, error) {
		branch, _ := gitBranchShowCurrent(ctx)
		pr, _ := ghPRViewJSON(ctx)
		url := os.Getenv("CC_JIRA_URL")
		if url == "" {
			return "", nil
		}
		code, ok := extractJiraCode(branch, pr.HeadRefName, pr.Title)
		if !ok {
			return "", nil
		}
		return style.Underline(url + "/" + code), nil
	}
}

func extractJiraCode(s ...string) (string, bool) {
	for _, s := range s {
		match := jiraCodeRegex.FindString(s)
		if match != "" {
			return match, true
		}
	}
	return "", false
}
