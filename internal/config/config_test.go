package config_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/config"
)

func TestValidate_TooFewFiles(t *testing.T) {
	cfg := &config.Config{Files: []string{"a.env"}}
	if err := config.Validate(cfg); err == nil {
		t.Fatal("expected error for fewer than 2 files")
	}
}

func TestValidate_DefaultsFormatToText(t *testing.T) {
	cfg := &config.Config{Files: []string{"a.env", "b.env"}}
	if err := config.Validate(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OutputFormat != config.FormatText {
		t.Errorf("expected default format %q, got %q", config.FormatText, cfg.OutputFormat)
	}
}

func TestValidate_ExplicitJSONFormat(t *testing.T) {
	cfg := &config.Config{
		Files:        []string{"a.env", "b.env"},
		OutputFormat: config.FormatJSON,
	}
	if err := config.Validate(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_UnknownFormat(t *testing.T) {
	cfg := &config.Config{
		Files:        []string{"a.env", "b.env"},
		OutputFormat: config.Format("xml"),
	}
	if err := config.Validate(cfg); err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct {
		input    string
		want     config.Format
	}{
		{"text", config.FormatText},
		{"TEXT", config.FormatText},
		{"json", config.FormatJSON},
		{"JSON", config.FormatJSON},
		{"env", config.FormatEnv},
		{"ENV", config.FormatEnv},
	}
	for _, tc := range cases {
		got, err := config.ParseFormat(tc.input)
		if err != nil {
			t.Errorf("ParseFormat(%q) unexpected error: %v", tc.input, err)
			continue
		}
		if got != tc.want {
			t.Errorf("ParseFormat(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	if _, err := config.ParseFormat("yaml"); err == nil {
		t.Fatal("expected error for unknown format 'yaml'")
	}
}
