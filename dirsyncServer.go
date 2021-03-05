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
	// for flag
	rootDir       string  // 需要传输的文件夹路径

	totalFileNum  uint64
	sendedFileNum uint64
	usedTime      uint64 // 传输文件的时间

	dirMapFilePath string
	dirMapFilePtr  *os.File
	dirMapFileBuf  *bufio.Writer
	recordNum      uint64

	dirFlag  string
	fileFlag string

	listenAddr string
)

func init() {
	recordNum = 0
	dirFlag = "D"
	fileFlag = "F"

	listenAddr = "127.0.0.1:13344"

	dirMapFilePath = "./dirMap.txt"
	// TODO: os.ModePerm 不合适，需要换成普通文件权限
	dirMapFilePtr, err := os.OpenFile(dirMapFilePath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil { // TODO 在init函数里打开文件是否合适
		fmt.Printf("open dirMapFile Failed. %v\n", err)
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
			fmt.Printf("send failed: %s error:%v\n", path, err)
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
		fmt.Println("error:", err)
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

func checkFileExist(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Println("no such file or dir:", path)
		return false
	}
	if fileInfo.IsFile() {
		return true
	} else {
		fmt.Println("path is a dir:", path)
		return false
	}
}

// 给客户端发送 dirMap 文件
func sendDirMap(conn net.Conn) {
	fmt.Println("sending dirMap file:", dirMapFilePath)

	bResult := sendFile(dirMapFilePath, conn)
	return bResult
}

// 等待客户端请求文件夹内容
func processDirRequest(conn net.Conn) bool {
	buf := make([]byte, 4096)
	for {
		// 等待客户端发过来文件请求或者请求完成信息
		n, err := conn.Read(buf)
		if err != nil {  // TODO: 服务端程序直接退出不合适，可以考虑记录所有传输失败的文件，稍后进行重传尝试
			fmt.Println("conn read err:", err)
			return false
		}

		path := string(buf[:n])

		// 收到客户端传递的传输完成信息
		if path == "IFinishPleaseCloseMe" {
			return true
		}
		
		// 校验文件是否合法
		if bOk := checkFileExist(path); !bResult {
			return false
		}

		// 传输文件
		if bOk := SendFile(path, conn); !bOk {
			return false
		}
	}
}

// 启动TCP服务器，等待客户端的连接
func startTCPServer() bool {
	fmt.Printf("Starting TCP server...")

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Println("create listener failed. err:", err)
		return false
	}

	fmt.Println("waiting client...")

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("accept connection failed. err:", err)
		return false
	}
	defer conn.Close()

	sendDirMap(conn)
	fmt.Println("send dirMap finish. waiting client request files...")

	bOk = processDirRequest(conn)
	
	return bOk
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
	usage := "Usage: dirsyncServer.exe -dir need/send/dir"
	fmt.Println(usage)
}

func main() {
	flag.Parse()

	if rootDir == "" {
		PrintUasge()
		return
	}

	// debug
	fmt.Printf("rootDir: %s\n", rootDir)

	bResult := createDirMap()
	if bResult == false {
		return
	}

	startTCPServer()
}
