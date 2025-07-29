package discovery

import (
	"fmt"
	"net"
	"strings"
)

const (
	broadcastAddr = "255.255.255.255"
	discoveryMsg  = "ravly::hello"
)

type Peer struct {
	Name string
	Ip   string
}

// BroadcastPresence sends a UDP broadcast packet.
func BroadcastPresence(name string, port int) error {
	raddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", broadcastAddr, port))
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	hostName := name
	if name == "" {
		myIP := strings.Split(conn.LocalAddr().String(), ":")[0]
		host, err := net.LookupAddr(myIP)
		if err != nil {
			return err
		}
		hostName = host[0]
	}

	if _, err := conn.Write([]byte(discoveryMsg + ">" + hostName)); err != nil {
		return err
	}
	return nil
}

// ListenForPeers listens for UDP discovery messages.
func ListenForPeers(port int, callback func(peer Peer)) error {
  ladd, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
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
			callback(Peer{
				Name: msg[1],
				Ip:   remoteAddr.IP.String(),
			})
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
