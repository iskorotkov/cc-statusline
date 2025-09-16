package transcript

import (
	"iter"
	"time"
)

type TimeModel struct {
	Date  time.Time
	Model string
}

type Usage struct {
	InputTokens      int
	OutputTokens     int
	CacheWriteTokens int
	CacheReadTokens  int
}

func (u Usage) Total() int {
	return u.InputTokens + u.OutputTokens + u.CacheWriteTokens + u.CacheReadTokens
}

func (u *Usage) Add(e EventUsage) {
	u.InputTokens += e.InputTokens
	u.OutputTokens += e.OutputTokens
	u.CacheWriteTokens += e.CacheCreationInputTokens
	u.CacheReadTokens += e.CacheReadInputTokens
}

func DateUsage(transcripts []Transcript, from, to time.Time) map[string]Usage {
	usages := make(map[string]Usage)
	for _, t := range transcripts {
		for e := range deduplicateEvents(t.Events) {
			if e.Timestamp.Before(from) || e.Timestamp.Equal(to) || e.Timestamp.After(to) {
				continue
			}
			usage := usages[e.Message.Model]
			usage.Add(e.Message.Usage)
			usages[e.Message.Model] = usage
		}
	}
	return usages
}

func UsageByDate(transcripts []Transcript) map[TimeModel]Usage {
	usages := make(map[TimeModel]Usage)
	for _, t := range transcripts {
		for e := range deduplicateEvents(t.Events) {
			key := TimeModel{
				Date:  e.Timestamp.Truncate(24 * time.Hour),
				Model: e.Message.Model,
			}
			usage := usages[key]
			usage.InputTokens += e.Message.Usage.InputTokens
			usage.OutputTokens += e.Message.Usage.OutputTokens
			usage.CacheWriteTokens += e.Message.Usage.CacheCreationInputTokens
			usage.CacheReadTokens += e.Message.Usage.CacheReadInputTokens
			usages[key] = usage
		}
	}
	return usages
}

func deduplicateEvents(events []Event) iter.Seq[Event] {
	return func(yield func(Event) bool) {
		seen := make(map[string]bool)
		for _, e := range events {
			if e.Message.ID == "" {
				if !yield(e) {
					return
				}
				continue
			}
			if !seen[e.Message.ID] {
				seen[e.Message.ID] = true
				if !yield(e) {
					return
				}
			}
		}
	}
}
