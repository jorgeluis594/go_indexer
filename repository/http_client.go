package repository

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

type Http interface {
	Post(path string, data []byte) ([]byte, bool)
	Get(path string) ([]byte, bool)
}

type HttpClient struct {
	Host     string
	username string
	password string
	client   http.Client
}

func InitHttpClient(host string, username string, password string) *HttpClient {
	httpClient := HttpClient{
		Host:     host,
		username: username,
		password: password,
	}
	httpClient.client = http.Client{}
	return &httpClient
}

func (c *HttpClient) Post(path string, data []byte) ([]byte, bool) {
	var json *bytes.Buffer
	if data != nil {
		json = bytes.NewBuffer(data)
	}

	req, err := http.NewRequest("POST", c.Host+path, json)
	if err != nil {
		log.Fatal("cannot make request with url: ", c.Host+path)
	}
	req.Header.Set("Content-Type", "application/json")
	c.setBasicAuth(req)
	return c.sendRequest(req)
}

func (c *HttpClient) Get(path string) ([]byte, bool) {
	req, err := http.NewRequest("GET", c.Host+path, nil)
	if err != nil {
		log.Fatal("cannot make request with url: ", c.Host+path)
	}
	req.Header.Set("Content-Type", "application/json")
	c.setBasicAuth(req)
	return c.sendRequest(req)
}

func (c *HttpClient) sendRequest(req *http.Request) ([]byte, bool) {
	resp, err := c.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return body, true
	} else {
		log.Printf("Get response error from %s, error message: %s", req.URL, string(body))
		return body, false
	}
}

func (c *HttpClient) setBasicAuth(req *http.Request) {
	req.SetBasicAuth(c.username, c.password)
}
