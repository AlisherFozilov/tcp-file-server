package main

import (
	"bufio"
	gc "github.com/AlisherFozilov/file-server/pkg/globconst"
	sc "github.com/AlisherFozilov/file-server/pkg/status-codes"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const serverDirName = "serverdata"

func init() {
	err := os.Mkdir(serverDirName, 0666)
	if err != nil {
		if !os.IsExist(err) {
			log.Fatalf("can't create directory: %s", err)
		}
	}
}

// 1) как отделить логи от сообщений, выводимых для пользователя
// 2) как замечать серьёзные ошибки в логах

func main() {
	listener, err := net.Listen("tcp", gc.Address)
	if err != nil {
		log.Fatal("can't listen: ", err)
	}
	defer func() {
		err := listener.Close()
		if err != nil {
			log.Printf("can't close Listener: %s", err)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("can't accept connection: ", err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Printf("can't close connection %s", err)
		}
	}()

	command := make([]byte, 1)
	_, err := conn.Read(command)
	if err != nil {
		log.Println("can't get command: ", err)
		return
	}

	switch command[0] {
	case gc.Upload:
		log.Println("upload start")
		upload(conn)
		log.Println("upload end")
	case gc.Download:
		log.Println("download start")
		download(conn)
		log.Println("download end")
	case gc.List:
		log.Println("list start")
		list(conn)
		log.Println("list end")
	default:
		log.Println("unknown command")
		return
	}
}

func upload(conn net.Conn) {

	reader := bufio.NewReader(conn)
	filename, err := reader.ReadString('\n')
	if err != nil {
		log.Println("can't get filename: ", err)
		return
	}

	filename = filename[:len(filename)-1]

	file, err := os.OpenFile(serverDirName+"/"+filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("can't open file %v: %v", filename, err)
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("can't close file %s: %s", filename, err)
		}
	}()

	readBytes, err := reader.ReadBytes('\n')
	if err != nil {
		log.Fatalf("can't read fileSize from connection: %s", err)
	}
	fileSizeStr := string(readBytes[:len(readBytes)-1])
	fileSize, err := strconv.Atoi(fileSizeStr)
	if err != nil {
		log.Fatalf("can't convert %s to int: %s", fileSizeStr, err)
	}

	fileWriter := bufio.NewWriter(file)
	_, err = io.CopyN(fileWriter, reader, int64(fileSize))
	if err != nil {
		log.Printf("can't write to file %s: %s", filename, err)
		return
	}

	err = fileWriter.Flush()
	if err != nil {
		log.Printf("can't flush to file %s: %s", filename, err)
		return
	}
}

func download(conn net.Conn) {

	reader := bufio.NewReader(conn)
	filename, err := reader.ReadString('\n')
	if err != nil {
		log.Println("can't get filename: ", err)
	}

	filename = filename[:len(filename)-1]

	file, err := os.Open(serverDirName+"/"+filename)
	if err != nil {
		log.Printf("can't open file %s: %s", filename, err)
		if os.IsNotExist(err) {
			_, err := conn.Write(sc.FILE_NOT_EXISTS)
			if err != nil {
				log.Printf("can't write %v to connection: %s",
					sc.FILE_NOT_EXISTS, err)
				return
			}
		}
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("can't close file %s: %s", filename, err)
		}
	}()

	_, err = conn.Write(sc.OK)
	if err != nil {
		log.Printf("can't write %v to connection: %s", sc.OK, err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("can't get %s's info: %s", filename, err)
		return
	}

	fileSize := fileInfo.Size()
	fileSizeStr := strconv.Itoa(int(fileSize))

	_, err = conn.Write([]byte(fileSizeStr))
	if err != nil {
		log.Printf("can't write fileSize to connection: %s", err)
		return
	}
	_, err = conn.Write([]byte{'\n'})
	if err != nil {
		log.Printf(`can't write '\n' to connection: %s`, err)
		return
	}

	_, err = io.Copy(conn, file)
	if err != nil {
		log.Printf("error while sending file: %s", err)
		return
	}
}
func list(conn net.Conn) {
	dirents, err := ioutil.ReadDir(serverDirName)
	if err != nil {
		log.Printf("can't get dirents for path %s/: %s", serverDirName, err)
		return
	}

	filenames := make([]string, 0)
	for _, entry := range dirents {
		if !entry.IsDir() {
			filenames = append(filenames, entry.Name())
		}
	}

	filenamesStr := strings.Join(filenames, " ")
	_, err = conn.Write([]byte(filenamesStr))
	if err != nil {
		log.Printf("can't write filenames (%s) to connection: %s", filenamesStr, err)
		return
	}
}
