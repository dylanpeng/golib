package http

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type ClientProxy struct {
	client   *http.Client
	proxyUrl string
}

func (c *ClientProxy) Get(url string, header map[string]string, body []byte) (rspCode int, rspBody []byte, err error) {
	return c.Do(http.MethodGet, url, header, body)
}

func (c *ClientProxy) Post(url string, header map[string]string, body []byte) (rspCode int, rspBody []byte, err error) {
	return c.Do(http.MethodPost, url, header, body)
}

func (c *ClientProxy) Do(method, url string, header map[string]string, body []byte) (rspCode int, rspBody []byte, err error) {
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

	defer response.Body.Close()

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

func NewClientProxy(proxyUrl string, timeout time.Duration) *Client {
	cookie, _ := cookiejar.New(nil)

	proxy, _ := url.Parse(proxyUrl)
	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &Client{client: &http.Client{Jar: cookie, Transport: tr, Timeout: timeout}}
}
