package parts

import (
	"cmp"
	"context"
	"slices"
	"strings"
)

func Fixed(s string) Part {
	return func(ctx context.Context, h CCHook) (string, error) {
		return s, nil
	}
}

type pair struct {
	k string
	v int
}

func count(s string) []pair {
	counts := make(map[string]int)
	for line := range strings.Lines(s) {
		key, _, ok := strings.Cut(strings.TrimSpace(line), " ")
		if ok {
			counts[key]++
		}
	}
	pairs := make([]pair, 0, len(counts))
	for k, v := range counts {
		pairs = append(pairs, pair{k: k, v: v})
	}
	slices.SortStableFunc(pairs, func(p1, p2 pair) int {
		return cmp.Compare(p1.k, p2.k)
	})
	return pairs
}

func limit(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return strings.TrimSpace(s[:n-3]) + "..."
}
