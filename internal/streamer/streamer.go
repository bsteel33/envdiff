// Package streamer provides line-by-line streaming comparison of .env files,
// emitting diff results incrementally rather than loading everything into memory.
package streamer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Event is emitted for each key encountered during streaming.
type Event struct {
	Result diff.Result
	Err    error
}

// Options controls streamer behaviour.
type Options struct {
	// BufferSize sets the scanner buffer size in bytes (default: 64 KiB).
	BufferSize int
}

var defaultOptions = Options{BufferSize: 64 * 1024}

// Stream reads leftPath and rightPath line by line and sends diff Events on ch.
// The channel is closed when both files have been fully processed.
func Stream(leftPath, rightPath string, opts *Options, ch chan<- Event) {
	if opts == nil {
		opts = &defaultOptions
	}

	go func() {
		defer close(ch)

		left, err := parseStream(leftPath, opts.BufferSize)
		if err != nil {
			ch <- Event{Err: fmt.Errorf("left file: %w", err)}
			return
		}

		right, err := parseStream(rightPath, opts.BufferSize)
		if err != nil {
			ch <- Event{Err: fmt.Errorf("right file: %w", err)}
			return
		}

		results := diff.Compare(left, right)
		for _, r := range results {
			ch <- Event{Result: r}
		}
	}()
}

// parseStream reads key=value pairs from a file path into a map.
func parseStream(path string, bufSize int) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return scanEnv(f, bufSize)
}

// scanEnv reads key=value lines from r, skipping comments and blank lines.
func scanEnv(r io.Reader, bufSize int) (map[string]string, error) {
	env := make(map[string]string)
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, bufSize), bufSize)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 1 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.Trim(strings.TrimSpace(line[idx+1:]), `"'`)
		env[key] = val
	}
	return env, scanner.Err()
}
