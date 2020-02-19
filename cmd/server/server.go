package main

import (
	"bufio"
	gc "github.com/AlisherFozilov/file-server/pkg/globconst"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", gc.Address)
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

	switch command {
	case gc.Upload:
		filename, err := reader.ReadString('\n')
		if err != nil {
			log.Println("can't get filename: ", err)
		}
		filename = filename[:len(filename)-1]
		file, err := os.OpenFile("copy-"+filename,
			os.O_CREATE | os.O_WRONLY, 0666)

		if err != nil {
			log.Printf("can't open file %v: %v", filename, err)
			return
		}
		defer file.Close()
		fileWriter := bufio.NewWriter(file)

		_, err = io.Copy(fileWriter, reader)
		if err != nil {
			log.Println("can't get file: ", err)
			return
		}
		err = fileWriter.Flush()
		if err != nil {
			log.Println("can't write file: ", file)
		}
	case gc.Download:
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
		defer file.Close()

		open, err := os.Open(filename)
		if err != nil {
			log.Println(err)
			return
		}

		stat, err := open.Stat()
		if err != nil {
			log.Println(err)
			return
		}
		fileSize := stat.Size()
		fileSizeStr := strconv.Itoa(int(fileSize))

		conn.Write([]byte(fileSizeStr))
		conn.Write([]byte{'\n'})

		_, err = io.Copy(conn, file)
		if err != nil {
			log.Println(err)
			return
		}
	case gc.List:
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
