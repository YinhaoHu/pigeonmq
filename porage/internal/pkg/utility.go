package pkg

import (
	"bytes"
	"encoding/binary"
)

// Int64ToBytes converts an int64 to a byte slice (Big Endian).
func Int64ToBytes(num int64) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)
	if err != nil {
		return nil, err
	}
	bufBytes := buf.Bytes()
	return bufBytes, nil
}

// BytesToInt64 converts a byte slice (Big Endian) back to an int64.
func BytesToInt64(b []byte) (int64, error) {
	buf := bytes.NewReader(b)
	var num int64
	err := binary.Read(buf, binary.BigEndian, &num)
	if err != nil {
		return 0, err
	}
	return num, nil
}
