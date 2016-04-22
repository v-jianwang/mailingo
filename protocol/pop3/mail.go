package pop3

type Mail struct {
	Number int
	Size int64
	Deleted bool
}


func (m Mail) Remove() error {
	println("remove mail")
	return nil
}