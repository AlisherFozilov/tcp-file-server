package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

const address = "0.0.0.0:9999"

func main() {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("can't listen: ", err)
	}
	defer listener.Close()
	
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("can't accept connection: ", err)
			continue
		}
		handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	command, err := reader.ReadByte()
	if err != nil {
		log.Println("can't get command: ", err)
		return
	}
	//command := commands[0]
	switch command {
	case '0':
		filename, err := reader.ReadString('\n')
		if err != nil {
			log.Println("can't get filename: ", err)
		}
		filename = filename[:len(filename)-1]
		file, err := os.OpenFile("copy-"+filename, os.O_CREATE | os.O_WRONLY, 0666)
		if err != nil {
			log.Printf("can't open file %v: %v", filename, err)
			return
		}
		defer file.Close()
		fileWriter := bufio.NewWriterSize(file, 1024*1024)

		_, err = io.Copy(fileWriter, reader)
		if err != nil {
			log.Println("can't get file: ", err)
			return
		}
		fileWriter.Flush()
	case '1':
		filename, err := reader.ReadString('\n')
		if err != nil {
			log.Println("can't get filename: ", err)
		}
		filename = filename[:len(filename)-1]

		file, err := os.Open(filename)
		if err != nil {
			log.Println(err)
			return
		}

		reader := bufio.NewReaderSize(file, 1024*1024)
		_, err = io.Copy(conn, reader)
		if err != nil {
			log.Println(err)
			return
		}
	case '2':
		log.Println("list start")
		dirents, err := ioutil.ReadDir(".")
		if err != nil {
			log.Println(err)
			return
		}

		filenames := make([]string, 0)
		log.Println("for loop start")
		for _, entry := range dirents {
			if !entry.IsDir() {
				filenames = append(filenames, entry.Name())
			}
		}
		log.Println("for loop end")
		filenamesStr := strings.Join(filenames, " ")
		_, err = conn.Write([]byte(filenamesStr))
		if err != nil {
			log.Println(err)
			return
		}
		//fmt.Println(filenamesStr)
		log.Println("list end")
	default:
		return
	}
}
