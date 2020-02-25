package download

import (
	"bufio"
	"fmt"
	"github.com/AlisherFozilov/file-server/cmd/client/internal/base"
	"github.com/AlisherFozilov/file-server/cmd/client/internal/execute"
	sch "github.com/AlisherFozilov/file-server/cmd/client/internal/status-codes-handling"
	gc "github.com/AlisherFozilov/file-server/pkg/globconst"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

var CmdDownload = &base.Command{
	Name:             "download",
	Run:              runDownload,
	UsageLine:        "client download 'filename'",
	ShortDescription: "download file from server",
}

func runDownload(cmd *base.Command, args []string) {

	if len(args) < 1 {
		fmt.Println("Usage:")
		fmt.Println(cmd.UsageLine)
		fmt.Println("Description")
		fmt.Println(cmd.ShortDescription)
		return
	}

	path := args[0]
	filename := filepath.Base(path)
	if filename == "" {
		log.Print("not a file: ", path)
		return
	}

	execute.HandleConn(gc.Address, func(conn net.Conn) {
		writer := bufio.NewWriter(conn)
		err := writer.WriteByte(gc.Download)
		if err != nil {
			log.Printf("can't write %v to connection: %s", gc.Download, err)
			return
		}

		_, err = writer.WriteString(filename + "\n")
		if err != nil {
			log.Print("can't send filename: ", err)
			return
		}
		err = writer.Flush()
		if err != nil {
			log.Print("can't flush to connection: ", err)
			return
		}

		reader := bufio.NewReader(conn)
		//
		statusCode := make([]byte, gc.BytesInStatusCode)
		_, err = reader.Read(statusCode)
		if err != nil {
			log.Printf("server does not response: %s", err)
			return
		}

		if sch.IsStatusCodeError(statusCode) {
			err := conn.Close()
			if err != nil {
				log.Printf("can't close connection: %s", err)
			}
			fmt.Println(sch.HandleStatusCode(statusCode))
			os.Exit(0)//return or os.Exit(0) ?
		}
		//
		readBytes, err := reader.ReadBytes('\n')
		if err != nil {
			log.Printf("can't read fileSize from connection: %s", err)
			return
		}
		fileSizeStr := string(readBytes[:len(readBytes)-1])
		fileSize, err := strconv.Atoi(fileSizeStr)
		if err != nil {
			log.Printf("can't convert %s to int: %s", fileSizeStr, err)
			return
		}

		file, err := os.OpenFile(filename, os.O_CREATE, 0666)
		if err != nil {
			log.Printf("can't open file %s: %s", filename, err)
			return
		}
		defer func() {
			err := file.Close()
			if err != nil {
				log.Print("can't close file: ", err)
			}
		}()

		_, err = io.CopyN(file, conn, int64(fileSize))
		if err != nil {
			log.Printf("can't download file %s: %s", filename, err)
			return
		}
	})

	return
}
