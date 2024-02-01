package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/pkg/errors"
)

var SmtpExpressURL string = os.Getenv("BASE_URL")

type Config struct {
	BaseUrl    *url.URL
	HttpClient *http.Client
}

type service struct {
	client *APIClient
}

type APIClient struct {
	projectSecret string
	config        *Config
	common        service

	SendApi *SendService
}

func NewAPIClient(projectSecret string, cfg *Config) *APIClient {
	cfg.BaseUrl = buildBaseURL(cfg)

	if cfg.HttpClient == nil {
		cfg.HttpClient = &http.Client{Timeout: 20 * time.Second}
	}

	c := &APIClient{projectSecret: projectSecret}
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.projectSecret))

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
