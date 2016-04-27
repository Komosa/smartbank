package main

import (
	"crypto/rand"
	"math/big"
)

// CryptoRandFloater provides more secure float64's (more than rand.Float64())
// please note, that returned floats have 64 bits anyway...
type CryptoRandFloater struct {
	err error
	mx  *big.Int
}

// NewRand returns generator of random values
// will use given number bits for generation
func NewRand(bits uint) *CryptoRandFloater {
	mx := big.NewInt(1)
	return &CryptoRandFloater{
		mx: mx.Lsh(mx, bits),
	}
}

// Err keeps first occured erorr
// if not nil, returned value(s) were not random
func (c *CryptoRandFloater) Err() error {
	return c.err
}

// Float64() returns next random value
// Check Err() for errors
func (c *CryptoRandFloater) Float64() float64 {
	if c.err != nil {
		return 0
	}
	var rnd *big.Int
	rnd, c.err = rand.Int(rand.Reader, c.mx)
	if c.err != nil {
		return 0
	}
	rat := &big.Rat{}
	rat.SetFrac(rnd, c.mx)
	f, _ := rat.Float64()
	return f
}
