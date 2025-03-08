package triage

import "testing"

func TestNetworkTriage(t *testing.T) {
	runNetworkTriage()
}

func TestUsersTriage(t *testing.T) {
	runUsersTriage()
}

func TestOSVersionTriage(t *testing.T) {
	runOSVersionTriage()
}

func TestFirewallTriage(t *testing.T) {
	runFirewallTriage()
}
