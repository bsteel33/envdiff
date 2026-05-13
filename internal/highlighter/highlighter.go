// Package highlighter provides diff-aware value highlighting,
// marking changed, added, or removed segments within env values.
package highlighter

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Style controls how highlights are rendered.
type Style string

const (
	StyleANSI  Style = "ansi"
	StyleMarkdown Style = "markdown"
	StylePlain Style = "plain"
)

// Highlight represents a single highlighted result entry.
type Highlight struct {
	Key    string
	Left   string
	Right  string
	Status diff.Status
}

// Options configures the highlighter.
type Options struct {
	Style Style
}

// Apply returns highlighted representations of the given diff results.
func Apply(results []diff.Result, opts Options) []Highlight {
	out := make([]Highlight, 0, len(results))
	for _, r := range results {
		out = append(out, Highlight{
			Key:    r.Key,
			Left:   renderValue(r.Left, r.Status, sideLeft, opts.Style),
			Right:  renderValue(r.Right, r.Status, sideRight, opts.Style),
			Status: r.Status,
		})
	}
	return out
}

type side int

const (
	sideLeft  side = iota
	sideRight
)

func renderValue(val string, status diff.Status, s side, style Style) string {
	switch status {
	case diff.StatusMatch:
		return val
	case diff.StatusMissingInLeft:
		if s == sideRight {
			return colorize("+"+val, "green", style)
		}
		return colorize("(missing)", "red", style)
	case diff.StatusMissingInRight:
		if s == sideLeft {
			return colorize("-"+val, "red", style)
		}
		return colorize("(missing)", "yellow", style)
	case diff.StatusMismatch:
		return colorize(val, "yellow", style)
	}
	return val
}

func colorize(text, color string, style Style) string {
	switch style {
	case StyleANSI:
		codes := map[string]string{"red": "31", "green": "32", "yellow": "33"}
		if code, ok := codes[color]; ok {
			return fmt.Sprintf("\033[%sm%s\033[0m", code, text)
		}
	case StyleMarkdown:
		return fmt.Sprintf("**%s**", text)
	case StylePlain:
		return fmt.Sprintf("[%s]", strings.ToUpper(color)+":"+text)
	}
	return text
}
