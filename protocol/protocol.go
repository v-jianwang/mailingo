package protocol

import (
	"bufio"
	"fmt"
	"log"
	"net"
)


type Handler interface {
	Handle(ReadWriter) error
}

type Protocol struct {
	Conn *net.Conn
}

type ReadWriter interface {
	Read() (string, error)
	Write(msg string) error
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

func newProtocol(c *net.Conn) Protocol {
	return Protocol{
		Conn: c,
	}
}

func newHandler(name string) Handler {
	var h Handler
	switch name {
		case "pop3":
			h = NewHandlerPop3()
		case "imap4":
			h = HandlerImap4{}
	}

	return h
}

func Serve(c *net.Conn, name string) {
	p := newProtocol(c)
	handler := newHandler(name)

	err := handler.Handle(p)
	if err != nil {
		log.Fatalf("Handle %s error: %v", name, err)
	}
}
