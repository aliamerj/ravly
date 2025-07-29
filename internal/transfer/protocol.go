package transfer

import (
	"encoding/binary"
	"io"
)

// fileHeader is sent before the file bytes
type fileHeader struct {
	Filename string
	Filesize int64
}

func writeHeader(w io.Writer, header fileHeader) error {
	nameBytes := []byte(header.Filename)

	if err := binary.Write(w, binary.BigEndian, uint32(len(nameBytes))); err != nil {
		return err
	}

	if _, err := w.Write(nameBytes); err != nil {
		return err
	}

	return binary.Write(w, binary.BigEndian, header.Filesize)
}

func readHeader(stream io.Reader) (*fileHeader, error) {
    var nameLen uint32
    if err := binary.Read(stream, binary.BigEndian, &nameLen); err != nil {
        return nil, err
    }
    nameBytes := make([]byte, nameLen)
    if _, err := io.ReadFull(stream, nameBytes); err != nil {
        return nil, err
    }
    var fileSize int64
    if err := binary.Read(stream, binary.BigEndian, &fileSize); err != nil {
        return nil, err
    }
    return &fileHeader{
        Filename: string(nameBytes),
        Filesize: fileSize,
    }, nil
}
