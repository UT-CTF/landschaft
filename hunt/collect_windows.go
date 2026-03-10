package hunt

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/UT-CTF/landschaft/util"
)

func getDetections(since time.Duration) ([]Detection, error) {
	return collectWindows(since)
}

// collectWindows runs Get-WinEvent and returns normalized detections.
func collectWindows(since time.Duration) ([]Detection, error) {
	minutes := int(since.Minutes())
	if minutes < 1 {
		minutes = 30
	}
	out, err := util.RunScriptQuiet("hunt/events.ps1", strconv.Itoa(minutes))
	if err != nil {
		return nil, err
	}
	return parseWindowsEvents(out)
}

type winEvent struct {
	TimeCreated string `json:"TimeCreated"`
	Id          int    `json:"Id"`
	Message     string `json:"Message"`
	LogName     string `json:"LogName"`
	Properties  []interface{} `json:"Properties"`
}

func parseWindowsEvents(text string) ([]Detection, error) {
	text = strings.TrimSpace(text)
	if text == "" || text == "null" {
		return nil, nil
	}
	var raw []winEvent
	if err := json.Unmarshal([]byte(text), &raw); err != nil {
		var single winEvent
		if err2 := json.Unmarshal([]byte(text), &single); err2 != nil {
			return nil, err
		}
		raw = []winEvent{single}
	}
	host, _ := os.Hostname()
	var out []Detection
	for _, e := range raw {
		d := Detection{
			Timestamp: e.TimeCreated,
			Host:      host,
			OS:        "windows",
			Source:    e.LogName,
			EventID:   strconv.Itoa(e.Id),
			Message:   e.Message,
		}
		TagSuspicious(&d)
		out = append(out, d)
	}
	return out, nil
}
