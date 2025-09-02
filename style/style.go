package style

import "fmt"

func RGB(text string, r, g, b int) string {
	return apply(text, 38, 2, r, g, b)
}

func Normal(text string) string {
	return apply(text, 0)
}

func Bold(text string) string {
	return apply(text, 1)
}

func Dim(text string) string {
	return apply(text, 2)
}

func Italic(text string) string {
	return apply(text, 3)
}

func Underline(text string) string {
	return apply(text, 4)
}

func apply(text string, codes ...int) string {
	if len(codes) == 0 {
		return text
	}
	codeStr := ""
	for i, code := range codes {
		if i > 0 {
			codeStr += ";"
		}
		codeStr += fmt.Sprintf("%d", code)
	}
	return fmt.Sprintf("\x1b[0m\x1b[%sm%s\x1b[0m", codeStr, text)
}
