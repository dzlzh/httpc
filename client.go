package httpc

import (
	"net/http"
	"net/url"
	"time"
)

type ClientOption func(*Client)

type Client struct {
	client    *http.Client
	transport *http.Transport
}

func NewClient(options ...ClientOption) *Client {
	tr := &http.Transport{}

	c := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}

	client := &Client{client: c, transport: tr}

	for _, option := range options {
		option(client)
	}

	return client
}

func (client *Client) SetClient(c *http.Client) *Client {
	client.client = c

	return client
}

func (client *Client) SetTransport(tr *http.Transport) *Client {
	client.transport = tr
	client.client.Transport = tr

	return client
}

func (client *Client) SetProxy(proxyUrl string) *Client {
	proxy, _ := url.Parse(proxyUrl)
	client.transport.Proxy = http.ProxyURL(proxy)

	return client
}

func (client *Client) SetTimeout(t time.Duration) *Client {
	client.client.Timeout = t
	return client
}

func Proxy(proxyUrl string) ClientOption {
	return func(client *Client) {
		client.SetProxy(proxyUrl)
	}
}

func Timeout(t time.Duration) ClientOption {
	return func(client *Client) {
		client.SetTimeout(t)
	}
}
