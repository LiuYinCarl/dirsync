package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"dirsync/dirproto"
	"dirsync/dirutil"
)

var (
	rootDir string // 保存传输的文件夹的根目录

	totalFileNum    uint64
	recvivedFileNum uint64
	usedTime        uint64

	dirMapFilePath string
	dirMapFilePtr  *os.File
	dirMapFileBuf  *bufio.Writer

	readBuf  *bufio.Reader
	writeBuf *bufio.Writer

	lastFinishRecord uint64 // 上一个传输完成的文件的序号

	dirFlag  string
	fileFlag string

	serverAddr string
)

func init() {
	dirFlag = "D"
	fileFlag = "F"
	serverAddr = ""

	//readBuf = bufio.NewReader()
	//writeBuf = bufio.NewWriter()

	dirMapFilePath = "./dirMap.txt"
	flag.StringVar(&rootDir, "dir", "", "root dir to save transform dir")
	flag.StringVar(&serverAddr, "serverAddr", "", "server address")
}

// 打开 dirMap 文件并设置读缓冲区
func openDirMapFile() bool {
	var err error
	dirMapFilePtr, err = os.OpenFile(dirMapFilePath, os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println("open dirMapFile Failed. err:", err)
		return false
	}
	dirMapFileBuf = bufio.NewReader(dirMapFilePtr)
	return true
}

// 拼接路径
func joinPath(path string) string {
	var fullPath string

	if len(path) == 0 {
		return fullPath
	}

	if rootDir[len(rootDir)-1] == '/' && path[0] == '/' {
		fullPath = rootDir + path[1:]
	} else {
		fullPath = rootDir + path
	}
	return fullPath
}

// 返回下一条文件记录
func getNextRecord() string {
	// 循环直到找到下一条文件记录
	for {
		line, err := dirMapFileBuf.ReadString('\n')
		if err == io.EOF {
			return ""
		}
		if err != nil {
			fmt.Println("read Record failed. err:", err)
			return ""
		}

		array := strings.Split(line, " ")

		// 如果该条记录是目录，则创建目录
		if array[1] == dirFlag {
			dirPath := joinPath(array[2])
			err = os.MkDirAll(dirPath)
			continue
		}
		return array[2]
	}
}

// 清空TCP连接的多余内容
func clearConnReadBuf(conn net.Conn) bool {
	tempBuf := make([]byte, 128)
	for {
		n, err := conn.Read(tempBuf)
		if err != nil {
			fmt.Println("clear conn read buffer failed")
			return false
		}
		if n == 0 {
			return true
		}
	}
}

// 接收单个文件
func recvFile(conn net.Conn) (result bool, path string) {
	var fileLen int32
	path = getNextRecord()
	if path == "" {
		return true, path
	}
	// 1 给服务端发送需要传输的文件的路径
	bOk, proto := dirproto.CreateProto1(path)
	if !bOk {
		return false, path
	}

	n, err := conn.Write(proto)
	if err != nil {
		fmt.Println("send path to server failed. err:", err)
		return false, path
	}
	if n != len(path) {
		fmt.Printf("send path failed. err:", err)
		return false, path
	}

	// 2 等待服务端返回消息
	readBuf = bufio.NewReader(conn)
	for {
		_, err := conn.Read(readBuf) //TODO 这里读取会把未读的缓冲区冲掉吗
		if err != nil {
			fmt.Println("readBuf err:", err)
			return false, path
		}
		if len(readBuf.buf) >= 4 {
			byte4Buf := make([]byte, 4)
			_, err := readBuf.Read(byte4Buf)
			if err != nil {
				fmt.Println("read proto header failed.")
				return false, path
			}
			protoID := dirproto.ParseProtoHeader(byte4Buf)
			if protoID != 2 {
				fmt.Println("read error protoID.")
				clearConnReadBuf(conn)
				return false, path
			}

			var isFileExist uint8
			_, err := readBuf.Read(&isFileExist)
			if err != nil {
				return false, path
			}
			if isFileExist != 'Y' {
				return false, path
			}

			_, err := readBuf.Read(byte4Buf)
			fileLen = dirutil.BytesToInt32(byte4Buf)
			break
		}
	}

	// 创建文件
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Println("create file err:", err)
		return false, path
	}
	defer file.Close()

	// 3 给服务端发送自己已经准备好接受文件的信号
	_, proto = dirproto.CreateProto3('Y')
	_, err = conn.Write(proto)
	if err != nil {
		fmt.Println(err)
	}

	// 4 等待服务端发来的文件
	for {
		n, err = conn.Read(readBuf)
		if err != nil {
			fmt.Println(err)
			// TODO 传输过程中出现错徐，需要告诉服务端终止本次传输
			return false, path
		}
		if n == 0 && fileLen == 0 {
			//TODO: 使用更合理的方式判断文件传输结束
			break
		}
		n, err = file.Write(readBuf)
		fileLen -= n
	}

	// 5 告诉服务端接收文件完成，完成本次文件传输
	_, proto = dirproto.CreateProto4('Y')
	_, err = conn.Write(proto)
	if err != nil {
		fmt.Println(err)
	}
	return true, path
}

// 接收文件
func recvFiles() {

}

// 接收DirMap
func recvDirMap(conn net.Conn) bool {
	fmt.Println("start recv dirMap file...")

	buf := make([]byte, 4096)

	for {
		n, err := conn.Read(buf)
		if n == 0 {
			// TODO: 需要将 buff 内容写入文件
			return true
		}

		if err != nil {
			fmt.Println("recv dirMap failed. err:", err)
			return false
		}

	}

}

// 启动TCP客户端
func startTCPClient() bool {
	fmt.Println("start TCP client...")

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("net.Dial failed. err:", err)
		return false
	}

	recvDirMap(conn)

}

func printUsage() {
	usage := "Usage: dirsyncClient.exe -serverAddr ip:port -dir need/recv/dir"
	fmt.Println(usage)
}

func main() {
	flag.Parse()

	if rootDir == "" || serverAddr == "" {
		printUsage()
		return
	}

	startTCPClient()
}
