package main

import (
	"fmt"
	"util"
	"net"
	"os"
)

var (
	rootDir string  // ���洫����ļ��еĸ�Ŀ¼

	totalFileNum uint64
	recvivedFileNum uint64
	usedTime uint64

	dirMapFilePath string
	dirMapFilePtr *os.File
	dirMapFileBuf *bufio.Writer

	lastFinishRecord uint64 // ��һ��������ɵ��ļ������

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

// �� dirMap �ļ������ö�������
func openDirMapFile() bool {
	dirMapFilePtr, err := os.OpenFile(dirMapFilePath, os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println("open dirMapFile Failed. err:", err)
		return false
	}
	dirMapFileBuf = bufio.NewReader()
	return true
}

// ƴ��·��
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

// ������һ���ļ���¼
func getNextRecord() string {
	// ѭ��ֱ���ҵ���һ���ļ���¼
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
		
		// ���������¼��Ŀ¼���򴴽�Ŀ¼
		if array[1] == dirFlag {
			dirPath := joinPath(array[2])
			err := os.MkDirAll(dirPath)
			continue
		}
		return array[2]
	}
}


// ���յ����ļ�
func recvFile() {

}

// �����ļ�
func recvFiles() {

}

// ����DirMap
func recvDirMap() bool {
	fmt.Println("start recv dirMap file...")

	buf := make([]byte, 4096)

	for {
		n, err := conn.Read(buf)
		if n == 0 {
			// TODO: ��Ҫ�� buff ����д���ļ�
			return true
		}
		
		if err != nil {
			fmt.Println("recv dirMap failed. err:", err)
			return false
		}

	}

}

// ����TCP�ͻ���
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