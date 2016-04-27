package main

import (
	"io/ioutil"

	"github.com/robertkrimen/otto"
)

func encPass(pass string) (string, error) {
	buf, err := ioutil.ReadFile("login.js")
	if err != nil {
		return "", err
	}

	rnd := NewRand(1024)
	vm := otto.New()
	vm.SetRandomSource(func() float64 { return rnd.Float64() })

	val, err := vm.Run(`
		var navigator = {appName:"x", appVersion:"3"};
		var window = {_init:"x"};
		` + string(buf) + `
		var a="` + pass + `",b=new RSAKey;
b.setPublic("b703da6dfaf6a07f34bd936158c5dc8c42e0601caa5b4beeba1412a0d3f898c1655d14f90dda1deea2ec2aac89fee7fd867634327e03f3f737b779118ad01115c5344e981c5bc169fa11bbaad8efcd5f28639066719e6d403a7fae13e8ba88ee0a88edbd28b8802b5e514bdcef837a379f3268ee190b80a3fd7cd484314fa045","10001");
		_encode(b.encrypt(a));`,
	)
	if rnd.Err() != nil {
		return "", WrapErr(rnd.Err())
	}
	return val.String(), WrapErr(err)
}
