package score

import (
	"fmt"
	"strings"
)

// RunList discovers listeners and prints them (text).
func RunList() error {
	listeners, err := discoverListeners()
	if err != nil {
		return err
	}
	for _, l := range listeners {
		exp := l.Explain
		if exp == "" {
			exp = "?"
		}
		fmt.Printf("%d/tcp %s %s %s\n", l.Port, l.Bind, l.Process, exp)
	}
	return nil
}

// RunExplain discovers listeners and prints with explanations.
func RunExplain() error {
	listeners, err := discoverListeners()
	if err != nil {
		return err
	}
	fmt.Println("Candidate scored services (auto-discovered listeners). Scoring may check different endpoints.")
	fmt.Println(strings.Repeat("-", 60))
	for _, l := range listeners {
		exp := l.Explain
		if exp == "" {
			exp = "(unknown)"
		}
		fmt.Printf("Port %d/%s bind=%s process=%s -> %s\n", l.Port, l.Proto, l.Bind, l.Process, exp)
	}
	return nil
}
