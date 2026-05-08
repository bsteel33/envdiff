package auditor

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// ReportText writes a human-readable audit report to w.
func ReportText(w io.Writer, report Report) {
	events := sortedEvents(report.Events)

	if len(events) == 0 {
		fmt.Fprintln(w, "No audit events recorded.")
		return
	}

	fmt.Fprintf(w, "Audit report — %s\n", report.GeneratedAt.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintln(w, "---")

	for _, ev := range events {
		switch ev.Kind {
		case EventIntroduced:
			fmt.Fprintf(w, "  [+] %-30s introduced (%s)\n", ev.Key, ev.Current.Status)
		case EventResolved:
			fmt.Fprintf(w, "  [-] %-30s resolved\n", ev.Key)
		case EventChanged:
			fmt.Fprintf(w, "  [~] %-30s changed: %s → %s\n", ev.Key, ev.Previous.Status, ev.Current.Status)
		case EventPersisted:
			fmt.Fprintf(w, "  [=] %-30s persisted (%s)\n", ev.Key, ev.Current.Status)
		}
	}
}

// ReportJSON writes a JSON-encoded audit report to w.
func ReportJSON(w io.Writer, report Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}

func sortedEvents(events []Event) []Event {
	out := make([]Event, len(events))
	copy(out, events)
	sort.Slice(out, func(i, j int) bool {
		if out[i].Kind != out[j].Kind {
			return kindOrder(out[i].Kind) < kindOrder(out[j].Kind)
		}
		return out[i].Key < out[j].Key
	})
	return out
}

func kindOrder(k EventKind) int {
	switch k {
		case EventIntroduced:
			return 0
		case EventChanged:
			return 1
		case EventResolved:
			return 2
		default:
			return 3
	}
}
