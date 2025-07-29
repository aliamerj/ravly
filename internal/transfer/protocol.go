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

	if err := binary.Write(w, binary.LittleEndian, uint16(len(nameBytes))); err != nil {
		return err
	}

	if _, err := w.Write(nameBytes); err != nil {
		return err
	}

	return binary.Write(w, binary.LittleEndian, uint16(header.Filesize))
}

func readHeader(r io.Reader) (fileHeader, error) {
	var nameLen uint16
	if err := binary.Read(r, binary.LittleEndian, &nameLen); err != nil {
		return fileHeader{}, err
	}

	nameBytes := make([]byte, nameLen)
	if _, err := io.ReadFull(r, nameBytes); err != nil {
		return fileHeader{}, err
	}

	var size int64
	if err := binary.Read(r, binary.LittleEndian, &size); err != nil {
		return fileHeader{}, err
	}

	return fileHeader{
		Filename: string(nameBytes),
		Filesize: size,
	}, nil
}
