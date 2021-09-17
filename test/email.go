package main

import (
	"fmt"

	"ginfra/utils"
)

func main() {

	// Sender data.
	from := "xxxx@163.com"
	password := "xxxx"

	// Receiver email address.
	to := []string{
		"xxxx@163.com",
	}

	// smtp server configuration.
	smtpHost := "smtp.163.com"
	smtpPort := "25"

	// Message.
	message := []byte("This is a test email message.")

	err := utils.SendMail(from, password, smtpHost, smtpPort, message, to)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}
