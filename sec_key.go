package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http/httputil"
	"strings"

	"github.com/Komosa/errors"
)

func (b *bank) grabSecKey() error {
	resp, err := b.doReq(
		"POST",
		AppURL+"?repaintAll=1&wsver=9.9.9",
		strings.NewReader("init"),
	)
	if err != nil {
		return err
	}
	defer Close(resp.Body, &err)

	data, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return WrapErr(err)
	}
	err = ioutil.WriteFile("out", data, 0644)
	if err != nil {
		return WrapErr(err)
	}

	_, err = io.ReadFull(resp.Body, []byte("for(;;);[")) // skip prefix
	if err != nil {
		return WrapErr(err)
	}

	var msg struct {
		SecKey string `json:"Vaadin-Security-Key"`
	}
	err = json.NewDecoder(resp.Body).Decode(&msg)
	if err != nil {
		return WrapErr(err)
	}

	if len(msg.SecKey) != 36 {
		return errors.Errorf("invalid security key length: %d", len(msg.SecKey))
	}

	b.secKey = msg.SecKey
	return nil
}
