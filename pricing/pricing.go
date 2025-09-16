package pricing

var pricingByModel = map[string]Pricing{
	"claude-opus-4-1-20250805": {
		InputTokens:      15e-6,
		OutputTokens:     75e-6,
		CacheWriteTokens: 18.75e-6,
		CacheReadTokens:  1.5e-6,
	},
	"claude-sonnet-4-20250514": {
		InputTokens:      3e-6,
		OutputTokens:     15e-6,
		CacheWriteTokens: 3.75e-6,
		CacheReadTokens:  0.3e-6,
	},
}

type Pricing struct {
	InputTokens      float64
	OutputTokens     float64
	CacheWriteTokens float64
	CacheReadTokens  float64
}

func ModelPricing(model string) (Pricing, bool) {
	p, ok := pricingByModel[model]
	return p, ok
}
