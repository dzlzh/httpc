package httpc

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	client   *Client
	request  *http.Request
	response *http.Response
	method   string
	url      string
	headers  map[string]string
	cookies  *[]*http.Cookie
	body     io.Reader
	query    url.Values
	form     url.Values
	json     []byte
	debug    bool
	err      error
	ch       chan struct{}
}

func NewRequest(c *Client) *Request {
	return &Request{
		client:  c,
		method:  "GET",
		headers: map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36"},
		cookies: new([]*http.Cookie),
		query:   url.Values{},
		form:    url.Values{},
	}
}

func (r *Request) SetMethod(method string) *Request {
	r.method = strings.ToUpper(method)
	return r
}

func (r *Request) SetURL(url string) *Request {
	r.url = url
	return r
}

func (r *Request) SetHeader(key, value string) *Request {
	r.headers[key] = value
	return r
}

func (r *Request) SetHeaders(headers map[string]string) *Request {
	r.headers = headers
	return r
}

func (r *Request) SetCookies(cookies *[]*http.Cookie) *Request {
	r.cookies = cookies
	return r
}

func (r *Request) SetDebug(debug bool) *Request {
	r.debug = debug
	return r
}

func (r *Request) SetQuery(key, value string) *Request {
	r.query.Set(key, value)
	return r
}

func (r *Request) SetBody(body io.Reader) *Request {
	r.body = body
	return r
}

func (r *Request) SetForm(key, value string) *Request {
	r.form.Set(key, value)
	return r
}

func (r *Request) SetJson(json []byte) *Request {
	r.json = json
	return r
}

func (r *Request) Reset() *Request {
	r.json = []byte("")
	r.form = url.Values{}
	r.body = nil
	return r
}

func (r *Request) Send() *Request {
	var err error
	url, err := url.Parse(r.url)
	if err != nil {
		r.err = err
		return r
	}
	url.RawQuery = r.query.Encode()

	body := r.body
	if body == nil {
		if r.json != nil {
			r.headers["Content-Type"] = "application/json"
			body = bytes.NewReader(r.json)
		} else if len(r.form) > 0 {
			r.headers["Content-Type"] = "application/x-www-form-urlencoded"
			body = strings.NewReader(r.form.Encode())
		}
	}

	r.request, err = http.NewRequest(r.method, url.String(), body)
	if err != nil {
		r.err = err
		return r
	}

	for k, v := range r.headers {
		r.request.Header.Set(k, v)
	}

	for _, v := range *r.cookies {
		r.request.AddCookie(v)
	}

	r.ch = make(chan struct{}, 1)
	go func() {
		r.response, r.err = r.client.client.Do(r.request)
		r.ch <- struct{}{}
		close(r.ch)
	}()
	return r
}

func (r *Request) End() (*http.Response, []byte, error) {
	<-r.ch
	if r.err != nil {
		return nil, []byte(""), r.err
	}
	var bodyByte []byte

	if r.response.Header.Get("Content-Encoding") == "gzip" {
		reader, _ := gzip.NewReader(r.response.Body)
		defer reader.Close()
		bodyByte, _ = ioutil.ReadAll(reader)
	} else {
		bodyByte, _ = ioutil.ReadAll(r.response.Body)
	}
	_ = r.response.Body.Close()
	return r.response, bodyByte, nil
}
