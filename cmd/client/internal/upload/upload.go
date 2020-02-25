package upload

import (
	"bufio"
	"fmt"
	"github.com/AlisherFozilov/file-server/cmd/client/internal/base"
	"github.com/AlisherFozilov/file-server/cmd/client/internal/execute"
	gc "github.com/AlisherFozilov/file-server/pkg/globconst"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

var CmdUpload = &base.Command{
	Name:             "upload",
	Run:              runUpload,
	UsageLine:        "client upload 'filepath'",
	ShortDescription: "upload specified file to server",
}

func runUpload(cmd *base.Command, args []string) {

	if len(args) < 1 {
		fmt.Println("Usage:")
		fmt.Println(cmd.UsageLine)
		fmt.Println("Description")
		fmt.Println(cmd.ShortDescription)
		return
	}

	path := args[0]
	file, err := os.Open(path)
	if err != nil {
		log.Printf("can't open file %s: %s", path, err)
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("can't close file %s: %s", path, err)
		}
	}()

	filename := filepath.Base(path)
	if filename == "" {
		log.Print("not a file: ", path)
		return
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("can't get %s's info: %s", filename, err)
		return
	}
	fileSize := fileInfo.Size()
	fileSizeStr := strconv.Itoa(int(fileSize))

	execute.HandleConn(gc.Address, func(conn net.Conn) {
		_, err = conn.Write([]byte{gc.Upload})
		if err != nil {
			log.Printf("can't write Upload %v to connection: %s", gc.Upload, err)
			return
		}

		writer := bufio.NewWriter(conn)
		_, err = writer.WriteString(filename + "\n")
		if err != nil {
			log.Printf("can't write filename %s to buffer: %s", filename, err)
			return
		}
		err = writer.Flush()
		if err != nil {
			log.Printf("can't flush to connection: %s", err)
			return
		}

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

		_, err = io.Copy(writer, file)
		if err != nil {
			log.Printf("can't sent file %s to server: %s", filename, err)
			return
		}
	})

}
