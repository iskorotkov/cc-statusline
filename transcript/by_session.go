package transcript

import (
	"maps"
	"slices"
)

func Sessions(transcripts []Transcript) []string {
	sessions := make(map[string]bool)
	for _, t := range transcripts {
		for _, e := range t.Events {
			if e.SessionID != "" {
				sessions[e.SessionID] = true
			}
		}
	}
	return slices.Collect(maps.Keys(sessions))
}

func SessionUsage(transcripts []Transcript, sessionID string) map[string]Usage {
	usages := make(map[string]Usage)
	for _, t := range transcripts {
		for e := range deduplicateEvents(t.Events) {
			if e.SessionID != sessionID {
				continue
			}
			usage := usages[e.Message.Model]
			usage.Add(e.Message.Usage)
			usages[e.Message.Model] = usage
		}
	}
	return usages
}
