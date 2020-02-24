package upload

import (
	"bufio"
	"github.com/AlisherFozilov/file-server/cmd/client/internal/base"
	gc "github.com/AlisherFozilov/file-server/pkg/globconst"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
)

var CmdUpload = &base.Command{
	Name:             "upload",
	Run:              runUpload,
	UsageLine:        "client upload 'filepath'",
	ShortDescription: "upload specified file to server",
}

func runUpload(_ *base.Command, args []string) {

	conn, err := net.Dial("tcp", gc.Address)
	if err != nil {
		log.Fatalf("can't listen %v: %v", gc.Address, err)
	}

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
