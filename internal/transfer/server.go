package transfer

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/quic-go/quic-go"
)

func StartServer(addr string) error {
	tlsConfig, err := generateTLSConfig()
	if err != nil {
		return err
	}

	listener, err := quic.ListenAddr(addr, tlsConfig, nil)
	if err != nil {
		return err
	}
	fmt.Println("🚀 Listening for files on", addr)

	for {
		conn, err := listener.Accept(context.TODO())
		if err != nil {
			return err
		}
		go handleConn(conn)
	}
}

func handleConn(conn *quic.Conn) {

	stream, err := conn.AcceptStream(context.TODO())
	if err != nil {
		fmt.Println("stream accept error:", err)
		return
	}
	defer stream.Close()

	header, err := readHeader(stream)
	if err != nil {
		fmt.Println("read header error:", err)
		return
	}
	fmt.Printf("📥 Receiving file: %s (%d bytes)\n", header.Filename, header.Filesize)

	f, err := os.Create(header.Filename)
	if err != nil {
		fmt.Println("file create error:", err)
		return
	}
	defer f.Close()

	written, err := io.CopyN(f, stream, header.Filesize)
	if err != nil {
		fmt.Println("copy error:", err)
		return
	}
	fmt.Printf("✅ Done. Received %d bytes\n", written)
}
