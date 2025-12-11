package main

import (
	"fmt"
	"taskflow/cmd"
)

func main() {
	fmt.Println("hello world")
	smtpAcccount := cmd.SMTPAccount{
		Name:     "google",
		Host:     "smtp.gmail.com",
		Port:     587,
		Auth:     "PLAIN",
		User:     "nepalsaurav123@gmail.com",
		Password: "eyce vgah wnqi ugzo",
		From:     "nepalsaurav123@gmail.com",
	}
	cmd.SetPostfixConfig(smtpAcccount)
}
