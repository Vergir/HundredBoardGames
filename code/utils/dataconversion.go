package utils

import (
	"bytes"
	"io"
)

func ReaderToBytes(reader io.Reader) ([]byte, error) {
	bytesBuffer := new(bytes.Buffer)

	_, err := bytesBuffer.ReadFrom(reader)
	if err != nil {
		return nil, err
	}

	return bytesBuffer.Bytes(), nil
}
