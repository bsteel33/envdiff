package colorizer_test

import (
	"strings"
	"testing"

	"github.com/envdiff/envdiff/internal/colorizer"
	"github.com/envdiff/envdiff/internal/diff"
)

func TestForStatus_NoColor_ReturnsPlain(t *testing.T) {
	c := colorizer.New(colorizer.Options{NoColor: true})
	result := c.ForStatus("hello", diff.StatusMismatch)
	if result != "hello" {
		t.Errorf("expected plain string, got %q", result)
	}
}

func TestForStatus_MissingInRight_ContainsRed(t *testing.T) {
	c := colorizer.New(colorizer.Options{})
	result := c.ForStatus("KEY", diff.StatusMissingInRight)
	if !strings.Contains(result, "\033[31m") {
		t.Errorf("expected red ANSI code in %q", result)
	}
	if !strings.Contains(result, "KEY") {
		t.Errorf("expected original text in %q", result)
	}
}

func TestForStatus_MissingInLeft_ContainsGreen(t *testing.T) {
	c := colorizer.New(colorizer.Options{})
	result := c.ForStatus("KEY", diff.StatusMissingInLeft)
	if !strings.Contains(result, "\033[32m") {
		t.Errorf("expected green ANSI code in %q", result)
	}
}

func TestForStatus_Mismatch_ContainsYellow(t *testing.T) {
	c := colorizer.New(colorizer.Options{})
	result := c.ForStatus("KEY", diff.StatusMismatch)
	if !strings.Contains(result, "\033[33m") {
		t.Errorf("expected yellow ANSI code in %q", result)
	}
}

func TestForStatus_Match_ContainsCyan(t *testing.T) {
	c := colorizer.New(colorizer.Options{})
	result := c.ForStatus("KEY", diff.StatusMatch)
	if !strings.Contains(result, "\033[36m") {
		t.Errorf("expected cyan ANSI code in %q", result)
	}
}

func TestBold_NoColor_ReturnsPlain(t *testing.T) {
	c := colorizer.New(colorizer.Options{NoColor: true})
	if got := c.Bold("title"); got != "title" {
		t.Errorf("expected plain, got %q", got)
	}
}

func TestBold_WithColor_ContainsBoldCode(t *testing.T) {
	c := colorizer.New(colorizer.Options{})
	result := c.Bold("title")
	if !strings.Contains(result, "\033[1m") {
		t.Errorf("expected bold ANSI code in %q", result)
	}
}

func TestLabel_NoColor_PlainText(t *testing.T) {
	c := colorizer.New(colorizer.Options{NoColor: true})
	if got := c.Label(diff.StatusMismatch); got != "MISMATCH" {
		t.Errorf("expected MISMATCH, got %q", got)
	}
}

func TestLabel_UnknownStatus(t *testing.T) {
	c := colorizer.New(colorizer.Options{NoColor: true})
	if got := c.Label(diff.Status("weird")); got != "UNKNOWN" {
		t.Errorf("expected UNKNOWN, got %q", got)
	}
}
