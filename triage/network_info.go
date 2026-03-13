package triage

import (
	"fmt"
	"net"
	"os"
	"sort"

	"github.com/cakturk/go-netstat/netstat"
	"github.com/charmbracelet/log"
)

func runNetworkTriage() (string, string) {
	hostname := getAndPrintHostname()
	csv := printDNSName(hostname) + "\t"
	csv = csv + printIPAddrs() + "\t"
	csv = csv + printNetstat()
	fmt.Println()
	return hostname, csv
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

func printDNSName(hostname string) string {
	var result string
	addrs, err := net.LookupHost(hostname)
	if err != nil {
		fmt.Println("error looking up hostname")
		return "err"
	}
	var check = false
	for _, addr := range addrs {
		names, err := net.LookupAddr(addr)
		if err != nil {
			return "N/A"
		}
		for _, name := range names {
			if name == "localhost" || addr == "127.0.1.1" || addr == "127.0.0.1" {
				continue
			}
			fmt.Printf("FQDN for %s: %s\n", addr, name)
			result = result + addr + ": " + name + "\n"
			check = true
		}
	}
	if !check {
		return "N/A"
	}
	return "\"" + result + "\""
}

func printIPAddrs() string {
	var result string
	var interfaces, err = net.Interfaces()
	if err != nil {
		log.Error("Failed to get network interfaces", "err", err)
		return "err"
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
			return "err"
		}
		var ipv4Addrs, ipv6Addrs []string
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				log.Error("Failed to cast address to IPNet", "addr", addr)
				return "err"
			}
			if ipNet.IP.To4() != nil {
				ipv4Addrs = append(ipv4Addrs, ipNet.IP.String())
			} else {
				ipv6Addrs = append(ipv6Addrs, ipNet.IP.String())
			}
		}
		temp := printAddrs(ipv4Addrs, "  IPv4 Addresses:") + "\n\n"
		if len(temp) > 2 {
			result += iface.Name + temp
		}
		//result = result + "\t" + printAddrs(ipv6Addrs, "  IPv6 Addresses:")
		printAddrs(ipv6Addrs, "  IPv6 Addresses:")
	}
	return "\"" + result + "\""
}

func printAddrs(list []string, msg string) string {
	var result string
	if len(list) > 0 {
		fmt.Println("  ", msg)
		for _, ip := range list {
			fmt.Println("   -", ip)
			result += fmt.Sprintf("\n\t%s", ip)
		}
	}
	return result
}

func printSockets(title string, sockets []netstat.SockTabEntry) string {
	type entry struct {
		port    uint16
		process string
	}

	var result string
	seen := make(map[uint16]bool)
	var entries []entry

	for _, e := range sockets {
		if e.State.String() == "LISTEN" && !e.LocalAddr.IP.IsLoopback() {
			port := e.LocalAddr.Port

			if !seen[port] {
				seen[port] = true

				// Guard against nil Process which can cause a panic when calling String()
				var procName string
				if e.Process != nil {
					procName = e.Process.String()
				} else {
					procName = "N/A"
				}

				entries = append(entries, entry{
					port:    port,
					process: procName,
				})

				fmt.Printf("%s %s %d %s\n", e.LocalAddr.String(), e.State.String(), e.UID, procName)
			}
		}
	}

	if len(entries) > 0 {
		fmt.Print(title)

		sort.Slice(entries, func(i, j int) bool {
			return entries[i].port < entries[j].port
		})

		for _, e := range entries {
			result += fmt.Sprintf("%d\t%s\n", e.port, e.process)
		}
	}

	if len(result) == 0 {
		result = "NONE"
	}

	return result
}

func printNetstat() string {
	var result string
	// Get TCP IPv4 sockets
	tcpSocks, err := netstat.TCPSocks(netstat.NoopFilter)
	if err != nil {
		fmt.Print(err)
		result += "err"
	} else {
		result += "## TCPv4 ##\n" + printSockets("\nTCP IPv4 Sockets:", tcpSocks)
	}
	// Get UDP IPv4 sockets
	udpSocks, err := netstat.UDPSocks(netstat.NoopFilter)
	if err != nil {
		fmt.Print(err)
		result += "err"
	} else {
		result += "\n\n## UDPv4 ##\n" + printSockets("\nUDP IPv4 Sockets:", udpSocks)
	}
	// Get TCP IPv6 sockets
	tcp6Socks, err := netstat.TCP6Socks(netstat.NoopFilter)
	if err != nil {
		fmt.Print(err)
		result += "err"
	} else {
		result += "\n\n## TCPv6 ##\n" + printSockets("\nTCP IPv6 Sockets:", tcp6Socks)
	}
	// Get UDP IPv6 sockets
	udp6Socks, err := netstat.UDP6Socks(netstat.NoopFilter)
	if err != nil {
		fmt.Print(err)
		result += "err"
	} else {
		result += "\n\n## UDPv6 ##\n" + printSockets("\nUDP IPv6 Sockets:", udp6Socks)
	}
	return "\"" + result + "\"\t"
}
