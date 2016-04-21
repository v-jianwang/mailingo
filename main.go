package main

import (
	"time"
	"github.com/v-jianwang/mailingo/usage"
)

func main() {

	mailingo := usage.NewUsageContainer()
	mailingo.NewUsage("pop3", 110, 30 * time.Second)
	mailingo.NewUsage("imap4", 143, time.Minute)

	println("pop3 and imap4 servers're running...")

	// time.Sleep(60 * time.Second)
	c := make(chan bool)
	<-c
}