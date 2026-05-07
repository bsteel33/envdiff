package templater_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/templater"
)

func TestGenerate_BasicOutput(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"PORT":     "8080",
	}
	var sb strings.Builder
	if err := templater.Generate(&sb, env, templater.Options{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "APP_NAME=<value>") {
		t.Errorf("expected APP_NAME=<value>, got:\n%s", out)
	}
	if !strings.Contains(out, "PORT=<value>") {
		t.Errorf("expected PORT=<value>, got:\n%s", out)
	}
}

func TestGenerate_SensitiveKeysRedacted(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "hunter2",
		"API_TOKEN":   "abc123",
		"APP_ENV":     "production",
	}
	var sb strings.Builder
	if err := templater.Generate(&sb, env, templater.Options{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "DB_PASSWORD=<secret>") {
		t.Errorf("expected DB_PASSWORD=<secret>, got:\n%s", out)
	}
	if !strings.Contains(out, "API_TOKEN=<secret>") {
		t.Errorf("expected API_TOKEN=<secret>, got:\n%s", out)
	}
	if !strings.Contains(out, "APP_ENV=<value>") {
		t.Errorf("expected APP_ENV=<value>, got:\n%s", out)
	}
}

func TestGenerate_CustomPlaceholder(t *testing.T) {
	env := map[string]string{"HOST": "localhost"}
	var sb strings.Builder
	opts := templater.Options{Placeholder: "CHANGEME"}
	if err := templater.Generate(&sb, env, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "HOST=CHANGEME") {
		t.Errorf("expected HOST=CHANGEME, got: %s", sb.String())
	}
}

func TestGenerate_CommentPrefix(t *testing.T) {
	env := map[string]string{"LOG_LEVEL": "debug"}
	var sb strings.Builder
	opts := templater.Options{CommentPrefix: "Set "}
	if err := templater.Generate(&sb, env, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "# Set LOG_LEVEL") {
		t.Errorf("expected comment line, got:\n%s", out)
	}
}

func TestGenerate_EmptyEnv(t *testing.T) {
	var sb strings.Builder
	if err := templater.Generate(&sb, map[string]string{}, templater.Options{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sb.Len() != 0 {
		t.Errorf("expected empty output, got: %q", sb.String())
	}
}

func TestGenerate_SortedOutput(t *testing.T) {
	env := map[string]string{
		"ZEBRA": "1",
		"ALPHA": "2",
		"MANGO": "3",
	}
	var sb strings.Builder
	if err := templater.Generate(&sb, env, templater.Options{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(sb.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "ALPHA") {
		t.Errorf("expected first line to be ALPHA, got: %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "ZEBRA") {
		t.Errorf("expected last line to be ZEBRA, got: %s", lines[2])
	}
}
