package parts

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/iskorotkov/cc-statusline/shell"
	"github.com/iskorotkov/cc-statusline/style"
)

var ghPRViewJSON = func() func(ctx context.Context) (GHPR, error) {
	var pr GHPR
	var err error
	var once sync.Once
	return func(ctx context.Context) (GHPR, error) {
		once.Do(func() {
			pr, err = shell.JSON[GHPR](
				ctx,
				"gh",
				"pr",
				"view",
				"--json",
				"number,url,title,mergeable,additions,deletions,changedFiles,baseRefName,headRefName",
			)
		})
		return pr, err
	}
}()

type GHPR struct {
	Number       int    `json:"number"`
	URL          string `json:"url"`
	Title        string `json:"title"`
	Mergeable    bool   `json:"mergeable"`
	Additions    int    `json:"additions"`
	Deletions    int    `json:"deletions"`
	ChangedFiles int    `json:"changedFiles"`
	BaseRefName  string `json:"baseRefName"`
	HeadRefName  string `json:"headRefName"`
}

func GHPRNumber() Part {
	return func(ctx context.Context, h CCHook) (string, error) {
		pr, _ := ghPRViewJSON(ctx)
		if pr == (GHPR{}) {
			return "", nil
		}
		return fmt.Sprintf(style.Bold("#%d"), pr.Number), nil
	}
}

func GHPRTitle() Part {
	return func(ctx context.Context, h CCHook) (string, error) {
		pr, _ := ghPRViewJSON(ctx)
		if pr == (GHPR{}) {
			return "", nil
		}
		return style.Italic(limit(pr.Title, 60)), nil
	}
}

func GHPRStats() Part {
	return func(ctx context.Context, h CCHook) (string, error) {
		pr, _ := ghPRViewJSON(ctx)
		if pr == (GHPR{}) {
			return "", nil
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, style.RGB("+%dL", 127, 255, 127), pr.Additions)
		sb.WriteString(" ")
		fmt.Fprintf(&sb, style.RGB("-%dL", 255, 127, 127), pr.Deletions)
		sb.WriteString(" ")
		fmt.Fprintf(&sb, "~%dF", pr.ChangedFiles)
		sb.WriteString(" ")
		if pr.Mergeable {
			fmt.Fprint(&sb, style.RGB("M", 127, 255, 127))
		} else {
			fmt.Fprint(&sb, style.RGB("NM", 255, 127, 127))
		}
		return sb.String(), nil
	}
}

func GHPRURL() Part {
	return func(ctx context.Context, h CCHook) (string, error) {
		pr, _ := ghPRViewJSON(ctx)
		if pr == (GHPR{}) {
			return "", nil
		}
		return style.Underline(pr.URL), nil
	}
}
