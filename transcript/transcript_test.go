package transcript_test

import (
	"os"
	"testing"

	"github.com/iskorotkov/cc-statusline/transcript"
)

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}

func TestParseTranscripts(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skipf("Tests depend on local configuration, skipping in CI")
	}
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
