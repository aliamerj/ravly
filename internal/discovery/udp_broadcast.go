package internal

import (
	"fmt"
	"net"
	"strings"
)

const (
	broadcastAddr = "255.255.255.255:9999"
	discoveryMsg  = "ravly::hello"
)

// BroadcastPresence sends a UDP broadcast packet.
func BroadcastPresence() error {
	raddr, err := net.ResolveUDPAddr("udp", broadcastAddr)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	myIP := strings.Split(conn.LocalAddr().String(), ":")[0]
	hostName, err := net.LookupAddr(myIP)

	if _, err := conn.Write([]byte(discoveryMsg + ">" + hostName[0])); err != nil {
		return err
	}
	return nil
}

// ListenForPeers listens for UDP discovery messages.
func ListenForPeers(callback func(addr, host string)) error {
	ladd, err := net.ResolveUDPAddr("udp", ":9999")
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", ladd)
	if err != nil {
		return err
	}
	defer conn.Close()

	myIP, err := externalIP()
	if err != nil {
		return err
	}

	buf := make([]byte, 1024)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		msg := strings.Split(string(buf[:n]), ">")
		ip := remoteAddr.IP.String()
		if len(msg) == 2 && msg[0] == discoveryMsg && ip != myIP {
			callback(remoteAddr.IP.String(), msg[1])
		}
	}
}
func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", fmt.Errorf("can't not find ip")
}
