package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Komosa/errors"
)

const AppURL = BaseURL + "/bim-webapp/bi/UIDL"

var seemsDone = errors.New("seems done")

func (b *bank) csv() error {
	err := b.doCSVInitPage()
	if err != nil {
		return err
	}

	output := &bytes.Buffer{}
	var nextHdr []byte
	var hdrLen int

	defer func(err *error) {
		e2 := ioutil.WriteFile("output.csv", output.Bytes(), 0644)
		if *err == nil {
			*err = e2
		}
	}(&err)

	doYear := func(y int, lastMonth time.Month) error {
		for m := lastMonth; m >= time.January; m-- {
			err := b.doMonth(y, m)
			if err != nil {
				return err
			}
			resp, err := b.doCSV()
			if err != nil {
				return err
			}
			io.ReadFull(resp.Body, nextHdr) // skip header
			n, err := io.Copy(output, resp.Body)
			if err != nil {
				return WrapErr(err)
			}
			if n < 5 {
				// arbitrary 'too small' treshold
				return seemsDone
			}
			err = resp.Body.Close()
			if err != nil {
				return WrapErr(err)
			}
			if !bytes.Equal(output.Bytes()[:hdrLen], nextHdr) {
				return errors.Errorf("headers mismatch: %q vs %q", output.Bytes()[:hdrLen], nextHdr)
			}
		}
		return nil
	}

	y, m := time.Now().Year(), time.Now().Month()
	// get first month data
	err = b.doMonth(y, m)
	if err != nil {
		return err
	}
	var resp *http.Response
	resp, err = b.doCSV()
	if err != nil {
		return err
	}

	_, err = io.Copy(output, resp.Body)
	if err != nil {
		return WrapErr(err)
	}
	err = resp.Body.Close()
	if err != nil {
		return WrapErr(err)
	}

	// and recognize csv header
	hdrLen = 1 + bytes.IndexByte(output.Bytes(), '\n')
	if hdrLen == 0 {
		return WrapErr(io.ErrUnexpectedEOF)
	}
	nextHdr = make([]byte, hdrLen)

	err = doYear(y, m-1)
	if err != nil {
		return err
	}

	for y--; y > 2010; y-- {
		fmt.Fprintf(os.Stderr, "%d\n", y)
		err = doYear(y, time.December)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *bank) doCSV() (*http.Response, error) {
	return b.doReq(
		"GET",
		BaseURL+"/bim-webapp/bi/APP/5/",
		nil,
	)
}

func (b *bank) doMonth(y int, m time.Month) error {
	resp, err := b.doReq(
		"POST",
		AppURL+"?windowName=1",
		strings.NewReader(prepQuery(y, m, b.secKey)),
	)
	Close(resp.Body, &err)
	return err
}

func prepQuery(y int, m time.Month, secKey string) string {
	last := time.Date(y, m, 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, -1).Day()
	format1 := "\x1d%d\x1fPID_Stext_date_%s\x1fyear\x1fi\x1e%d\x1fPID_Stext_date_%s\x1fmonth\x1fi\x1e%d\x1fPID_Stext_date_%s\x1fday\x1fi"
	return fmt.Sprintf(
		"%s"+format1+format1,
		secKey,
		y, "from", m, "from", 1, "from",
		y, "to", m, "to", last, "to",
	)
}

func (b *bank) doCSVInitPage() error {
	resp, err := b.doReq(
		"POST",
		AppURL+"?windowName=1",
		strings.NewReader(b.secKey+
			"\x1d911\x1fPID0\x1fheight\x1fi\x1e961\x1fPID0\x1fwidth\x1fi\x1e961\x1fPID0\x1fbrowserWidth"+
			"\x1fi\x1e911\x1fPID0\x1fbrowserHeight\x1fi\x1esmouseDetails\x1c1,715,46,"+
			"false,false,false,false,8,113,26\x1cpcomponent\x1cPID15\x1fPID_Stab_history\x1flayout_click\x1fm"),
	)
	Close(resp.Body, &err)
	return err
}
