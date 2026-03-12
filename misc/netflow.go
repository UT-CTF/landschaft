package misc

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/spf13/cobra"
)

type FlowKey struct {
	SrcIP   string
	SrcPort int
	DstIP   string
	DstPort int
	Proto   string
}

type FilteredFlowKey struct {
	IP    string
	Port  int
	Proto string
}

type FlowStat struct {
	Packets int
	Bytes   int
}

type PacketSummary struct {
	SrcIP   string
	DstIP   string
	SrcPort int
	DstPort int
	Proto   string
	Length  int
}

var knownPorts = map[int]string{
	21:   "FTP",
	22:   "SSH",
	25:   "SMTP",
	53:   "DNS",
	80:   "HTTP",
	88:   "Kerberos",
	110:  "POP3",
	123:  "NTP",
	135:  "RPC",
	139:  "NetBIOS",
	143:  "IMAP",
	161:  "SNMP",
	389:  "LDAP",
	443:  "HTTPS",
	445:  "SMB",
	465:  "SMTPS",
	587:  "SMTP",
	636:  "LDAPS",
	993:  "IMAPS",
	995:  "POP3S",
	1433: "MS SQL",
	3268: "LDAP GC",
	3269: "LDAPS GC",
	3306: "SQL",
	3389: "RDP",
	5900: "VNC",
	5985: "WinRM",
	5986: "WinRM SSL",
	6001: "MAPI",
	8000: "HTTP Alt",
	8008: "HTTP Alt",
	8080: "HTTP Alt",
	8081: "HTTP Alt",
	8443: "HTTPS Alt",
	8888: "HTTP Alt",
	9389: "AD WS",
}

var netflowFlags struct {
	duration    int
	subnetStr   string
	ifaceName   string
	forceBackup bool
}

func setupNetflowCommand(cmd *cobra.Command) {
	netflowCmd := &cobra.Command{
		Use:   "netflow",
		Short: "Capture and analyze network flows",
		Run: func(cmd *cobra.Command, args []string) {
			analyzeNetflow(netflowFlags.duration, netflowFlags.subnetStr, netflowFlags.ifaceName, netflowFlags.forceBackup)
		},
	}

	netflowCmd.Flags().IntVarP(&netflowFlags.duration, "duration", "d", 60, "Capture duration in seconds")
	netflowCmd.Flags().StringVarP(&netflowFlags.subnetStr, "subnet", "s", "0.0.0.0/0", "Subnet(s) in CIDR format (comma-separated, e.g., 192.168.1.0/24,10.0.0.0/8)")
	netflowCmd.Flags().StringVarP(&netflowFlags.ifaceName, "iface", "i", "", "Interface to capture on (optional)")
	netflowCmd.Flags().BoolVarP(&netflowFlags.forceBackup, "backup", "b", false, "Force use of backup capture method (only on Windows)")

	cmd.AddCommand(netflowCmd)
}

func getLocalIPs() map[string]bool {
	localIPs := make(map[string]bool)

	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err == nil {
				localIPs[ip.String()] = true
			}
		}
	}
	return localIPs
}

func selectInterface(provided string) string {
	if provided != "" {
		return provided
	}

	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	if len(devices) == 0 {
		log.Fatal("No network devices found")
	}

	fmt.Println("Available interfaces:")
	for i, dev := range devices {
		desc := dev.Description
		if desc == "" {
			desc = "(no description)"
		}
		fmt.Printf("[%d] %s — %s\n", i+1, dev.Name, desc)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter interface number: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read input: %v", err)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		n, err := strconv.Atoi(line)
		if err != nil || n < 1 || n > len(devices) {
			fmt.Println("Invalid selection, try again.")
			continue
		}
		return devices[n-1].Name
	}
}

func parseSubnets(s string) ([]*net.IPNet, error) {
	if strings.TrimSpace(s) == "" {
		return nil, fmt.Errorf("empty subnet list")
	}
	parts := strings.Split(s, ",")
	var out []*net.IPNet
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		_, ipnet, err := net.ParseCIDR(p)
		if err != nil {
			return nil, fmt.Errorf("invalid CIDR '%s': %v", p, err)
		}
		out = append(out, ipnet)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no valid CIDRs provided")
	}
	return out, nil
}

func containsAny(nets []*net.IPNet, ip net.IP) bool {
	for _, n := range nets {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}

func defaultPacketCapture(ifaceName string, duration int) ([]PacketSummary, error) {
	iface := selectInterface(ifaceName)
	fmt.Printf("Starting packet capturing on interface %s...\n", iface)
	h, err := pcap.OpenLive(iface, 65535, true, pcap.BlockForever)
	if err != nil {
		return nil, err
	}
	defer h.Close()
	if err := h.SetBPFFilter("tcp or udp"); err != nil {
		return nil, err
	}

	packetSource := gopacket.NewPacketSource(h, h.LinkType())
	out := make([]PacketSummary, 0)
	timeout := time.After(time.Duration(duration) * time.Second)

	for {
		select {
		case packet := <-packetSource.Packets():
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

		case <-timeout:
			return out, nil
		}
	}
}

func filterAndAggregate(packets []PacketSummary, subnets []*net.IPNet, localIPs map[string]bool) (map[FlowKey]*FlowStat, map[FlowKey]*FlowStat) {
	inbound := make(map[FlowKey]*FlowStat)
	outbound := make(map[FlowKey]*FlowStat)

	for _, p := range packets {
		src := net.ParseIP(p.SrcIP)
		dst := net.ParseIP(p.DstIP)
		if src == nil || dst == nil {
			continue
		}

		if !((containsAny(subnets, src) || localIPs[src.String()]) && (containsAny(subnets, dst) || localIPs[dst.String()])) {
			continue
		}

		key := FlowKey{
			SrcIP:   p.SrcIP,
			SrcPort: p.SrcPort,
			DstIP:   p.DstIP,
			DstPort: p.DstPort,
			Proto:   p.Proto,
		}

		if localIPs[dst.String()] {
			if inbound[key] == nil {
				inbound[key] = &FlowStat{}
			}
			inbound[key].Packets++
			inbound[key].Bytes += p.Length
		} else if localIPs[src.String()] {
			if outbound[key] == nil {
				outbound[key] = &FlowStat{}
			}
			outbound[key].Packets++
			outbound[key].Bytes += p.Length
		}
	}

	return inbound, outbound
}

func analyzeNetflow(duration int, subnetStr string, ifaceName string, forceBackup bool) {
	if subnetStr == "" {
		fmt.Println("Usage: netflow -duration 60 -subnet 192.168.1.0/24[,10.0.0.0/8] [-iface eth0]")
		os.Exit(1)
	}

	subnets, err := parseSubnets(subnetStr)
	if err != nil {
		log.Fatal("Invalid subnet(s):", err)
	}

	localIPs := getLocalIPs()

	var packets []PacketSummary

	if !forceBackup && verifyPacketCapture() {
		fmt.Println("Default packet capture available")
		packets, err = defaultPacketCapture(ifaceName, duration)
		if err != nil {
			log.Fatalf("Default packet capture failed: %v", err)
		}
	} else {
		fmt.Println("Default packet capture not available, using backup method")
		packets, err = backupPacketCapture(duration)
		if err != nil {
			log.Fatalf("Backup packet capture failed: %v", err)
		}
	}

	inbound, outbound := filterAndAggregate(packets, subnets, localIPs)

	printTable("Inbound Flows", inbound)
	printTable("Outbound Flows", outbound)
	printFiltered("Filtered Inbound (known dest ports)", inbound, true)
	printFiltered("Filtered Outbound (known dest ports)", outbound, false)
}

func printTable(title string, flows map[FlowKey]*FlowStat) {
	fmt.Println("\n====", title, "====")

	type row struct {
		Key  FlowKey
		Stat *FlowStat
	}

	var rows []row
	for k, v := range flows {
		rows = append(rows, row{k, v})
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Stat.Bytes > rows[j].Stat.Bytes
	})

	fmt.Printf("%-15s %-6s %-15s %-6s %-5s %-8s %-8s\n",
		"SRC IP", "SPORT", "DST IP", "DPORT", "PROTO", "PACKETS", "BYTES")

	for _, r := range rows {
		fmt.Printf("%-15s %-6d %-15s %-6d %-5s %-8d %-8d\n",
			r.Key.SrcIP,
			r.Key.SrcPort,
			r.Key.DstIP,
			r.Key.DstPort,
			r.Key.Proto,
			r.Stat.Packets,
			r.Stat.Bytes,
		)
	}
}

func printFiltered(title string, flows map[FlowKey]*FlowStat, inbound bool) {
	filteredFlows := make(map[FilteredFlowKey]*FlowStat)
	for k, v := range flows {
		var ip string
		if inbound {
			ip = k.SrcIP
		} else {
			ip = k.DstIP
		}
		if _, ok := knownPorts[k.DstPort]; ok {
			fk := FilteredFlowKey{
				IP:    ip,
				Port:  k.DstPort,
				Proto: k.Proto,
			}
			if filteredFlows[fk] == nil {
				filteredFlows[fk] = &FlowStat{}
			}
			filteredFlows[fk].Packets += v.Packets
			filteredFlows[fk].Bytes += v.Bytes
		}
	}
	fmt.Println("\n====", title, "====")
	fmt.Printf("%-15s %-6s %-20s %-8s %-8s\n", "IP", "PORT", "SERVICE", "PACKETS", "BYTES")
	for k, v := range filteredFlows {
		svc := knownPorts[k.Port]
		fmt.Printf("%-15s %-6d %-20s %-8d %-8d\n", k.IP, k.Port, svc, v.Packets, v.Bytes)
	}
}
