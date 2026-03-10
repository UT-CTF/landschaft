package score

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/UT-CTF/landschaft/util"
)

func discoverListeners() ([]Listener, error) {
	out, err := util.RunScriptQuiet("score/listeners.ps1")
	if err != nil {
		return nil, err
	}
	return parseWindowsListeners(out)
}

type winListener struct {
	LocalPort     int    `json:"LocalPort"`
	OwningProcess int    `json:"OwningProcess"`
	LocalAddress  string `json:"LocalAddress"`
}

func parseWindowsListeners(text string) ([]Listener, error) {
	text = strings.TrimSpace(text)
	var raw []winListener
	if err := json.Unmarshal([]byte(text), &raw); err != nil {
		var single winListener
		if err2 := json.Unmarshal([]byte(text), &single); err2 != nil {
			return nil, err
		}
		raw = []winListener{single}
	}
	var out []Listener
	for _, r := range raw {
		port := uint16(r.LocalPort)
		out = append(out, Listener{
			Port:    port,
			Proto:   "tcp",
			Process: strconv.Itoa(r.OwningProcess),
			Bind:    r.LocalAddress,
			Explain: PortExplain(port),
		})
	}
	return out, nil
}
