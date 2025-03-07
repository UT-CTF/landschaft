package triage

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/cakturk/go-netstat/netstat"
)

func print_network_info() {
	// var hostname = print_hostname()
	// print_dnsname(hostname)
	// print_ip_addrs()
	print_netstat()
}

func print_hostname() string {
	var hostname, name_err = os.Hostname()
	if name_err == nil {
		fmt.Printf("Host Name: %s\n", hostname)
	}
	return hostname
}

func print_dnsname(hostname string) {
	addrs, err := net.LookupAddr(hostname)
	if err != nil || len(addrs) == 0 {
		fmt.Println("Error getting FQDN:", err)
		fmt.Printf("DNS Name: %s\n", hostname)
		return
	}
	fmt.Printf("DNS Name: %s\n", strings.TrimSuffix(addrs[0], "."))
}

func print_ip_addrs() {
	var interfaces, err = net.Interfaces()
	if err != nil {
		fmt.Printf("Error obtaining interfaces.\n")
		return
	}
	for _, iface := range interfaces {
		// skip loopback address
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		fmt.Printf("Interface: %s\n", iface.Name)
		var addrs, ip_err = iface.Addrs()
		if ip_err != nil {
			fmt.Printf("Error obtaining IP addresses")
			return
		}
		var ipv4_addrs, ipv6_addrs []string
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				fmt.Printf("Error converting IP Addr")
				return
			}
			if ipNet.IP.To4() != nil {
				ipv4_addrs = append(ipv4_addrs, ipNet.IP.String())
			} else {
				ipv6_addrs = append(ipv6_addrs, ipNet.IP.String())
			}
		}
		print_addrs(ipv4_addrs, "  IPv4 Addresses:")
		print_addrs(ipv6_addrs, "  IPv6 Addresses:")
	}
}

func print_addrs(list []string, msg string) {
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
			if e.State.String() == "LISTEN" && !strings.Contains(e.LocalAddr.String(), "127.0.0") {
				fmt.Printf("%s %s %d %s\n", e.LocalAddr.String(), e.State.String(), e.UID, e.Process)
			}
		}
	}
}

func print_netstat() {
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
