package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type MailSender struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type MailRecipient = MailSender

type SendMailOptions struct {
	Subject    string          `json:"subject,omitempty"`
	Message    string          `json:"message,omitempty"`
	Template   MailTemplate    `json:"template,omitempty"`
	Sender     MailSender      `json:"sender,omitempty"`
	Recipients []MailRecipient `json:"recipients,omitempty"`
}

type MailTemplate struct {
	ID        string            `json:"id"`
	Variables map[string]string `json:"variables"`
}

type SendMailResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	Data       struct {
		Ref string `json:"ref"`
	} `json:"data"`
}

type SendService service

func (s *SendService) SendMail(ctx context.Context, options SendMailOptions) (SendMailResponse, error) {
	var resp SendMailResponse

	url := s.client.config.BaseUrl.JoinPath("/send")

	jsonBody, _ := json.Marshal(options)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), bytes.NewBuffer(jsonBody))
	if err != nil {
		return resp, err
	}

	_, err = s.client.sendRequest(req, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
