package main

import (
	"github.com/v-jianwang/mailingo/usage"
)

func main() {

	mailingo := usage.NewUsageContainer()
	mailingo.NewUsage("pop3", 110)
	mailingo.NewUsage("imap4", 143)

	println("pop3 and imap4 servers're running...")

	c := make(chan bool)
	<-c
}