package pop3


type User struct {
	Username string
	Password string
}

func (u *User) Authenticate() bool {
	if u.Username == "jiang" && u.Password == "Password" {
		return true
	}
	return false
}