package usage

import (
	"log"
	"net"
	"strconv"
	"time"
	
	"github.com/v-jianwang/mailingo/protocol"
)


type Usage struct {
	ID string
	Protocol string
	Port int
	InactiveTimeout time.Duration
	// context map[string]interface{}
	// ctxMux *sync.Mutex
}

// var (
// 	context map[string][]byte
// )

// func Context(key string, defaul []byte) ([]byte, bool) {
// 	b, ok := context[key]
// 	if !ok {
// 		context[key] = defaul
// 		b := defaul
// 	}
// 	return b, ok
// }

// func SetContext(key string, b []byte) {
// 	context[key] = b
// }

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

		go func() {
			log.Println("Hello to a new client", conn.RemoteAddr())

			// create inactive timer for this connection
			protocol.Serve(&conn, u.ID, u.Protocol, u.InactiveTimeout)
			
			log.Println("Goodbye to the client", conn.RemoteAddr())
		}()
	}
}

// func (u Usage) Context(key string) (interface{}, bool) {
// 	v, ok := u.context[key]
// 	return v, ok
// }

// func (u Usage) SetContext(key string, ctx interface{}) {
// 	u.context[key] = ctx
// }

// func (u Usage) Lock() {
// 	u.ctxMux.Lock()
// }

// func (u Usage) Unlock() {
// 	u.ctxMux.Unlock()
// }
