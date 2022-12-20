package indexer

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Http interface {
	Post(path string, object interface{}) ([]byte, bool)
	Get(path string) ([]byte, bool)
}

type HttpClient struct {
	Host string
	username string
	password string
	client http.Client
}

func InitHttpClient(host string, username string, password string) *HttpClient {
	httpClient := HttpClient{
	  Host: host,
	  username: username,
	  password: password,
	}
	httpClient.client = http.Client{}
	return &httpClient
}

func (c *HttpClient) Post(path string, object interface{}) ([]byte, bool) {
	var json *bytes.Buffer
	if object != nil {
		jsonData := toJson(&object)
		json = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest("POST", c.Host + path, json)
	if err != nil {
		log.Fatal("cannot make request with url: ", c.Host + path)
	}
	req.Header.Set("Content-Type", "application/json")
	c.setBasicAuth(req)
	return c.sendRequest(req)
}

func (c *HttpClient) Get(path string) ([]byte, bool) {
	req, err := http.NewRequest("GET", c.Host + path, nil)
	if err != nil {
		log.Fatal("cannot make request with url: ", c.Host + path)
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
		return body, false
	}
}

func (c *HttpClient) setBasicAuth(req *http.Request) {
	req.SetBasicAuth(c.username, c.password)
}

func toJson(object *interface{}) []byte {
	jsonData, err := json.Marshal(object)
	if err != nil {
		log.Fatal(err)
	}
	return jsonData
}