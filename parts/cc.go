package parts

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/iskorotkov/cc-statusline/style"
)

type CCHook struct {
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

func CCVersion() Part {
	return func(ctx context.Context, h CCHook) (string, error) {
		return fmt.Sprintf(style.Dim("v%s"), h.Version), nil
	}
}

func CCModel() Part {
	return func(ctx context.Context, h CCHook) (string, error) {
		return style.Bold(h.Model.DisplayName), nil
	}
}

func CCOutputStyle() Part {
	return func(ctx context.Context, h CCHook) (string, error) {
		return style.Dim(h.OutputStyle.Name), nil
	}
}

func CCDir() Part {
	return func(ctx context.Context, h CCHook) (string, error) {
		dir, err := filepath.Rel(filepath.Dir(h.Workspace.ProjectDir), h.Workspace.CurrentDir)
		if err != nil {
			return "", err
		}
		return style.Italic(limit(dir, 20)), nil
	}
}

func CCStats() Part {
	return func(ctx context.Context, h CCHook) (string, error) {
		duration := time.Millisecond * time.Duration(h.Cost.TotalAPIDurationMS)
		return fmt.Sprintf("%s %s %.1fm %s",
			fmt.Sprintf(style.RGB("+%dL", 127, 255, 127), h.Cost.TotalLinesAdded),
			fmt.Sprintf(style.RGB("-%dL", 255, 127, 127), h.Cost.TotalLinesRemoved),
			duration.Minutes(),
			fmt.Sprintf(style.RGB("$%.1f", 127, 255, 127), h.Cost.TotalCostUSD)), nil
	}
}

func CC200KContextBadge() Part {
	return func(ctx context.Context, h CCHook) (string, error) {
		if h.Exceeds200KTokens {
			return style.Bold("200K+"), nil
		}
		return "", nil
	}
}

func CCTranscriptPath() Part {
	return func(ctx context.Context, h CCHook) (string, error) {
		return style.Dim(h.TranscriptPath), nil
	}
}
