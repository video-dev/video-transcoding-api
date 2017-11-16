// Package jwt provides a simplified and secure API for encoding, decoding and
// verifying JSON Web Tokens (JWT).
//
// See https://jwt.io/ and https://tools.ietf.org/html/rfc7519
//
// The API is designed to be instantly familiar to users of the standard crypto
// and json packages:
//
//		// create jwt.Signer from the RSA key on disk
//		rs384 := jwt.RS384.New(jwt.PEM{"myrsakey.pem"})
//
//		// create claims
//		claims := jwt.Claims{
//			Issuer: "user@example.com",
//		}
//
//		// encode claims as a JWT:
//		buf, err := rs384.Encode(&claims)
//		if err != nil {
//			log.Fatalln(err)
//		}
//		fmt.Printf("token: %s\n", string(buf))
//
// 		// decode and verify claims:
//		cl2 := jwt.Claims{}
//		err = rs384.Decode(buf, &cl2)
//		if err == jwt.ErrInvalidSignature {
//			// invalid signature
//		} else if err != nil {
//			// handle general error
//		}
//
//		fmt.Printf("decoded claims: %+v\n", cl2)
//
package jwt

import (
	"bytes"
	"encoding/json"
	"reflect"
)

var (
	// tokenSep is the token separator.
	tokenSep = []byte{'.'}
)

// Decode decodes a serialized JWT in buf into obj, and verifies the JWT
// signature using the Algorithm and Signer.
//
// If the token or signature is invalid, ErrInvalidToken or ErrInvalidSignature
// will be returned, respectively. Otherwise, any other errors encountered
// during token decoding will be returned.
func Decode(alg Algorithm, signer Signer, buf []byte, obj interface{}) error {
	var err error

	// split token
	ut := UnverifiedToken{}
	err = DecodeUnverifiedToken(buf, &ut)
	if err != nil {
		return err
	}

	// verify signature
	sig, err := signer.Verify(buf[:len(ut.Header)+len(tokenSep)+len(ut.Payload)], ut.Signature)
	if err != nil {
		return ErrInvalidSignature
	}

	// b64 decode header
	headerBuf, err := b64.DecodeString(string(ut.Header))
	if err != nil {
		return err
	}

	// json decode header
	header := Header{}
	err = json.Unmarshal(headerBuf, &header)
	if err != nil {
		return err
	}

	// verify alg matches header algorithm
	if alg != header.Algorithm {
		return ErrInvalidAlgorithm
	}

	// set header in the provided obj
	err = decodeToObjOrFieldWithTag(headerBuf, obj, "header", &header)
	if err != nil {
		return err
	}

	// b64 decode payload
	payloadBuf, err := b64.DecodeString(string(ut.Payload))
	if err != nil {
		return err
	}

	// json decode payload
	payload := Claims{}
	err = json.Unmarshal(payloadBuf, &payload)
	if err != nil {
		return err
	}

	// set payload in the provided obj
	err = decodeToObjOrFieldWithTag(payloadBuf, obj, "payload", &payload)
	if err != nil {
		return err
	}

	// set sig in the provided obj
	field := getFieldWithTag(obj, "signature")
	if field != nil {
		field.Set(reflect.ValueOf(sig))
	}

	return nil
}

// Encode encodes a JWT using the Algorithm and Signer, returning the URL-safe
// encoded token or any errors encountered during encoding.
func Encode(alg Algorithm, signer Signer, obj interface{}) ([]byte, error) {
	var err error

	// grab encode targets
	headerObj, payloadObj, err := encodeTargets(alg, obj)
	if err != nil {
		return nil, err
	}

	// json encode header
	header, err := json.Marshal(headerObj)
	if err != nil {
		return nil, err
	}

	// b64 encode playload
	headerEnc := make([]byte, b64.EncodedLen(len(header)))
	b64.Encode(headerEnc, header)

	// json encode payload
	payload, err := json.Marshal(payloadObj)
	if err != nil {
		return nil, err
	}

	// b64 encode playload
	payloadEnc := make([]byte, b64.EncodedLen(len(payload)))
	b64.Encode(payloadEnc, payload)

	// allocate result
	var buf bytes.Buffer

	// add header
	_, err = buf.Write(headerEnc)
	if err != nil {
		return nil, err
	}

	// add 1st separator
	_, err = buf.Write(tokenSep)
	if err != nil {
		return nil, err
	}

	// add payload
	_, err = buf.Write(payloadEnc)
	if err != nil {
		return nil, err
	}

	// sign
	sig, err := signer.Sign(buf.Bytes())
	if err != nil {
		return nil, err
	}

	// add 2nd separator
	_, err = buf.Write(tokenSep)
	if err != nil {
		return nil, err
	}

	// add sig
	_, err = buf.Write(sig)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// PeekHeaderField extracts the specified field from the serialized JWT buf's
// header. An error will be returned if the field is not present in the decoded
// header.
func PeekHeaderField(buf []byte, field string) (string, error) {
	return peekField(buf, field, tokenPositionHeader)
}

// PeekPayloadField extracts the specified field from the serialized JWT buf's
// payload (ie, the token claims). An error will be returned if the field is
// not present in the decoded payload.
func PeekPayloadField(buf []byte, field string) (string, error) {
	return peekField(buf, field, tokenPositionPayload)
}

// PeekAlgorithm extracts the signing algorithm listed in the "alg" field of
// the serialized JWT buf's header and attempts to unmarshal it into an
// Algorithm. An error will be returned if the alg field is not specified in
// the JWT header, or is otherwise invalid.
func PeekAlgorithm(buf []byte) (Algorithm, error) {
	alg := NONE

	// get alg
	algVal, err := PeekHeaderField(buf, "alg")
	if err != nil {
		return NONE, err
	}

	// decode alg
	err = (&alg).UnmarshalText([]byte(algVal))
	if err != nil {
		return NONE, err
	}

	return alg, nil
}

// PeekAlgorithmAndIssuer extracts the signing algorithm listed in the "alg"
// field and the issuer from the "iss" field of the serialized JWT buf's header
// and payload, attempting to unmarshal alg to Algorithm and iss to a string.
// An error will be returned if the Algorithm or Issuer fields are not
// specified in the JWT header and payload, or are otherwise invalid.
func PeekAlgorithmAndIssuer(buf []byte) (Algorithm, string, error) {
	var err error

	// get algorithm
	alg, err := PeekAlgorithm(buf)
	if err != nil {
		return NONE, "", err
	}

	// get issuer
	issuer, err := PeekPayloadField(buf, "iss")
	if err != nil {
		return NONE, "", err
	}

	return alg, issuer, nil
}
