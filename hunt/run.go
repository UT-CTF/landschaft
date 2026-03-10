package hunt

import (
	"fmt"
	"time"
)

// Run collects events for the given duration, appends to detectionsLogPath, and prints highlights.
func Run(since time.Duration, detectionsLogPath string) error {
	list, err := getDetections(since)
	if err != nil {
		return err
	}
	for i := range list {
		TagSuspicious(&list[i])
	}
	// Append each to JSONL
	for _, d := range list {
		if err := AppendDetection(detectionsLogPath, d); err != nil {
			// non-fatal
			continue
		}
	}
	// Print suspicious or high severity
	fmt.Println("Suspicious / notable events (also appended to", detectionsLogPath, "):")
	for _, d := range list {
		if d.Severity == "high" || d.Severity == "medium" || len(d.Tags) > 0 {
			fmt.Printf("  [%s] %s %s\n", d.Severity, d.EventID, d.Message)
			if d.Explain != "" {
				fmt.Printf("    -> %s\n", d.Explain)
			}
		}
	}
	return nil
}
