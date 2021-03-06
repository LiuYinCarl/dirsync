package dirproto

import (
	"dirsync/dirutil"
)

const (
	//InvalidProtoID 非法协议ID
	InvalidProtoID uint8 = 0
	//MinProtoID 当前最小的协议ID号
	MinProtoID uint8 = 1
	//MaxProtoID 当前最大的协议ID号
	MaxProtoID uint8 = 4
	//ProtoHeaderFlag 每个协议头的开始部分，用来标记这是一条dir协议
	ProtoHeaderFlag string = "dir" //TODO: 定义为字节数组
)

var (
	byte1Buf = make([]byte, 1)
)

//isProtoValid 检查协议ID是否合法
func isProtoValid(protoID uint8) bool {
	if protoID >= MinProtoID && protoID <= MaxProtoID {
		return true
	}
	return false
}

//createCommonHeader 创建公共协议头
func createCommonHeader(protoID uint8) (bool, []byte) {
	if !isProtoValid(protoID) {
		return false, []byte("")
	}
	header := make([]byte, 4)
	copy(header[0:3], []byte(ProtoHeaderFlag))
	header[3] = protoID
	return true, header
}

//ParseProtoHeader 根据传入的协议头返回这个协议的ID
func ParseProtoHeader(header []byte) uint8 {
	if len(header) != 4 {
		return InvalidProtoID
	}
	if string(header[:3]) != ProtoHeaderFlag {
		return InvalidProtoID
	}
	if !isProtoValid(header[3]) {
		return InvalidProtoID
	}
	return header[3]
}

//CreateProto1 创建协议1
func CreateProto1(path string) (bool, []byte) {
	_, header := createCommonHeader(1)
	pathLen := int16(len(path))
	proto := dirutil.BytesCombine(header, dirutil.Int16ToBytes(pathLen), []byte(path))
	return true, proto
}

//CreateProto2 创建协议2
// func CreateProto2() (bool, []byte) {

// }

//CreateProto3 创建协议3
func CreateProto3(c uint8) (bool, []byte) {
	if c != 'Y' && c != 'N' {
		return false, []byte{}
	}
	_, header := createCommonHeader(3)
	byte1Buf[0] = c
	proto := dirutil.BytesCombine(header, byte1Buf)
	return true, proto
}

//CreateProto4 创建协议4
func CreateProto4(c uint8) (bool, []byte) {
	if c != 'Y' && c != 'N' {
		return false, []byte{}
	}
	_, header := createCommonHeader(4)
	byte1Buf[0] = c
	proto := dirutil.BytesCombine(header, byte1Buf)
	return true, proto
}
