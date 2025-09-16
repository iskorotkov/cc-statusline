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
		var sb strings.Builder
		s, err := rows[0](ctx, h)
		if err != nil {
			return "", err
		}
		sb.WriteString(s)
		for _, r := range rows[1:] {
			s, err := printWithSeparator(ctx, h, r, rowSeparator)
			if err != nil {
				return "", err
			}
			sb.WriteString(s)
		}
		return sb.String(), nil
	}
}

func Row(prefix string, row ...Part) Part {
	return func(ctx context.Context, h CCHook) (string, error) {
		if len(row) == 0 {
			return "", nil
		}
		var sb strings.Builder
		sb.WriteString(prefix)
		for _, c := range row {
			s, err := printWithSeparator(ctx, h, c, partSeparator)
			if err != nil {
				return "", err
			}
			sb.WriteString(s)
		}
		if sb.Len() <= len(prefix) {
			return "", nil
		}
		return sb.String(), nil
	}
}

func printWithSeparator(ctx context.Context, h CCHook, p Part, separator string) (string, error) {
	s, err := p(ctx, h)
	if err != nil {
		return "", err
	}
	if s == "" {
		return "", nil
	}
	return separator + s, nil
}
