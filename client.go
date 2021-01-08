package httpc

import (
	"net/http"
	"net/url"
	"time"
)

type HTTPClientOption func(*HTTPClient)

type HTTPClient struct {
	client    *http.Client
	transport *http.Transport
}

func NewHTTPClient(options ...HTTPClientOption) *HTTPClient {
	tr := &http.Transport{}

	client := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}

	httpClient := &HTTPClient{client: client, transport: tr}

	for _, option := range options {
		option(httpClient)
	}

	return httpClient
}

func (httpClient *HTTPClient) SetClient(client *http.Client) *HTTPClient {
	httpClient.client = client

	return httpClient
}

func (httpClient *HTTPClient) SetTransport(tr *http.Transport) *HTTPClient {
	httpClient.transport = tr
	httpClient.client.Transport = tr

	return httpClient
}

func (httpClient *HTTPClient) SetProxy(proxyUrl string) *HTTPClient {
	proxy, _ := url.Parse(proxyUrl)
	httpClient.transport.Proxy = http.ProxyURL(proxy)

	return httpClient
}

func (httpClient *HTTPClient) SetTimeout(t time.Duration) *HTTPClient {
	httpClient.client.Timeout = t
	return httpClient
}

func Proxy(proxyUrl string) HTTPClientOption {
	return func(httpClient *HTTPClient) {
		httpClient.SetProxy(proxyUrl)
	}
}

func Timeout(t time.Duration) HTTPClientOption {
	return func(httpClient *HTTPClient) {
		httpClient.SetTimeout(t)
	}
}
