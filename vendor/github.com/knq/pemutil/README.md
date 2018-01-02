# About pemutil [![Build Status](https://travis-ci.org/knq/pemutil.svg)](https://travis-ci.org/knq/pemutil) [![Coverage Status](https://coveralls.io/repos/knq/pemutil/badge.svg?branch=master&service=github)](https://coveralls.io/github/knq/pemutil?branch=master) #

A [Golang](https://golang.org/project) package that provides a light wrapper to
load PEM-encoded data, meant to ease the loading, parsing and decoding of PEM
data into standard [crypto](https://golang.org/pkg/crypto/) primitives.

## Installation ##

Install the package via the following:

    go get -u github.com/knq/pemutil

## Usage ##

Please see [the GoDoc API page](http://godoc.org/github.com/knq/pemutil) for a
full API listing.

The pemutil package can be used similarly to the following:

```go
// example/main.go
package main

//go:generate openssl genrsa -out rsa-private.pem 2048
//go:generate openssl rsa -in rsa-private.pem -outform PEM -pubout -out rsa-public.pem

import (
	"log"
	"os"

	"github.com/knq/pemutil"
)

func main() {
	var err error

	// create store and load our private key
	keyset, err := pemutil.LoadFile("rsa-private.pem")
	if err != nil {
		log.Fatal(err)
	}

	// do something with keyset.RSAPrivateKey()

	// get pem data and write to disk
	buf, err := keyset.Bytes()
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(buf)
}
```
