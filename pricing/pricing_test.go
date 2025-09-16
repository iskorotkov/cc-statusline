package pricing_test

import (
	"os"
	"testing"

	"github.com/iskorotkov/cc-statusline/pricing"
	"github.com/iskorotkov/cc-statusline/transcript"
)

func TestModelPricing(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skipf("Tests depend on local configuration, skipping in CI")
	}
	transcripts, err := transcript.ParseTranscripts()
	if err != nil {
		t.Fatalf("Failed to parse transcripts: %v", err)
	}
	usages := transcript.UsageByDate(transcripts)
	for k := range usages {
		p, ok := pricing.ModelPricing(k.Model)
		if !ok {
			t.Errorf("No pricing found for model %q", k.Model)
		}
		_ = p
	}
}
