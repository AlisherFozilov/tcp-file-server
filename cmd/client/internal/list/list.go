package list

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/AlisherFozilov/file-server/cmd/client/internal/base"
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

func runList(_ *base.Command, args []string) {

	conn, err := net.Dial("tcp", gc.Address)
	if err != nil {
		log.Fatalf("can't listen %v: %v", gc.Address, err)
	}

	writer := bufio.NewWriter(conn)

	err = writer.WriteByte(gc.List)
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
