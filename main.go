package main

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := run(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	var hook Hook
	if err := json.NewDecoder(os.Stdin).Decode(&hook); err != nil {
		return err
	}

	dir, err := filepath.Rel(
		filepath.Dir(hook.Workspace.ProjectDir),
		hook.Workspace.CurrentDir,
	)
	if err != nil {
		return err
	}

	branch, _ := shellString(ctx, "git", "branch", "--show-current")

	files, _ := shellString(ctx, "git", "status", "--porcelain")
	fileCount := count(files)

	var taskReferences []string

	var gitFormat string
	if branch != "" {
		var sb strings.Builder
		if len(fileCount) > 0 {
			sb.WriteString(" /")
		}

		for _, p := range fileCount {
			sb.WriteString(fmt.Sprintf(" %s:%d", p.k, p.v))
		}

		gitFormat = fmt.Sprintf(
			"\nGIT %s%s",
			limit(branch, 60),
			sb.String(),
		)

		taskReferences = append(taskReferences, branch)
	}

	pr, _ := shellJSON[GHPR](
		ctx,
		"gh",
		"pr",
		"view",
		"--json",
		"number,url,title,mergeable,additions,deletions,changedFiles,baseRefName,headRefName",
	)

	var prFormat string
	if pr != (GHPR{}) {
		var mergeableFormat string
		if pr.Mergeable {
			mergeableFormat = "M"
		} else {
			mergeableFormat = "NM"
		}

		prFormat = fmt.Sprintf(
			"\nPR #%d / %s / +%dL -%dL ~%dF %s\nPR %s",
			pr.Number,
			limit(pr.Title, 60),
			pr.Additions,
			pr.Deletions,
			pr.ChangedFiles,
			mergeableFormat,
			pr.URL,
		)

		taskReferences = append(taskReferences, pr.Title)
	}

	var taskFormat string
	if taskServer := os.Getenv("CC_TASK_SERVER"); taskServer != "" && len(taskReferences) > 0 {
		taskCode, ok := extractTaskCode(taskReferences...)
		if ok {
			taskFormat = fmt.Sprintf("\nTASK %s", buildTaskLink(taskServer, taskCode))
		}
	}

	duration := time.Millisecond * time.Duration(hook.Cost.TotalAPIDurationMS)

	var badge200k string
	if hook.Exceeds200KTokens {
		badge200k = " / 200K+"
	}

	fmt.Printf(
		"CC v%s / %s / %s / %s / +%dL -%dL %.1fm $%.1f%s%s%s%s",
		hook.Version,
		hook.Model.DisplayName,
		hook.OutputStyle.Name,
		limit(dir, 20),
		hook.Cost.TotalLinesAdded,
		hook.Cost.TotalLinesRemoved,
		duration.Minutes(),
		hook.Cost.TotalCostUSD,
		badge200k,
		gitFormat,
		prFormat,
		taskFormat,
	)

	return nil
}

type pair struct {
	k string
	v int
}

func count(s string) []pair {
	counts := make(map[string]int)
	for line := range strings.Lines(s) {
		key, _, ok := strings.Cut(strings.TrimSpace(line), " ")
		if ok {
			counts[key]++
		}
	}

	pairs := make([]pair, 0, len(counts))
	for k, v := range counts {
		pairs = append(pairs, pair{k: k, v: v})
	}

	slices.SortStableFunc(pairs, func(p1, p2 pair) int {
		return cmp.Compare(p1.v, p2.v)
	})

	return pairs
}

func limit(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}

func shellString(ctx context.Context, s ...string) (string, error) {
	cmd := exec.CommandContext(ctx, s[0], s[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w: \n\n %s", err, strings.TrimSpace(string(output)))
	}

	return strings.TrimSpace(string(output)), nil
}

func shellJSON[T any](ctx context.Context, s ...string) (T, error) {
	var result T

	cmd := exec.CommandContext(ctx, s[0], s[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return result, fmt.Errorf("%w: \n\n %s", err, strings.TrimSpace(string(output)))
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return result, err
	}

	return result, nil
}

func buildTaskLink(server, task string) string {
	return server + "/" + task
}

var taskRegex = regexp.MustCompile(`(\w+\-?|#)?\d+`)

func extractTaskCode(s ...string) (string, bool) {
	for _, s := range s {
		match := taskRegex.FindString(s)
		if match != "" {
			return match, true
		}
	}

	return "", false
}

type Hook struct {
	SessionID      string `json:"session_id"`
	TranscriptPath string `json:"transcript_path"`
	CWD            string `json:"cwd"`
	Model          struct {
		ID          string `json:"id"`
		DisplayName string `json:"display_name"`
	} `json:"model"`
	Workspace struct {
		CurrentDir string `json:"current_dir"`
		ProjectDir string `json:"project_dir"`
	} `json:"workspace"`
	Version     string `json:"version"`
	OutputStyle struct {
		Name string `json:"name"`
	} `json:"output_style"`
	Cost struct {
		TotalCostUSD       float64 `json:"total_cost_usd"`
		TotalDurationMS    int     `json:"total_duration_ms"`
		TotalAPIDurationMS int     `json:"total_api_duration_ms"`
		TotalLinesAdded    int     `json:"total_lines_added"`
		TotalLinesRemoved  int     `json:"total_lines_removed"`
	} `json:"cost"`
	Exceeds200KTokens bool `json:"exceeds_200k_tokens"`
}

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
