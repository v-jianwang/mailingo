package usage

import (
	"log"
	"net"
	"strconv"
	"time"
	
	"github.com/v-jianwang/mailingo/protocol"
)

type Usage struct {
	Protocol string
	Port int
	InactiveTimeout time.Duration
}


func (u Usage) Launch() {
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

		go func(c *net.Conn, p string, t time.Duration) {
			log.Println("Hello to a new client", (*c).RemoteAddr())

			// create inactive timer for this connection
			protocol.Serve(c, p, t)
			
			log.Println("Goodbye to the client", (*c).RemoteAddr())
		}(&conn, u.Protocol, u.InactiveTimeout)
	}
}
