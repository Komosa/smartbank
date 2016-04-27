package main

import (
	"io"
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/publicsuffix"
)

const BaseURL = "https://online.banksmart.pl"

func (b *bank) initClient() error {
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	b.client = &http.Client{Jar: jar}
	return WrapErr(err)
}

func (b *bank) doReq(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, WrapErr(err)
	}
	setPlainHeaders(req)
	return b.client.Do(req)
}

func Close(c io.Closer, err *error) {
	e := c.Close()
	if *err != nil {
		*err = e
	}
}

func setPlainHeaders(req *http.Request) {
	req.Header.Set("Content-type", "text/plain;charset=UTF-8")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36")
}
