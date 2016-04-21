package protocol

import (
	"fmt"
	"log"
	"strconv"

	"github.com/v-jianwang/mailingo/protocol/pop3"
)


type State string
const (
	StateNotSet State = "NotSet"
	StateAuthorization = "Authorization"
	StateTransaction = "Transaction"
	StateUpdate = "Update"
)

const (
	crlf = "\r\n"
	plusOK = "+OK"
	negaERR = "-ERR"
)


type HandlerPop3 struct {
	UsageID string
	CurrentState State
	CurrentCommand *pop3.Command
	CanQuit bool

	AcceptableCommands map[State][]string

	User *pop3.User
	Maildrop *pop3.Maildrop
}

func NewHandlerPop3(usageID string) HandlerPop3 {
	h := HandlerPop3{
		UsageID: usageID,
		CurrentState: StateNotSet,
		CurrentCommand: nil,
		CanQuit: false,

		AcceptableCommands: make(map[State][]string),

		User: new(pop3.User),
		Maildrop: new(pop3.Maildrop),
	}

	// define acceptable commands by state
	acceptCommand(&h, StateAuthorization, "USER", "PASS", "QUIT")
	acceptCommand(&h, StateTransaction, "STAT", "LIST", "RETR", "DELE", "NOOP", "RSET", "QUIT","TOP")

	return h
}


func acceptCommand(h *HandlerPop3, s State, cmds ...string) {
	h.AcceptableCommands[s] = cmds
}

func (h HandlerPop3) Handle(base BaseHandler) error {
	var msg, line, keyword string
	var ok bool
	var err error
	// var readErrCount, writeErrCount int = 0, 0

	for {

		keyword = "GREET"
		log.Printf("Current State: %v\n", h.CurrentState)

		if h.CurrentState != StateNotSet {
			line, err = base.Read()

			if err != nil {
				if err.Error() == "EOF" || base.IsClosed() {
					break
				}
				log.Printf("Read in pop3.Handler error[%t]: %v\n", err, err)
				break
			} else {
				base.Active()
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
			case "USER":
				msg = h.user()
			case "PASS":
				msg = h.pass()
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

		err = base.Write(msg)
		if err != nil {
			log.Printf("Write in pop3.Handler %q error: %v\n", msg, err)
			break
		}		
	
		if h.CanQuit {
			break
		}
	}

	h.dispose()
	base.Dispose()

	return err
}

func (h *HandlerPop3) dispose() {
	h.Maildrop.Unlock()
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
	return fmt.Sprint(plusOK, " server is ready", crlf)
}


func (h *HandlerPop3) quit() string {
	h.CanQuit = true
	return fmt.Sprint(plusOK, " server is signing off", crlf)
}


func (h *HandlerPop3) user() string {
	cmd := h.CurrentCommand
	if len(cmd.Args) != 1 {
		return fmt.Sprint(negaERR, " argument username is invalid", crlf)
	}
	h.User.Username = cmd.Args[0]
	return fmt.Sprint(plusOK, crlf)
}


func (h *HandlerPop3) pass() string {
	cmd := h.CurrentCommand
	if len(cmd.Args) != 1 {
		return fmt.Sprint(negaERR, " argument password is invalid", crlf)
	}
	h.User.Password = cmd.Args[0]
	if ok := h.User.Authenticate(); !ok {
		return fmt.Sprint(negaERR, " permission denied", crlf)
	}
	
	username := h.User.Username
	maildrop := &pop3.Maildrop{
		Username: username,
		UsageID: h.UsageID,
	}
	if err := maildrop.Lock(); err != nil {
		log.Printf("lock %s's maildrop error: %v\n", username, err)
		return fmt.Sprint(negaERR, " maildrop is being used", crlf)
	}

	if err := maildrop.Open(); err != nil {
		log.Printf("Maildrop.Open error: %v", err, crlf)
		maildrop.Unlock()
		return fmt.Sprint(negaERR, " failed to open maildrop")
	}

	h.Maildrop = maildrop
	h.CurrentState = StateTransaction

	count := maildrop.Count()
	size := maildrop.Size()
	msg := fmt.Sprintf(" %s's maildrop has %d message(s) (%d octets)", username, count, size)
	return fmt.Sprint(plusOK, msg, crlf)
}


func (h HandlerPop3) stat() string {
	if h.Maildrop == nil {
		return fmt.Sprint(negaERR, " client must identify itself at first", crlf)
	}
	count := h.Maildrop.Count()
	size := h.Maildrop.Size()
	return fmt.Sprint(plusOK, " ", count, size, crlf)
}


func (h HandlerPop3) list() string {
	if h.Maildrop == nil {
		return fmt.Sprint(negaERR, " client must identify itself at first", crlf)
	}

	cmd := h.CurrentCommand
	maildrop := h.Maildrop
	count := maildrop.Count()

	if len(cmd.Args) == 0 {
		var msg string
		if count > 0 {
			msg = fmt.Sprintf(" %d message(s) (%d octets)", maildrop.Count(), maildrop.Size())
			msg = fmt.Sprint(plusOK, msg, crlf)
			for _, mail := range maildrop.Mails {
				if !mail.Deleted {
					msg += fmt.Sprint(mail.Number, mail.Size, crlf)
				}
			}
		} else {
			msg = fmt.Sprint(plusOK, crlf)
		}
		// dot ending
		msg += fmt.Sprint(".", crlf)
		return msg
	}

	n, err := strconv.Atoi(cmd.Args[1])
	if err != nil {
		return fmt.Sprint(negaERR, " argument is invalid")
	}

	var index int = -1
	for i, m := range maildrop.Mails {
		if !m.Deleted && m.Number == n {
			index = i
		}
	}

	if index > -1 {
		return fmt.Sprint(plusOK, " ", n, maildrop.Mails[index].Size)
	} else {
		return fmt.Sprint(negaERR, " no such message")
	}
}


func (h HandlerPop3) retr() string {
	return ""
}


func (h HandlerPop3) dele() string {
	return ""
}


func (h HandlerPop3) noop() string {
	return fmt.Sprint(plusOK, crlf)
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
	msg := fmt.Sprintf(" command %s is unrecognized", command)
	return fmt.Sprint(negaERR, msg, crlf)
}


func (h HandlerPop3) incorrect_stat() string {
	command := h.CurrentCommand.Name
	state := h.CurrentState

	msg := fmt.Sprintf(" command %s is in incorrect state %s", command, state)
	return fmt.Sprint(negaERR, msg, crlf)
}

