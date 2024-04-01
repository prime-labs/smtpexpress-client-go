package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	lib "github.com/prime-labs/smtpexpress-client-go/lib"
)

func main() {
	projectSecret := os.Getenv("PROJECT_SECRET")
	client := lib.CreateClient(projectSecret, &lib.Config{})
	ctx := context.Background()
	opts := lib.SendMailOptions{
		Message: fmt.Sprintf("<h1> Welcome to the future of Email Delivery - message sent at %s </h1>", time.Now().String()),
		Subject: "golang-sdk test subject",
		Sender: lib.MailSender{
			Email: os.Getenv("SENDER_EMAIL"),
			Name:  "smtpexpress-client-go",
		},
		Recipients: []lib.MailRecipient{
			{
				Email: os.Getenv("RECIPIENT_EMAIL"),
				Name:  "end user",
			},
		},
	}
	res, err := client.Send.SendMail(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n send mail was a success: ", res)
}
