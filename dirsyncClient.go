package main

import (
	"fmt"
	"util"
	"net"
	"os"
)

var (
	rootDir string  // 保存传输的文件夹的根目录

	totalFileNum uint64
	recvivedFileNum uint64
	usedTime uint64

	dirMapFilePath string
	dirMapFilePtr *os.File
	dirMapFileBuf *bufio.Writer

	lastFinishRecord uint64 // 上一个传输完成的文件的序号

	dirFlag  string
	fileFlag string

	serverAddr string
)

func init() {
	dirFlag = "D"
	fileFlag = "F"
	serverAddr = ""

	dirMapFilePath = "./dirMap.txt"
	flag.StringVar(&rootDir, "dir", "", "root dir to save transform dir")
	flag.StringVar(&serverAddr, "serverAddr", "", "server address")
}

// 打开 dirMap 文件并设置读缓冲区
func openDirMapFile() bool {
	dirMapFilePtr, err := os.OpenFile(dirMapFilePath, os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println("open dirMapFile Failed. err:", err)
		return false
	}
	dirMapFileBuf = bufio.NewReader()
	return true
}

// 拼接路径
func joinPath(path string) string {
	if len(path) == 0 {
		return fullPath
	}

	var fullPath string
	
	if rootDir[len(rootDir-1)] == "/" && path[0] == "/" {
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
			err := os.MkDirAll(dirPath)
			continue
		}
		return array[2]
	}
}


// 接收单个文件
func recvFile() {

}

// 接收文件
func recvFiles() {

}

// 接收DirMap
func recvDirMap() bool {
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
	
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("net.Dial failed. err:", err)
		return false
	}

	recvDirMap()


}

func printUsage() {
	usage := "Usage: dirsyncClient.exe -serverAddr ip:port -dir need/recv/dir"
	fmt.Println(usage)
}

func main() {
	flag.Parse()

	if rootDir == "" || serverAddr == "" {
		PrintUsage()
		return
	}

	startTCPClient()
}