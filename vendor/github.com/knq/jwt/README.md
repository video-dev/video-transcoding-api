# About jwt [![Build Status](https://travis-ci.org/knq/jwt.svg)](https://travis-ci.org/knq/jwt) [![Coverage Status](https://coveralls.io/repos/knq/jwt/badge.svg?branch=master&service=github)](https://coveralls.io/github/knq/jwt?branch=master) #

A [Golang](https://golang.org/project) package that provides a simple and
secure way to encode and decode [JWT](https://jwt.io/) tokens.

## Installation ##

Install the package via the following:

    go get -u github.com/knq/jwt

Additionally, if you need to do command line encoding/decoding of JWTs, there
is a functional command line tool available:

    go get -u github.com/knq/jwt/cmd/jwt

## Usage ##

Please see [the GoDoc API page](http://godoc.org/github.com/knq/jwt) for a
full API listing.

The jwt package can be used similarly to the following:

```go
// example/main.go
package main

//go:generate openssl genrsa -out rsa-private.pem 2048
//go:generate openssl rsa -in rsa-private.pem -outform PEM -pubout -out rsa-public.pem

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/knq/jwt"
	"github.com/knq/pemutil"
)

func main() {
	var err error

	// load key
	keyset, err := pemutil.LoadFile("rsa-private.pem")
	if err != nil {
		log.Fatal(err)
	}

	// create PS384 using keyset
	// in addition, there are the other standard JWT encryption implementations:
	// HMAC:         HS256, HS384, HS512
	// RSA-PKCS1v15: RS256, RS384, RS512
	// ECC:          ES256, ES384, ES512
	// RSA-SSA-PSS:  PS256, PS384, PS512
	ps384, err := jwt.PS384.New(keyset)
	if err != nil {
		log.Fatal(err)
	}

	// calculate an expiration time
	expr := time.Now().Add(14 * 24 * time.Hour)

	// create claims using provided jwt.Claims
	c0 := jwt.Claims{
		Issuer:     "user@example.com",
		Audience:   "client@example.com",
		Expiration: json.Number(strconv.FormatInt(expr.Unix(), 10)),
	}
	fmt.Printf("Claims: %+v\n\n", c0)

	// encode token
	buf, err := ps384.Encode(&c0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Encoded token:\n\n%s\n\n", string(buf))

	// decode the generated token and verify
	c1 := jwt.Claims{}
	err = ps384.Decode(buf, &c1)
	if err != nil {
		// if the signature was bad, the err would not be nil
		log.Fatal(err)
	}
	if reflect.DeepEqual(c0, c1) {
		fmt.Printf("Claims Match! Decoded claims: %+v\n\n", c1)
	}

	fmt.Println("----------------------------------------------")

	// use custom claims
	c3 := map[string]interface{}{
		"aud": "my audience",
		"http://example/api/write": true,
	}
	fmt.Printf("My Custom Claims: %+v\n\n", c3)

	buf, err = ps384.Encode(&c3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Encoded token with custom claims:\n\n%s\n\n", string(buf))

	// decode custom claims
	c4 := myClaims{}
	err = ps384.Decode(buf, &c4)
	if err != nil {
		log.Fatal(err)
	}
	if c4.Audience == "my audience" {
		fmt.Printf("Decoded custom claims: %+v\n\n", c1)
	}
	if c4.WriteScope {
		fmt.Println("myClaims custom claims has write scope!")
	}
}

// or use any type that the standard encoding/json library recognizes as a
// payload (you can also "extend" jwt.Claims in this fashion):
type myClaims struct {
	jwt.Claims
	WriteScope bool `json:"http://example/api/write"`
}
```

The command line tool can be used as follows (assuming jwt is somewhere on $PATH):

```sh
# encode arbitrary JSON as payload (ie, claims)
echo '{"iss": "issuer", "nbf": '$(date +%s)'}' | jwt -k ./testdata/rsa.pem -enc

# quick encode name/value pairs from command line
jwt -k ./testdata/rsa.pem -enc iss=issuer nbf=$(date +%s)

# decode (and verify) token
echo "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg" | jwt -k ./testdata/rsa.pem -dec

# encode and decode in one sweep:
jwt -k ./testdata/rsa.pem -enc iss=issuer nbf=$(date +%s) | jwt -k ./testdata/rsa.pem -dec

# specify algorithm -- this will error since the token here is encoded using RS256, not RS384
echo "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg" | jwt -k ./testdata/rsa.pem -dec -alg RS384
```
