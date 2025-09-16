package transcript_test

import (
	"testing"
	"time"

	"github.com/iskorotkov/cc-statusline/transcript"
)

func TestDateUsage(t *testing.T) {
	transcripts, err := transcript.ParseTranscripts()
	if err != nil {
		t.Fatal(err)
	}
	sessions := transcript.UsageByDate(transcripts)
	if len(sessions) == 0 {
		t.Skipf("no sessions found")
	}
	for session := range sessions {
		usage := transcript.DateUsage(transcripts, session.Date, session.Date.Add(24*time.Hour))
		if len(usage) == 0 {
			t.Errorf("no usage found for date %s", session.Date)
		}
		for model, usage := range usage {
			if model == "" {
				t.Errorf("model is empty for date %s", session.Date)
			}
			if usage == (transcript.Usage{}) {
				t.Errorf("usage for model %s is empty for date %s", model, session.Date)
			}
		}
	}
}

func TestUsageByDate(t *testing.T) {
	transcripts, err := transcript.ParseTranscripts()
	if err != nil {
		t.Fatalf("ParseTranscripts() error: %v", err)
	}
	usages := transcript.UsageByDate(transcripts)
	for k, usage := range usages {
		if !k.Date.Equal(k.Date.Truncate(24 * time.Hour)) {
			t.Errorf("date %v is not truncated to day", k.Date)
		}
		if k.Date.Before(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)) {
			t.Errorf("date %v is before 2023-01-01", k.Date)
		}
		if k.Date.After(time.Now()) {
			t.Errorf("date %v is in the future", k.Date)
		}
		if k.Model == "" {
			t.Errorf("model is empty for date %v", k.Date)
		}
		if usage == (transcript.Usage{}) {
			t.Errorf("usage for date %v is empty", k.Date)
		}
	}
}
