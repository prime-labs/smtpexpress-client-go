# SMTP Express client go
A golang package to send emails using SMTP Express


## Installation

```bash
go get github.com/prime-labs/smtpexpress-client-go/lib
```

## Usage

```go
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
	projectSecret := "<<REPLACE_WITH_YOUR_PROJECT_SECRET>>"
	client := lib.CreateClient(projectSecret, &lib.Config{})
	ctx := context.Background()
	opts := lib.SendMailOptions{
		Message: fmt.Sprintf("<h1> Welcome to the future of Email Delivery - message sent at %s </h1>", time.Now().String()),
		Subject: "smtpexpress-client-go Mail Subject",
		Sender: lib.MailSender{
			Email: "<<REPLACE_WITH_YOUR_SENDER_EMAIL>>",
			Name:  "smtpexpress-client-go",
		},
		Recipients: []lib.MailRecipient{
			{
				Email: "<<REPLACE_WITH_YOUR_RECIPIENT_EMAIL>>",
				Name:  "Jane Doe",
			},
		},
	}
	res, err := client.Send.SendMail(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n send mail was a success: ", res)
}
```


## Bugs/Contributions/Requests
If you encounter any problems using this package, please feel free to open an [issue](https://github.com/prime-labs/smtpexpress-client-go/issues).

If you'd like to contribute to this package, kindly open a pull request [here](https://github.com/prime-labs/smtpexpress-client-go)


## TODO

- [ ] Add more tests
- [ ] Add more examples
- [ ] Add more documentation
- [ ] add Attachments
- [ ] add Calendar invites
