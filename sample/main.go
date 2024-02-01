package main

import (
	"context"
	"fmt"
	"log"
	"os"

	lib "github.com/olad5/smtpexpress-client-go/lib"
)

func main() {
	projectId := os.Getenv("PROJECT_ID")
	client := lib.NewAPIClient(projectId, &lib.Config{})
	ctx := context.Background()
	opts := lib.SendMailOptions{
		Message: "<h1> Welcome to the future of Email Delivery - message 34</h1>",
		Subject: "golang-sdk test subject - 1",
		Sender: lib.MailSender{
			Email: os.Getenv("SENDER_EMAIL"),
			Name:  "smtpexpress-client-go",
		},
		Recipients: os.Getenv("RECIPIENT_EMAIL"),
	}
	res, err := client.SendApi.SendMail(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nin the main function: ", res)
}
