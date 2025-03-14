package triage

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/cakturk/go-netstat/netstat"
	"github.com/charmbracelet/log"
)

func runNetworkTriage() {
	var hostname = getAndPrintHostname()
	printDNSName(hostname)
	printIPAddrs()
	printNetstat()
	fmt.Println()
}

func getAndPrintHostname() string {
	var hostname, nameErr = os.Hostname()
	if nameErr != nil {
		log.Error("Failed to get hostname", "err", nameErr)
		return "Error getting hostname"
	}
	fmt.Printf("Host Name: %s\n", hostname)
	return hostname
}

func printDNSName(hostname string) {
	addrs, err := net.LookupAddr(hostname)
	if err != nil || len(addrs) == 0 {
		log.Error("Failed to get FQDN", "err", err)
		fmt.Printf("DNS Name: %s\n", hostname)
		return
	}
	fmt.Printf("DNS Name: %s\n", strings.TrimSuffix(addrs[0], "."))
}

func printIPAddrs() {
	var interfaces, err = net.Interfaces()
	if err != nil {
		log.Error("Failed to get network interfaces", "err", err)
		return
	}
	for _, iface := range interfaces {
		// skip loopback address
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		fmt.Printf("Interface: %s\n", iface.Name)
		var addrs, ipErr = iface.Addrs()
		if ipErr != nil {
			log.Error("Failed to get IP addresses", "err", ipErr)
			return
		}
		var ipv4Addrs, ipv6Addrs []string
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				log.Error("Failed to cast address to IPNet", "addr", addr)
				return
			}
			if ipNet.IP.To4() != nil {
				ipv4Addrs = append(ipv4Addrs, ipNet.IP.String())
			} else {
				ipv6Addrs = append(ipv6Addrs, ipNet.IP.String())
			}
		}
		printAddrs(ipv4Addrs, "  IPv4 Addresses:")
		printAddrs(ipv6Addrs, "  IPv6 Addresses:")
	}
}

func printAddrs(list []string, msg string) {
	if len(list) > 0 {
		fmt.Println("  ", msg)
		for _, ip := range list {
			fmt.Println("   -", ip)
		}
	}
}

func printSockets(title string, sockets []netstat.SockTabEntry) {
	if len(sockets) > 0 {
		fmt.Println(title)
		for _, e := range sockets {
			if e.State.String() == "LISTEN" && !e.LocalAddr.IP.IsLoopback() {
				fmt.Printf("%s %s %d %s\n", e.LocalAddr.String(), e.State.String(), e.UID, e.Process)
			}
		}
	}
}

func printNetstat() {
	// Get TCP IPv4 sockets
	tcpSocks, err := netstat.TCPSocks(netstat.NoopFilter)
	if err != nil {
		fmt.Print(err)
		return
	}
	printSockets("\nTCP IPv4 Sockets:", tcpSocks)

	// Get UDP IPv4 sockets
	udpSocks, err := netstat.UDPSocks(netstat.NoopFilter)
	if err != nil {
		fmt.Print(err)
		return
	}
	printSockets("\nUDP IPv4 Sockets:", udpSocks)

	// Get TCP IPv6 sockets
	tcp6Socks, err := netstat.TCP6Socks(netstat.NoopFilter)
	if err != nil {
		fmt.Print(err)
		return
	}
	printSockets("\nTCP IPv6 Sockets:", tcp6Socks)

	// Get UDP IPv6 sockets
	udp6Socks, err := netstat.UDP6Socks(netstat.NoopFilter)
	if err != nil {
		fmt.Print(err)
		return
	}
	printSockets("\nUDP IPv6 Sockets:", udp6Socks)
}
