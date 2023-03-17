package http

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type Client struct {
	client *http.Client
}

func (c *Client) Get(url string, header map[string]string, body []byte) (rspCode int, rspBody []byte, err error) {
	return c.Do(http.MethodGet, url, header, body)
}

func (c *Client) Post(url string, header map[string]string, body []byte) (rspCode int, rspBody []byte, err error) {
	return c.Do(http.MethodPost, url, header, body)
}

func (c *Client) Do(method, url string, header map[string]string, body []byte) (rspCode int, rspBody []byte, err error) {
	var req *http.Request

	if len(body) > 0 {
		req, err = http.NewRequest(method, url, bytes.NewReader(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	response, err := c.client.Do(req)

	if err != nil {
		return
	}

	rspCode = response.StatusCode

	if rspCode != http.StatusOK {
		err = fmt.Errorf("error http code %d", rspCode)
		return
	}

	rBody := response.Body

	if response.Header.Get("Content-Encoding") == "gzip" {
		rBody, err = gzip.NewReader(response.Body)
		if err != nil {
			fmt.Println("http resp unzip is failed,err: ", err)
		}
	}

	rspBody, err = io.ReadAll(rBody)

	return
}

func NewClient(timeout time.Duration) *Client {
	cookie, _ := cookiejar.New(nil)
	return &Client{client: &http.Client{Jar: cookie, Timeout: timeout}}
}
