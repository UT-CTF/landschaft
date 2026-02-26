package misc

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
)

func verifyPacketCapture() bool {
	devices, err := pcap.FindAllDevs()
	return err == nil && len(devices) > 0
}

func backupPacketCapture(duration int) ([]PacketSummary, error) {
	fmt.Println("Starting pktmon capture...")

	etl := "PktMon.etl"
	cmdStart := exec.Command("pktmon", "start", "--capture", "--pkt-size", "0", "--file-size", "4096", "--file-name", etl)
	cmdStart.Stdout = io.Discard
	cmdStart.Stderr = io.Discard
	if err := cmdStart.Run(); err != nil {
		return nil, fmt.Errorf("Failed to start pktmon capture: %v", err)
	}

	time.Sleep(time.Duration(duration) * time.Second)

	cmdStop := exec.Command("pktmon", "stop")
	cmdStop.Stdout = io.Discard
	cmdStop.Stderr = io.Discard
	if err := cmdStop.Run(); err != nil {
		return nil, fmt.Errorf("Failed to stop pktmon capture: %v", err)
	}

	pcapOut := "pktmon.pcap"
	cmdFmt := exec.Command("pktmon", "etl2pcap", etl, "-o", pcapOut)
	cmdFmt.Stdout = io.Discard
	cmdFmt.Stderr = io.Discard
	if err := cmdFmt.Run(); err != nil {
		return nil, fmt.Errorf("Failed to format pktmon ETL to pcap: %v", err)
	}
	defer os.Remove(etl)
	defer os.Remove(pcapOut)

	f, err := os.Open(pcapOut)
	if err != nil {
		return nil, fmt.Errorf("Failed to open pktmon pcap file: %v", err)
	}
	defer f.Close()

	reader, err := pcapgo.NewNgReader(f, pcapgo.DefaultNgReaderOptions)
	if err != nil {
		return nil, fmt.Errorf("Failed to create pcapng reader: %v", err)
	}
	ps := gopacket.NewPacketSource(reader, reader.LinkType())

	out := make([]PacketSummary, 0)

	for packet := range ps.Packets() {
		if packet == nil {
			continue
		}
		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if ipLayer == nil {
			continue
		}
		ip := ipLayer.(*layers.IPv4)

		var srcPort, dstPort int
		var proto string
		if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
			tcp := tcpLayer.(*layers.TCP)
			srcPort = int(tcp.SrcPort)
			dstPort = int(tcp.DstPort)
			proto = "TCP"
		} else if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
			udp := udpLayer.(*layers.UDP)
			srcPort = int(udp.SrcPort)
			dstPort = int(udp.DstPort)
			proto = "UDP"
		} else {
			continue
		}

		out = append(out, PacketSummary{
			SrcIP:   ip.SrcIP.String(),
			DstIP:   ip.DstIP.String(),
			SrcPort: srcPort,
			DstPort: dstPort,
			Proto:   proto,
			Length:  len(packet.Data()),
		})
	}

	return out, nil
}
