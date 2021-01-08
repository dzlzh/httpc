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

func (c *Client) SetClient(client *http.Client) *Client {
	c.client = client

	return c
}

func (c *Client) SetTransport(tr *http.Transport) *Client {
	c.transport = tr
	c.client.Transport = tr

	return c
}

func (c *Client) SetProxy(proxyUrl string) *Client {
	proxy, _ := url.Parse(proxyUrl)
	c.transport.Proxy = http.ProxyURL(proxy)

	return c
}

func (c *Client) SetTimeout(t time.Duration) *Client {
	c.client.Timeout = t
	return c
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
