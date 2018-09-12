package jwt

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
)

var (
	// b64 is the base64 encoding config used for encoding/decoding jwt
	// parts.
	b64 = base64.URLEncoding.WithPadding(base64.NoPadding)
)

// getFieldWithTag lookups jwt tag, with specified tagName on obj, returning
// its reflected value.
func getFieldWithTag(obj interface{}, tagName string) *reflect.Value {
	objVal := reflect.ValueOf(obj)
	if objVal.Kind() != reflect.Struct {
		objVal = objVal.Elem()
	}

	for i := 0; i < objVal.NumField(); i++ {
		fieldType := objVal.Type().Field(i)
		if tagName == fieldType.Tag.Get("jwt") {
			field := objVal.Field(i)
			return &field
		}
	}

	return nil
}

// decodeToObjOrFieldWithTag decodes the buf into obj's field having the
// specified jwt tagName. If the provided obj's has the same type as
// defaultObj, then the obj is set to the defaultObj, otherwise an attempt is
// made to json.Decode the buf into obj.
func decodeToObjOrFieldWithTag(buf []byte, obj interface{}, tagName string, defaultObj interface{}) error {
	// reflect values
	objValElem := reflect.ValueOf(obj).Elem()
	defaultObjValElem := reflect.ValueOf(defaultObj).Elem()

	// first check type, if same type, then set
	if objValElem.Type() == defaultObjValElem.Type() {
		objValElem.Set(defaultObjValElem)
		return nil
	}

	// get field with specified jwt tagName (if any)
	fieldVal := getFieldWithTag(obj, tagName)
	if fieldVal != nil {
		// check field type and defaultObj type, if same, set
		if fieldVal.Type() == defaultObjValElem.Type() {
			fieldVal.Set(defaultObjValElem)
			return nil
		}

		// otherwise, assign obj address of field
		obj = fieldVal.Addr().Interface()
	}

	// decode json
	d := json.NewDecoder(bytes.NewBuffer(buf))
	d.UseNumber()
	return d.Decode(obj)
}

// grabEncodeTargets grabs the fields for the obj.
func grabEncodeTargets(alg Algorithm, obj interface{}) (interface{}, interface{}, error) {
	var headerObj, payloadObj interface{}

	// get header
	if headerVal := getFieldWithTag(obj, "header"); headerVal != nil {
		headerObj = headerVal.Interface()
	}
	if headerObj == nil {
		headerObj = alg.Header()
	}

	// get payload
	if payloadVal := getFieldWithTag(obj, "payload"); payloadVal != nil {
		payloadObj = payloadVal.Interface()
	}
	if payloadObj == nil {
		payloadObj = obj
	}

	return headerObj, payloadObj, nil
}

// encodeTargets determines what to encode.
func encodeTargets(alg Algorithm, obj interface{}) (interface{}, interface{}, error) {
	// determine what to encode
	switch val := obj.(type) {
	case *Token:
		return val.Header, val.Payload, nil
	}

	objVal := reflect.ValueOf(obj)
	objKind := objVal.Kind()
	if objKind == reflect.Struct || (objKind == reflect.Ptr && objVal.Elem().Kind() == reflect.Struct) {
		return grabEncodeTargets(alg, obj)
	}

	return alg.Header(), obj, nil
}

// tokenPosition is the different positions of the constituent JWT parts.
//
// Used in conjunction with peekField.
type tokenPosition int

const (
	tokenPositionHeader tokenPosition = iota
	tokenPositionPayload

//	tokenPositionSignature
)

// peekField looks at an undecoded JWT, JSON decoding the data at pos, and
// returning the specified field's value as string.
//
// If the fieldName is not present, then an error will be returned.
func peekField(buf []byte, fieldName string, pos tokenPosition) (string, error) {
	var err error

	// split token
	ut := UnverifiedToken{}
	err = DecodeUnverifiedToken(buf, &ut)
	if err != nil {
		return "", err
	}

	// determine position decode
	var typ string
	var b []byte
	switch pos {
	case tokenPositionHeader:
		typ = "header"
		b = ut.Header
	case tokenPositionPayload:
		typ = "payload"
		b = ut.Payload

	default:
		return "", fmt.Errorf("invalid field %d", pos)
	}

	// b64 decode
	dec, err := b64.DecodeString(string(b))
	if err != nil {
		return "", fmt.Errorf("could not decode token %s", typ)
	}

	// json decode
	m := make(map[string]interface{})
	err = json.Unmarshal(dec, &m)
	if err != nil {
		return "", err
	}

	if val, ok := m[fieldName]; ok {
		return fmt.Sprintf("%v", val), nil
	}

	return "", fmt.Errorf("token %s field %s not present or invalid", typ, fieldName)
}
