package transfer

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/quic-go/quic-go"
)

func SendFile(ip, filePath string, port int) error {
	tlsConfig, err := generateTLSConfig()
	if err != nil {
		return err
	}
  addr := fmt.Sprintf("%s:%d", ip, port)
  
	conn, err := quic.DialAddr(context.TODO(), addr, tlsConfig, nil)
	if err != nil {
		return err
	}

	steam, err := conn.OpenStreamSync(context.TODO())
	if err != nil {
		return err
	}
	defer steam.Close()

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	// send header
	if err := writeHeader(steam, fileHeader{
		Filename: stat.Name(),
		Filesize: stat.Size(),
	}); err != nil {
		return err
	}

  // send file
	written, err := io.Copy(steam, file)
	if err != nil {
		return err
	}

	fmt.Printf("📤 Sent %s (%d bytes)\n", filePath, written)
	return nil
}
