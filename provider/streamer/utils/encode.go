package utils

import (
	"encoding/base64"
	"encoding/json"
)

/*
	@function: EncodeBase64
	@description:
		encode any input object into base64-encoded string
*/
func EncodeBase64(obj interface{}) (string, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

/*
	@function: DecodeBase64
	@description:
		decode base64-encoded string into an object
*/
func DecodeBase64(in string, obj interface{}) error {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, obj)
	if err != nil {
		return err
	}

	return nil
}
