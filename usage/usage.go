package usage

import (
	"log"
	"net"
	"strconv"
	"github.com/v-jianwang/mailingo/protocol"
)

type Usage struct {
	Protocol string
	Port int
}


func (u *Usage) Launch() {

	addr := ":" + strconv.Itoa(u.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Listen to %s error %v", addr, err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Accept from %s error %v", addr, err)
			continue
		}

		
		go func(c net.Conn, p string) {
			log.Println("Receive a connection from client", c.RemoteAddr())

			protocol.Serve(&c, p)
			c.Close()
			
			log.Println("Close a connection from client", c.RemoteAddr())
		}(conn, u.Protocol)
	}
}
