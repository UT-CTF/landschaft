package score

import (
	"github.com/cakturk/go-netstat/netstat"
)

func discoverListeners() ([]Listener, error) {
	var out []Listener
	tcp, err := netstat.TCPSocks(netstat.NoopFilter)
	if err != nil {
		return nil, err
	}
	for _, e := range tcp {
		if e.State.String() != "LISTEN" {
			continue
		}
		port := uint16(e.LocalAddr.Port)
		procName := ""
		if e.Process != nil {
			procName = e.Process.Name
		}
		out = append(out, Listener{
			Port:    port,
			Proto:   "tcp",
			Process: procName,
			Bind:    e.LocalAddr.IP.String(),
			Explain: PortExplain(port),
		})
	}
	return out, nil
}
