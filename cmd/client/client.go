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


func main() {
	flag.Parse()
	args := flag.Args()

	//args := []string{"list"} // for debugging

	switch len(args) {
	case 1:
		oneArg(args)
	case 2:
		twoArg(args)
	case 0:
		fmt.Println(commandsList)
		return
	default:
		fmt.Println("too many arguments")
		return
	}

}

func oneArg(args []string) {
	conn, err := net.Dial("tcp", gc.Address)
	if err != nil {
		log.Fatalf("can't listen %v: %v", gc.Address, err)
	}
	defer conn.Close()

	command := args[0]
	writer := bufio.NewWriter(conn)

	switch command {
	case "list":
		err = writer.WriteByte(gc.List)
		if err != nil {
			log.Fatal(err)
		}
		err := writer.Flush()
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
	default:
		fmt.Println(`wrong command`)
		fmt.Println(commandsList)
	}
}

func twoArg(args []string) {
	conn, err := net.Dial("tcp", gc.Address)
	if err != nil {
		log.Fatalf("can't listen %v: %v", gc.Address, err)
	}
	defer conn.Close()

	command := args[0]
	path := args[1]

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)
	switch command {
	case "upload":
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
	case "download":
		err = writer.WriteByte(gc.Download)
		if err != nil {
			log.Fatal(err)
		}
		_, err = writer.WriteString(filepath.Base(path) + "\n")
		if err != nil {
			log.Fatal("can't send filename: ", err)
		}
		err := writer.Flush()
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

	default:
		fmt.Println(`wrong command`)
		fmt.Println(commandsList)
	}

}
