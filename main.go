package main

import (
	"fmt"
	"os"

	"github.com/Komosa/errors"
)

func WrapErr(err error) error {
	return errors.Wrap(err, "")
}

func main() {
	err := run()
	if err != nil && err != seemsDone {
		fmt.Fprintf(os.Stderr, "err: %+v\n", err)
		os.Exit(1)
	}
}

func run() error {
	b, err := newBank()
	if err != nil {
		return err
	}

	cred, err := askCredentials()
	if err != nil {
		return err
	}

	err = b.login(cred)
	if err != nil {
		return err
	}

	err = b.grabSecKey()
	if err != nil {
		return err
	}

	return b.csv()
}
