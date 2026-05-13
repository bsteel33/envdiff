package tagger_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/tagger"
)

func sampleTagged() []tagger.TaggedResult {
	return []tagger.TaggedResult{
		{
			Result: diff.Result{Key: "API_TOKEN", Status: diff.StatusMismatch},
			Tags:   []tagger.Tag{tagger.TagSecret},
		},
		{
			Result: diff.Result{Key: "APP_ENV", Status: diff.StatusMatch},
			Tags:   []tagger.Tag{tagger.TagUnknown},
		},
	}
}

func TestReportText_NoTagged(t *testing.T) {
	var buf bytes.Buffer
	tagger.ReportText(&buf, nil)
	if !strings.Contains(buf.String(), "No tagged results") {
		t.Errorf("expected no-results message, got: %s", buf.String())
	}
}

func TestReportText_WithTagged(t *testing.T) {
	var buf bytes.Buffer
	tagger.ReportText(&buf, sampleTagged())
	out := buf.String()
	if !strings.Contains(out, "API_TOKEN") {
		t.Error("expected API_TOKEN in output")
	}
	if !strings.Contains(out, "secret") {
		t.Error("expected 'secret' tag in output")
	}
	if !strings.Contains(out, "APP_ENV") {
		t.Error("expected APP_ENV in output")
	}
}

func TestReportJSON_Structure(t *testing.T) {
	var buf bytes.Buffer
	err := tagger.ReportJSON(&buf, sampleTagged())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0]["key"] != "API_TOKEN" {
		t.Errorf("expected key API_TOKEN, got %v", out[0]["key"])
	}
	tags, ok := out[0]["tags"].([]interface{})
	if !ok || len(tags) == 0 {
		t.Error("expected non-empty tags array")
	}
}
