package transcript_test

import (
	"testing"

	"github.com/iskorotkov/cc-statusline/transcript"
)

func TestSessions(t *testing.T) {
	transcripts, err := transcript.ParseTranscripts()
	if err != nil {
		t.Fatal(err)
	}
	sessions := transcript.Sessions(transcripts)
	visited := make(map[string]bool)
	for _, session := range sessions {
		if visited[session] {
			t.Errorf("duplicate session found: %s", session)
		}
		if session == "" {
			t.Errorf("empty session found")
		}
		visited[session] = true
	}
}

func TestSessionUsage(t *testing.T) {
	transcripts, err := transcript.ParseTranscripts()
	if err != nil {
		t.Fatal(err)
	}
	sessions := transcript.Sessions(transcripts)
	if len(sessions) == 0 {
		t.Skipf("no sessions found")
	}
	for _, session := range sessions {
		usage := transcript.SessionUsage(transcripts, session)
		if len(usage) == 0 {
			t.Errorf("no usage found for session %s", session)
		}
		for model, usage := range usage {
			if model == "" {
				t.Errorf("model is empty for session %s", session)
			}
			if usage == (transcript.Usage{}) {
				t.Errorf("usage for model %s is empty for session %s", model, session)
			}
		}
	}
}
