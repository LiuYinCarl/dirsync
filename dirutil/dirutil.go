package dirutil

import (
	"bytes"
	"encoding/binary"
)

//Int16ToBytes int16转[]byte
func Int16ToBytes(n int16) []byte {
	bytesBuf := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuf, binary.BigEndian, n)
	return bytesBuf.Bytes()
}

//Int32ToBytes int32转[]byte
func Int32ToBytes(n int16) []byte {
	bytesBuf := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuf, binary.BigEndian, n)
	return bytesBuf.Bytes()
}

//Int64ToBytes int64转[]byte
func Int64ToBytes(n int64) []byte {
	bytesBuf := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuf, binary.BigEndian, n)
	return bytesBuf.Bytes()
}

// BytesToInt16 []byte转int16
func BytesToInt16(b []byte) int16 {
	bytesBuf := bytes.NewBuffer(b)
	var n int16
	binary.Read(bytesBuf, binary.BigEndian, &n)
	return n
}

// BytesToInt32 []byte转int32
func BytesToInt32(b []byte) int32 {
	bytesBuf := bytes.NewBuffer(b)
	var n int32
	binary.Read(bytesBuf, binary.BigEndian, &n)
	return n
}

// BytesToInt64 []byte转int64
func BytesToInt64(b []byte) int64 {
	bytesBuf := bytes.NewBuffer(b)
	var n int64
	binary.Read(bytesBuf, binary.BigEndian, &n)
	return n
}

//BytesCombine 多个[]byte合并成一个[]byte
func BytesCombine(b ...[]byte) []byte {
	return bytes.Join(b, []byte(""))
}
