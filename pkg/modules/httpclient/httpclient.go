package httpclient

import (
	"net/http"
)

type Client interface {
	Send(*http.Request) (*http.Response, error)
}

type httpClient struct {
	cl *http.Client
}

func NewClient() Client {
	return &httpClient{
		cl: &http.Client{},
	}
}

func (cl *httpClient) Send(req *http.Request) (*http.Response, error) {
	resp, err := cl.cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp, nil
}
