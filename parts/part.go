package parts

import (
	"context"
)

type Part func(ctx context.Context, h CCHook) (string, error)
