package list

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/AlisherFozilov/file-server/cmd/client/internal/base"
	"github.com/AlisherFozilov/file-server/cmd/client/internal/execute"
	gc "github.com/AlisherFozilov/file-server/pkg/globconst"
	"io"
	"log"
	"net"
	"strings"
)

var CmdList = &base.Command{
	Name:             "list",
	Run:              runList,
	UsageLine:        "client list",
	ShortDescription: "show all files on server",
}

func runList(_ *base.Command, _ []string) {

	execute.HandleConn(gc.Address, func(conn net.Conn) {
		writer := bufio.NewWriter(conn)

		err := writer.WriteByte(gc.List)
		if err != nil {
			log.Printf("can't write %v to writer: %s", gc.List, err)
			return
		}
		err = writer.Flush()
		if err != nil {
			log.Printf("can't flush to connection: %s", err)
			return
		}

		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, conn)
		if err != nil {
			log.Printf("can't read from %v: %v", gc.Address, err)
			return
		}

		fileNames := strings.Split(buf.String(), " ")
		for _, fileName := range fileNames {
			fmt.Println(fileName)
		}
	})
}
