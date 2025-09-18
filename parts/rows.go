package parts

import (
	"context"
	"strings"

	"github.com/iskorotkov/cc-statusline/style"
)

var (
	rowSeparator  = "\n"
	partSeparator = style.Dim(" / ")
)

func Rows(rows ...Part) Part {
	return func(ctx context.Context, h CCHook) (string, error) {
		if len(rows) == 0 {
			return "", nil
		}
		results := make([]string, 0, len(rows))
		for _, r := range rows {
			s, err := r(ctx, h)
			if err != nil {
				return "", err
			}
			if s != "" {
				results = append(results, s)
			}
		}
		return strings.Join(results, rowSeparator), nil
	}
}

func Row(prefix string, row ...Part) Part {
	return func(ctx context.Context, h CCHook) (string, error) {
		if len(row) == 0 {
			return "", nil
		}
		results := make([]string, 0, len(row))
		for _, c := range row {
			s, err := c(ctx, h)
			if err != nil {
				return "", err
			}
			if s != "" {
				results = append(results, s)
			}
		}
		if len(results) == 0 {
			return "", nil
		}
		return prefix + " " + strings.Join(results, partSeparator), nil
	}
}

