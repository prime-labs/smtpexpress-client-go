package main

import (
	"context"
	"fmt"
	"log"
	"os"

	lib "github.com/olad5/smtpexpress-client-go/lib"
)

func main() {
	projectSecret := os.Getenv("PROJECT_SECRET")
	client := lib.NewAPIClient(projectSecret, &lib.Config{})
	ctx := context.Background()
	opts := lib.SendMailOptions{
		Message: "<h1> Welcome to the future of Email Delivery - message 35</h1>",
		Subject: "golang-sdk test subject - 1",
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
	res, err := client.SendApi.SendMail(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n send mail was a success: ", res)
}
