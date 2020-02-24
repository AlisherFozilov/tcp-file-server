package download

import (
	"bufio"
	"github.com/AlisherFozilov/file-server/cmd/client/internal/base"
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

func runDownload(_ *base.Command, args []string) {

	conn, err := net.Dial("tcp", gc.Address)
	if err != nil {
		log.Fatalf("can't listen %v: %v", gc.Address, err)
	}

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	path := args[0]

	err = writer.WriteByte(gc.Download)
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
