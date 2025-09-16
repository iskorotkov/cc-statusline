package transcript_test

import (
	"testing"

	"github.com/iskorotkov/cc-statusline/transcript"
)

func TestParseTranscripts(t *testing.T) {
	transcripts, err := transcript.ParseTranscripts()
	if err != nil {
		t.Fatalf("ParseTranscripts() error: %v", err)
	}
	for _, tr := range transcripts {
		if tr.File == "" {
			t.Errorf("transcript file is empty")
		}
		if len(tr.Events) == 0 {
			t.Errorf("transcript %q has no events", tr.File)
		}
		for _, e := range tr.Events {
			if e.Message.Usage == (transcript.EventUsage{}) {
				t.Errorf("event in transcript %q has empty usage", tr.File)
			}
		}
	}
}
