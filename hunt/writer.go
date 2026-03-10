package hunt

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// AppendDetection appends a detection to the JSONL file. Creates dirs if needed.
func AppendDetection(path string, d Detection) error {
	if path == "" {
		return nil
	}
	_ = os.MkdirAll(filepath.Dir(path), 0755)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	return enc.Encode(d)
}
