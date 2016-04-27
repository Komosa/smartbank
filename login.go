package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/Komosa/errors"
)

const LoginURL = BaseURL + "/bim-webapp/smart/login"

func (b *bank) login(cred credentials) error {
	pass, err := encPass(cred.pass)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		LoginURL,
		strings.NewReader("&password="+pass+"&nik="+cred.user),
	)
	if err != nil {
		return WrapErr(err)
	}
	setPlainHeaders(req)
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := b.client.Do(req)
	if err != nil {
		return err
	}
	defer Close(resp.Body, &err)

	if !jsonStatus(resp.Body) {
		return errors.New("invalid password")
	}
	return err
}

func jsonStatus(stream io.Reader) bool {
	var status struct {
		Status string `json:"status"`
	}
	err := json.NewDecoder(stream).Decode(&status)
	return err == nil && status.Status == "OK"
}
