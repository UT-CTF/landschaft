package misc

import (
	"fmt"

	"github.com/google/gopacket/pcap"
)

func verifyPacketCapture() bool {
	devices, err := pcap.FindAllDevs()
	return err == nil && len(devices) > 0
}

func backupPacketCapture(duration int) ([]PacketSummary, error) {
	return nil, fmt.Errorf("pktmon is only available on Windows")
}
