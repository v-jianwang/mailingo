package pop3

import (
	// "bytes"
	"os"
	"bufio"
)

type Mail struct {
	Number int
	Title string
	Size int64
	Deleted bool
	Fullname func(string) string
}


func (m Mail) Delete() error {
	fullname := m.Fullname(m.Title)
	return os.Remove(fullname)
}


func (m Mail) Head() (string, error) {
	head := ""
	fullname := m.Fullname(m.Title)
	file, err := os.Open(fullname)
	if err != nil {
		return head, err
	}
	defer file.Close()

	crlf := "\r\n"
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			return head, err
		}
		if len(line) == 0 {
			break
		}
		head += (string(line) + crlf)
	}

	return head, err
}


func (m Mail) Body(n int) (string, error) {
	body := ""
	fullname := m.Fullname(m.Title)
	file, err := os.Open(fullname)
	if err != nil {
		return body, err
	}
	defer file.Close()

	crlf := "\r\n"
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			return body, err
		}
		if len(line) == 0 {
			break
		}
	}


	for i:=0; (i < n || n < 0); i++ {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				err = nil
			}
			return body, err
		}
		body += (string(line) + crlf)
	}

	return body, err	
}