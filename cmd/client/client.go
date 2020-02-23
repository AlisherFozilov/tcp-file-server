package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	gc "github.com/AlisherFozilov/file-server/pkg/globconst"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const commandsList = `Here's list of gc:
upload
download
list`

//var command = flag.Bool("download", false, "command to do")


func main() {

	flag.Parse()
//	fmt.Println(*command)
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println(commandsList)
		return
	}

	command := args[0]
	var function func(net.Conn, []string)

	switch command {
	case "upload":
		function = upload
	case "download":
		function = download
	case "list":
		function = list
	default:
		fmt.Println(`wrong command`)
		fmt.Println(commandsList)
		return
	}

	conn, err := net.Dial("tcp", gc.Address)
	if err != nil {
		log.Fatalf("can't listen %v: %v", gc.Address, err)
	}

	function(conn, args[1:])
}

func upload(conn net.Conn, args []string) {
	writer := bufio.NewWriter(conn)
	//reader := bufio.NewReader(conn)

	path := args[0]
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("can't open file: ", err)
	}

	defer file.Close()
	//reader := bufio.NewReaderSize(file, 1024*1024*100)

	err = writer.WriteByte(gc.Upload)
	if err != nil {
		log.Fatal("can't write gc to buffer: ", err)
	}
	err = writer.Flush()
	if err != nil {
		log.Fatal("can't send gc: ", err)
	}
	_, err = writer.WriteString(filepath.Base(path) + "\n")
	if err != nil {
		log.Fatal("can't send filename: ", err)
	}
	writer.Flush()
	_, err = io.Copy(writer, file)
	if err != nil {
		log.Fatalf("can't send file: %v", err)
	}
}
func download(conn net.Conn, args []string) {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	path := args[0]

	err := writer.WriteByte(gc.Download)
	if err != nil {
		log.Fatal(err)
	}
	_, err = writer.WriteString(filepath.Base(path) + "\n")
	if err != nil {
		log.Fatal("can't send filename: ", err)
	}
	err = writer.Flush()
	if err != nil {
		log.Fatal(err)
	}

	readBytes, err := reader.ReadBytes('\n')
	fileSize, err := strconv.Atoi(string(readBytes[:len(readBytes)-1]))
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile("downloaded.txt", os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("can't open file %v: %v", path, err)
	}
	defer file.Close()
	fileWriter := bufio.NewWriterSize(file, 1024*1024)

	_, err = io.CopyN(fileWriter, conn, int64(fileSize))
	if err != nil {
		log.Fatal("can't download: ", err)
	}
	fileWriter.Flush()
}
func list(conn net.Conn, _ []string) {
	writer := bufio.NewWriter(conn)

	err := writer.WriteByte(gc.List)
	if err != nil {
		log.Fatal(err)
	}
	err = writer.Flush()
	if err != nil {

	}
	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, conn)
	if err != nil {
		log.Fatalf("can't read from %v: %v", gc.Address, err)
	}

	fileNames := strings.Split(buf.String(), " ")
	for _, fileName := range fileNames {
		fmt.Println(fileName)
	}
}