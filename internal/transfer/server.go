package transfer

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/quic-go/quic-go"
)

type ReceivedFile struct {
	Name   string
	Sender string
	Size   int64
}

type FileTracker struct {
	sync.Mutex
	Files []ReceivedFile
}

func (ft *FileTracker) Add(file ReceivedFile) {
	ft.Lock()
	defer ft.Unlock()
	ft.Files = append(ft.Files, file)
}

func (ft *FileTracker) List() []ReceivedFile {
	ft.Lock()
	defer ft.Unlock()
	filesCopy := make([]ReceivedFile, len(ft.Files))
	copy(filesCopy, ft.Files)
	return filesCopy
}

func StartServer(ctx context.Context, addr, receiveDir string, tracker *FileTracker) error {
	tlsConfig, err := generateTLSConfig()
	if err != nil {
		return err
	}

	listener, err := quic.ListenAddr(addr, tlsConfig, nil)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept(context.TODO())
		if err != nil {
			return err
		}
		go handleConn(receiveDir, tracker, conn)
	}
}

func handleConn(receiveDir string, tracker *FileTracker, conn *quic.Conn) {
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
	// show it to user
	// if autoAccept is true then start getting the file

	fmt.Printf("📥 Receiving file: %s (%d bytes)\n", header.Filename, header.Filesize)

	path := filepath.Join(receiveDir, header.Filename)
	f, err := os.Create(path)
	if err != nil {
		fmt.Println("file create error:", err)
		return
	}

	written, err := io.Copy(f, stream)
	if err != nil {
		fmt.Println("copy error:", err)
		return
	}
	f.Close()

	fmt.Printf("✅ Done. Received %d bytes\n", written)
	ip := strings.Split(conn.LocalAddr().String(), ":")[0]

	tracker.Add(ReceivedFile{
		Name:   header.Filename,
		Size:   header.Filesize,
		Sender: ip,
	})
}
