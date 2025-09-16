package transcript

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

type Event struct {
	SessionID string       `json:"sessionId"`
	Timestamp time.Time    `json:"timestamp"`
	Type      string       `json:"type"`
	Cwd       string       `json:"cwd"`
	GitBranch string       `json:"gitBranch"`
	Message   EventMessage `json:"message"`
}

type EventMessage struct {
	ID    string     `json:"id"`
	Type  string     `json:"type"`
	Role  string     `json:"role"`
	Model string     `json:"model"`
	Usage EventUsage `json:"usage"`
}

type EventUsage struct {
	InputTokens              int `json:"input_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens"`
	OutputTokens             int `json:"output_tokens"`
}

type Transcript struct {
	File   string
	Events []Event
}

func ParseTranscripts() ([]Transcript, error) {
	path, err := transcriptPath()
	if err != nil {
		return nil, fmt.Errorf("get transcript path: %w", err)
	}
	dir := os.DirFS(path)
	var transcripts []Transcript
	if err := fs.WalkDir(dir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("access path %q: %w", path, err)
		}
		if d.IsDir() || filepath.Ext(path) != ".jsonl" {
			return nil
		}
		f, err := dir.Open(path)
		if err != nil {
			return fmt.Errorf("open file %q: %w", path, err)
		}
		defer func() {
			_ = f.Close()
		}()
		events, err := parseEvents(f)
		if err != nil {
			return fmt.Errorf("parse events in file %q: %w", path, err)
		}
		if len(events) > 0 {
			transcripts = append(transcripts, Transcript{
				File:   path,
				Events: events,
			})
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("walk dir %q: %w", dir, err)
	}
	return transcripts, nil
}

func parseEvents(f io.Reader) ([]Event, error) {
	dec := json.NewDecoder(bufio.NewReader(f))
	var events []Event
	for {
		var e Event
		if err := dec.Decode(&e); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("decode JSON: %w", err)
		}
		if e.Message.Usage == (EventUsage{}) {
			continue
		}
		events = append(events, e)
	}
	return events, nil
}

func transcriptPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get user home dir: %w", err)
	}
	return filepath.Join(home, ".claude", "projects"), nil
}
