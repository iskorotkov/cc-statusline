package style

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type theme int

const (
	themeDark theme = iota
	themeLight
)

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

func Blue(text string) string {
	return themeBlue(text)
}

func Red(text string) string {
	return themeRed(text)
}

func Green(text string) string {
	return themeGreen(text)
}

func themeBlue(text string) string {
	theme := detectTheme()
	switch theme {
	case themeLight:
		return rgb(text, 0, 64, 160)
	default:
		return rgb(text, 192, 192, 255)
	}
}

func themeRed(text string) string {
	theme := detectTheme()
	switch theme {
	case themeLight:
		return rgb(text, 180, 0, 0)
	default:
		return rgb(text, 255, 127, 127)
	}
}

func themeGreen(text string) string {
	theme := detectTheme()
	switch theme {
	case themeLight:
		return rgb(text, 0, 128, 0)
	default:
		return rgb(text, 127, 255, 127)
	}
}

func detectTheme() theme {
	if ccTheme := os.Getenv("CC_THEME"); ccTheme != "" {
		switch strings.ToLower(ccTheme) {
		case "light":
			return themeLight
		case "dark":
			return themeDark
		case "auto":
			return detectTerminalTheme()
		default:
			return themeDark
		}
	}
	return detectTerminalTheme()
}

func detectTerminalTheme() theme {
	colorFgBg := os.Getenv("COLORFGBG")
	if colorFgBg == "" {
		return themeDark
	}
	parts := strings.Split(colorFgBg, ";")
	if len(parts) < 2 {
		return themeDark
	}
	bg, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return themeDark
	}
	if bg >= 8 {
		return themeLight
	}
	return themeDark
}

func rgb(text string, r, g, b int) string {
	return apply(text, 38, 2, r, g, b)
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
