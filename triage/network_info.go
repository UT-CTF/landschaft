package triage

import (
	"fmt"
	"net"
	"os"

	"github.com/cakturk/go-netstat/netstat"
	"github.com/charmbracelet/log"
)

func runNetworkTriage() string {
	var hostname = getAndPrintHostname()
	var csv = hostname + "\t"
	csv = csv + printDNSName(hostname) + "\t"
	csv = csv + printIPAddrs() + "\t"
	csv = csv + printNetstat()
	fmt.Println()
	return csv
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
	for _, addr := range addrs {
		names, err := net.LookupAddr(addr)
		if err != nil {
			return "N/A"
		}
		for _, name := range names {
			if name == "localhost" {
				continue
			}
			fmt.Printf("FQDN for %s: %s\n", addr, name)
			result = result + addr + ": " + name + ", "
		}
	}
	return result
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
		result = printAddrs(ipv4Addrs, "  IPv4 Addresses:")
		result = result + "\t" + printAddrs(ipv6Addrs, "  IPv6 Addresses:")
		return result
	}
	return result
}

func printAddrs(list []string, msg string) string {
	var result string
	if len(list) > 0 {
		fmt.Println("  ", msg)
		for _, ip := range list {
			fmt.Println("   -", ip)
			result = result + ip + ", "
		}
	}
	return result
}

func printSockets(title string, sockets []netstat.SockTabEntry) string {
	var result = ""
	if len(sockets) > 0 {
		fmt.Print(title)
		for _, e := range sockets {
			if e.State.String() == "LISTEN" && !e.LocalAddr.IP.IsLoopback() {
				fmt.Printf("%s %s %d %s\n", e.LocalAddr.String(), e.State.String(), e.UID, e.Process)
				result += fmt.Sprintf("%s %s %d %s;", e.LocalAddr.String(), e.State.String(), e.UID, e.Process)
			}
		}
	}

	if len(result) == 0 {
		result = "NONE"
	}

	return result + "\t"
}

func printNetstat() string {
	var result string
	// Get TCP IPv4 sockets
	tcpSocks, err := netstat.TCPSocks(netstat.NoopFilter)
	if err != nil {
		fmt.Print(err)
		result += "err"
	} else {
		result += printSockets("\nTCP IPv4 Sockets:", tcpSocks)
	}
	// Get UDP IPv4 sockets
	udpSocks, err := netstat.UDPSocks(netstat.NoopFilter)
	if err != nil {
		fmt.Print(err)
		result += "err"
	} else {
		result += printSockets("\nUDP IPv4 Sockets:", udpSocks)
	}
	// Get TCP IPv6 sockets
	tcp6Socks, err := netstat.TCP6Socks(netstat.NoopFilter)
	if err != nil {
		fmt.Print(err)
		result += "err"
	} else {
		result += printSockets("\nTCP IPv6 Sockets:", tcp6Socks)
	}
	// Get UDP IPv6 sockets
	udp6Socks, err := netstat.UDP6Socks(netstat.NoopFilter)
	if err != nil {
		fmt.Print(err)
		result += "err"
	} else {
		result += printSockets("\nUDP IPv6 Sockets:", udp6Socks)
	}
	return result
}
