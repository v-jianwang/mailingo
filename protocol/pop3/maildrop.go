package pop3

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/v-jianwang/mailingo/utils"
)

const (
	maildropRoot = "F:\\Root\\temp\\email"
	mailLockingKey = "locking_maildrops"
)

type Maildrop struct {
	Username string
	UsageID string
	Mails []Mail
}


func (m *Maildrop) Open() error {
	path := maildropRoot + "\\" + m.Username
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	
	var mails []Mail
	for _, file := range files {
		if dir := file.IsDir(); !dir {
			mail := Mail{
				Number: len(mails) + 1,
				Size: file.Size(),
				Deleted: false,
			}
			mails = append(mails, mail)
		}
	}
	m.Mails = mails
	return nil
}


func (m Maildrop) Count() int {
	count := 0
	for _, mail := range m.Mails {
		if !mail.Deleted {
			count++
		}
	}
	return count
}


func (m Maildrop) Size() int64 {
	var totalSize int64
	for _, mail := range m.Mails {
		if !mail.Deleted {
			totalSize += mail.Size	
		}
	}
	return totalSize
}


func (m Maildrop) Lock() error {
	var locking []string
	var user string = m.Username

	b, err := json.Marshal(locking)
	if err != nil {
		return err
	}

	state, _ := utils.Stated(m.UsageID)
	state.Locker.Lock()
	defer state.Locker.Unlock()

	if item, ok := state.Item(mailLockingKey, b); ok {
		err := json.Unmarshal(item, &locking)
		if err != nil {
			return err
		}

		for _, name := range locking {
			if name == user {
				// the maildrop has been locked
				return errors.New("maildrop is being used")
			}
		}
	}

	locking = append(locking, user)
	b, err = json.Marshal(locking)
	if err != nil {
		return err
	}

	state.SetItem(mailLockingKey, b)
	return nil
}


func (m Maildrop) Unlock() error {
	var locking []string
	var user string = m.Username

	state, _ := utils.Stated(m.UsageID)
	state.Locker.Lock()
	defer state.Locker.Unlock()

	if item, ok := state.Item(mailLockingKey, nil); ok {
		err := json.Unmarshal(item, &locking)
		if err != nil {
			return err
		}

		var index int
		for i, name := range locking {
			if name == user {
				index = i
			}
		}
		locking = append(locking[:index], locking[index+1:]...)
	}

	b, err := json.Marshal(locking)
	if err != nil {
		return err
	}

	state.SetItem(mailLockingKey, b)
	return nil
}