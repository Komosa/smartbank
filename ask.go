package main

import (
	"os"

	"github.com/Komosa/go-input"
)

type credentials struct {
	user, pass string
}

func askCredentials() (credentials, error) {
	var c credentials
	var err error
	ui := &input.UI{Writer: os.Stdout, Reader: os.Stdin}
	c.user, err = ui.Ask("NIK", &input.Options{Required: true, HideOrder: true, Loop: true})
	if err != nil {
		return c, WrapErr(err)
	}

	c.pass, err = ui.Ask("Password", &input.Options{Required: true, HideOrder: true, Loop: true, Mask: true})
	return c, WrapErr(err)
}
