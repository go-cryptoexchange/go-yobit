package yobit

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type client struct {
	apiKey      string
	apiSecret   string
	httpClient  *http.Client
	httpTimeout time.Duration
}

// NewClient return a new Bittrex HTTP client
func NewClient(apiKey, apiSecret string) (c *client) {
	return &client{apiKey, apiSecret, &http.Client{}, 30 * time.Second}
}

// NewClientWithCustomHttpConfig returns a new Bittrex HTTP client using the predefined http client
func NewClientWithCustomHttpConfig(apiKey, apiSecret string, httpClient *http.Client) (c *client) {
	timeout := httpClient.Timeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &client{apiKey, apiSecret, httpClient, timeout}
}

// NewClient returns a new Bittrex HTTP client with custom timeout
func NewClientWithCustomTimeout(apiKey, apiSecret string, timeout time.Duration) (c *client) {
	return &client{apiKey, apiSecret, &http.Client{}, timeout}
}

// doTimeoutRequest do a HTTP request with timeout
func (c *client) doTimeoutRequest(timer *time.Timer, req *http.Request) (*http.Response, error) {
	// Do the request in the background so we can check the timeout
	type result struct {
		resp *http.Response
		err  error
	}
	done := make(chan result, 1)
	go func() {
		resp, err := c.httpClient.Do(req)
		done <- result{resp, err}
	}()
	// Wait for the read or the timeout
	select {
	case r := <-done:
		return r.resp, r.err
	case <-timer.C:
		return nil, errors.New("timeout on reading data from Bittrex API")
	}
}

// do prepare and process HTTP request to Bittrex API
func (c *client) do(method string, ressource string, payload string) (response []byte, err error) {
	connectTimer := time.NewTimer(c.httpTimeout)

	var rawurl string
	if strings.HasPrefix(ressource, "http") {
		rawurl = ressource
	} else {
		rawurl = fmt.Sprintf("%s/%s/%s", API_BASE, API_VERSION, ressource)
	}

	req, err := http.NewRequest(method, rawurl, strings.NewReader(payload))
	if err != nil {
		return
	}

	resp, err := c.doTimeoutRequest(connectTimer, req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	response, err = ioutil.ReadAll(resp.Body)
	// fmt.Println(fmt.Sprintf("reponse %s", response), err)
	if err != nil {
		return response, err
	}
	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
	}
	return response, err
}
