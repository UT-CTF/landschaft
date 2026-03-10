package util

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// ActionLogEntry is one JSONL record for landschaft invocations.
type ActionLogEntry struct {
	Timestamp string   `json:"timestamp"`
	Hostname  string   `json:"hostname"`
	User      string   `json:"user"`
	GOOS      string   `json:"goos"`
	Command   string   `json:"command"`
	Args      []string `json:"args"`
	ExitCode  int      `json:"exit_code"`
	DurationMs int64   `json:"duration_ms"`
}

// AppendActionLog appends a single JSONL line to path. Creates parent dirs if needed. Does not fail the process on error.
func AppendActionLog(path string, entry ActionLogEntry) {
	if path == "" {
		return
	}
	dir := filepath.Dir(path)
	_ = os.MkdirAll(dir, 0755)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(entry)
}

// NewActionLogEntry builds an entry from argv and metadata.
func NewActionLogEntry(args []string, exitCode int, start time.Time) ActionLogEntry {
	command := ""
	if len(args) > 0 {
		command = args[0]
	}
	argList := args
	if len(args) > 1 {
		argList = args[1:]
	} else {
		argList = nil
	}
	hostname, _ := os.Hostname()
	user := os.Getenv("USER")
	if runtime.GOOS == "windows" && user == "" {
		user = os.Getenv("USERNAME")
	}
	return ActionLogEntry{
		Timestamp:  start.UTC().Format(time.RFC3339),
		Hostname:   hostname,
		User:       user,
		GOOS:       runtime.GOOS,
		Command:    command,
		Args:       argList,
		ExitCode:   exitCode,
		DurationMs: time.Since(start).Milliseconds(),
	}
}

// DefaultActionLogPath returns a default path for the action log (current dir or home).
func DefaultActionLogPath() string {
	dir, err := os.Getwd()
	if err != nil {
		dir = os.TempDir()
	}
	return filepath.Join(dir, "landschaft-actions.jsonl")
}

// ParseActionLogPath returns path from flag or env LANDSCHAFT_ACTION_LOG, or default.
func ParseActionLogPath(flagPath string) string {
	if flagPath != "" {
		return strings.TrimSpace(flagPath)
	}
	if p := os.Getenv("LANDSCHAFT_ACTION_LOG"); p != "" {
		return strings.TrimSpace(p)
	}
	return DefaultActionLogPath()
}
