package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"net"
)

var (
	rootDir       string
	totalFileNum  uint64
	sendedFileNum uint64
	usedTime      uint64 // 传输文件的时间

	dirMapFilePath string
	dirMapFilePtr  *os.File
	dirMapFileBuf  *bufio.Writer
	recordNum      uint64

	dirFlag  string
	fileFlag string
)

func init() {
	recordNum = 0
	dirFlag = "D"
	fileFlag = "F"

	dirMapFilePath = "./dirMap.txt"
	// TODO: os.ModePerm 不合适，需要换成普通文件权限
	dirMapFilePtr, err := os.OpenFile(dirMapFilePath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil { // TODO 在init函数里打开文件是否合适
		fmt.Printf("open dirMapFile Failed. %v", err)
		return
	}

	dirMapFileBuf = bufio.NewWriter(dirMapFilePtr)

	flag.StringVar(&rootDir, "dir", "", "root dir need send")
}

func getRecordNum() uint64 {
	recordNum++
	return recordNum
}

// 展示文件夹的传输进度
func ShowProgress() {
	fmt.Print("\x1b7")   // 保存光标位置 保存光标和Attrs <ESC> 7
	fmt.Print("\x1b[2k") // 清空当前行的内容 擦除线<ESC> [2K
	fmt.Printf("sendeFile/totalFile: %s/%s\n", sendedFileNum, totalFileNum)
	fmt.Print("\x1b8") // 恢复光标位置 恢复光标和Attrs <ESC> 8
}

// 给客户端传输单个文件
func SendFile(path string, conn net.Conn) bool {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("load file failed. path:%s\n", path)
		return false
	}
	defer file.Close()

	buf := make([]byte, 4096)

	// 循环读取文件内容，写入客户端
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			fmt.Printf("Send finish: %s\n", path)
			return true
		}

		_, err = conn.Write(buf[:n])
		if err != nil {
			fmt.Printf("send failed: %s error:%v", path, err)
			return false
		}
	}

	return true
}

func writeRecord(record string) {
	dirMapFileBuf.WriteString(record)
}

// 处理文件遍历过程中的单条记录
func processRecord(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Printf("error:%v", err)
		return err
	}

	var record string

	if info.IsDir() {
		record = fmt.Sprintf("%d %s %s\n", getRecordNum(), dirFlag, path)
	} else {
		record = fmt.Sprintf("%d %s %s\n", getRecordNum(), fileFlag, path)
	}
	writeRecord(record)

	return nil
}

// 给客户端发送 dirMap 文件
func SendDirMap(conn net.Conn) {
	fmt.Printf("sending dirMap file: %s", dirMapFilePath)

	_, err := co
}

// 启动TCP服务器，等待客户端的连接
func startTcpServer() bool {
	fmt.Printf("Starting TCP server...")
	addr := "127.0.0.1:13344"

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("create listener failed. err:", err)
		return false
	}~

	fmt.Println("waiting client...")

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("accept connection failed. err:", err)
		return false
	}

	
	



	fmt.Println("send dirMap finish. waiting client request files...")

	// 等待客户端请求文件

	

}

// 根据传入的目录，创建一个包含该目录下所有文件和目录的路径的文件
func createDirMap() bool {
	fmt.Println("collecting files information. please wait...")

	err := filepath.Walk(rootDir, processRecord)
	if err != nil {
		fmt.Println("generate dir map file failed.")
		return false
	}

	dirMapFileBuf.Flush() // 清空缓冲区
	return true
}

func PrintUasge() {
	usage := "Usage: dirsync.exe -dir need/send/dir"
	fmt.Println(usage)
}

func main() {
	flag.Parse()

	if rootDir == "" {
		PrintUasge()
		return
	}

	// debug
	fmt.Printf("rootDir: %s", rootDir)

	bResult := createDirMap()
	if bResult == false {
		return
	}

}
