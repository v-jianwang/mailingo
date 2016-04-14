package protocol

import (
	"fmt"
	"log"	
	"github.com/v-jianwang/mailingo/protocol/pop3"
)


type State string
const (
	StateNotSet State = "NotSet"
	StateAuthorization State = "Authorization"
	StateTransaction State = "Transaction"
	StateUpdate State = "Update"
)


type HandlerPop3 struct {
	CurrentState State
	CurrentCommand *pop3.Command
	CanQuit bool

	AcceptableCommands map[State][]string
}

func NewHandlerPop3() HandlerPop3 {
	h := HandlerPop3{
		CurrentState: StateNotSet,
		CurrentCommand: nil,
		CanQuit: false,

		AcceptableCommands: make(map[State][]string),
	}

	h.acceptCommand(StateAuthorization, "USER", "PASS", "QUIT")
	h.acceptCommand(StateTransaction, "STAT", "LIST", "RETR", "DELE", "NOOP", "RSET", "QUIT","TOP")

	return h
}

func (h HandlerPop3) acceptCommand(s State, cmds ...string) {
	h.AcceptableCommands[s] = cmds
}

func (h HandlerPop3) Handle(rw ReadWriter) error {
	var msg, line, keyword string
	var ok bool
	var err error

	for {

		keyword = "GREET"
		log.Printf("Current State: %v\n", h.CurrentState)
		if h.CurrentState != StateNotSet {
			line, err = rw.Read()
			if err != nil {
				log.Fatalf("Read in pop3.Handler error: %v", err)
			}

			h.CurrentCommand = pop3.NewCommand(line)
			keyword = h.CurrentCommand.Name
			// Check if keyword of commmand is unrecognized or unimplemented
			if keyword, ok = h.checkCommand(keyword); ok {
				// Check if keyword of commmand is in correct state
				keyword, _ = h.checkCommandState(keyword)
			}
		}
		println("keyword: " + keyword)
		switch keyword {
			case "GREET":
				msg = h.greet()
			case "QUIT":
				msg = h.quit()
			case "STAT":
				msg = h.stat()
			case "LIST":
				msg = h.list()
			case "RETR":
				msg = h.retr()
			case "DELE":
				msg = h.dele()
			case "NOOP":
				msg = h.noop()
			case "RSET":
				msg = h.rset()
			case "TOP":
				msg = h.top()
			case "INVALID_CMD":
				msg = h.invalid_command()
			case "INCORRECT_STAT":
				msg = h.incorrect_stat()
		}		

		err = rw.Write(msg)
		if err != nil {
			log.Fatalf("Write in pop3.Handler %q error: %v", msg, err)
		}		
	
		if h.CanQuit {
			break
		}
	}
	return nil	
}


func (h HandlerPop3) checkCommand(k string) (string, bool) {
	for _, commands := range h.AcceptableCommands {
		for _, cmd := range commands {
			if cmd == k {
				return k, true
			}
		}
	}

	return "INVALID_CMD", false
}


func (h HandlerPop3) checkCommandState(k string) (string, bool) {
	commands := h.AcceptableCommands[h.CurrentState]
	for _, cmd := range commands {
		if cmd == k {
			return k, true
		}
	}
	return "INCORRECT_STAT", false
}

func (h *HandlerPop3) greet() string {
	h.CurrentState = StateAuthorization
	return fmt.Sprintf("%s %s\r\n", "+OK", "server is ready")
}

func (h *HandlerPop3) quit() string {
	h.CanQuit = true
	return fmt.Sprintf("%s %s\r\n", "+OK", "server is signing off")
}

func (h HandlerPop3) stat() string {
	return ""
}

func (h HandlerPop3) list() string {
	return ""
}

func (h HandlerPop3) retr() string {
	return ""
}

func (h HandlerPop3) dele() string {
	return ""
}

func (h HandlerPop3) noop() string {
	return fmt.Sprintf("%s\r\n", "+OK")
}

func (h HandlerPop3) rset() string {
	return ""
}

func (h HandlerPop3) top() string {
	return ""
}

func (h HandlerPop3) invalid_command() string {
	command := h.CurrentCommand.Name
	// response to unrecognized, unimplemented or invalid command
	msg := fmt.Sprintf("command %s is unrecognized", command)
	return fmt.Sprintf("%s %s\r\n", "-ERR", msg)
}

func (h HandlerPop3) incorrect_stat() string {
	command := h.CurrentCommand.Name
	state := h.CurrentState

	msg := fmt.Sprintf("command %s is in incorrect state %s", command, state)
	return fmt.Sprintf("%s %s\r\n", "-ERR", msg)
}

