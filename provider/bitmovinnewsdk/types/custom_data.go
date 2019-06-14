package types

import "errors"

// CustomData holds data to set on encodings and configurations to hold relationship state
type CustomData *map[string]map[string]interface{}

// CustomDataStringValAtKeys is a helper for pulling values out of custom data
func CustomDataStringValAtKeys(data CustomData, primary, secondary string) (string, error) {
	if data == nil {
		return "", errors.New("custom data is nil")
	}

	d := *data

	raw, ok := d[primary][secondary]
	if !ok {
		return "", errors.New("no value found")
	}

	strVal, ok := raw.(string)
	if !ok {
		return "", errors.New("value is not a string")
	}

	return strVal, nil
}
