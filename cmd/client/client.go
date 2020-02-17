package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const commandsList = `Here's list of commands:
upload
download
list`

func main() {
	//flag.Parse()
	//args := flag.Args()

	args := []string{"list"}
	if len(args) > 2 {
		fmt.Println("too many arguments")
		return
	}

	if len(args) == 0 {
		fmt.Println(commandsList)
		return
	}

	command := args[0]
	path := ""
	if len(args) > 1 {
		path = args[1]
	}

	const address = "localhost:9999"
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("can't listen %v: %v", address, err)
	}
	defer conn.Close()
	writer := bufio.NewWriter(conn)
	switch command {
	case "upload":
		file, err := os.Open(path)
		if err != nil {
			log.Fatal("can't open file: ", err)
		}
		defer file.Close()
		reader := bufio.NewReaderSize(file, 1024*1024)

		_, err = writer.WriteString("0")
		if err != nil {
			log.Fatal(err)
		}
		err = writer.Flush()
		if err != nil {
			log.Fatal(err)
		}
		_, err = writer.WriteString(filepath.Base(path) + "\n")
		if err != nil {
			log.Fatal("can't send filename: ", err)
		}
		writer.Flush()
		_, err = io.Copy(writer, reader)
		if err != nil {
			log.Fatalf("can't send file: %v", err)
		}
	case "download":
		err = writer.WriteByte('1')
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

		file, err := os.OpenFile("downloaded.txt", os.O_CREATE, 0666)
		if err != nil {
			log.Fatalf("can't open file %v: %v", path, err)
		}
		defer file.Close()
		fileWriter := bufio.NewWriterSize(file, 1024*1024)

		_, err = io.Copy(fileWriter, conn)
		if err != nil {
			log.Fatal("can't download: ", err)
		}
		fileWriter.Flush()
	case "list":
		err = writer.WriteByte('2')
		if err != nil {
			log.Fatal(err)
		}
		writer.Flush()
		buf := &bytes.Buffer{}
		_, err := io.Copy(buf, conn)
		if err != nil {
			log.Fatalf("can't read from %v: %v", address, err)
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
