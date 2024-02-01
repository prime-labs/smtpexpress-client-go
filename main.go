package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/pkg/errors"
)

var SmtpExpressURL string = os.Getenv("BASE_URL")

type MailSender struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type SendMailOptions struct {
	Subject    string       `json:"subject"`
	Message    string       `json:"message"`
	Template   MailTemplate `json:"template"`
	Sender     MailSender   `json:"sender"`
	Recipients string       `json:"recipients"`
}

type MailTemplate struct {
	ID        string            `json:"id"`
	Variables map[string]string `json:"variables"`
}

type MailAttachment string

type Data struct {
	Ref        string `json:"ref"`
	StatusCode int    `json:"number"`
	Name       string `json:"name"`
}

type SendMailResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	Data       struct {
		Ref string `json:"ref"`
	} `json:"data"`
}
type Config struct {
	BaseUrl    *url.URL
	HttpClient *http.Client
}

type APIClient struct {
	projectId string
	config    *Config
	common    service

	SendApi *SendService
}

type service struct {
	client *APIClient
}

func NewAPIClient(projectId string, cfg *Config) *APIClient {
	cfg.BaseUrl = buildBaseURL(cfg)

	if cfg.HttpClient == nil {
		cfg.HttpClient = &http.Client{Timeout: 20 * time.Second}
	}

	c := &APIClient{projectId: projectId}
	c.config = cfg
	c.common.client = c

	c.SendApi = (*SendService)(&c.common)
	return c
}

func MustParseURL(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}

	return u
}

func buildBaseURL(cfg *Config) *url.URL {
	if cfg.BaseUrl == nil {
		return MustParseURL(SmtpExpressURL)
	}

	return cfg.BaseUrl.JoinPath("")
}

func (c APIClient) sendRequest(req *http.Request, resp interface{}) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.projectId))

	res, err := c.config.HttpClient.Do(req)
	if err != nil {
		return res, errors.Wrap(err, "failed to execute request")
	}

	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	if res.StatusCode >= http.StatusMultipleChoices {
		return res, errors.Errorf(
			`request was not successful, status code %d, %s`, res.StatusCode,
			string(body),
		)
	}

	if string(body) == "" {
		resp = map[string]string{}
		return res, nil
	}

	err = c.decode(&resp, body)
	if err != nil {
		return res, errors.Wrap(err, "unable to unmarshal response body")
	}

	return res, nil
}

func (c APIClient) decode(v interface{}, b []byte) (err error) {
	if err = json.Unmarshal(b, v); err != nil {
		return err
	}
	return nil
}

type SendService service

type SendMailRequest struct {
	Subject    string     `json:"subject"`
	Message    string     `json:"message"`
	Sender     MailSender `json:"sender"`
	Recipients string     `json:"recipients"`
}

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

func NewConfiguration() *Config {
	cfg := &Config{
		BaseUrl:    MustParseURL(SmtpExpressURL),
		HttpClient: &http.Client{Timeout: 20 * time.Second},
	}
	return cfg
}

func main() {
	cfg := NewConfiguration()
	projectId := os.Getenv("PROJECT_ID")
	client := NewAPIClient(projectId, cfg)
	ctx := context.Background()
	opts := SendMailOptions{
		Message: "<h1> Welcome to the future of Email Delivery - message 34</h1>",
		Subject: "golang-sdk test subject - 1",
		Sender: MailSender{
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
