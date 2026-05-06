package exporter_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/exporter"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Status: diff.Missing, LeftValue: "localhost", RightValue: ""},
		{Key: "API_KEY", Status: diff.Mismatched, LeftValue: "abc", RightValue: "xyz"},
		{Key: "PORT", Status: diff.Extra, LeftValue: "", RightValue: "8080"},
	}
}

func TestExport_EnvFormat_Stdout(t *testing.T) {
	// Just ensure no error is returned when writing to stdout (path="")
	err := exporter.Export(sampleResults(), "", exporter.Options{Format: exporter.FormatEnv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExport_EnvFormat_ToFile(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "out.env")
	err := exporter.Export(sampleResults(), tmp, exporter.Options{Format: exporter.FormatEnv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(tmp)
	content := string(data)

	if !strings.Contains(content, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST=localhost in output, got:\n%s", content)
	}
	if !strings.Contains(content, "API_KEY=abc") {
		t.Errorf("expected API_KEY=abc in output, got:\n%s", content)
	}
}

func TestExport_EnvFormat_OnlyMissing(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "out.env")
	err := exporter.Export(sampleResults(), tmp, exporter.Options{
		Format:      exporter.FormatEnv,
		OnlyMissing: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(tmp)
	content := string(data)

	if !strings.Contains(content, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST in output")
	}
	if strings.Contains(content, "API_KEY") {
		t.Errorf("did not expect API_KEY (mismatched) in only-missing output")
	}
}

func TestExport_JSONFormat_ToFile(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "out.json")
	err := exporter.Export(sampleResults(), tmp, exporter.Options{Format: exporter.FormatJSON})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(tmp)
	var entries []map[string]string
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}

	found := false
	for _, e := range entries {
		if e["key"] == "API_KEY" && e["status"] == "mismatched" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected API_KEY mismatched entry in JSON output")
	}
}
