package pop3

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/v-jianwang/mailingo/utils"
)

const (
	maildropRoot = "F:\\Root\\temp\\email"
	mailLockingKey = "locking_maildrops"
	pathSeparator = string(os.PathSeparator)
)

type Maildrop struct {
	Username string
	UsageID string
	Mails []*Mail
}


func (md *Maildrop) Open() error {
	dirname := maildropRoot + pathSeparator + md.Username
	f, err := os.Open(dirname)
	if err != nil {
		return err
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return err
	}
	
	var mails []*Mail
	for _, file := range files {
		if dir := file.IsDir(); !dir {
			mail := &Mail{
				Number: len(mails) + 1,
				Title: file.Name(),
				Size: file.Size(),
				Deleted: false,
				Fullname: func(name string) string {
					return dirname + pathSeparator + name 
				},
			}
			mails = append(mails, mail)
		}
	}
	md.Mails = mails
	return nil
}


func (md Maildrop) Count() int {
	count := 0
	for _, mail := range md.Mails {
		if !mail.Deleted {
			count++
		}
	}
	return count
}


func (md Maildrop) Size() int64 {
	var totalSize int64
	for _, mail := range md.Mails {
		if !mail.Deleted {
			totalSize += mail.Size	
		}
	}
	return totalSize
}


func (md Maildrop) Lock() error {
	var locking []string
	var user string = md.Username

	b, err := json.Marshal(locking)
	if err != nil {
		return err
	}

	state, _ := utils.Stated(md.UsageID)
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


func (md Maildrop) Unlock() error {
	var locking []string
	var user string = md.Username

	state, _ := utils.Stated(md.UsageID)
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


func (md Maildrop) GetMail(num int, ignoreDeleted bool) *Mail {
	for _, mail := range md.Mails {
		if mail.Number == num && 
			(ignoreDeleted || !mail.Deleted) {
			return mail
		}
	}
	return nil
}


func (md Maildrop) ResetMails() {
	for _, mail := range md.Mails {
		mail.Deleted = false
	}	
}


func (md Maildrop) RemoveMails(ignoreDeleted bool) {
	for _, mail := range md.Mails {
		if ignoreDeleted || mail.Deleted {
			err := mail.Delete()
			if err != nil {
				println(err.Error())
			}
		}
	}
}