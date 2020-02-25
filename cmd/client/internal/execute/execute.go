package execute

import (
	"log"
	"net"
)

func HandleConn(addr string, ex func(conn net.Conn)) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("can't connect to %s: %s", addr, err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Print("can't close connection: ", err)
		}
	}()
	ex(conn)
}
