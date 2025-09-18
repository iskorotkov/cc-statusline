package parts

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/iskorotkov/cc-statusline/pricing"
	"github.com/iskorotkov/cc-statusline/style"
	"github.com/iskorotkov/cc-statusline/transcript"
)

var parsedTranscripts = func() func(ctx context.Context) ([]transcript.Transcript, error) {
	var transcripts []transcript.Transcript
	var err error
	var once sync.Once
	return func(ctx context.Context) ([]transcript.Transcript, error) {
		once.Do(func() {
			transcripts, err = transcript.ParseTranscripts()
		})
		return transcripts, err
	}
}()

func CCSessionUsage() Part {
	return func(ctx context.Context, h CCHook) (string, error) {
		transcripts, err := parsedTranscripts(ctx)
		if err != nil {
			return "", err
		}
		usage := transcript.SessionUsage(transcripts, h.SessionID)
		return formatUsage("session", usage), nil
	}
}

func CCHourUsage() Part {
	return func(ctx context.Context, h CCHook) (string, error) {
		transcripts, err := parsedTranscripts(ctx)
		if err != nil {
			return "", err
		}
		from := time.Now().Truncate(time.Hour)
		to := from.Add(time.Hour)
		usage := transcript.DateUsage(transcripts, from, to)
		return formatUsage("hour", usage), nil
	}
}

func CCDayUsage() Part {
	return func(ctx context.Context, h CCHook) (string, error) {
		transcripts, err := parsedTranscripts(ctx)
		if err != nil {
			return "", err
		}
		from := time.Now().Truncate(24 * time.Hour)
		to := from.Add(24 * time.Hour)
		usage := transcript.DateUsage(transcripts, from, to)
		return formatUsage("day", usage), nil
	}
}

func CCWeekUsage() Part {
	return func(ctx context.Context, h CCHook) (string, error) {
		transcripts, err := parsedTranscripts(ctx)
		if err != nil {
			return "", err
		}
		to := time.Now()
		from := to.Add(-7 * 24 * time.Hour)
		usage := transcript.DateUsage(transcripts, from, to)
		return formatUsage("week", usage), nil
	}
}

func formatUsage(title string, usage map[string]transcript.Usage) string {
	var combinedTokens int
	var combinedPrice float64
	for model, usage := range usage {
		combinedTokens += usage.Total()
		price, ok := pricing.ModelPricing(model)
		if ok {
			combinedPrice += totalPrice(usage, price)
		}
	}
	return fmt.Sprintf("%s %s%s",
		title,
		formatTokens(combinedTokens),
		fmt.Sprintf(style.Green(" $%.1f"), combinedPrice))
}

func formatTokens(tokens int) string {
	if tokens < 1000 {
		return fmt.Sprintf("%dt", tokens)
	}
	if tokens < 1000_000 {
		return fmt.Sprintf("%.1fKt", float64(tokens)/1000)
	}
	if tokens < 1000_000_000 {
		return fmt.Sprintf("%.1fMt", float64(tokens)/1000_000)
	}
	return fmt.Sprintf("%.1fBt", float64(tokens)/1000_000_000)
}

func totalPrice(usage transcript.Usage, price pricing.Pricing) float64 {
	return float64(usage.InputTokens)*price.InputTokens +
		float64(usage.OutputTokens)*price.OutputTokens +
		float64(usage.CacheWriteTokens)*price.CacheWriteTokens +
		float64(usage.CacheReadTokens)*price.CacheReadTokens
}
