package internal

import "net"

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

	if _, err := conn.Write([]byte(discoveryMsg)); err != nil {
		return err
	}
	return nil
}

// ListenForPeers listens for UDP discovery messages.
func ListenForPeers(callback func(addr string)) error {
	ladd, err := net.ResolveUDPAddr("udp", ":9999")
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", ladd)
	if err != nil {
		return err
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		msg := string(buf[:n])

		if msg == discoveryMsg {
			callback(remoteAddr.IP.String())
		}
	}
}
