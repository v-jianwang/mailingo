package protocol

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)


type Protocol struct {
	Conn *net.Conn
	Inactive time.Duration
	Timer *time.Timer
	Closed *bool

}

type BaseHandler interface {
	Read() (string, error)
	Write(string) error

	Pulse() bool
	IsClosed() bool
	Dispose() bool
}

type Handler interface {
	Handle(BaseHandler) error
}


func (p Protocol) Read() (string, error) {
	return func(c *net.Conn) (string, error) {
		line, _, err := bufio.NewReader((*c)).ReadLine()
		return string(line), err
	}(p.Conn)
}

func (p Protocol) Write(message string) error {
	return func(c *net.Conn, msg string) error {
		_, err := fmt.Fprint((*c), msg)
		return err
	}(p.Conn, message)
}

func (p Protocol) Pulse() bool {
	return p.Timer.Reset(p.Inactive)
}

func (p Protocol) IsClosed() bool {
	return (*p.Closed)
}

func (p Protocol) Dispose() bool {
	var active bool
	if active := p.Timer.Stop(); active {
		c := p.Conn
		(*c).Close()
		(*p.Closed) = true
	}
	return active
}


func newBaseHandler(c *net.Conn, inactive time.Duration) BaseHandler {
	cb := false
	cp := &cb
	t := time.AfterFunc(inactive, func() {
			log.Println("inactive expired")
			(*c).Close()
			(*cp) = true
		})
	p := Protocol{
		Conn: c,
		Inactive: inactive,
		Timer: t,
		Closed: cp,
	}
	return p
}

func newHandler(usageID string, name string) Handler {
	var h Handler
	switch name {
		case "pop3":
			h = NewHandlerPop3(usageID)
		case "imap4":
			h = HandlerImap4{}
	}

	return h
}


func Serve(c *net.Conn, usageID string, name string, inactive time.Duration) {
	base := newBaseHandler(c, inactive)
	handler := newHandler(usageID, name)
	err := handler.Handle(base)
	if err != nil {
		log.Printf("Handle %s error: %v\n", name, err)
	}

}
