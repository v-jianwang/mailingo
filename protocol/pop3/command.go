package pop3

import (
	"log"
	"strings"
)


type Command struct {
	Name string
	Args []string
}

func NewCommand(cmd string) *Command {
	//"STAT"
	//"LIST 1"
	//"TOP 2 10"
	var array = strings.Split(cmd, " ")
	if len(array) <= 0 {
		log.Panic("Command parsed error: %v", cmd)
	}

	var length = len(array) - 1
	var command = &Command{
		Name: cmd,
		Args: make([]string, length),
	}

	for i, v := range array {
		if i == 0 {
			command.Name = v
		} else {
			command.Args[i-1] = v
		}
	}
	return command
}
