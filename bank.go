package main

import "net/http"

type bank struct {
	client *http.Client
	secKey string
}

func newBank() (*bank, error) {
	b := &bank{}
	return b, b.initClient()
}
